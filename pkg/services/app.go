package service

import "gopkg.in/gomail.v2"

// UserRepository ...
type UserRepository interface {
	CreateUser(usr User) error
	GetResponse() string
	GetEmailList() ([]User, error)
}

// RecipeRepository ...
type RecipeRepository interface {
	CreateRecipe(usr Recipe) error
	GetDailyRecipes() ([]EmailRecipe, error)
}

// App holds dependencies for the services
type App struct {
	uR           UserRepository
	rR           RecipeRepository
	mailProvider *gomail.Dialer
}

// NewApp ...
func NewApp(ur UserRepository, rr RecipeRepository, mp *gomail.Dialer) *App {
	return &App{
		ur,
		rr,
		mp,
	}
}
