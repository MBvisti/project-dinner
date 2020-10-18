package app

import (
	"html/template"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repository struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Email string `gorm:"not null;unique_index"`
	Name  string
}

// Recipe ...
type Recipe struct {
	gorm.Model
	Type           string
	Description    string
	RecipeCategory string
	RecipeCuisine  string
	RecipeYield    string
}

type RecipeIngredient struct {
	gorm.Model
	Ingredient string
	RecipeID   uint
	Recipe     Recipe
}

type RecipeKeyword struct {
	gorm.Model
	Keyword  string
	RecipeID uint
	Recipe   Recipe
}

type RecipeImage struct {
	gorm.Model
	Image    string
	RecipeID uint
	Recipe   Recipe
}

type RecipeInstruction struct {
	gorm.Model
	Type     string
	Name     string
	URL      string
	Text     string
	RecipeID uint
	Recipe   Recipe
	// ItemListElement map[string]string
}

type RatingSection struct {
	gorm.Model
	RatingCount string
	RatingValue string
	RecipeID    uint
	Recipe      Recipe
}

type DailyRecipes struct {
	gorm.Model
	Name         string        `json:"title"`
	Image        string        `json:"image"`
	Description  template.HTML `json:"summary"`
	Source       string        `json:"sourceUrl"`
	Instructions string        `json:"instructions"`
}

type EmailList struct {
	Email string
	Name  string
}

func NewRepository(db *gorm.DB) *Repository {

	return &Repository{
		db: db,
	}
}

func (r *Repository) DestructiveReset() error {
	err := r.db.DropTableIfExists(&User{}, &RecipeInstruction{},
		&RatingSection{}, &RecipeIngredient{}, &RecipeKeyword{}, &RecipeImage{}, &Recipe{}, &DailyRecipes{}).Error
	if err != nil {
		return err
	}

	err = r.AutoMigrate()
	if err != nil {
		return err
	}

	morten := User{
		Email: "mbv1406@gmail.com",
		Name:  "Morten",
	}

	err = r.db.Create(&morten).Error

	if err != nil {
		return err
	}

	javiera := User{
		Email: "j.camuslaso@gmail.com",
		Name:  "Javiera",
	}

	err = r.db.Create(&javiera).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateRecipe(scrapedRecipe *ScrapedRecipeSection) error {

	newRecipe := Recipe{
		Description: scrapedRecipe.Description,
		// RecipeCategory: scrapedRecipe.RecipeCategory[0],
		// RecipeCuisine: scrapedRecipe.RecipeCuisine[0],
		// RecipeYield:   scrapedRecipe.RecipeYield[0],
		Type: scrapedRecipe.Type,
	}

	err := r.db.Table("recipes").Create(&newRecipe).Error

	if err != nil {
		return err
	}

	for _, image := range scrapedRecipe.Image {
		newImage := RecipeImage{
			Image:    image,
			RecipeID: newRecipe.ID,
		}
		err = r.db.Table("recipe_images").Create(&newImage).Error

		if err != nil {
			return err
		}
	}

	for _, instruction := range scrapedRecipe.RecipeInstructions {
		newInstruction := RecipeInstruction{
			Type:     instruction.Type,
			Name:     instruction.Name,
			URL:      instruction.URL,
			Text:     instruction.Text,
			RecipeID: newRecipe.ID,
		}
		err = r.db.Table("recipe_instructions").Create(&newInstruction).Error

		if err != nil {
			return err
		}
	}

	for _, ingredient := range scrapedRecipe.RecipeIngredients {
		newIngredient := RecipeIngredient{
			Ingredient: ingredient,
			RecipeID:   newRecipe.ID,
		}
		err = r.db.Table("recipe_ingredients").Create(&newIngredient).Error

		if err != nil {
			return err
		}
	}

	newRating := RatingSection{
		RatingCount: scrapedRecipe.AggregatedRating.RatingCount,
		RatingValue: scrapedRecipe.AggregatedRating.RatingValue,
		RecipeID:    newRecipe.ID,
	}

	err = r.db.Table("rating_sections").Create(&newRating).Error

	if err != nil {
		return err
	}

	newKeyword := RecipeKeyword{
		Keyword:  scrapedRecipe.Keywords,
		RecipeID: newRecipe.ID,
	}
	err = r.db.Table("recipe_keywords").Create(&newKeyword).Error

	if err != nil {
		return err
	}

	return err
}

func (r *Repository) AutoMigrate() error {
	if err := r.db.AutoMigrate(&User{}, &RecipeInstruction{},
		&RatingSection{}, &RecipeKeyword{}, &RecipeIngredient{}, &RecipeImage{}, &Recipe{}, &DailyRecipes{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetEmailList() ([]EmailList, error) {
	var emailList []EmailList
	err := r.db.Table("users").Select("email, name").Scan(&emailList).Error

	if err != nil {
		return nil, err
	}

	return emailList, nil
}

func (r *Repository) TodaysRecipes() ([]DailyRecipes, error) {
	var selectedRecipes []DailyRecipes
	err := r.db.Table("daily_recipes").Find(&selectedRecipes).Error

	if err != nil {
		return nil, err
	}

	return selectedRecipes, nil
}
