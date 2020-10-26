package repositorytest

import (
	"github.com/jinzhu/gorm"
	"log"
	"project-dinner/pkg/api"
	"project-dinner/pkg/repository"
	"testing"
)

// UserTable ...
type user struct {
	gorm.Model
	Email string `gorm:"not null;unique_index"`
	Name  string
}

func TestCreateUser(t *testing.T) {
	testDB, teardown := NewTestDatabase(&user{})

	db := repository.NewStorage(testDB)

	newUser := api.User{
		Name:  "Morten",
		Email: "vistisen@live.dk",
	}
	err := db.CreateUser(newUser)

	if err != nil {
		log.Printf("this is the test create user error: %s", err.Error())
	}

	log.Printf("this is the err variable from create user test: %s", err)

	defer teardown()
}
