package utils

import (
	"crypto/rand"
	"encoding/hex"
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

func GenerateVerificationToken() string {
	bytes := make([]byte, 16) // 16 bytes = 128 bit
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
