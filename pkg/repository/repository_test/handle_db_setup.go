package repositorytest

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

func NewTestDatabase(tables ...interface{}) (*gorm.DB, func()) {
	connectionString := os.Getenv("DATABASE_URL")
	whatEnv := os.Getenv("WHAT_ENVIRONMENT_IS_THIS")
	testDatabase, err := gorm.Open("postgres", connectionString)

	log.Printf("this is the connection string: %s", connectionString)
	if err != nil {
		log.Fatal("could not create test database")
	}

	if whatEnv == "test" {
		testDatabase.LogMode(true)
	}

	testDatabase.AutoMigrate(tables...)

	return testDatabase, func() {
		err := Teardown(testDatabase, tables...)

		if err != nil {
			log.Fatal("could not tear down test database")
		}

		testDatabase.Close()
	}

}

func Teardown(db *gorm.DB, tables ...interface{}) error {
	err := db.DropTable(tables...).Error

	if err != nil {
		return err
	}

	return nil
}
