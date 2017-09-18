package notification

import (
	"database/sql"
	"fmt"

	"encoding/json"
	"log"

	"github.com/autlamps/delay-backend-notification/input"
	"github.com/autlamps/delay-backend-notification/notify"
	"github.com/autlamps/delay-backend-notification/static"
)

// Conf stores our initial string connection values before being turned into services
type Conf struct {
	DBURL string
	MQURL string
}

// Env stores our services to be used
type Env struct {
	Notification notify.Notifier
	Trips        static.TripStore
	StopTimes    static.StopTimeStore
	Routes       static.RouteStore
}

// EnvFromConf
func EnvFromConf(c Conf) (Env, error) {
	db, err := sql.Open("postgres", c.DBURL)

	if err != nil {
		return Env{}, fmt.Errorf("Failed to open db connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return Env{}, fmt.Errorf("Failed to ping db connection: %v", err)
	}

	n, err := notify.InitService(c.MQURL)

	if err != nil {
		return Env{}, err
	}

	return Env{
		Notification: n,
		Trips:        static.TripServiceInit(db),
		StopTimes:    static.StopTimeServiceInit(db),
		Routes:       static.RouteServiceInit(db),
	}, nil
}

// Start contains our main loop
func (e *Env) Start() {
	nc, err := e.Notification.Receive()

	if err != nil {
		log.Fatal(err)
	}

	for nm := range nc {
		b := nm.Body
		n := input.Notification{}

		err := json.Unmarshal(b, &n)

		if err != nil {
			log.Fatal(err)
		}

		// Do work

		nm.Ack(true)
	}
}
