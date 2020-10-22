package main

import (
	"fmt"
	"math/rand"
	"os"
	app "project-dinner/pkg"
	repository "project-dinner/pkg/storage"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("PORT")
	whatEnv := os.Getenv("WHAT_ENVIRONMENT_IS_THIS")
	sendGridUser := os.Getenv("SEND_GRID_USER")
	sendGridPassword := os.Getenv("SEND_GRID_API_KEY")
	mailHost := os.Getenv("HOST")
	connectionString := os.Getenv("DATABASE_URL")
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))

	gin.SetMode(gin.ReleaseMode)

	if whatEnv == "development" {
		port = "5000"
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	r.Use(cors.Default())

	database, err := setupDatabase(connectionString, whatEnv)

	if err != nil {
		return err
	}

	s := repository.NewStorage(database)

	defer database.Close()
	s.MigrateTables()

	t, err := time.LoadLocation("Europe/Copenhagen")

	if err != nil {
		return err
	}

	c := cron.New(cron.WithLocation(t))

	m := gomail.NewDialer(mailHost, mailPort, sendGridUser, sendGridPassword)

	server := app.NewServer(s, r, c, m)

	err = server.Run(":" + port)

	if err != nil {
		return err
	}

	// seed random generator used to pick the emails
	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func setupDatabase(connectionInfo string, environment string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connectionInfo)

	if err != nil {
		return nil, err
	}

	if environment == "development" {
		db.LogMode(true)
	}

	return db, nil
}
