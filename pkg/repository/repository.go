package repository

import (
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

func (r *Repository) CreateRecipe(nr Recipe) error {

	// newRecipe := Recipe{
	// 	Description: scrapedRecipe.Description,
	// 	// RecipeCategory: scrapedRecipe.RecipeCategory[0],
	// 	// RecipeCuisine: scrapedRecipe.RecipeCuisine[0],
	// 	// RecipeYield:   scrapedRecipe.RecipeYield[0],
	// 	Type: scrapedRecipe.Type,
	// }

	// err := r.db.Table("recipes").Create(&newRecipe).Error

	// if err != nil {
	// 	return err
	// }

	// for _, image := range scrapedRecipe.Image {
	// 	newImage := RecipeImage{
	// 		Image:    image,
	// 		RecipeID: newRecipe.ID,
	// 	}
	// 	err = r.db.Table("recipe_images").Create(&newImage).Error

	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// for _, instruction := range scrapedRecipe.RecipeInstructions {
	// 	newInstruction := RecipeInstruction{
	// 		Text:     instruction.Text,
	// 		RecipeID: newRecipe.ID,
	// 		Step:     2,
	// 	}
	// 	err = r.db.Table("recipe_instructions").Create(&newInstruction).Error

	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// for _, ingredient := range scrapedRecipe.RecipeIngredients {
	// 	newIngredient := Ingredient{
	// 		Ingredient: ingredient,
	// 		RecipeID:   newRecipe.ID,
	// 	}
	// 	err = r.db.Table("recipe_ingredients").Create(&newIngredient).Error

	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// newRating := Rating{
	// 	Votes:    scrapedRecipe.AggregatedRating.RatingCount,
	// 	Score:    scrapedRecipe.AggregatedRating.RatingValue,
	// 	RecipeID: newRecipe.ID,
	// }

	// err = r.db.Table("rating_sections").Create(&newRating).Error

	// if err != nil {
	// 	return err
	// }

	// newKeyword := Keyword{
	// 	Keyword:  scrapedRecipe.Keywords,
	// 	RecipeID: newRecipe.ID,
	// }
	// err = r.db.Table("recipe_keywords").Create(&newKeyword).Error

	// if err != nil {
	// 	return err
	// }

	return nil
}

// Utility functinos
func (r *Repository) DestructiveReset() error {
	err := r.db.DropTableIfExists(&UserTable{}, &RecipeCategoryTable{},
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
	if err := r.db.AutoMigrate(&UserTable{}, &RecipeCategoryTable{},
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
