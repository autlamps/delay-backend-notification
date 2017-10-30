package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Subscription contains our sub info
type Subscription struct {
	ID              string
	TripID          string
	StopTimeID      string
	UserID          string
	Archived        bool
	Created         time.Time
	Monday          bool
	Tuesday         bool
	Wednesday       bool
	Thursday        bool
	Friday          bool
	Saturday        bool
	Sunday          bool
	NotificationIDs []string
}

// SubStore is our interface defining methods for a concrete implementation
type SubscriptionStore interface {
	GetSubsByStopTimeID(string) ([]Subscription, error)
	GetSubsByTripID(id string) ([]Subscription, error)
	Notified(s Subscription) error
	RecentlyNotified(id string) (bool, error)
}

// SubService is our psql implementation of SubStore
type SubscriptionService struct {
	db *sql.DB
}

func InitSubscriptionService(db *sql.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

// GetSubsByStopTimeID gets all non archived subscriptions for a given stoptime
// GetAll returns all subscriptions belonging to a single user
func (ss *SubscriptionService) GetSubsByStopTimeID(stid string) ([]Subscription, error) {
	subs := []Subscription{}

	rows, err := ss.db.Query("SELECT sub_id, trip_id, stoptime_id, user_id, archived, date_created, monday, tuesday, wednesday, thursday, friday, saturday, sunday FROM subscription WHERE stoptime_id = $1", stid)

	if err != nil {
		return []Subscription{}, fmt.Errorf("subscription - GetAll: Failed to get subs from db: %v", err)
	}

	for rows.Next() {
		var s Subscription

		err := rows.Scan(&s.ID, &s.TripID, &s.StopTimeID, &s.UserID, &s.Archived, &s.Created, &s.Monday, &s.Tuesday, &s.Wednesday, &s.Thursday, &s.Friday, &s.Saturday, &s.Sunday)

		if err != nil {
			return []Subscription{}, fmt.Errorf("subscription - GetAll: Failed to scan for individual subscription: %v", err)
		}

		s.Created = s.Created.Local()

		notifyRows, err := ss.db.Query("SELECT notification_id from sub_notification WHERE sub_id = $1", s.ID)

		if err != nil {
			return []Subscription{}, fmt.Errorf("subscription - GetAll: Failed get notification ids: %v", err)
		}

		for notifyRows.Next() {
			var id string

			err := notifyRows.Scan(&id)

			if err != nil {
				return []Subscription{}, fmt.Errorf("subscription - GetAll: Failed to read individual notification id: %v", err)
			}

			s.NotificationIDs = append(s.NotificationIDs, id)
		}

		subs = append(subs, s)
	}

	return subs, nil
}

// GetSubsByTripID gets all subscriptions associated with the given trip id
func (ss *SubscriptionService) GetSubsByTripID(id string) ([]Subscription, error) {
	subs := []Subscription{}

	rows, err := ss.db.Query("SELECT sub_id, trip_id, stoptime_id, user_id, archived, date_created, monday, tuesday, wednesday, thursday, friday, saturday, sunday FROM subscription WHERE trip_id = $1", id)

	if err != nil {
		return []Subscription{}, fmt.Errorf("subscription - GetSubsByTripID: Failed to get subs from db: %v", err)
	}

	for rows.Next() {
		var s Subscription

		err := rows.Scan(&s.ID, &s.TripID, &s.StopTimeID, &s.UserID, &s.Archived, &s.Created, &s.Monday, &s.Tuesday, &s.Wednesday, &s.Thursday, &s.Friday, &s.Saturday, &s.Sunday)

		if err != nil {
			return []Subscription{}, fmt.Errorf("subscription - GetSubsByTripID: Failed to scan for individual subscription: %v", err)
		}

		s.Created = s.Created.Local()

		notifyRows, err := ss.db.Query("SELECT notification_id from sub_notification WHERE sub_id = $1", s.ID)

		if err != nil {
			return []Subscription{}, fmt.Errorf("subscription - GetSubsByTripID: Failed get notification ids: %v", err)
		}

		for notifyRows.Next() {
			var id string

			err := notifyRows.Scan(&id)

			if err != nil {
				return []Subscription{}, fmt.Errorf("subscription - GetSubsByTripID: Failed to read individual notification id: %v", err)
			}

			s.NotificationIDs = append(s.NotificationIDs, id)
		}

		subs = append(subs, s)
	}

	return subs, nil
}

var ErrNoNotificationMethods = errors.New("users - No notification methods specificed.")

// Day is one of our three letter day codes
type Day string

// Defines our three letter day codes
const (
	MONDAY    Day = "Mon"
	TUESDAY   Day = "Tue"
	WEDNESDAY Day = "Wed"
	THURSDAY  Day = "Thur"
	FRIDAY    Day = "Fri"
	SATURDAY  Day = "Sat"
	SUNDAY    Day = "Sun"
)

// NewSubscription is received from called and transformed into a db backed Subscription
type NewSubscription struct {
	TripID          string
	StopTimeID      string
	Days            []Day
	NotificationIDs []string
	UserID          string
}

// New creates a new database backed Subscription, or returns an error - Used for testing in this service
func (ss *SubscriptionService) new(ns NewSubscription) (Subscription, error) {
	// If no notification methods are specified then we send an error back
	if len(ns.NotificationIDs) == 0 {
		return Subscription{}, ErrNoNotificationMethods
	}

	id, err := uuid.NewRandom()

	if err != nil {
		return Subscription{}, fmt.Errorf("subscriptions - New: failed to generate uuid: %v", err)
	}

	s := Subscription{
		ID:              id.String(),
		TripID:          ns.TripID,
		StopTimeID:      ns.StopTimeID,
		Archived:        false,
		Created:         time.Now().Round(time.Second),
		NotificationIDs: ns.NotificationIDs,
		UserID:          ns.UserID,
	}

	// Setup subscribed days, days not present will remain false as per golang default
	for _, d := range ns.Days {
		switch d {
		case MONDAY:
			s.Monday = true
		case TUESDAY:
			s.Tuesday = true
		case WEDNESDAY:
			s.Wednesday = true
		case THURSDAY:
			s.Thursday = true
		case FRIDAY:
			s.Friday = true
		case SATURDAY:
			s.Saturday = true
		case SUNDAY:
			s.Sunday = true
		}
	}

	_, err = ss.db.Exec("INSERT INTO subscription (sub_id, trip_id, stoptime_id, user_id, archived, date_created, monday, tuesday, wednesday, thursday, friday, saturday, sunday) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		s.ID,
		s.TripID,
		s.StopTimeID,
		s.UserID,
		s.Archived,
		s.Created,
		s.Monday,
		s.Tuesday,
		s.Wednesday,
		s.Thursday,
		s.Friday,
		s.Saturday,
		s.Sunday,
	)

	if err != nil {
		return Subscription{}, fmt.Errorf("users - Subscription: Failed to insert subscription into db: %v", err)
	}

	for _, sn := range s.NotificationIDs {
		_, err = ss.db.Exec("INSERT INTO sub_notification (sub_id, notification_id) VALUES ($1, $2)",
			s.ID,
			sn,
		)

		if err != nil {
			return Subscription{}, fmt.Errorf("users - New: failed to link notification methods and subscription: %v", err)
		}
	}

	return s, nil
}

// RecentlyNotified returns true if a user was recently notified by this subscription
func (ss *SubscriptionService) RecentlyNotified(id string) (bool, error) {
	row := ss.db.QueryRow("SELECT COUNT(*) from notification_event WHERE sub_id = $1 AND date_created > NOW() - INTERVAL '15 minutes'", id)

	var count int

	err := row.Scan(&count)

	if err != nil {
		return false, fmt.Errorf("subscriptions - RecentlyNotified: Failed to get if subscription notified recently: %v", err)
	}

	if count < 1 {
		return false, nil
	}

	return true, nil
}

// Notified takes in a subscription and adds a record to indicate it has been notified
func (ss *SubscriptionService) Notified(s Subscription) error {
	id, err := uuid.NewRandom()

	if err != nil {
		return fmt.Errorf("subscriptions - Notified: failed to generate id: %v", err)
	}

	_, err = ss.db.Exec("INSERT INTO notification_event (notification_event_id, sub_id, date_created) VALUES ($1, $2, $3)",
		id,
		s.ID,
		time.Now().Round(time.Second),
	)

	if err != nil {
		return fmt.Errorf("subscriptions - Notified: failed to insert event: %v", err)
	}

	return nil
}

// SubscribedForToday returns whether or not the user wishes to receive notifications for this subscription today
func (s Subscription) SubscribedForToday() bool {
	day := Day(time.Now().Format("Mon"))

	switch day {
	case MONDAY:
		return s.Monday == true
	case TUESDAY:
		return s.Tuesday == true
	case WEDNESDAY:
		return s.Wednesday == true
	case THURSDAY:
		return s.Thursday == true
	case FRIDAY:
		return s.Friday == true
	case SATURDAY:
		return s.Saturday == true
	case SUNDAY:
		return s.Sunday == true
	}

	return false
}