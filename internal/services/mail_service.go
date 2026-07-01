package services

import (
	"net/smtp"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
)

func sendEmail(serverConfig *config.ServerConfig, to string, subject string, body string) error {
	from := serverConfig.MailFrom
	password := serverConfig.MailPassword
	smtpHost := serverConfig.MailSmtpHost
	smtpPort := serverConfig.MailSmtpPort

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}
