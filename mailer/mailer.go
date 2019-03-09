package mailer

import (
	"fmt"
	"github.com/A9u/function_junction/config"
	"net/smtp"
)

func NotifyAll() {
	smtpConfig := config.Smtp()

	fmt.Println("inside notify member")
	var (
		from       = "anusha@joshsoftware.com"
		recipients = []string{"anusha+test@joshsoftware.com"}
	)

	msg := "From: " + from + "\n" +
		"To: " + "anusha+test@joshsoftware.com" + "\n" +
		"Subject: Hello there\n\n" +
		"This is first email from Golang"

	hostname := smtpConfig.Domain()
	fmt.Println("inside notify member 2")

	auth := smtp.PlainAuth("", smtpConfig.Username(), smtpConfig.Password(), hostname)
	fmt.Println("inside notify member 3")

	err := smtp.SendMail(hostname+":"+smtpConfig.Port(), auth, from, recipients, []byte(msg))
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
}
