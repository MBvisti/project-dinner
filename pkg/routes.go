package app

import "github.com/gin-gonic/gin"

func (s *Server) routes() *gin.Engine {
	r := s.router

	v1 := r.Group("/v1/api")
	{
		v1.GET("/status", s.ApiStatus())
		v1.POST("/reset", s.ResetDatabase())
		v1.POST("/stop-cron", s.StopCronJob())
	}

	return r
}
