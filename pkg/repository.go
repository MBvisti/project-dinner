package app

import (
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

type Recipe struct {
	gorm.Model
	Name        string `json:"title"`
	Image       string `json:"image"`
	Description string `json:"summary"`
	Source      string `json:"sourceUrl"`
	Rating      string
	Instructions string `json:"instructions"`
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
	err := r.db.DropTableIfExists(&User{}, &Recipe{}).Error
	if err != nil {
		return err
	}

	err = r.AutoMigrate()
	if err != nil {
		return err
	}


	user := User{
		Email: "mbv1406@gmail.com",
		Name:  "Morten",
	}

	err = r.db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateRecipe(recipe *Recipe) error {
	err := r.db.Create(&recipe).Error

	return err
}

func (r *Repository) AutoMigrate() error {
	if err := r.db.AutoMigrate(&User{}, &Recipe{}).Error; err != nil {
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

func (r *Repository) TodaysRecipes() ([]Recipe, error) {
	var recipes []Recipe
	var count int
	err := r.db.Model(&recipes).Count(&count).Error
	if err != nil {
		return nil, err
	}

	var selectedRecipes []Recipe
	err = r.db.Find(&selectedRecipes, []int{1, 2, 3, 4}).Error

	if err != nil {
		return nil, err
	}

	return selectedRecipes, nil
}
