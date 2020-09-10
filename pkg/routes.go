package app

import "github.com/gin-gonic/gin"

func (s *Server) routes() *gin.Engine {
	r := s.router

	v1 := r.Group("/v1/api")
	{
		v1.GET("/status", s.ApiStatus())
	}

	return r
}
