package repository

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Services ...
type Services struct {
	Recipe RecipeService
	User   UserService
	db     *gorm.DB
}

// NewStorage returns a new repository
func NewStorage(db *gorm.DB) *Services {
	return &Services{
		Recipe: NewRecipeService(db),
		User:   NewUserService(db),
		db:     db,
	}
}

var (
	// ErrNoData is returned when there is no data in a given table
	ErrNoData = errors.New("repo: no data for the requested resource")
)

// DestructiveReset resets the database and and creates two users
func (r *Services) DestructiveReset() error {
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

// MigrateTables migrates all tables in definitions
func (r *Services) MigrateTables() error {
	if err := r.db.AutoMigrate(&UserTable{}, &RecipeCategoryTable{}, &RecipeInstructionTable{},
		&RatingTable{}, &IngredientTable{}, &KeywordTable{}, &RecipeImageTable{}, &RecipeTable{}, &DailyRecipes{}).Error; err != nil {
		return err
	}
	return nil
}
