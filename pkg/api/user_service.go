package api

import (
	"strings"
)

// UserService is the user service's interface
type UserService interface {
	CreateUser(u User) error
}

// Repository ...
type UserRepository interface {
	CreateUser(u User) error
}

type userService struct {
	storage UserRepository
}

// NewUserService ...
func NewUserService(r UserRepository) UserService {
	return &userService{
		r,
	}
}

func (s *userService) CreateUser(u User) error {
	var newUser User

	newUser.Email = strings.ToLower(u.Email)
	newUser.Name = strings.ToLower(u.Name)
	newUser.RecipeType = strings.ToLower(u.RecipeType)
	newUser.DietaryType = strings.ToLower(u.DietaryType)

	err := s.storage.CreateUser(newUser)

	if err != nil {
		return err
	}

	return nil

}
