package services

import (
	"bytes"
	"ecommerce-price-tracker/internal/models"
	"ecommerce-price-tracker/internal/queue"
	"encoding/json"
	"html/template"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
	"gorm.io/gorm"
)

const sendMailMaxAttempts = 3

type EmailAlert struct {
	Recipient string       `json:"recipient"`
	Product   ProductAlert `json:"product"`
}

type ProductAlert struct {
	Name      string  `json:"name"`
	Url       string  `json:"url"`
	ImgLink   string  `json:"img_link"`
	Price     float64 `json:"price"`
	Direction string  `json:"direction"`
}

func sendMail(recipient string, body string) error {
	message := gomail.NewMessage()

	message.SetHeader("From", os.Getenv("ALERT_EMAIL_ADDRESS"))

	if gin.Mode() == "debug" {
		slog.Info("forwarding mail in debug mode", "recipient", recipient, "forward_to", os.Getenv("MAILTRAP_RECIPIENT"))
		message.SetHeader("To", os.Getenv("MAILTRAP_RECIPIENT"))
	} else {
		message.SetHeader("To", recipient)
	}
	message.SetHeader("Subject", "Price Alert")

	message.SetBody("text/html", body)
	api := os.Getenv("MAILTRAP_API")
	if api == "" {
		slog.Warn("MAILTRAP_API not set")
	}
	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", api)

	if err := dialer.DialAndSend(message); err != nil {
		return err
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

func getDirectionColor(direction string) string {
	switch direction {
	case models.BelowRecentAvg:
		return "green"
	// case models.BelowStart:
	// 	return "red"
	case models.Decreased:
		return "orange"
	default:
		return "black"
	}
}

const defaultAlertTemplate = `<html>
<body>
	<p>Hello {{.Recipient}},</p>
	<p>Here is an update on the products you are tracking:</p>
	<table border="1" cellpadding="5" cellspacing="0">
		<tr>
			<th>Image</th>
			<th>Product</th>
			<th>Price</th>
			<th>Status</th>
		</tr>
		<tr>
			<td><img src="{{.Product.ImgLink}}" alt="{{.Product.Name}}" width="100"></td>
			<td><a href="{{.Product.Url}}">{{.Product.Name}}</a></td>
			<td>₹{{printf "%.2f" .Product.Price}}</td>
			<td style="color: {{directionColor .Product.Direction}};"><strong>{{.Product.Direction}}</strong></td>
		</tr>
	</table>
	<p>Best Regards,<br>Your Price Tracker</p>
</body>
</html>`

var (
	alertTmpl     *template.Template
	alertTmplOnce sync.Once
)

func emailTemplate() *template.Template {
	alertTmplOnce.Do(func() {
		fns := template.FuncMap{"directionColor": getDirectionColor}

		path := os.Getenv("EMAIL_TEMPLATE_PATH")
		if path == "" {
			path = "templates/price_alert.html"
		}

		if raw, err := os.ReadFile(path); err != nil {
			slog.Warn("email template not found, using built-in fallback", "path", path, "err", err)
		} else if t, err := template.New("price_alert").Funcs(fns).Parse(string(raw)); err != nil {
			slog.Warn("email template invalid, using built-in fallback", "path", path, "err", err)
		} else {
			alertTmpl = t
			return
		}

		alertTmpl = template.Must(template.New("price_alert").Funcs(fns).Parse(defaultAlertTemplate))
	})
	return alertTmpl
}

func renderEmailBody(recipient string, product ProductAlert) (string, error) {
	data := struct {
		Recipient string
		Product   ProductAlert
	}{Recipient: recipient, Product: product}

	var buf bytes.Buffer
	if err := emailTemplate().Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func EmailAll(db *gorm.DB) error {
	slog.Info("publishing email alerts")
	var products []models.Product
	if err := db.Preload("UsersTracking").
		Where("direction IN ?", []models.PriceState{models.BelowRecentAvg, models.Decreased}).
		Find(&products).Error; err != nil {
		return err
	}

	for _, product := range products {
		for _, user := range product.UsersTracking {
			alert := EmailAlert{
				Recipient: user.Email,
				Product: ProductAlert{
					Name:      product.Name,
					Url:       product.Url,
					ImgLink:   product.ImgLink,
					Price:     product.Price,
					Direction: string(product.Direction),
				},
			}

			body, err := json.Marshal(alert)
			if err != nil {
				return err
			}

			if err := queue.Publish(queue.EmailAlertsQueue, body); err != nil {
				return err
			}
		}
	}

	return nil
}

func ConsumeEmailAlerts() error {
	ch, msgs, err := queue.Consume(queue.EmailAlertsQueue)
	if err != nil {
		return err
	}
	defer ch.Close()

	for d := range msgs {
		var alert EmailAlert
		if err := json.Unmarshal(d.Body, &alert); err != nil {
			slog.Error("email consumer: dropping unreadable message", "err", err)
			_ = d.Nack(false, false)
			continue
		}

		body, err := renderEmailBody(alert.Recipient, alert.Product)
		if err != nil {
			slog.Error("email consumer: cannot render template", "err", err)
			_ = d.Nack(false, false)
			continue
		}

		var sendErr error
		for attempt := 1; attempt <= sendMailMaxAttempts; attempt++ {
			if sendErr = sendMail(alert.Recipient, body); sendErr == nil {
				break
			}
			slog.Warn("email consumer: send attempt failed", "attempt", attempt, "max", sendMailMaxAttempts, "recipient", alert.Recipient, "err", sendErr)
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		if sendErr != nil {
			slog.Error("email consumer: giving up, dead-lettering", "recipient", alert.Recipient)
			_ = d.Nack(false, false)
			continue
		}

		_ = d.Ack(false)
	}

	return nil
}
