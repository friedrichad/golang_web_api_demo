package service

import (
	"bytes"
	"html/template"
	"log"

	"github.com/friedrichad/golang_web_api_demo/backend/configs/mail"
	"github.com/friedrichad/golang_web_api_demo/backend/shared"
	"gopkg.in/gomail.v2"
)
type IMailService interface {
	SendWelcomeEmail(userID int, email string) error
}
func NewMailService() IMailService {
	return &MailService{}
}
type MailService struct{
}

func (m *MailService) SendWelcomeEmail(userID int, email string) error{
	tmpl, err := template.ParseFiles("template/welcome/welcome.html")
	if err != nil {
		log.Print("Error parsing template: ", err)
		return err
	}
	var body bytes.Buffer
	err = tmpl.Execute(&body, shared.WelcomeMailData{
		UserID: userID,
		Email: email,
	})
	if err != nil {
		log.Print("Error executing template: ", err)
		return err
	}
	message := gomail.NewMessage()
	message.SetHeader(
		"From",
		mail.Mail.From,
	)
	message.SetHeader(
		"To",
		email,
	)
	message.SetHeader(
		"Subject",
		"Welcome to our system",
	)
	message.SetBody(
	"text/html",
	body.String(),
)
	dialer := gomail.NewDialer(
		mail.Mail.Host,
		mail.Mail.Port,
		mail.Mail.Username,
		mail.Mail.Password,
	)
	return dialer.DialAndSend(message)
}