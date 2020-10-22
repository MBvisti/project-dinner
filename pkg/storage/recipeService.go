package repository

import (
	"math/rand"
	"strings"

	"github.com/jinzhu/gorm"
)

// RecipeService ...
type RecipeService interface {
	GetRandomRecipes() ([]EmailRecipe, error)
	CreateScrapedRecipe(nR Recipe) error
}

type recipeService struct {
	rS RecipeService
	db *gorm.DB
}

// NewRecipeService ...
func NewRecipeService(db *gorm.DB) RecipeService {
	return &recipeService{
		db: db,
	}
}

// GetRandomRecipes ...
func (r *recipeService) GetRandomRecipes() ([]EmailRecipe, error) {
	numberOfEntries := 0
	err := r.db.Raw("select count(*) from recipe_tables").Count(&numberOfEntries).Error

	if err != nil {
		return nil, err
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
		return nil, err
	}
	err = r.db.First(&recipeTwo).Error

	if err != nil {
		return nil, err
	}
	err = r.db.First(&recipeThree).Error

	if err != nil {
		return nil, err
	}
	err = r.db.First(&recipeFour).Error

	if err != nil {
		return nil, err
	}

	recipeOneImage := recipeImage{}
	recipeTwoImage := recipeImage{}
	recipeThreeImage := recipeImage{}
	recipeFourImage := recipeImage{}

	err = r.db.Where("recipe_id = ?", recipeOneID).Last(&recipeOneImage).Error

	if err != nil {
		return nil, err
	}
	err = r.db.Where("recipe_id = ?", recipeTwoID).Last(&recipeTwoImage).Error

	if err != nil {
		return nil, err
	}
	err = r.db.Where("recipe_id = ?", recipeThreeID).Last(&recipeThreeImage).Error

	if err != nil {
		return nil, err
	}
	err = r.db.Where("recipe_id = ?", recipeFourID).Last(&recipeFourImage).Error

	if err != nil {
		return nil, err
	}

	rOne := EmailRecipe{
		Name:        recipeOne.Name,
		Description: recipeOne.Description,
		Category:    recipeOne.Category,
		Cuisine:     recipeOne.Cuisine,
		ThumbNail:   recipeOneImage.Image,
		FoundOn:     recipeOne.FoundOn,
	}

	rTwo := EmailRecipe{
		Name:        recipeTwo.Name,
		Description: recipeTwo.Description,
		Category:    recipeTwo.Category,
		Cuisine:     recipeTwo.Cuisine,
		ThumbNail:   recipeTwoImage.Image,
		FoundOn:     recipeTwo.FoundOn,
	}

	rThree := EmailRecipe{
		Name:        recipeThree.Name,
		Description: recipeThree.Description,
		Category:    recipeThree.Category,
		Cuisine:     recipeThree.Cuisine,
		ThumbNail:   recipeThreeImage.Image,
		FoundOn:     recipeThree.FoundOn,
	}

	rFour := EmailRecipe{
		Name:        recipeFour.Name,
		Description: recipeFour.Description,
		Category:    recipeFour.Category,
		Cuisine:     recipeFour.Cuisine,
		ThumbNail:   recipeFourImage.Image,
		FoundOn:     recipeFour.FoundOn,
	}

	selectedRecipes := []EmailRecipe{
		rOne,
		rTwo,
		rThree,
		rFour,
	}

	return selectedRecipes, nil
}

// CreateScrapedRecipe saves a recipe from a scraped site
func (r *recipeService) CreateScrapedRecipe(nR Recipe) error {

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
		return err
	}

	for _, img := range nR.Images {
		newImage := recipeImage{
			Image:    img,
			RecipeID: newRecipe.ID,
		}

		err = r.db.Create(&newImage).Error

		if err != nil {
			return err
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
			return err
		}
	}

	for _, ingredient := range nR.Ingredients {
		newIngredient := recipeIngredient{
			Ingredient: ingredient,
			RecipeID:   newRecipe.ID,
		}
		err = r.db.Create(&newIngredient).Error

		if err != nil {
			return err
		}
	}

	newRating := rating{
		Votes:    nR.Score.Votes,
		Score:    nR.Score.Score,
		RecipeID: newRecipe.ID,
	}

	err = r.db.Create(&newRating).Error

	if err != nil {
		return err
	}

	for _, keyWord := range nR.Keywords {
		newKeyword := recipeKeyword{
			Keyword:  strings.TrimSpace(keyWord),
			RecipeID: newRecipe.ID,
		}
		err = r.db.Create(&newKeyword).Error

		if err != nil {
			return err
		}
	}

	return nil
}
