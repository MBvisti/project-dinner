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

func RenderSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		var signupView *views.View

		type TestData struct {
			Name string
			Msg  string
		}
		data := TestData{
			Name: "MBV",
			Msg:  "Its working! Its working!!",
		}

		wd, _ := os.Getwd()
		signupView = views.NewView("base", wd+"/pkg/views/signup.gohtml")

		views.Must(signupView.Render(c.Writer, data))
	}
}
