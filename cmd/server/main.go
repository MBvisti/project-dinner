package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"project-dinner/pkg/api"
	"project-dinner/pkg/repository"
	"project-dinner/pkg/rest"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
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
	sentryDNS := os.Getenv("LOGGING_SERVICE_DNS")

	gin.SetMode(gin.ReleaseMode)

	if whatEnv == "development" {
		port = "5000"
		gin.SetMode(gin.DebugMode)
	}

	if whatEnv != "development" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDNS,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
	}

	database, err := setupDatabase(connectionString, whatEnv)
	defer database.Close()

	if err != nil {
		return err
	}

	s := repository.NewStorage(database)
	err = s.MigrateTables()

	if err != nil {
		log.Fatal(err)
	}

	m := gomail.NewDialer(mailHost, mailPort, sendGridUser, sendGridPassword)

	usrService := api.NewUserService(s)
	emailService := api.NewEmailService(m, s)
	spiderService := api.NewSpiderService(m, setupCrawler())

	// start up the every day mailer
	job, err := emailService.EveryDayMailer()
	go startMailer(job, err)

	router := rest.Routes(usrService, emailService, spiderService, s)
	log.Printf("starting server on: " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))

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

func setupCrawler() *colly.Collector {
	c := colly.NewCollector()

	return c
}

func startMailer(j cron.Job, err error) {
	if err != nil {
		log.Printf("this is the mailer job err: %v", err)
	}

	j.Run()
}
