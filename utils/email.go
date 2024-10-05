package utils

import (
	"fmt"
	"net/smtp"

	"github.com/thoriqdharmawan/be-question-generator/config"
)

func SendEmail(to []string, subject, body string) error {

	from := config.Conf.SmtpEmail
	password := config.Conf.SmtpPassword

	smtpHost := config.Conf.SmtpHost
	smtpPort := config.Conf.SmtpPort

	// Setup authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
}
