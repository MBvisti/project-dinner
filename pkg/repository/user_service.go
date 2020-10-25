package repository

import (
	"github.com/jinzhu/gorm"
	service "project-dinner/pkg/services"
	"regexp"
)

// UserService ...
type UserService interface {
	GetEmailList() ([]service.User, error)
	CreateUser(usr service.User) error
	GetResponse() string
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

func (r *userService) GetResponse() string {
	return "HELLO HELLO"
}

// GetEmailList returns the email list - TODO: maybe separate the user list and emailing list into two different tables?
func (r *userService) GetEmailList() ([]service.User, error) {
	var emailList []service.User
	err := r.db.Table("users").Select("email, name").Scan(&emailList).Error

	if err != nil {
		return nil, ErrNoResourceFound
	}

	return emailList, nil
}

var isMailValid = regexp.MustCompile(
	`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)

// CreateUser ...
func (r *userService) CreateUser(usr service.User) error {
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
