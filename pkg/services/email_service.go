package service

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

type emailService struct {
	app *App
}

// NewEmailService ...
func NewEmailService(a *App) EmailService {
	return &emailService{
		app: a,
	}
}

func (e *emailService) SendRecipes() error {
	mailTemplate, err := template.ParseFiles("./template/daily_recipe_email.html")
	if err != nil {
		return err
	}

	dailyRecipes, err := e.app.rR.GetDailyRecipes()

	if err != nil {
		return err
	}

	emailList, err := e.app.uR.GetEmailList()

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

		err = e.app.mailProvider.DialAndSend(mail)

		if err != nil {
			return err
		}
	}

	return nil
}
