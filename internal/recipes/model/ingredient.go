package model

import (
	"encoding/json"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

type Ingredient struct {
	id   uuid.UUID
	name string
	kind *string
}

func (i Ingredient) ID() uuid.UUID {
	return i.id
}

func (i Ingredient) Name() string {
	return i.name
}

func (i Ingredient) Kind() *string {
	return i.kind
}

func NewIngredient(uuid uuid.UUID, name string, kind *string) Ingredient {
	return Ingredient{
		id:   uuid,
		name: name,
		kind: kind,
	}
}

func (i Ingredient) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"id":   i.id.String(),
		"name": i.name,
	}

	if i.kind != nil {
		data["kind"] = *i.kind
	}

	res, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling ingredient: %w", err)
	}

	return res, nil
}

func (i Ingredient) Fake(faker *gofakeit.Faker) (any, error) {
	name := faker.AdjectiveDescriptive() + " " + faker.NounConcrete()
	kind := faker.NounAbstract()

	return NewIngredient(uuid.New(), name, &kind), nil
}
