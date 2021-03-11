package mailer

import (
	"fmt"

	"github.com/A9u/function_junction/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailerService interface {
	Send(to []string, from, subject, body string) (err error)
}

type sendgridMailer struct{}

func NewMailer() MailerService {
	return &sendgridMailer{}
}

func (sm *sendgridMailer) Send(to []string, from, subject, body string) (err error) {
	email := mail.NewV3Mail()
	p := mail.NewPersonalization()

	tos := getTos(to)
	p.AddTos(tos...)

	p.Subject = subject

	email.AddPersonalizations(p)
	fromEmail := mail.NewEmail("", from)

	htmlBody := "<p> Hi, </p>" + body
	content := mail.NewContent("text/html", htmlBody)
	email.AddContent(content)
	email.SetFrom(fromEmail)

	client := sendgrid.NewSendClient(config.SmtpApiKey())

	response, err := client.Send(email)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return
}

func getTos(emails []string) []*mail.Email {
	var tos = make([]*mail.Email, len(emails))

	fmt.Println(len(emails))
	// TODO: for i, e := range emails{}

	for i, email := range emails {
		tos[i] = mail.NewEmail("", email)
	}

	return tos
}

// TODO: always send err from a method if there are any
// caller can always decide whether to use them or not
