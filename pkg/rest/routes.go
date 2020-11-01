package rest

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	service "project-dinner/pkg/api"
)

// Routes set up the routes
func Routes(userService service.UserService, emailService service.EmailService, spiderService service.SpiderService) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	// All api endpoints here
	v1 := r.Group("/v1/api")
	{
		v1.GET("/status", APIStatus())
		v1.POST("/sign-up", SignupUser(userService, emailService))
		v1.GET("/send-mails", SendMails(emailService))
		v1.GET("/start-scraping-procedure", StartSpider(spiderService))
	}

	// All view endpoints here
	r.GET("/", RenderHome())
	r.GET("/signup", RenderSignup())

	return r
}
