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
		subscribeView = views.NewView("base", wd+"/pkg/views/todays_recipes.gohtml")

		views.Must(subscribeView.Render(c.Writer, data))
	}
}
