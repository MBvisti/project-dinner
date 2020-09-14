package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
)

type Server struct {
	storage *Repository
	crawler *colly.Collector
	router  *gin.Engine
}

func NewServer(s *Repository, r *gin.Engine, c *colly.Collector) Server {
	return Server{
		storage: s,
		router:  r,
		crawler: c,
	}
}

func (s *Server) Run(addr string) error {
	log.Printf("Starting the server on: %v", addr)
	err := http.ListenAndServe(addr, s.routes())

	if err != nil {
		return err
	}

	return nil
}
