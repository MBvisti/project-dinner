package repository

import (
	"regexp"

	"github.com/jinzhu/gorm"
)

// UserService ...
type UserService interface {
	GetEmailList() ([]UserEmail, error)
	CreateUser(usr NewUser) error
}

type userService struct {
	uS UserService
	db *gorm.DB
}

// NewUserService ...
func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

// GetEmailList returns the email list - TODO: maybe separate the user list and emailing list into two different tables?
func (r *userService) GetEmailList() ([]UserEmail, error) {
	var emailList []UserEmail
	err := r.db.Table("users").Select("email, name").Scan(&emailList).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}

	return emailList, nil
}

var isMailValid = regexp.MustCompile(
	`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)

type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	TimeZone string `json:"time_zone"`
}

// CreateUser ...
func (r *userService) CreateUser(usr NewUser) error {
	if usr.Email == "" {
		return ErrEmailRequired
	}

	if !isMailValid.MatchString(usr.Email) {
		return ErrEmailInvalid
	}

	nU := user{
		Name:  usr.Name,
		Email: usr.Email,
	}

	err := r.db.Create(&nU).Error

	if err != nil {
		return ErrNoCreate
	}

	return nil
}
