package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	service "project-dinner/pkg/api"
)

type Repository interface {
	GetDailyRecipes() ([]service.EmailRecipe, error)
	CreateScrapedRecipe(nR service.Recipe) error
	CreateRecipe(usr service.Recipe) error
	GetEmailList() ([]service.User, error)
	CreateUser(usr service.User) error
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
	ErrNoResourceFound = errors.New("repo - resource doesn't exists")
	// ErrNoCreate ...
	ErrNoCreate = errors.New("repo - couldn't create the resource")
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
	).Error
	if err != nil {
		return err
	}

	err = r.MigrateTables()
	if err != nil {
		return err
	}

	morten := user{
		Email: "mbv1406@gmail.com",
		Name:  "Morten",
	}

	err = r.db.Create(&morten).Error

	if err != nil {
		return err
	}

	javiera := user{
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
	).Error; err != nil {
		return err
	}
	return nil
}
