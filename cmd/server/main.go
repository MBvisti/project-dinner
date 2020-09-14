package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/jinzhu/gorm"
	"os"
	app "project-dinner/pkg"
	"strconv"
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

	c := colly.NewCollector()

	server := app.NewServer(s, r, c)

	err = server.Run(":" + port)

	if err != nil {
		return err
	}

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
