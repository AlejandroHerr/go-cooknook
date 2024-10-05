package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlejandroHerr/cook-book-go/internal/core"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IngredientsRepository struct {
	pool *pgxpool.Pool
}

func NewIngredientsRepository(pool *pgxpool.Pool) *IngredientsRepository {
	return &IngredientsRepository{
		pool: pool,
	}
}

func (r IngredientsRepository) Find(id uuid.UUID) (*model.Ingredient, error) {
	return r.findBy("id", id.String())
}

func (r IngredientsRepository) FindByName(name string) (*model.Ingredient, error) {
	return r.findBy("name", name)
}

func (r IngredientsRepository) findBy(field string, value string) (*model.Ingredient, error) {
	row := r.pool.QueryRow(context.TODO(), "SELECT * from ingredients WHERE "+field+" = $1", value)
	ingredientDB := struct {
		ID   uuid.UUID
		Name string
		Kind *string
	}{}

	err := row.Scan(&ingredientDB.ID, &ingredientDB.Name, &ingredientDB.Kind)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &core.NotFoundError{Field: field, Value: value, Err: err}
		}

		return nil, fmt.Errorf("error scanning ingredient: %w", err)
	}

	ingredient := model.NewIngredient(ingredientDB.ID, ingredientDB.Name, ingredientDB.Kind)

	return &ingredient, nil
}

func (r IngredientsRepository) Save(ingredient model.Ingredient) error {
	values := pgx.NamedArgs{
		"id":   ingredient.ID(),
		"name": ingredient.Name(),
		"kind": ingredient.Kind(),
	}

	_, err := r.pool.Exec(
		context.TODO(),
		"INSERT INTO ingredients (id, name, kind) VALUES (@id,@name,@kind)",
		values,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return &core.DuplicateKeyError{Key: pgErr.ConstraintName, Err: err}
			case "23503": // foreign_key_violation
				return &core.ConstraintError{Constraint: pgErr.ConstraintName, Err: err}
			}
		}

		return fmt.Errorf("error saving ingredient: %w", err)
	}

	return nil
}
