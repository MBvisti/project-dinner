package repository

import (
	"math/rand"
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

func (r *Repository) GetRandomRecipes() ([]EmailRecipe, error) {
	numberOfEntries := 0
	err := r.db.Raw("select count(*) from recipe_tables").Count(&numberOfEntries).Error

	if err != nil {
		return nil, err
	}

	recipeOneID := uint(rand.Intn(numberOfEntries-1) + 1)
	recipeTwoID := uint(rand.Intn(numberOfEntries-1) + 1)
	recipeThreeID := uint(rand.Intn(numberOfEntries-1) + 1)
	recipeFourID := uint(rand.Intn(numberOfEntries-1) + 1)

	recipeOne := RecipeTable{
		Model: gorm.Model{
			ID: recipeOneID,
		},
	}
	recipeTwo := RecipeTable{
		Model: gorm.Model{
			ID: recipeTwoID,
		},
	}
	recipeThree := RecipeTable{
		Model: gorm.Model{
			ID: recipeThreeID,
		},
	}
	recipeFour := RecipeTable{
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

	recipeOneImage := RecipeImageTable{}
	recipeTwoImage := RecipeImageTable{}
	recipeThreeImage := RecipeImageTable{}
	recipeFourImage := RecipeImageTable{}

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

func (r *Repository) CreateScrapedRecipe(nR Recipe) error {

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

func (r *Repository) GetEmailList() ([]EmailList, error) {
	var emailList []EmailList
	err := r.db.Table("user_tables").Select("email, name").Scan(&emailList).Error

	if err != nil {
		return nil, err
	}

	return emailList, nil
}
