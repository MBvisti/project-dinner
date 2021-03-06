package repository

import (
	"math/rand"
	"project-dinner/pkg/api"
	"strings"

	"github.com/jinzhu/gorm"
)

func (r repoService) CreateRecipe(usr api.Recipe) error {
	panic("implement me")
}

// GetRandomRecipes ...
func (r *repoService) GetDailyRecipes() ([]api.EmailRecipe, error) {
	numberOfEntries := 0
	err := r.db.Raw("select count(*) from recipes").Count(&numberOfEntries).Error

	if err != nil {
		return nil, ErrNoData
	}

	if numberOfEntries == 0 {
		return nil, ErrNoData
	}

	recipeOneID := uint(rand.Intn(numberOfEntries-1) + 1)
	recipeTwoID := uint(rand.Intn(numberOfEntries-1) + 1)
	recipeThreeID := uint(rand.Intn(numberOfEntries-1) + 1)
	recipeFourID := uint(rand.Intn(numberOfEntries-1) + 1)

	recipeOne := recipe{
		Model: gorm.Model{
			ID: recipeOneID,
		},
	}
	recipeTwo := recipe{
		Model: gorm.Model{
			ID: recipeTwoID,
		},
	}
	recipeThree := recipe{
		Model: gorm.Model{
			ID: recipeThreeID,
		},
	}
	recipeFour := recipe{
		Model: gorm.Model{
			ID: recipeFourID,
		},
	}

	err = r.db.First(&recipeOne).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}
	err = r.db.First(&recipeTwo).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}
	err = r.db.First(&recipeThree).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}
	err = r.db.First(&recipeFour).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}

	recipeOneImage := recipeImage{}
	recipeTwoImage := recipeImage{}
	recipeThreeImage := recipeImage{}
	recipeFourImage := recipeImage{}

	err = r.db.Where("recipe_id = ?", recipeOneID).Last(&recipeOneImage).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}
	err = r.db.Where("recipe_id = ?", recipeTwoID).Last(&recipeTwoImage).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}
	err = r.db.Where("recipe_id = ?", recipeThreeID).Last(&recipeThreeImage).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}
	err = r.db.Where("recipe_id = ?", recipeFourID).Last(&recipeFourImage).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}

	rOne := api.EmailRecipe{
		Name:        recipeOne.Name,
		Description: recipeOne.Description,
		Category:    recipeOne.Category,
		Cuisine:     recipeOne.Cuisine,
		ThumbNail:   recipeOneImage.Image,
		FoundOn:     recipeOne.FoundOn,
	}

	rTwo := api.EmailRecipe{
		Name:        recipeTwo.Name,
		Description: recipeTwo.Description,
		Category:    recipeTwo.Category,
		Cuisine:     recipeTwo.Cuisine,
		ThumbNail:   recipeTwoImage.Image,
		FoundOn:     recipeTwo.FoundOn,
	}

	rThree := api.EmailRecipe{
		Name:        recipeThree.Name,
		Description: recipeThree.Description,
		Category:    recipeThree.Category,
		Cuisine:     recipeThree.Cuisine,
		ThumbNail:   recipeThreeImage.Image,
		FoundOn:     recipeThree.FoundOn,
	}

	rFour := api.EmailRecipe{
		Name:        recipeFour.Name,
		Description: recipeFour.Description,
		Category:    recipeFour.Category,
		Cuisine:     recipeFour.Cuisine,
		ThumbNail:   recipeFourImage.Image,
		FoundOn:     recipeFour.FoundOn,
	}

	selectedRecipes := []api.EmailRecipe{
		rOne,
		rTwo,
		rThree,
		rFour,
	}

	return selectedRecipes, nil
}

// CreateScrapedRecipe saves a recipe from a scraped site
func (r *repoService) CreateScrapedRecipe(nR api.Recipe) error {

	newRecipe := recipe{
		Category:    nR.Category,
		Cuisine:     nR.Cuisine,
		Description: nR.Description,
		Name:        nR.Name,
		Yield:       nR.Yield,
		FoundOn:     nR.FoundOn,
	}

	err := r.db.Create(&newRecipe).Error

	if err != nil {
		return ErrNoCreate
	}

	for _, img := range nR.Images {
		newImage := recipeImage{
			Image:    img,
			RecipeID: newRecipe.ID,
		}

		err = r.db.Create(&newImage).Error

		if err != nil {
			return ErrNoCreate
		}
	}

	for _, instruction := range nR.Instructions {
		newInstruction := recipeInstruction{
			Text:     instruction.Text,
			RecipeID: newRecipe.ID,
			Step:     instruction.Step,
		}
		err = r.db.Create(&newInstruction).Error

		if err != nil {
			return ErrNoCreate
		}
	}

	for _, ingredient := range nR.Ingredients {
		newIngredient := recipeIngredient{
			Ingredient: ingredient,
			RecipeID:   newRecipe.ID,
		}
		err = r.db.Create(&newIngredient).Error

		if err != nil {
			return ErrNoCreate
		}
	}

	newRating := rating{
		Votes:    nR.Score.Votes,
		Score:    nR.Score.Score,
		RecipeID: newRecipe.ID,
	}

	err = r.db.Create(&newRating).Error

	if err != nil {
		return ErrNoCreate
	}

	for _, keyWord := range nR.Keywords {
		newKeyword := recipeKeyword{
			Keyword:  strings.TrimSpace(keyWord),
			RecipeID: newRecipe.ID,
		}
		err = r.db.Create(&newKeyword).Error

		if err != nil {
			return ErrNoCreate
		}
	}

	return nil
}
