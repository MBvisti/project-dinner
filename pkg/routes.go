package app

import "github.com/gin-gonic/gin"

func (s *Server) routes() *gin.Engine {
	r := s.router

	v1 := r.Group("/v1/api")
	{
		v1.GET("/status", s.ApiStatus())
		v1.POST("/reset", s.ResetDatabase())
		v1.POST("/stop-cron", s.StopCronJob())
		v1.POST("/recipe", s.CreateRecipe())
		v1.GET("/wakey-wakey", s.WakeDyno())
		v1.GET("/random-recipes", s.GetFourRandomRecipes())
		v1.GET("/users", s.EmailList())
	}

	return r
}
