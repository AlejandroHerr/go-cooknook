package repositories_test

import (
	"testing"

	"github.com/AlejandroHerr/cook-book-go/internal/core"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/model"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/repositories"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRecipesRepository_Find_ReturnsRecipe(t *testing.T) {
	recipe := testFixtures.Recipes[0]

	repo := repositories.NewRecipesRepository(testPool)

	foundRecipe, err := repo.Find(recipe.ID())
	assert.NoError(t, err, "should not return an error")

	assert.Equal(t, recipe, *foundRecipe, "should return the same recipe")
}

func TestRecipesRepository_Find_ReturnsNotFound(t *testing.T) {
	repo := repositories.NewRecipesRepository(testPool)

	foundRecipe, err := repo.Find(uuid.New())
	assert.Nil(t, foundRecipe, "should return nil")

	var notFoundError *core.NotFoundError

	assert.ErrorAs(t, err, &notFoundError, "should return a NotFoundError")
}

func TestRecipesRepository_FindByName_ReturnsRecipe(t *testing.T) {
	recipe := testFixtures.Recipes[0]

	repo := repositories.NewRecipesRepository(testPool)

	foundRecipe, err := repo.FindByName(recipe.Name())
	assert.NoError(t, err, "should not return an error")

	assert.Equal(t, recipe, *foundRecipe, "should return the same recipe")
}

func TestRecipesRepository_FindByName_ReturnsNotFound(t *testing.T) {
	repo := repositories.NewRecipesRepository(testPool)

	foundRecipe, err := repo.FindByName(gofakeit.Word())
	assert.Nil(t, foundRecipe, "should return nil")

	var notFoundError *core.NotFoundError

	assert.ErrorAs(t, err, &notFoundError, "should return a NotFoundError")
}

func TestRecipesRepository_Save_SavesRecipe(t *testing.T) {
	recipe := model.NewRecipe(uuid.New(), gofakeit.Adjective()+" "+gofakeit.Word(), nil, nil)

	recipe.AddIngredient(testFixtures.Ingredients[5])
	recipe.AddIngredient(testFixtures.Ingredients[6])

	repo := repositories.NewRecipesRepository(testPool)

	err := repo.Save(recipe)

	assert.NoError(t, err, "no error")

	foundRecipe, err := repo.Find(recipe.ID())
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, recipe, *foundRecipe, "should return the same recipe")
}

func TestRecipesRepository_Save_Failes0(t *testing.T) {
	recipe := model.NewRecipe(uuid.New(), testFixtures.Recipes[0].Name(), nil, nil)

	repo := repositories.NewRecipesRepository(testPool)

	err := repo.Save(recipe)

	var duplicateKeyError *core.DuplicateKeyError

	assert.ErrorAs(t, err, &duplicateKeyError, "should return a DuplicateKeyError")
}

func TestRecipesRepository_Save_Failes(t *testing.T) {
	recipe := model.NewRecipe(uuid.New(), gofakeit.Adjective()+" "+gofakeit.Word(), nil, nil)

	recipe.AddIngredient(testFixtures.Ingredients[5])
	recipe.AddIngredient(testFixtures.Ingredients[6])
	recipe.AddIngredient(model.NewIngredient(uuid.New(), gofakeit.Word(), nil))

	repo := repositories.NewRecipesRepository(testPool)

	err := repo.Save(recipe)

	var constraintError *core.ConstraintError

	assert.ErrorAs(t, err, &constraintError, "should return a ConstraintError")
}
