package app

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"project-dinner/pkg/repository"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/gin-gonic/gin"
)

type handlerResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// APIStatus returns the status of the api
func (s *Server) APIStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "project dinner api running smoothly"})
	}
}

// ResetDatabase endpoint used to reset DB in development
func (s *Server) ResetDatabase() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		whatEnv := os.Getenv("WHAT_ENVIRONMENT_IS_THIS")

		if whatEnv == "production" {
			c.JSON(http.StatusBadRequest, handlerResponse{Status: "failure", Data: "don't reset the database in production, idiot"})
			return
		}

		err := s.storage.DestructiveReset()

		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "successfully reset the database"})
	}
}

// StopCronJob stops cron object running on the server struct
func (s *Server) StopCronJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		s.cron.Stop()

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "the cron job is stopped"})
	}
}

// AllRecipes ...
type AllRecipes struct {
	Recipes []repository.Recipe `json:"recipes"`
}

// UserRecipe ... TODO: move this in to a user service at one point
type UserRecipe struct {
	UserName string
	Recipes  []repository.EmailRecipe
}

// SendMails send out recipe emails
func (s *Server) SendMails() gin.HandlerFunc {
	return func(c *gin.Context) {
		mailTemplate, err := template.ParseFiles("./template/daily_recipe_email.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		var emailList []*gomail.Message

		userList, err := s.storage.User.GetEmailList()

		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		recipes, err := s.storage.Recipe.GetRandomRecipes()

		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		for _, user := range userList {
			usrRecipe := UserRecipe{
				UserName: user.Name,
				Recipes:  recipes,
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

		err = s.mailer.DialAndSend(emailList...)

		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "all emails sent"})

	}
}

// CrawlSite endpint gets all links to a recipe on the clean eating couple, scrape the recipe and stores it in db
func (s *Server) CrawlSite() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		pages := 10
		var links []string
		baseURL := "https://thecleaneatingcouple.com/category/recipes/lunch-dinner/page/"

		for page := 1; page < pages; page++ {
			newLinks := s.CrawlUrls(baseURL + strconv.Itoa(page))
			links = append(links, newLinks...)
		}

		var returnedData []repository.Recipe
		for _, link := range links {
			res, err := s.Crawler(link)
			if err != nil {
				log.Printf("this is the crawler error: %v", err.Error())
				continue
			}

			if res.Description == "" {
				continue
			}
			returnedData = append(returnedData, res)
			err = s.storage.Recipe.CreateScrapedRecipe(res)

			if err != nil {
				log.Printf("this is the storage error: %v", err.Error())
				continue
			}
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "website crawled"})
	}
}

// NewUser ....
type NewUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	TimeZone string `json:"time_zone"`
}

// SignupUser endpoint
func (s *Server) SignupUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var user repository.NewUser

		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, handlerResponse{Status: "failure", Data: "user data no good"})
			return
		}

		user.Email = strings.ToLower(user.Email)
		user.Name = strings.ToLower(user.Name)

		err = s.storage.User.CreateUser(user)

		if err != nil {
			c.JSON(http.StatusBadRequest, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "user created"})
	}
}
