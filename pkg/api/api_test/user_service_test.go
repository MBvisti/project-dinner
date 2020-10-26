package api_test

import (
	"project-dinner/pkg/api"
	"project-dinner/pkg/repository"
	"regexp"
	"testing"
)

type mockRepo struct{}

var isMailValid = regexp.MustCompile(
	`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)

func (mR *mockRepo) CreateUser(usr api.User) error {
	if usr.Email == "" {
		return repository.ErrEmailRequired
	}

	if !isMailValid.MatchString(usr.Email) {
		return repository.ErrEmailInvalid
	}

	return nil
}

func TestCreateUser(t *testing.T) {
	mR := mockRepo{}
	t.Run("should create a user and return nil", func(t *testing.T) {

		userService := api.NewUserService(&mR)

		newUser := api.User{
			Email: "jonsnow@gmail.com",
			Name:  "Testesen",
		}

		err := userService.CreateUser(newUser)

		if err != nil {
			t.Errorf("Create user test failed. Got err: %s", err.Error())
		}
	})

	t.Run("should not create a user and return 'email required' error", func(t *testing.T) {

		userService := api.NewUserService(&mR)

		newUser := api.User{
			Email: "",
			Name:  "Testesen",
		}

		err := userService.CreateUser(newUser)

		if err != repository.ErrEmailRequired {
			t.Errorf("Create user test failed. Got err: %s", err.Error())
		}
	})

	t.Run("should not create a user and return 'email not valid' error", func(t *testing.T) {

		userService := api.NewUserService(&mR)

		newUser := api.User{
			Email: "!4131mbv@gmail.com",
			Name:  "Testesen",
		}

		err := userService.CreateUser(newUser)

		if err != repository.ErrEmailInvalid {
			t.Errorf("Create user test failed. Got err: %s", err)
		}
	})
}
