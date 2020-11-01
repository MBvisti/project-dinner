package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project-dinner/pkg/api"
)

// SendMails send out recipe emails
func SendMails(u api.EmailService) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := u.SendRecipes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "all emails sent"})

	}
}
