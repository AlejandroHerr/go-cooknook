package fixtures

import (
	"context"
	"fmt"

	"github.com/AlejandroHerr/cook-book-go/internal/recipes/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Fixtures struct {
	Ingredients []model.Ingredient
	Recipes     []model.Recipe
}

func FakeIngredient() (model.Ingredient, error) {
	var ingredient model.Ingredient

	err := gofakeit.Struct(&ingredient)
	if err != nil {
		return model.Ingredient{}, fmt.Errorf("could not generate a fake ingredient: %w", err)
	}

	return ingredient, nil
}

func FakeRecipe() (model.Recipe, error) {
	var recipe model.Recipe

	err := gofakeit.Struct(&recipe)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("could not generate a fake recipe: %w", err)
	}

	fmt.Println(recipe.Name())

	return recipe, nil
}

func generateIngredientFixtures(number int) ([]model.Ingredient, error) {
	ingredients := make([]model.Ingredient, number)

	for i := 0; i < number; i++ {
		ingredient, err := FakeIngredient()
		if err != nil {
			return nil, err
		}

		ingredients[i] = ingredient
	}

	return ingredients, nil
}

func GenerateFixtures(number int) (Fixtures, error) {
	ingredients, err := generateIngredientFixtures(number * 2)
	if err != nil {
		return Fixtures{}, fmt.Errorf("could not generate the ingredients fixtures: %w", err)
	}

	recipesList := make([]model.Recipe, number)

	for i := 0; i < number; i++ {
		recipe, err := FakeRecipe()
		if err != nil {
			return Fixtures{}, err
		}

		if recipe.Ingredients() != nil {
			recipe.AddIngredients([]model.Ingredient{
				ingredients[(i*2)%len(ingredients)],
				ingredients[(i*2+1)%len(ingredients)],
			})
		}

		recipesList[i] = recipe
	}

	return Fixtures{
		Ingredients: ingredients,
		Recipes:     recipesList,
	}, nil
}

func LoadFixtures(pool *pgxpool.Pool, fixtures Fixtures) error {
	for _, ingredient := range fixtures.Ingredients {
		_, err := pool.Exec(
			context.Background(),
			"INSERT INTO ingredients (id, name, kind) VALUES ($1, $2, $3)",
			ingredient.ID(),
			ingredient.Name(),
			ingredient.Kind(),
		)
		if err != nil {
			return fmt.Errorf("could not insert the ingredient fixture: %w", err)
		}
	}

	for _, recipe := range fixtures.Recipes {
		_, err := pool.Exec(
			context.Background(),
			"INSERT INTO recipes (id, name, description) VALUES ($1, $2, $3)",
			recipe.ID(),
			recipe.Name(),
			recipe.Description(),
		)
		if err != nil {
			return fmt.Errorf("could not insert the recipe fixture: %w", err)
		}

		for _, ingredient := range recipe.Ingredients() {
			_, err = pool.Exec(
				context.Background(),
				"INSERT INTO recipe_ingredients (recipe_id, ingredient_id) VALUES ($1, $2)",
				recipe.ID(),
				ingredient.ID(),
			)
			if err != nil {
				return fmt.Errorf("could not insert the recipe_ingredient fixture: %w", err)
			}
		}
	}

	return nil
}

func GenerateAndLoadFixtures(pool *pgxpool.Pool, number int) (Fixtures, error) {
	fixtures, err := GenerateFixtures(number)
	if err != nil {
		return Fixtures{}, err
	}

	err = LoadFixtures(pool, fixtures)
	if err != nil {
		return Fixtures{}, err
	}

	return fixtures, nil
}

func Cleanup(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), "DELETE FROM recipe_ingredients")
	if err != nil {
		return fmt.Errorf("could not clean the recipe_ingredients table: %w", err)
	}

	_, err = pool.Exec(context.Background(), "DELETE FROM recipes")
	if err != nil {
		return fmt.Errorf("could not clean the recipes table: %w", err)
	}

	_, err = pool.Exec(context.Background(), "DELETE FROM ingredients")
	if err != nil {
		return fmt.Errorf("could not clean the ingredients table: %w", err)
	}

	return nil
}
