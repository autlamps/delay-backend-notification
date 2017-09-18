// Notify implements recieving of late services from the message broker
package notify

import (
	"fmt"

	"github.com/streadway/amqp"
)

// ContentType is our type for specifying the type of message we are sending
type ContentType string

// JSON is one of our content types
const JSON = "text/json"

// Notifier interface specifies the methods needed for a concrete notifier service
type Notifier interface {
	Receive() (<-chan amqp.Delivery, error)
	Close()
}

// Service is our concrete implementation of Notifier using rabbit mq
type Service struct {
	con *amqp.Connection
	ch  *amqp.Channel
	q   *amqp.Queue
}

// InitService returns a new service
func InitService(url string) (*Service, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return &Service{}, fmt.Errorf("notify/InitService: failed to connect to rabbitmq: %v", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		return &Service{}, fmt.Errorf("notify/InitService: failed to create channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		"notify",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return &Service{}, fmt.Errorf("notify/InitService: failed to create queue: %v", err)
	}

	return &Service{conn, ch, &q}, nil
}

// Receive reads in a message from the notification service
func (s *Service) Receive() (<-chan amqp.Delivery, error) {
	return s.ch.Consume(
		s.q.Name, // queue
		"",       // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
}

// Close closes a services channel then connection
func (s *Service) Close() {
	s.ch.Close()
	s.con.Close()
}
