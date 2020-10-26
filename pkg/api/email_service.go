package api

import (
	"bytes"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

// EmailService is the email service's interface
type EmailService interface {
	SendRecipes() error
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

	log.Print(emailList)

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
