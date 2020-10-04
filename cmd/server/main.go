package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"math/rand"
	"os"
	app "project-dinner/pkg"
	"strconv"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("PORT")
	developmentMode, err := strconv.ParseBool(os.Getenv("DEVELOPMENT_MODE"))

	if err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)

	if developmentMode {
		port = "5000"
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	r.Use(cors.Default())

	connectionString := os.Getenv("DATABASE_URL")
	database, err := setupDatabase(connectionString, developmentMode)

	if err != nil {
		return err
	}

	s := app.NewRepository(database)

	defer database.Close()
	s.AutoMigrate()

	t, err := time.LoadLocation("Europe/Copenhagen")

	if err != nil {
		return err
	}

	c := cron.New(cron.WithLocation(t))

	sendGridUser := os.Getenv("SEND_GRID_USER")
	sendGridPassword := os.Getenv("SEND_GRID_API_KEY")
	mailHost := os.Getenv("HOST")
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))

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

func setupDatabase(connectionInfo string, logMode bool) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connectionInfo)

	if err != nil {
		return nil, err
	}

	db.LogMode(logMode)

	return db, nil
}
