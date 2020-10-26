package api

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"reflect"
	"time"
)

// UserService is the user service's interface
type SpiderService interface {
	// FindLinks searches a page for anchor tags matching a query pattern
	FindLinks(pageURL, querySelector string) []string
	// CollectRecipe takes in a url and collects the recipe by looking for json+ld
	CollectRecipe(URL string) (*RecipeSection, error)
	// Go ...
	Go(baseURL string) error
}

// Repository ...
type SpiderRepository interface {
}

type spiderService struct {
	storage SpiderRepository
	crawler *colly.Collector
}

// NewUserService ...
func NewSpiderService(r SpiderRepository, c *colly.Collector) SpiderService {
	return &spiderService{
		r,
		c,
	}
}

func (s *spiderService) Go(baseURL string) error {
	return nil
}

// CrawlUrls finds links and returns an array of them for a given page
func (s *spiderService) FindLinks(pageURL, querySelector string) []string {
	s.crawler.Limit(&colly.LimitRule{
		DomainGlob:  "https://thecleaneatingcouple.com/*",
		RandomDelay: 4 * time.Second,
	})

	// Before making a request print "Visiting ..."
	s.crawler.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	var linkList []string

	//a[class='entry-image-link']
	s.crawler.OnHTML(querySelector, func(e *colly.HTMLElement) {
		link := e.Attr("href")

		linkList = append(linkList, link)
	})

	s.crawler.Visit(pageURL)

	return linkList
}

// CollectRecipe actually scrapes a website
func (s *spiderService) CollectRecipe(url string) (*RecipeSection, error) {
	s.crawler.Limit(&colly.LimitRule{
		DomainGlob:  "https://thecleaneatingcouple.com/*",
		RandomDelay: 5 * time.Second,
	})

	var crawlerResult BaseSection
	var WholeBody map[string]interface{}

	s.crawler.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		err := json.Unmarshal([]byte(e.Text), &WholeBody)
		json.Unmarshal([]byte(e.Text), &crawlerResult)

		if err != nil {
			log.Printf("this is from the crawler error: %v", err.Error())
		}
	})

	err := s.crawler.Visit(url)

	if err != nil {
		log.Printf("this is from the crawler visit error: %v", err)
		return nil, err
	}

	var ScapedRecipe RecipeSection

	for _, section := range crawlerResult.Graph {
		if section.Type != "Recipe" {
			continue
		}

		if reflect.DeepEqual(section.Description, "") {
			fmt.Print("desc empty")
		}

		var recipeInstructions []Instruction

		if len(section.RecipeInstructions[0].ItemListElement) == 0 {
			for index, instruction := range section.RecipeInstructions {

				r := Instruction{
					Step: index,
					Text: instruction.Name,
				}

				recipeInstructions = append(recipeInstructions, r)
			}
		}

		if len(section.RecipeInstructions[0].ItemListElement) != 0 {
			for index, ingredient := range section.RecipeInstructions {

				r := Instruction{
					Step: index,
					Text: ingredient.ItemListElement[0].Text,
				}

				recipeInstructions = append(recipeInstructions, r)
			}
		}
		//
		//recipeKeywords := strings.Split(section.Keywords, ",")
		//
		//yield, err := strconv.Atoi(section.RecipeYield[0])
		//
		//if err != nil {
		//	yield = 0
		//}

		//recipe := repository.Recipe{
		//	Name:         section.Name,
		//	Category:     strings.ToLower(section.RecipeCategory[0]),
		//	Cuisine:      strings.ToLower(section.RecipeCuisine[0]),
		//	Description:  section.Description,
		//	Images:       section.Image,
		//	Ingredients:  section.RecipeIngredients,
		//	Instructions: recipeInstructions,
		//	Keywords:     recipeKeywords,
		//	Score: repository.Rating{
		//		Score: section.AggregatedRating.RatingValue,
		//		Votes: section.AggregatedRating.RatingCount,
		//	},
		//	FoundOn: url,
		//	Yield:   yield,
		//}
		//
		//ScapedRecipe = recipe
	}

	return &ScapedRecipe, nil
}
