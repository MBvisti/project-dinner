package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

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

type Section struct {
	Type              string        `json:"@type"`
	ID                string        `json:"@id"`
	AggregatedRating  RatingSection `json:"aggregateRating"`
	Description       string        `json:"description"`
	RecipeIngredients []string      `json:"recipeIngredient"`
	//RecipeIngredientsTwo []map[string]string `json:"recipeIngredient"`
}

type RatingSection struct {
	RatingCount string `json:"ratingCount"`
	RatingValue string `json:"ratingValue"`
}

type RecipeIngredient struct {
	Ingredient string
}

func (s *Server) Crawler(url string) (interface{}, error) {
	log.Printf("crawling this url: %v", url)
	crawler := colly.NewCollector()
	var crawlerResult map[string]json.RawMessage

	crawler.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		err := json.Unmarshal([]byte(e.Text), &crawlerResult)

		if err != nil {
			log.Printf("this is from the crawler error: %v", err.Error())
		}
	})

	err := crawler.Visit(url)

	if err != nil {
		log.Printf("this is from the crawler visit error: %v", err)
		return "", err
	}

	//log.Printf("this is the conversion test: %v", string(crawlerResult["@graph"]))

	//log.Printf("this is crawlerResult[@graph]: %v", crawlerResult["@graph"])

	var TestStruct []map[string]interface{}
	err = json.Unmarshal(crawlerResult["@graph"], &TestStruct)

	if err != nil {
		log.Printf("unmarshalling error one: %v", err)
	}

	var AllSections Section
	for _, section := range crawlerResult["@graph"] {
		log.Printf("this is the all sections loop value: %v", section)
		err = json.Unmarshal([]byte(string(section)), &AllSections)
	}

	if err != nil {
		log.Printf("unmarshalling error two: %v", err)
	}

	log.Printf("this is AllSections: %v", AllSections)

	//for _, element := range TestStruct {
	//	for _, elementInner := range element {
	//		//log.Printf("all inner elements: %v", elementInner)
	//		if elementInner == "Recipe" {
	//			log.Printf("This is the inner element: %v", elementInner)
	//			log.Printf("this is the wanted element: %v", element)
	//		}
	//	}
	//}

	//log.Printf("this is the test struct: %v", TestStruct)

	//log.Printf("this is the whole object: %v", crawlerResult)

	//for _, element := range crawlerResult["@graph"] {
	//	log.Printf("this is the element: %v", element)
	//}

	return AllSections, nil
}
