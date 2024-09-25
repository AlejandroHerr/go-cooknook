package ingredients

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *IngredientsRepository) SaveIngredient(ingredient *Ingredient) error {
	values := pgx.NamedArgs{
		"id":   ingredient.ID(),
		"name": ingredient.Name(),
		"kind": ingredient.Kind(),
	}

	_, err := r.pool.Exec(context.TODO(), "INSERT INTO ingredients (id, name, kind) VALUES (@id,@name,@kind)", values)
	if err != nil {
		return fmt.Errorf("could not save the ingredient: %w", err)
	}

	return nil
}

func (r *IngredientsRepository) GetIngredient(id uuid.UUID) (*Ingredient, error) {
	row := r.pool.QueryRow(context.TODO(), "SELECT * from ingredients WHERE id = $1", id)

	ingredientDB := struct {
		ID   uuid.UUID
		Name string
		Kind *string
	}{}

	err := row.Scan(&ingredientDB.ID, &ingredientDB.Name, &ingredientDB.Kind)
	if err != nil {
		return nil, fmt.Errorf("could not get the ingredient: %w", err)
	}

	ingredient := NewIngredient(ingredientDB.ID, ingredientDB.Name, ingredientDB.Kind)

	return ingredient, nil
}
