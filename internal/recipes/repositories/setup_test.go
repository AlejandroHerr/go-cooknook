package repositories_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/AlejandroHerr/cook-book-go/internal/infra"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/fixtures"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testFixtures fixtures.Fixtures
	testPool     *pgxpool.Pool
)

func TestMain(m *testing.M) {
	user := "tests"
	password := "123456"
	config := infra.Config{
		Host:     "localhost",
		Port:     5433,
		Database: "tests",
		User:     &user,
		Password: &password,
	}

	var err error

	testPool, err = infra.Connect(&config)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	defer teardown()

	err = fixtures.Cleanup(testPool)
	if err != nil {
		log.Fatalf("Could not clean up the database: %v", err)
	}

	testFixtures, err = fixtures.GenerateAndLoadFixtures(testPool, 20)
	if err != nil {
		log.Panicf("Could not load the fixtures: %v", err)
	}

	fmt.Println(testFixtures)
	code := m.Run()

	teardown()

	os.Exit(code)
}

func cleanUp() {
	fixtures.Cleanup(testPool)
}

func teardown() {
	cleanUp()
	testPool.Close()
}
