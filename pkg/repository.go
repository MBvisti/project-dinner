package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Repository struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Email    string    `gorm:"not null;unique_index"`
	MailTime time.Time `gorm:"not null"`
}

type Recipe struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Image        string
	Link         string `gorm:"not null"`
	Credits      string
	Instructions []Instruction `gorm:"not null"`
	Ingredients  []Ingredient  `gorm:"not null"`
}

type Ingredient struct {
	gorm.Model
	RecipeID        uint   `gorm:"not null"`
	Amount          string `gorm:"not null"`
	MeasurementType string
	What            string `gorm:"not null"`
}

type Instruction struct {
	gorm.Model
	IngredientID uint `gorm:"not null"`
	Name         string
	Step         int    `gorm:"not null"`
	Text         string `gorm:"not null"`
}

func NewRepository(db *gorm.DB) *Repository {

	return &Repository{
		db: db,
	}
}

func (r *Repository) DestructiveReset() error {
	err := r.db.DropTableIfExists(&User{}, &Recipe{}, &Ingredient{}, &Instruction{}).Error
	if err != nil {
		return err
	}

	return r.AutoMigrate()
}

func (r *Repository) AutoMigrate() error {
	if err := r.db.AutoMigrate(&User{}, &Recipe{}, &Ingredient{}, &Instruction{}).Error; err != nil {
		return err
	}
	return nil
}
