package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
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
	// err := s.CronMailer()

	// if err != nil {
	// 	log.Printf("this is err from cronjob: %v", err)
	// 	return err
	// }

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

func (s *Server) GetDailyRecipes() error {
	resp, err := http.Get("https://api.spoonacular.com/recipes/random?apiKey=5ce66a1c4dc546f2a512059d8df566f7&tags=vegetarian,dinner&number=4")

	if err != nil {
		return err
	}

	var recipe AllRecipes

	if err := json.NewDecoder(resp.Body).Decode(&recipe); err != nil {
		return err
	}

	// err = s.storage.CreateRecipe(recipe.Recipes)

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

	s.cron.AddFunc("0 12 * * *", func() {
		err := s.GetDailyRecipes()

		if err != nil {
			log.Printf("there was an error getting recipes: %v", err.Error())
		}
	})

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

type Base struct {
	Context string                 `json:"@context"`
	Graph   []ScrapedRecipeSection `json:"@graph"`
}

type ScrapedRecipeSection struct {
	Type               string                      `json:"@type"`
	ID                 string                      `json:"@id"`
	AggregatedRating   ScrapedRatingSection        `json:"aggregateRating"`
	Description        string                      `json:"description"`
	RecipeIngredients  []string                    `json:"recipeIngredient"`
	Image              []string                    `json:"image"`
	Nutrition          map[string]string           `json:"nutrition"`
	Keywords           string                      `json:"keywords"`
	RecipeCategory     []string                    `json:"recipeCategory"`
	RecipeCuisine      []string                    `json:"recipeCuisine"`
	RecipeInstructions []ScrapedRecipeInstructions `json:"recipeInstructions"`
	RecipeYield        []string                    `json:"recipeYield"`
}

type ScrapedRecipeInstructions struct {
	Type            string              `json:"@type"`
	Name            string              `json:"name,omitempty"`
	URL             string              `json:"url,omitempty"`
	Text            string              `json:"text,omitempty"`
	ItemListElement []map[string]string `json:"itemListElement,omitempty"`
}

type ScrapedRatingSection struct {
	RatingCount string `json:"ratingCount"`
	RatingValue string `json:"ratingValue"`
}

func (s *Server) CrawlUrls() []string {
	crawler := colly.NewCollector()

	crawler.Limit(&colly.LimitRule{
		DomainGlob: "https://thecleaneatingcouple.com/*",
		Delay:      4 * time.Second,
	})
	var linkList []string

	crawler.OnHTML("a[class='entry-image-link']", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		linkList = append(linkList, link)
	})

	crawler.Visit("https://thecleaneatingcouple.com/category/recipes/lunch-dinner/")

	return linkList
}

func (s *Server) Crawler(url string) (*ScrapedRecipeSection, error) {
	crawler := colly.NewCollector()
	crawler.Limit(&colly.LimitRule{
		DomainGlob: "https://thecleaneatingcouple.com/*",
		Delay:      2 * time.Second,
	})
	var crawlerResult Base
	var WholeBody map[string]interface{}

	crawler.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		err := json.Unmarshal([]byte(e.Text), &WholeBody)
		json.Unmarshal([]byte(e.Text), &crawlerResult)

		if err != nil {
			log.Printf("this is from the crawler error: %v", err.Error())
		}
	})

	err := crawler.Visit(url)

	if err != nil {
		log.Printf("this is from the crawler visit error: %v", err)
		return nil, err
	}

	var ScapedRecipe ScrapedRecipeSection

	for _, section := range crawlerResult.Graph {
		if section.Type != "Recipe" {
			continue
		}

		log.Printf("this is section: %v", section)

		ScapedRecipe = section
	}

	return &ScapedRecipe, nil
}
