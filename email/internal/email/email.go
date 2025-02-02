package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func Send(target string, orderID string) error {
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPassword := os.Getenv("SENDER_PASSWORD")

	recipientEmail := target

	message := []byte(fmt.Sprintf("Subject: Payment Proccessed! \n Process ID: %s", orderID))

	smtpServer := "smtp.gmail.com"
	smtpPort := 587

	creds := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	smtpAddress := fmt.Sprintf("%s:%d", smtpServer, smtpPort)

	err := smtp.SendMail(smtpAddress, creds, senderEmail, []string{recipientEmail}, message)
	if err != nil {
		return err
	}

	return nil
}
