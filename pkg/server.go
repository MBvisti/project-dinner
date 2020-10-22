package app

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"project-dinner/pkg/repository"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

// Server ...
type Server struct {
	router  *gin.Engine
	storage *repository.Services
	cron    *cron.Cron
	mailer  *gomail.Dialer
}

// NewServer returns a new server
func NewServer(s *repository.Services, r *gin.Engine, c *cron.Cron, m *gomail.Dialer) Server {
	return Server{
		storage: s,
		router:  r,
		cron:    c,
		mailer:  m,
	}
}

// Run starts the server
func (s *Server) Run(addr string) error {
	// TODO: change this to setup the main cronjob
	err := s.RecipeMailer()

	if err != nil {
		log.Printf("this is err from cronjob: %v", err)
		return err
	}

	log.Printf("Starting the server on: %v", addr)
	err = http.ListenAndServe(addr, s.routes())

	if err != nil {
		return err
	}

	return nil
}

// RecipeMailer sends out the daily recipes
func (s *Server) RecipeMailer() error {
	mailTemplate, err := template.ParseFiles("./template/daily_recipe_email.html")

	if err != nil {
		return err
	}

	var emailList []*gomail.Message

	userList, err := s.storage.User.GetEmailList()

	if err != nil {
		return err
	}

	dailyRecipes, err := s.storage.Recipe.GetRandomRecipes()

	for _, user := range userList {
		usrRecipe := UserRecipe{
			UserName: user.Name,
			Recipes:  dailyRecipes,
		}

		var t bytes.Buffer
		err = mailTemplate.Execute(&t, usrRecipe)

		mail := gomail.NewMessage()
		mail.SetAddressHeader("From", "noreply@mbvistisen.dk", "Morten's recipe service")
		mail.SetHeader("To", user.Email)
		mail.SetHeader("Subject", "Your daily recipes are here!")
		mail.SetBody("text/html", t.String())

		emailList = append(emailList, mail)
	}

	s.cron.AddFunc("0 14 * * *", func() {
		err := s.mailer.DialAndSend(emailList...)

		log.Printf("this is from the cron job")
		if err != nil {
			log.Printf("there was an error sending the mail: %v", err.Error())
		}
	})

	s.cron.Start()
	return nil
}
