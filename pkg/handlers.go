package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
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

		//err = s.storage.CreateRecipe(&recipe)

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

type AllRecipes struct {
	Recipes []Recipe `json:"recipes"`
}

func (s *Server) GetFourRandomRecipes() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		resp, err := http.Get("https://api.spoonacular.com/recipes/random?apiKey=5ce66a1c4dc546f2a512059d8df566f7&tags=vegetarian,dinner&number=4")

		if err != nil {
			response := map[string]string{
				"status": "failure",
				"data":   err.Error(),
			}
			c.JSON(http.StatusInternalServerError, response)
		}

		var recipe AllRecipes

		if err := json.NewDecoder(resp.Body).Decode(&recipe); err != nil {
			response := map[string]string{
				"status": "failure",
				"data":   err.Error(),
			}
			c.JSON(http.StatusInternalServerError, response)
		}

		log.Print(recipe)
		err = s.storage.CreateRecipe(recipe.Recipes)

		if err != nil {
			log.Printf("there was an error saving the recipe to the database: %v", err)
		}

		//response := map[string]string{
		//	"status": "success",
		//	"data":   "saved the recipes to the database",
		//}
		c.JSON(http.StatusOK, recipe.Recipes)
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

func (s *Server) EmailList() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		//users, err := s.storage.GetEmailList()

		recipes, err := s.storage.TodaysRecipes()

		if err != nil {
			response := map[string]string{
				"status": "failure",
				"data":   "couldn't retrive list",
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		c.JSON(http.StatusOK, recipes)
	}
}
