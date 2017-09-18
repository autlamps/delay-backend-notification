package main

import (
	"flag"
	"log"
	"os"
)

var mqurl string
var dburl string

func init() {
	flag.StringVar(&dburl, "DB_URL", "", "database url")
	flag.StringVar(&mqurl, "MQ_URL", "", "message queue url")

	if dburl == "" {
		dburl = os.Getenv("DB_URL")
	}

	if mqurl == "" {
		mqurl = os.Getenv("MQ_URL")
	}

	if mqurl == "" || dburl == "" {
		log.Fatal("DB url and/or message broker url not set by either env variable or flag.")
	}
}

func main() {

}
