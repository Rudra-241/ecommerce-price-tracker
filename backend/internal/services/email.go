package services

import (
	gomail "gopkg.in/mail.v2"
	"log"
	"os"
)

func sendMail(recipient string) {
	message := gomail.NewMessage()

	message.SetHeader("From", "pricetracker@demomailtrap.com")
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", "Price Alert")

	message.SetBody("text/plain", "Your tracked products prices are decreasing")
	api := os.Getenv("MAILTRAP_API")
	if api == "" {
		log.Print("MAILTRAP_API not set")
	}
	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", api)

	if err := dialer.DialAndSend(message); err != nil {
		log.Fatal("Error:", err)
	}
}

/*
 * ISSUE: This function currently sends emails every 2 hours to all users
 * where the price direction is either below the starting price or decreasing.
 * This is repetitive
 * Additionally, the "below start" condition is measured from when the product
 * was first tracked, making it non-user-specific.
 */
