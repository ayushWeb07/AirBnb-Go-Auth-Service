package services

import (
	"net/smtp"
)

func sendEmail(to, subject, body string) error {
	from := "bommanaayush07@gmail.com"
	password := "wjdk xcqc cnpz gamn"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}
