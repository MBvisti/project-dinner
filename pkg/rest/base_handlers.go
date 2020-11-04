package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project-dinner/pkg/api"
	"project-dinner/pkg/repository"
)

type handlerResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// APIStatus returns the status of the api
func APIStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "project dinner api running smoothly"})
	}
}

func SeedProdData(r repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		err := r.SeedProductionData()

		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: "could not seed database"})
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "production db seeded"})
	}
}

// // AllRecipes ...
// type AllRecipes struct {
// 	Recipes []repository.Recipe `json:"recipes"`
// }

// StartSpider endpoint
func StartSpider(s api.SpiderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		err := s.Go("")

		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: "could not crawl site"})
			return
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "user created"})
	}
}

// CrawlSite endpint gets all links to a recipe on the clean eating couple, scrape the recipe and stores it in db
//func CrawlSite() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.Header("Content-Type", "application/json")
//
//		pages := 10
//		var links []string
//		baseURL := "https://thecleaneatingcouple.com/category/recipes/lunch-dinner/page/"
//
//		for page := 1; page < pages; page++ {
//			newLinks := s.CrawlUrls(baseURL + strconv.Itoa(page))
//			links = append(links, newLinks...)
//		}
//
//		var returnedData []repository.Recipe
//		for _, link := range links {
//			res, err := s.Crawler(link)
//			if err != nil {
//				log.Printf("this is the crawler error: %v", err.Error())
//				continue
//			}
//
//			if res.Description == "" {
//				continue
//			}
//			returnedData = append(returnedData, res)
//			err = s.storage.Recipe.CreateScrapedRecipe(res)
//
//			if err != nil {
//				log.Printf("this is the storage error: %v", err.Error())
//				continue
//			}
//		}
//
//		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "website crawled"})
//	}
//}
