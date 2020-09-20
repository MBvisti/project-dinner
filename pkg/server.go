package app

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
	"strconv"
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
	// TODO: change this to setup the main cronjob
	err := s.CronMailer()
	if err != nil {
		log.Printf("this is err from cronjob: %v", err)
		return err
	}

	// TODO: change this when no longer needed
	isStaging, err := strconv.ParseBool(os.Getenv("IS_STAGING"))
	if err != nil {
		return err
	}

	if isStaging {
		err = s.storage.DestructiveReset()
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
	//address := []string{"morten@mbvistisen.dk", "vistisen@live.dk", "mbv1406@gmail.com"}
	mail := gomail.NewMessage()
	mail.SetAddressHeader("From", "noreply@mbvistisen.dk", "CronJob")
	mail.SetHeader("To", "mbv1406@gmail.com")
	mail.SetHeader("Subject", "test cron job at 4pm")
	mail.SetBody("text/html", "This is a test email sent everyday at 4pm by the cronjob")

	mails := []*gomail.Message{mail}

	s.cron.AddFunc("0 16 * * *", func() {
		err := s.mailer.DialAndSend(mails...)

		log.Printf("this is from the cron job")
		if err != nil {
			log.Printf("there was an error sending the mail: %v", err)
		}
	})

	s.cron.Start()
	return nil
}
