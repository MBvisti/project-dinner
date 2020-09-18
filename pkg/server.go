package app

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
)

type Server struct {
	storage *Repository
	router  *gin.Engine
	cron    *cron.Cron
	mailer  *gomail.Dialer
}

func NewServer(s *Repository, r *gin.Engine, c *cron.Cron, m *gomail.Dialer) Server {
	return Server{
		storage: s,
		router:  r,
		cron:    c,
		mailer:  m,
	}
}

func (s *Server) Run(addr string) error {
	log.Print("Calling cronjob")
	err := s.CronMailer()

	if err != nil {
		log.Printf("this is err from cronjob: %v", err)
		return err
	}

	log.Printf("Starting the server on: %v", addr)
	err = http.ListenAndServe(addr, s.routes())

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) CronMailer() error {
	mail := gomail.NewMessage()
	mail.SetAddressHeader("From", "noreply@mbvistisen.dk", "CronJob")
	mail.SetHeader("To", "morten@mbvistisen.dk")
	mail.SetHeader("Subject", "test cron job")
	mail.SetBody("text/html", "This is a test email sent every 5 minute by the cronjob")

	mailTwo := gomail.NewMessage()
	mailTwo.SetAddressHeader("From", "noreply@mbvistisen.dk", "CronJob")
	mailTwo.SetHeader("To", "mbv1406@gmail.com")
	mailTwo.SetHeader("Subject", "test cron job")
	mailTwo.SetBody("text/html", "This is a test email sent every 5 minute by the cronjob")

	s.cron.AddFunc("*/2 * * * *", func() {
		err := s.mailer.DialAndSend(mail, mailTwo)

		if err != nil {
			log.Printf("there was an error sending the mail: %v", err)
		}
	})

	s.cron.Start()
	return nil
}
