package repositories_test

import (
	"testing"

	"github.com/AlejandroHerr/cook-book-go/internal/core"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/fixtures"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/model"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/repositories"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq" // add this
)

func TestIngredientRepository_Find_ReturnsTheIngredient(t *testing.T) {
	t.Parallel()

	ingredient := testFixtures.Ingredients[0]

	repo := repositories.NewIngredientsRepository(testPool)

	foundIngredient, err := repo.Find(ingredient.ID())

	assert.NoError(t, err, "Could not get the ingredient")
	assert.Equal(t, ingredient, *foundIngredient, "should return the same ingredient")
}

func TestIngredientRepository_Find_ReturnsNotFoundError(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	repo := repositories.NewIngredientsRepository(testPool)

	_, err := repo.Find(id)

	var notFoundError *core.NotFoundError

	assert.ErrorAs(t, err, &notFoundError, "Expected a core.NotFoundError")
}

func TestIngredientRepository_FindByName_ReturnsTheIngredient(t *testing.T) {
	t.Parallel()

	ingredient := testFixtures.Ingredients[0]

	repo := repositories.NewIngredientsRepository(testPool)

	foundIngredient, err := repo.FindByName(ingredient.Name())

	assert.NoError(t, err, "Could not get the ingredient")
	assert.Equal(t, ingredient, *foundIngredient, "should return the same ingredient")
}

func TestIngredientRepository_FindByName_ReturnsNotFoundError(t *testing.T) {
	t.Parallel()

	repo := repositories.NewIngredientsRepository(testPool)

	_, err := repo.FindByName("Not an ingredient")

	var notFoundError *core.NotFoundError

	assert.ErrorAs(t, err, &notFoundError, "Expected a core.NotFoundError")
}

func TestIngredientRespository_Save_CreatesIngredientInDB(t *testing.T) {
	t.Parallel()

	ingredient, err := fixtures.FakeIngredient()
	if err != nil {
		t.Fatal(err)
	}

	repo := repositories.NewIngredientsRepository(testPool)

	err = repo.Save(ingredient)

	assert.NoError(t, err, "Could not save the ingredient")

	foundIngredient, err := repo.Find(ingredient.ID())

	assert.NoError(t, err, "Could not get the ingredient")
	assert.Equal(t, ingredient, *foundIngredient, "should return the same ingredient")
}

func TestIngredientRespositoryi_Save_CreatesIngredientWithRequiredFields(t *testing.T) {
	t.Parallel()

	ingredient := model.NewIngredient(uuid.New(), "Tomato", nil)

	repo := repositories.NewIngredientsRepository(testPool)

	err := repo.Save(ingredient)

	assert.NoError(t, err, "Could not save the ingredient")

	foundIngredient, err := repo.Find(ingredient.ID())

	assert.NoError(t, err, "Could not get the ingredient")

	assert.Equal(t, ingredient.ID(), foundIngredient.ID(), "The ID of the ingredient is not the same")
	assert.Equal(t, ingredient.Name(), foundIngredient.Name(), "The name of the ingredient is not the same")
	assert.Nil(t, foundIngredient.Kind(), "The kind of the ingredient is not the same")
}

func TestIngredientRespository_Save_ReturnsDuplicateKeyError(t *testing.T) {
	t.Parallel()

	repo := repositories.NewIngredientsRepository(testPool)
	ingredientA := model.NewIngredient(uuid.New(), "Broccoli", nil)

	err := repo.Save(ingredientA)

	assert.NoError(t, err, "Could not save the ingredient")

	ingredientB := model.NewIngredient(uuid.New(), "Broccoli", nil)

	err = repo.Save(ingredientB)

	var duplicateKeyError *core.DuplicateKeyError

	assert.ErrorAs(t, err, &duplicateKeyError, "Expected a core.ConstraintError")
}
