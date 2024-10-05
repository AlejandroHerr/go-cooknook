package services

import (
	"errors"
	"fmt"

	"github.com/AlejandroHerr/cook-book-go/internal/core"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/model"
	"github.com/AlejandroHerr/cook-book-go/internal/recipes/repositories"
	"github.com/google/uuid"
)

type RecipesServices struct {
	recipesRepository     *repositories.RecipesRepository
	ingredientsRepository *repositories.IngredientsRepository
}

func (s RecipesServices) getOrCreateIngredient(name string) (*model.Ingredient, error) {
	ingredient, err := s.ingredientsRepository.FindByName(name)
	if err != nil {
		var notFoundError *core.NotFoundError
		if as := errors.As(err, &notFoundError); as {
			newIngredient := model.NewIngredient(uuid.New(), name, nil)

			err := s.ingredientsRepository.Save(newIngredient)
			if err != nil {
				return nil, fmt.Errorf("could not save the ingredient: %w", err)
			}

			return &newIngredient, nil
		}

		return nil, fmt.Errorf("could not find the ingredient: %w", err)
	}

	return ingredient, nil
}

type CreateRecipeDto struct {
	Name        string
	Description *string
	URL         *string
	Ingredients []string
}
