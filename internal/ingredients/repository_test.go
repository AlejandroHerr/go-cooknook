package ingredients_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/AlejandroHerr/cook-book-go/internal/infra"
	"github.com/AlejandroHerr/cook-book-go/internal/ingredients"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	// "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq" // add this
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	user := "tests"
	password := "123456"
	pgConfig := infra.Config{
		Host:     "localhost",
		Port:     5433,
		Database: "tests",
		User:     &user,
		Password: &password,
	}

	pool, err := infra.Connect(&pgConfig)
	if err != nil {
		t.Fatalf("Could not connect to the database: %v", err)
	}
	// Setup the test database
	_, err = pool.Exec(context.TODO(), "CREATE TABLE IF NOT EXISTS ingredients (id UUID PRIMARY KEY, name TEXT NOT NULL, kind TEXT)")
	if err != nil {
		t.Fatalf("Could not create the table: %v", err)
	}

	return pool
}

func TestCreateIngredient(t *testing.T) {
	t.Parallel()
	// Setup the test database
	pool := setupTestDB(t)

	ingredient := ingredients.NewIngredient(uuid.New(), "Tomato", nil)

	repo := ingredients.NewIngredientsRepository(pool)

	err := repo.SaveIngredient(ingredient)

	assert.NoError(t, err, "Could not save the ingredient")

	foundIngredient, err := repo.GetIngredient(ingredient.ID())

	assert.NoError(t, err, "Could not get the ingredient")

	assert.Equal(t, ingredient.ID(), foundIngredient.ID(), "The ID of the ingredient is not the same")
	assert.Equal(t, ingredient.Name(), foundIngredient.Name(), "The name of the ingredient is not the same")
	assert.Equal(t, ingredient.Kind(), foundIngredient.Kind(), "The kind of the ingredient is not the same")
}

func TestIngredientNotFound(t *testing.T) {
	t.Parallel()
	pool := setupTestDB(t)

	repo := ingredients.NewIngredientsRepository(pool)

	_, err := repo.GetIngredient(uuid.New())
	fmt.Println(err)
	// var pgErr *pgconn.PgError
	assert.ErrorIs(t, err, pgx.ErrNoRows, fmt.Sprintf("Expected a pgx.ErrNoRows, got %v", err))
	// assert.ErrorAs(t, err, &pgErr, fmt.Sprintf("Expected a pgx.PgError, got %v", err))
}
