package notify

import (
	"flag"
	"os"
	"testing"
	"time"
)

// TODO: this testing file really needs fleshing out. Need to look into how to test it

var mqurl string

func init() {
	flag.StringVar(&mqurl, "MQ_URL", "", "message queue url for testing")
	flag.Parse()

	if mqurl == "" {
		mqurl = os.Getenv("MQ_URL")
	}
}

func TestService_Send(t *testing.T) {
	// TODO: review this test to make sure it's actually useful
	s, err := InitService(mqurl)
	defer s.Close()

	if err != nil {
		t.Errorf("Failed to create init service %v", err)
	}

	err = s.Send([]byte("hello"))

	if err != nil {
		t.Errorf("Failed to send message: %v", err)
	}

	msgs, err := s.ch.Consume(
		s.q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	select {
	case msg := <-msgs:
		if string(msg.Body) != "hello" {
			t.Errorf("Message recieved not the same as sent")
		}
	case <-time.After(time.Second * 10):
		t.Errorf("Failed to retrieve message after 10 seconds")
	}
}
