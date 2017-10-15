package data

import (
	"flag"
	"os"

	_ "github.com/lib/pq"
)

var dburl string

func init() {
	flag.StringVar(&dburl, "DATABASE_URL", "", "database url")
	flag.Parse()

	if dburl == "" {
		dburl = os.Getenv("DATABASE_URL")
	}
}
