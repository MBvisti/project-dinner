package repository

import (
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repository struct {
	db *gorm.DB
}

type EmailList struct {
	Email string
	Name  string
}

// NewRepository returns a new repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateRecipe(nR Recipe) error {

	newRecipe := RecipeTable{
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
		newImage := RecipeImageTable{
			Image:    img,
			RecipeID: newRecipe.ID,
		}

		err = r.db.Create(&newImage).Error

		if err != nil {
			return err
		}
	}

	for _, instruction := range nR.Instructions {
		newInstruction := RecipeInstructionTable{
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
		newIngredient := IngredientTable{
			Ingredient: ingredient,
			RecipeID:   newRecipe.ID,
		}
		err = r.db.Create(&newIngredient).Error

		if err != nil {
			return err
		}
	}

	newRating := RatingTable{
		Votes:    nR.Score.Votes,
		Score:    nR.Score.Score,
		RecipeID: newRecipe.ID,
	}

	err = r.db.Create(&newRating).Error

	if err != nil {
		return err
	}

	for _, keyWord := range nR.Keywords {
		newKeyword := KeywordTable{
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

// Utility functinos
func (r *Repository) DestructiveReset() error {
	err := r.db.DropTableIfExists(&UserTable{}, &RecipeCategoryTable{}, &RecipeInstructionTable{},
		&RatingTable{}, &IngredientTable{}, &KeywordTable{}, &RecipeImageTable{}, &RecipeTable{}, &DailyRecipes{}).Error
	if err != nil {
		return err
	}

	err = r.MigrateTables()
	if err != nil {
		return err
	}

	morten := UserTable{
		Email: "mbv1406@gmail.com",
		Name:  "Morten",
	}

	err = r.db.Create(&morten).Error

	if err != nil {
		return err
	}

	javiera := UserTable{
		Email: "j.camuslaso@gmail.com",
		Name:  "Javiera",
	}

	err = r.db.Create(&javiera).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) MigrateTables() error {
	if err := r.db.AutoMigrate(&UserTable{}, &RecipeCategoryTable{}, &RecipeInstructionTable{},
		&RatingTable{}, &IngredientTable{}, &KeywordTable{}, &RecipeImageTable{}, &RecipeTable{}, &DailyRecipes{}).Error; err != nil {
		return err
	}
	return nil
}

// func (r *Repository) GetEmailList() ([]EmailList, error) {
// 	var emailList []EmailList
// 	err := r.db.Table("users").Select("email, name").Scan(&emailList).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return emailList, nil
// }

// func (r *Repository) TodaysRecipes() ([]DailyRecipes, error) {
// 	var selectedRecipes []DailyRecipes
// 	err := r.db.Table("daily_recipes").Find(&selectedRecipes).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return selectedRecipes, nil
// }
