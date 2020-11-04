package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"project-dinner/pkg/api"
)

type Repository interface {
	GetDailyRecipes() ([]api.EmailRecipe, error)
	CreateScrapedRecipe(nR api.Recipe) error
	CreateRecipe(usr api.Recipe) error
	GetEmailList() ([]api.User, error)
	CreateUser(usr api.User) error
	MigrateTables() error
	SeedProductionData() error
	GetFeaturedRecipes() []api.FeaturedRecipe
}

type repoService struct {
	db *gorm.DB
}

// NewStorage returns a new repository
func NewStorage(db *gorm.DB) Repository {
	return &repoService{
		db: db,
	}
}

var (
	// ErrNoData is returned when there is no data in a given table
	ErrNoData = errors.New("repo - no data for the requested resource")
	// ErrEmailInvalid is returned when the provided email doesn't pass the regex validation
	ErrEmailInvalid = errors.New("repo - email not valid")
	// ErrEmailRequired is returned when the email is empty/non-existent
	ErrEmailRequired = errors.New("repo - email is required")
	// ErrNoResourceFound is returned when a query couldn't be performed
	ErrNoResourceFound = errors.New("repo - resource does not exists")
	// ErrNoCreate is returned when a resource couldn't be created
	ErrNoCreate = errors.New("repo - could not create the resource")
	// ErrNoMigrate is returned when migrations couldn't be run
	ErrNoMigrate = errors.New("repo - could not run migrations")
)

// DestructiveReset resets the database and and creates two users
func (r *repoService) DestructiveReset() error {
	err := r.db.DropTableIfExists(
		&recipe{},
		&category{},
		&cuisine{},
		&recipeCategory{},
		&recipeCuisine{},
		&recipeIngredient{},
		&recipeKeyword{},
		&recipeImage{},
		&recipeInstruction{},
		&rating{},
		&user{},
		&recipeType{},
		&dietaryType{},
	).Error
	if err != nil {
		return err
	}

	err = r.MigrateTables()
	if err != nil {
		return err
	}
	err = r.SeedProductionData()
	if err != nil {
		return err
	}

	morten := user{
		Email:         "mbv1406@gmail.com",
		Name:          "Morten",
		RecipeTypeID:  3,
		DietaryTypeID: 2,
	}

	err = r.db.Create(&morten).Error

	if err != nil {
		return err
	}

	javiera := user{
		Email:         "j.camuslaso@gmail.com",
		Name:          "Javiera",
		RecipeTypeID:  1,
		DietaryTypeID: 1,
	}

	err = r.db.Create(&javiera).Error

	if err != nil {
		return err
	}

	return nil
}

// MigrateTables migrates all tables in definitions
func (r *repoService) MigrateTables() error {
	if err := r.db.AutoMigrate(
		&recipe{},
		&category{},
		&cuisine{},
		&recipeCategory{},
		&recipeCuisine{},
		&recipeIngredient{},
		&recipeKeyword{},
		&recipeImage{},
		&recipeInstruction{},
		&rating{},
		&user{},
		&recipeType{},
		&dietaryType{},
	).Error; err != nil {
		log.Printf("repo - this is the migration error: %s", err.Error())
		return ErrNoMigrate
	}
	return nil
}

// SeedProductionData creates necessary data
func (r *repoService) SeedProductionData() error {
	fitnessType := recipeType{
		Type: "fitness",
	}

	regularType := recipeType{
		Type: "regular",
	}

	comfyType := recipeType{
		Type: "comfy",
	}

	err := r.db.Create(&fitnessType).Error
	if err != nil {
		return err
	}
	err = r.db.Create(&regularType).Error
	if err != nil {
		return err
	}
	err = r.db.Create(&comfyType).Error
	if err != nil {
		return err
	}

	vegetarian := dietaryType{
		Type: "vegetarian",
	}
	meatEater := dietaryType{
		Type: "meat",
	}

	err = r.db.Create(&vegetarian).Error
	if err != nil {
		return err
	}

	err = r.db.Create(&meatEater).Error
	if err != nil {
		return err
	}
	return nil
}
