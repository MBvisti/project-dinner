package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"status": "success",
			"data":   "project dinner api running smootly",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) ResetDatabase() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		developmentMode, err := strconv.ParseBool(os.Getenv("DEVELOPMENT_MODE"))

		if err != nil {
			response := map[string]string{
				"status": "failure",
				"data":   "couldn't determine the development mode",
			}

			c.JSON(http.StatusInternalServerError, response)
		}

		if !developmentMode {
			response := map[string]string{
				"status": "failure",
				"data":   "don't reset the db in production idiot",
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		err = s.storage.DestructiveReset()

		if err != nil {
			response := map[string]string{
				"status": "failure",
				"data":   "couldn't reset the database",
			}

			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := map[string]string{
			"status": "success",
			"data":   "reset the database",
		}
		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) StopCronJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		s.cron.Stop()

		response := map[string]string{
			"status": "success",
			"data":   "cron job stopped",
		}
		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) CreateRecipe() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var recipe Recipe

		err := c.ShouldBindJSON(&recipe)

		if err != nil {
			response := map[string]string{
				"status": "failure",
				"data":   "recipe not created",
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		err = s.storage.CreateRecipe(&recipe)

		if err != nil {
			response := map[string]string{
				"status": "failure",
				"data":   "recipe not created",
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := map[string]string{
			"status": "success",
			"data":   "recipe created",
		}
		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) WakeDyno() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"status": "success",
			"data":   "dyno is awake again",
		}
		c.JSON(http.StatusOK, response)
	}
}
