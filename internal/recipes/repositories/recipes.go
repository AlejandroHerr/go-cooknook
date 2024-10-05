package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlejandroHerr/cook-book-go/internal/core"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RecipesRepository struct {
	pool *pgxpool.Pool
}

func NewRecipesRepository(pool *pgxpool.Pool) *RecipesRepository {
	return &RecipesRepository{
		pool: pool,
	}
}

func (r RecipesRepository) Find(id uuid.UUID) (*model.Recipe, error) {
	return r.findBy("id", id.String())
}

func (r RecipesRepository) FindByName(name string) (*model.Recipe, error) {
	return r.findBy("name", name)
}

func (r RecipesRepository) findBy(field string, value string) (*model.Recipe, error) {
	row := r.pool.QueryRow(context.TODO(), "SELECT * from recipes WHERE "+field+" = $1", value)
	recipeDB := struct {
		ID          uuid.UUID
		Name        string
		Description *string
	}{}

	err := row.Scan(&recipeDB.ID, &recipeDB.Name, &recipeDB.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &core.NotFoundError{Field: field, Value: value, Err: err}
		}

		return nil, fmt.Errorf("error scanning recipe: %w", err)
	}

	ingredientsList, err := r.loadIngredients(recipeDB.ID)
	if err != nil {
		return nil, fmt.Errorf("error loading ingredients: %w", err)
	}

	recipe := model.NewRecipe(recipeDB.ID, recipeDB.Name, recipeDB.Description, ingredientsList)

	return &recipe, nil
}

func (r RecipesRepository) loadIngredients(recipeID uuid.UUID) ([]model.Ingredient, error) {
	query := `
    select i.id, i.name, i.kind 
    from recipe_ingredients ri
    inner JOIN ingredients i on ri.ingredient_id = i.id 
    where ri.recipe_id = $1
  `

	rows, err := r.pool.Query(
		context.TODO(),
		query,
		recipeID,
	)
	if err != nil {
		return nil, fmt.Errorf("error loading ingredients: %w", err)
	}

	ingredientsList := make([]model.Ingredient, 0)

	for rows.Next() {
		ingredientDB := struct {
			ID   uuid.UUID
			Name string
			Kind *string
		}{}

		err = rows.Scan(&ingredientDB.ID, &ingredientDB.Name, &ingredientDB.Kind)
		if err != nil {
			return nil, fmt.Errorf("error scanning ingredient: %w", err)
		}

		ingredient := model.NewIngredient(ingredientDB.ID, ingredientDB.Name, ingredientDB.Kind)
		ingredientsList = append(ingredientsList, ingredient)
	}

	return ingredientsList, nil
}

func (r RecipesRepository) Save(recipe model.Recipe) error {
	values := pgx.NamedArgs{
		"id":          recipe.ID(),
		"name":        recipe.Name(),
		"description": recipe.Description(),
	}

	_, err := r.pool.Exec(
		context.TODO(),
		"INSERT INTO recipes (id, name, description) VALUES (@id,@name,@description)",
		values,
	)
	if err != nil {
		return core.HandleError(err, "error saving recipe")
	}

	for _, ingredient := range recipe.Ingredients() {
		_, err := r.pool.Exec(
			context.TODO(),
			"INSERT INTO recipe_ingredients (recipe_id, ingredient_id) VALUES ($1, $2)",
			recipe.ID(),
			ingredient.ID(),
		)
		if err != nil {
			return core.HandleError(err, "error saving recipe ingredients")
		}
	}

	return nil
}

//
// func (r Repository) Save(recipe model.Recipe) error {
// 	values := pgx.NamedArgs{
// 		"id":   recipe.ID(),
// 		"name": recipe.Name(),
// 		"kind": recipe.Kind(),
// 	}
//
// 	_, err := r.pool.Exec(
// 		context.TODO(),
// 		"INSERT INTO recipes (id, name, kind) VALUES (@id,@name,@kind)",
// 		values,
// 	)
// 	if err != nil {
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			switch pgErr.Code {
// 			case "23505": // unique_violation
// 				return &core.DuplicateKeyError{Key: pgErr.ConstraintName, Err: err}
// 			case "23503": // foreign_key_violation
// 				return &core.ConstraintError{Constraint: pgErr.ConstraintName, Err: err}
// 			}
// 		}
//
// 		return fmt.Errorf("error saving recipe: %w", err)
// 	}
//
// 	return nil
// }
