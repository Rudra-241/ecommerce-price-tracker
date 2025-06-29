package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EmailAlertsQueue       = "email.alerts"
	EmailAlertsFailedQueue = "email.alerts.failed"
)

var conn *amqp.Connection

func Init(url string) error {
	c, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("connecting to rabbitmq: %w", err)
	}
	conn = c

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("opening channel: %w", err)
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(EmailAlertsFailedQueue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declaring %s: %w", EmailAlertsFailedQueue, err)
	}

	_, err = ch.QueueDeclare(EmailAlertsQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": EmailAlertsFailedQueue,
	})
	if err != nil {
		return fmt.Errorf("declaring %s: %w", EmailAlertsQueue, err)
	}

	return nil
}

func Publish(queueName string, body []byte) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("opening channel: %w", err)
	}
	defer ch.Close()

	return ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func Consume(queueName string) (*amqp.Channel, <-chan amqp.Delivery, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("opening channel: %w", err)
	}

	if err := ch.Qos(5, 0, false); err != nil {
		ch.Close()
		return nil, nil, fmt.Errorf("setting qos: %w", err)
	}

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		ch.Close()
		return nil, nil, fmt.Errorf("consuming %s: %w", queueName, err)
	}

	return ch, msgs, nil
}
