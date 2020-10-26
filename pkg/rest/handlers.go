package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	service "project-dinner/pkg/api"
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

// // StopCronJob stops cron object running on the server struct
// func (s *Server) StopCronJob() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Content-Type", "application/json")

// 		s.cron.Stop()

// 		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "the cron job is stopped"})
// 	}
// }

// // AllRecipes ...
// type AllRecipes struct {
// 	Recipes []repository.Recipe `json:"recipes"`
// }

// // UserRecipe ... TODO: move this in to a user service at one point
// type UserRecipe struct {
// 	UserName string
// 	Recipes  []repository.EmailRecipe
// }

// SendMails send out recipe emails
func SendMails(u service.EmailService) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := u.SendRecipes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "all emails sent"})

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

func GetRes(e service.EmailService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		text := e.SendRecipes()
		c.JSON(http.StatusOK, gin.H{"is it working": text})
	}
}

// SignupUser endpoint
func SignupUser(u service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var user service.User

		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, handlerResponse{Status: "failure", Data: "user data no good"})
			return
		}

		err = u.CreateUser(user)

		if err != nil {
			c.JSON(http.StatusBadRequest, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "user created"})
	}
}
