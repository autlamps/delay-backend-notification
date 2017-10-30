package static

import (
	"flag"
	"os"
)

var dburl string

// Grab postgres url
func init() {
	flag.StringVar(&dburl, "DB_URL", "", "database url for testing")
	flag.Parse()

	if dburl == "" {
		dburl = os.Getenv("DB_URL")
	}
}
