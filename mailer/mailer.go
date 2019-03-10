package mailer

import (
	"fmt"
	"github.com/A9u/function_junction/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	To      []string
	From    string
	Subject string
	Body    string
}

func getTos(emails []string) []*mail.Email {
	var tos = make([]*mail.Email, len(emails))

	fmt.Println(len(emails))
	for i := 0; i < len(emails); i++ {
		tos[i] = mail.NewEmail("", emails[i])
	}

	return tos
}

func (e *Email) Send() {
	from := mail.NewEmail("", e.From)
	tos := getTos(e.To)

	email := mail.NewV3Mail()
	p := mail.NewPersonalization()

	p.AddTos(tos...)
	p.Subject = e.Subject

	email.AddPersonalizations(p)

	body := "<p> Hi, </p>" + e.Body
	content := mail.NewContent("text/html", body)
	email.AddContent(content)
	email.SetFrom(from)

	client := sendgrid.NewSendClient(config.SmtpApiKey())

	response, err := client.Send(email)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
