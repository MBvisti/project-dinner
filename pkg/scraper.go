package app

import (
	"encoding/json"
	"fmt"
	"log"
	"project-dinner/pkg/repository"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// Base is the uppermost level of the recipe structure defined by google
type Base struct {
	Context string          `json:"@context"`
	Graph   []RecipeSection `json:"@graph"`
}

// RecipeSection is a scraped recipe defined in ld+json format
type RecipeSection struct {
	Type               string               `json:"@type"`
	Name               string               `json:"name"`
	ID                 string               `json:"@id"`
	AggregatedRating   RatingSection        `json:"aggregateRating"`
	Description        string               `json:"description"`
	RecipeIngredients  []string             `json:"recipeIngredient"`
	Image              []string             `json:"image"`
	Nutrition          map[string]string    `json:"nutrition"`
	Keywords           string               `json:"keywords"`
	RecipeCategory     []string             `json:"recipeCategory"`
	RecipeCuisine      []string             `json:"recipeCuisine"`
	RecipeInstructions []RecipeInstructions `json:"recipeInstructions"`
	RecipeYield        []string             `json:"recipeYield"`
	FoundOn            string
}

// RecipeInstructions is a scaped recipe's instructions defined in ld+json format
type RecipeInstructions struct {
	Type            string            `json:"@type"`
	Name            string            `json:"name,omitempty"`
	URL             string            `json:"url,omitempty"`
	Text            string            `json:"text,omitempty"`
	ItemListElement []ItemListElement `json:"itemListElement,omitempty"`
}

// ItemListElement ...
type ItemListElement struct {
	Type string `json:"@type"`
	Text string `json:"text"`
}

// RatingSection ...
type RatingSection struct {
	RatingCount string `json:"ratingCount"`
	RatingValue string `json:"ratingValue"`
}

// FindLinks finds links and returns an array of them for a given page
func (s *Server) FindLinks(pageURL, queryString, domainGlob string) []string {
	crawler := colly.NewCollector()

	crawler.Limit(&colly.LimitRule{
		DomainGlob:  domainGlob,
		RandomDelay: 4 * time.Second,
	})

	// Before making a request print "Visiting ..."
	crawler.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	var linkList []string

	crawler.OnHTML(queryString, func(e *colly.HTMLElement) {
		link := e.Attr("href")

		linkList = append(linkList, link)
	})

	crawler.Visit(pageURL)

	return linkList
}

// ScrapeRecipe collects data about a recipe
func (s *Server) ScrapeRecipe(url, queryString string) (*repository.Recipe, error) {
	crawler := colly.NewCollector()

	crawler.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
	})

	var crawlerResult Base
	var WholeBody map[string]interface{}

	crawler.OnHTML(queryString, func(e *colly.HTMLElement) {
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

	var ScapedRecipe repository.Recipe

	for _, section := range crawlerResult.Graph {
		if section.Type != "Recipe" {
			continue
		}

		if reflect.DeepEqual(section.Description, "") {
			fmt.Print("desc empty")
		}

		var recipeInstructions []repository.Instruction

		if len(section.RecipeInstructions[0].ItemListElement) == 0 {
			for index, instruction := range section.RecipeInstructions {

				r := repository.Instruction{
					Step: index,
					Text: instruction.Name,
				}

				recipeInstructions = append(recipeInstructions, r)
			}
		}

		if len(section.RecipeInstructions[0].ItemListElement) != 0 {
			for index, ingredient := range section.RecipeInstructions {

				r := repository.Instruction{
					Step: index,
					Text: ingredient.ItemListElement[0].Text,
				}

				recipeInstructions = append(recipeInstructions, r)
			}
		}

		recipeKeywords := strings.Split(section.Keywords, ",")

		yield, err := strconv.Atoi(section.RecipeYield[0])

		if err != nil {
			yield = 0
		}

		recipe := repository.Recipe{
			Name:         section.Name,
			Category:     strings.ToLower(section.RecipeCategory[0]),
			Cuisine:      strings.ToLower(section.RecipeCuisine[0]),
			Description:  section.Description,
			Images:       section.Image,
			Ingredients:  section.RecipeIngredients,
			Instructions: recipeInstructions,
			Keywords:     recipeKeywords,
			Score: repository.Rating{
				Score: section.AggregatedRating.RatingValue,
				Votes: section.AggregatedRating.RatingCount,
			},
			FoundOn: url,
			Yield:   yield,
		}

		ScapedRecipe = recipe
	}

	return &ScapedRecipe, nil
}
