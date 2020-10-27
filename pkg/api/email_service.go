package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"html/template"
	"time"
)

// EmailService is the email service's interface
type EmailService interface {
	SendRecipes() error
	EveryDayMailer() (cron.Job, error)
}

// EmailRepository ...
type EmailRepository interface {
	GetDailyRecipes() ([]EmailRecipe, error)
	GetEmailList() ([]User, error)
}

type emailService struct {
	mailProvider *gomail.Dialer
	storage      EmailRepository
}

// NewEmailService ...
func NewEmailService(mp *gomail.Dialer, r EmailRepository) EmailService {
	return &emailService{
		mp,
		r,
	}
}

var (
	ErrLoadLocation = errors.New("email service - could not load location")
	ErrCronJob      = errors.New("email service - could not run cron job")
)

// SendRecipes will email all users the daily recipes
func (e *emailService) SendRecipes() error {
	mailTemplate, err := template.ParseFiles("./template/daily_recipe_email.html")
	if err != nil {
		return err
	}

	dailyRecipes, err := e.storage.GetDailyRecipes()

	if err != nil {
		return err
	}

	emailList, err := e.storage.GetEmailList()

	if err != nil {
		return err
	}

	for _, user := range emailList {
		usrRecipe := UserRecipe{
			UserName: user.Name,
			Recipes:  dailyRecipes,
		}

		var t bytes.Buffer
		err = mailTemplate.Execute(&t, usrRecipe)
		mail := gomail.NewMessage()
		mail.SetAddressHeader("From", "noreply@mbvistisen.dk", "Morten's recipe service")
		mail.SetHeader("To", user.Email)
		mail.SetHeader("Subject", "Your daily recipes are here!")
		mail.SetBody("text/html", t.String())

		err = e.mailProvider.DialAndSend(mail)

		if err != nil {
			return err
		}
	}

	return nil
}

// EveryDayMailer calls SendMails everyday at 4pm
func (e *emailService) EveryDayMailer() (cron.Job, error) {
	fmt.Print("starting up everyday mailer...")

	// Set the location for the cron job
	t, err := time.LoadLocation("Europe/Copenhagen")

	if err != nil {
		return nil, ErrLoadLocation
	}

	// create instance of cron
	c := cron.New(cron.WithLocation(t))
	var cronJobError error

	// Send out recipes at 4pm Cph time
	c.AddFunc("0 14 * * *", func() {
		cronJobError = e.SendRecipes()
	})

	if cronJobError != nil {
		return nil, ErrCronJob
	}

	fmt.Print("everyday mailer started")

	return c, nil
}
