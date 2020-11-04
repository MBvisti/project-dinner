package rest

import (
	"github.com/gin-gonic/gin"
	"os"
	"project-dinner/pkg/views"
)

func RenderHome() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		var homeView *views.View

		type TestData struct {
			Name string
			Msg  string
		}
		data := TestData{
			Name: "MBV",
			Msg:  "Its working! Its working!!",
		}

		wd, _ := os.Getwd()
		homeView = views.NewView("base", wd+"/pkg/views/home.gohtml")

		views.Must(homeView.Render(c.Writer, data))
	}
}

func RenderSubscribe() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		var subscribeView *views.View

		type TestData struct {
			Name string
			Msg  string
		}
		data := TestData{
			Name: "MBV",
			Msg:  "Its working! Its working!!",
		}

		wd, _ := os.Getwd()
		subscribeView = views.NewView("base", wd+"/pkg/views/subscribe.gohtml")

		views.Must(subscribeView.Render(c.Writer, data))
	}
}

func RenderTodaysRecipes() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		var todaysRecipesView *views.View

		type TestData struct {
			Name string
			Msg  string
		}
		data := TestData{
			Name: "MBV",
			Msg:  "Its working! Its working!!",
		}

		wd, _ := os.Getwd()
		todaysRecipesView = views.NewView("base", wd+"/pkg/views/todays_recipes.gohtml")

		views.Must(todaysRecipesView.Render(c.Writer, data))
	}
}

func RenderAbout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		var aboutView *views.View

		type TestData struct {
			Name string
			Msg  string
		}
		data := TestData{
			Name: "MBV",
			Msg:  "Its working! Its working!!",
		}

		wd, _ := os.Getwd()
		aboutView = views.NewView("base", wd+"/pkg/views/about.gohtml")

		views.Must(aboutView.Render(c.Writer, data))
	}
}
