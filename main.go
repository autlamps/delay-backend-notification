package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/autlamps/delay-backend-notification/notification"
)

var mqurl string
var dburl string
var firebaseKey string

func init() {
	flag.StringVar(&dburl, "DATABASE_URL", "", "database url")
	flag.StringVar(&mqurl, "RABBITMQ_URL", "", "message queue url")
	flag.StringVar(&firebaseKey, "FIREBASE_KEY", "", "firebase url")
	flag.Parse()

	if firebaseKey == "" {
		firebaseKey = os.Getenv("FIREBASE_KEY")
	}

	if dburl == "" {
		dburl = os.Getenv("DATABASE_URL")
	}

	if mqurl == "" {
		mqurl = os.Getenv("RABBITMQ_URL")
	}

	if mqurl == "" || dburl == "" {
		log.Fatal("DB url and/or message broker url not set by either env variable or flag.")
	}
}

func main() {
	c := notification.Conf{DBURL: dburl, MQURL: mqurl, FirebaseToken: firebaseKey}
	e, err := notification.EnvFromConf(c)

	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	// Exit channel used to signal env.Start that we want to stop executing
	ec := make(chan bool)

	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)

	// This func on receiving the syscall signal from the channel sends true down our exit channel.
	go func() {
		<-sc
		fmt.Println("Exit signal recieved")
		ec <- true
	}()

	e.Start(ec)
}
