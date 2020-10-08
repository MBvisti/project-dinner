package app

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	storage *Repository
	router  *gin.Engine
	cron    *cron.Cron
	mailer  *gomail.Dialer
}

type UserRecipe struct {
	UserName string
	Recipes  []DailyRecipes
}

func NewServer(s *Repository, r *gin.Engine, c *cron.Cron, m *gomail.Dialer) Server {
	return Server{
		storage: s,
		router:  r,
		cron:    c,
		mailer:  m,
	}
}

func (s *Server) Run(addr string) error {
	// TODO: change this to setup the main cronjob
	err := s.CronMailer()
	if err != nil {
		log.Printf("this is err from cronjob: %v", err)
		return err
	}

	if err != nil {
		log.Printf("this is err from cronjob: %v", err)
		return err
	}

	// TODO: change this when no longer needed
	isStaging, err := strconv.ParseBool(os.Getenv("IS_STAGING"))
	if err != nil {
		return err
	}

	if isStaging {
		err = s.storage.DestructiveReset()
		if err != nil {
			return err
		}
	}

	log.Printf("Starting the server on: %v", addr)
	err = http.ListenAndServe(addr, s.routes())

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) GetInitialRecipes() error {
	resp, err := http.Get("https://api.spoonacular.com/recipes/random?apiKey=5ce66a1c4dc546f2a512059d8df566f7&tags=vegetarian,dinner&number=4")

	if err != nil {
		return err
	}

	var recipe AllRecipes

	if err := json.NewDecoder(resp.Body).Decode(&recipe); err != nil {
		return err
	}

	err = s.storage.CreateRecipe(recipe.Recipes)

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) CronMailer() error {
	mailTemplate, err := template.ParseFiles("../template/daily_recipe_email.html")
	if err != nil {
		return err
	}

	var emailList []*gomail.Message

	userList, err := s.storage.GetEmailList()

	if err != nil {
		return err
	}

	dailyRecipes, err := s.storage.TodaysRecipes()

	if err != nil {
		return err
	}

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

	s.cron.AddFunc("0 13 * * *", func() {
		err := s.mailer.DialAndSend(emailList...)

		log.Printf("this is from the cron job")
		if err != nil {
			log.Printf("there was an error sending the mail: %v", err.Error())
		}
	})

	s.cron.Start()
	return nil
}
