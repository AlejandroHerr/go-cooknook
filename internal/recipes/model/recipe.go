package model

import (
	"encoding/json"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

type Recipe struct {
	id          uuid.UUID
	name        string
	description *string
	ingredients []Ingredient
}

func NewRecipe(id uuid.UUID, name string, description *string, ingredients []Ingredient) Recipe {
	if ingredients == nil {
		ingredients = make([]Ingredient, 0)
	}

	return Recipe{
		id:          id,
		name:        name,
		description: description,
		ingredients: ingredients,
	}
}

func (r Recipe) ID() uuid.UUID {
	return r.id
}

func (r Recipe) Name() string {
	return r.name
}

func (r Recipe) Description() *string {
	return r.description
}

func (r Recipe) Ingredients() []Ingredient {
	return r.ingredients
}

func (r *Recipe) AddIngredient(ingredient Ingredient) {
	r.ingredients = append(r.ingredients, ingredient)
}

func (r *Recipe) AddIngredients(ingredients []Ingredient) {
	r.ingredients = append(r.ingredients, ingredients...)
}

func (r Recipe) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"id":          r.id.String(),
		"name":        r.name,
		"ingredients": r.ingredients,
	}

	res, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling recipe: %w", err)
	}

	return res, nil
}

func (r Recipe) Fake(faker *gofakeit.Faker) (any, error) {
	name := faker.Adjective() + " " + faker.Dinner()
	descrption := faker.LoremIpsumParagraph(1, 3, 6, " ")

	return NewRecipe(uuid.New(), name, &descrption, nil), nil
}
