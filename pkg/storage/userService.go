package repository

import "github.com/jinzhu/gorm"

// UserService ...
type UserService interface {
	GetEmailList() ([]UserEmail, error)
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
		return nil, err
	}

	return emailList, nil
}
