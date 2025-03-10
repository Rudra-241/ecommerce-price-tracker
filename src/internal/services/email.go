package services

import (
	"ecommerce-price-tracker/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
	"gorm.io/gorm"
	"log"
	"os"
)

func sendMail(recipient string, body string) error {
	message := gomail.NewMessage()

	message.SetHeader("From", os.Getenv("ALERT_EMAIL_ADDRESS"))

	if gin.Mode() == "debug" {
		log.Print("Mail for ", recipient, " but forwarding to ", os.Getenv("MAILTRAP_RECIPIENT"))
		message.SetHeader("To", os.Getenv("MAILTRAP_RECIPIENT"))
	} else {
		message.SetHeader("To", recipient)
	}
	message.SetHeader("Subject", "Price Alert")

	message.SetBody("text/html", body)
	api := os.Getenv("MAILTRAP_API")
	if api == "" {
		log.Print("MAILTRAP_API not set")
	}
	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", api)

	if err := dialer.DialAndSend(message); err != nil {
		log.Fatal("Error:", err)
	}
	return nil
}

/*
 * ISSUE: This function currently sends emails every 2 hours to all users
 * where the price direction is either below the starting price or decreasing.
 * This is repetitive
 * Additionally, the "below start" condition is measured from when the product
 * was first tracked, making it non-user-specific.
 */

func getDirectionColor(direction models.PriceState) string {
	switch direction {
	case models.BelowStart:
		return "red"
	case models.Decreased:
		return "orange"
	default:
		return "black"
	}
}

func EmailAll(db *gorm.DB) error {
	fmt.Println("EmailAll")
	var products []models.Product
	if err := db.Preload("UsersTracking").
		Where("direction IN ?", []models.PriceState{models.BelowStart, models.Decreased}).
		Find(&products).Error; err != nil {
		return err
	}
	for _, product := range products {
		productHTML := fmt.Sprintf(`
			<tr>
				<td><img src="%s" alt="%s" width="100"></td>
				<td><a href="%s">%s</a></td>
				<td>₹%.2f</td>
				<td style="color: %s;"><strong>%s</strong></td>
			</tr>`,
			product.ImgLink, product.Name, product.Url, product.Name, product.Price,
			getDirectionColor(product.Direction), product.Direction)

		for _, user := range product.UsersTracking {
			emailBody := fmt.Sprintf(`
				<html>
				<body>
					<p>Hello %s,</p>
					<p>Here is an update on the products you are tracking:</p>
					<table border="1" cellpadding="5" cellspacing="0">
						<tr>
							<th>Image</th>
							<th>Product</th>
							<th>Price</th>
							<th>Status</th>
						</tr>
						%s
					</table>
					<p>Best Regards,<br>Your Price Tracker</p>
				</body>
				</html>`, user.Email, productHTML)

			err := sendMail(user.Email, emailBody)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
