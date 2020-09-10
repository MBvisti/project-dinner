package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Server struct {
	storage *Repository
	router  *gin.Engine
}

func NewServer(s *Repository, r *gin.Engine) Server {
	return Server{
		storage: s,
		router:  r,
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
