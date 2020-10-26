package rest

import (
	service "project-dinner/pkg/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Routes set up the routes
func Routes(uS service.UserService, eS service.EmailService) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	v1 := r.Group("/v1/api")
	{
		v1.GET("/status", APIStatus())
		v1.GET("/test", GetRes(eS))
		// v1.POST("/stop-cron", StopCronJob())
		v1.GET("/send-mails", SendMails(eS))
		//v1.GET("/start-scraping-procedure", CrawlSite())
		v1.POST("/sign-up", SignupUser(uS))
	}

	return r
}
