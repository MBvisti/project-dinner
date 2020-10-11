package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	"log"
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
	Name         string `json:"title"`
	Image        string `json:"image"`
	Description  string `json:"summary"`
	Source       string `json:"sourceUrl"`
	Instructions string `json:"instructions"`
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
	err := r.db.DropTableIfExists(&User{}, &Recipe{}, &DailyRecipes{}).Error
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

func (r *Repository) CreateRecipe(recipe []Recipe) error {
	r.db.Exec("DELETE FROM daily_recipes")

	var err error
	for index, reci := range recipe {
		log.Printf("this is index: %v", index)
		err = r.db.Create(&reci).Error

		err = r.db.Table("daily_recipes").Create(&reci).Error
	}

	return err
}

func (r *Repository) AutoMigrate() error {
	if err := r.db.AutoMigrate(&User{}, &Recipe{}, &DailyRecipes{}).Error; err != nil {
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
