package data

import (
	"database/sql"
	"time"
)

// Subscription contains our sub info
type Subscription struct {
	ID         string
	TripID     string
	StopTimeID string
	UserID     string
	Mon        bool
	Tue        bool
	Wed        bool
	Thu        bool
	Fri        bool
	Sat        bool
	Sun        bool
	Created    time.Time
}

// Subscriptions is a slice of subscription
type Subscriptions []Subscription

// SubStore is our interface defining methods for a concrete implementation
type SubStore interface {
	GetSubsByStopTimeID(string) (Subscriptions, error)
}

// SubService is our psql implementation of SubStore
type SubService struct {
	db *sql.DB
}

// GetSubsByStopTimeID gets all non archived subscriptions for a given stoptime
func (ss *SubService) GetSubsByStopTimeID(stid string) (Subscriptions, error) {
	rows, err := ss.db.Query("SELECT sub_id, trip_id, stoptime_id, user_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday, date_created WHERE archived = FALSE AND stoptime_id = '$1'",
		stid,
	)

	if err != nil {
		return Subscriptions{}, err
	}

	subs := Subscriptions{}

	for rows.Next() {
		sub := Subscription{}

		err := rows.Scan(&sub.ID, &sub.TripID, &sub.UserID, &sub.Mon, &sub.Tue, &sub.Thu, &sub.Wed, &sub.Thu, &sub.Fri, &sub.Sat, &sub.Sun, &sub.Created)

		if err != nil {
			continue // TODO: Yeah lets maybe do something else here?
		}

		subs = append(subs, sub)
	}

	return subs, nil
}
