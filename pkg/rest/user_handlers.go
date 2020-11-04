package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project-dinner/pkg/api"
)

// SubscribeUser endpoint
func SubscribeUser(u api.UserService, e api.EmailService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var user api.User

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

		mail, err := e.CreateWelcomeMail(user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		err = e.MailSender(*mail)

		if err != nil {
			c.JSON(http.StatusInternalServerError, handlerResponse{Status: "failure", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, handlerResponse{Status: "success", Data: "user created"})
	}
}
