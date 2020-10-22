package app

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	repository "project-dinner/pkg/storage"
	"strconv"

	"gopkg.in/gomail.v2"

	"github.com/gin-gonic/gin"
)

// APIStatus returns the status of the api
func (s *Server) APIStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"status": "success",
			"data":   "project dinner api running super smootly",
		}

		c.JSON(http.StatusOK, response)
	}
}

// ResetDatabase endpoint used to reset DB in development
func (s *Server) ResetDatabase() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		developmentMode, err := strconv.ParseBool(os.Getenv("DEVELOPMENT_MODE"))

		if err != nil {
			response := map[string]string{
				"status": "failure",
				"data":   "couldn't determine the development mode",
			}

			c.JSON(http.StatusInternalServerError, response)
		}

		if !developmentMode {
			response := map[string]string{
				"status": "failure",
				"data":   "don't reset the db in production idiot",
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		err = s.storage.DestructiveReset()

		if err != nil {
			response := map[string]string{
				"status": "failure",
				"data":   "couldn't reset the database",
			}

			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := map[string]string{
			"status": "success",
			"data":   "reset the database",
		}
		c.JSON(http.StatusOK, response)
	}
}

// StopCronJob stops cron object running on the server struct
func (s *Server) StopCronJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		s.cron.Stop()

		response := map[string]string{
			"status": "success",
			"data":   "cron job stopped",
		}
		c.JSON(http.StatusOK, response)
	}
}

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
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		}

		var emailList []*gomail.Message

		userList, err := s.storage.User.GetEmailList()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		}

		recipes, err := s.storage.Recipe.GetRandomRecipes()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
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
			log.Printf("there was an error sending the mail: %v", err.Error())
		}

		response := map[string]interface{}{
			"status":   "success",
			"response": "all email sent",
		}
		c.JSON(http.StatusOK, response)

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

		response := map[string]interface{}{
			"status":   "success",
			"response": "site successfully crawled",
		}
		c.JSON(http.StatusOK, response)
	}
}
