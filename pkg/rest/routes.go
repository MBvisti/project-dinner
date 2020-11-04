package rest

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"project-dinner/pkg/api"
	"project-dinner/pkg/repository"
)

// Routes set up the routes
func Routes(userService api.UserService, emailService api.EmailService, spiderService api.SpiderService, recipeService api.RecipeService, repo repository.Repository) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.StaticFS("/static", http.Dir("static"))

	// All api endpoints here
	v1 := r.Group("/v1/api")
	{
		v1.GET("/status", APIStatus())
		v1.POST("/subscribe", SubscribeUser(userService, emailService))
		v1.GET("/send-mails", SendMails(emailService))
		v1.GET("/start-scraping-procedure", StartSpider(spiderService))
		v1.POST("/seed-prod-db", SeedProdData(repo))
	}

	// All view endpoints here
	r.GET("/", RenderHome(recipeService))
	// Subscribe
	r.GET("/subscribe", RenderSubscribe())
	r.GET("/subscribe/success", RenderSubscribe())
	r.GET("/subscribe/failure", RenderSubscribe())
	// Today's recipes
	r.GET("/todays-recipes", RenderTodaysRecipes())

	return r
}
