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
	Name    string `gorm:"not null"`
	Image   string
	Link    string `gorm:"not null"`
	Credits string
}

func NewRepository(db *gorm.DB) *Repository {

	return &Repository{
		db: db,
	}
}

func (r *Repository) DestructiveReset() error {
	err := r.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}

	return r.AutoMigrate()
}

func (r *Repository) AutoMigrate() error {
	if err := r.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
