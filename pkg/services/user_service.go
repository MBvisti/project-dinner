package service

import (
	"strings"
)

// UserService is the user service's interface
type UserService interface {
	CreateUser(u User) error
}

type userService struct {
	app *App
}

// NewUserService ...
func NewUserService(a *App) UserService {
	return &userService{
		a,
	}
}

func (s *userService) CreateUser(u User) error {
	var newUser User
	newUser.Email = strings.ToLower(u.Email)
	newUser.Name = strings.ToLower(u.Name)

	err := s.app.uR.CreateUser(newUser)

	if err != nil {
		return err
	}

	return nil

}
