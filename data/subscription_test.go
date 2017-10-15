package data

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

func subSetup() (NotifyInfo, User, *sql.DB, error) {
	db, err := sql.Open("postgres", dburl)

	if err != nil {
		return NotifyInfo{}, User{}, nil, fmt.Errorf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return NotifyInfo{}, User{}, nil, fmt.Errorf("Failed to ping db: %v", err)
	}

	us := InitUserService(db)

	u, err := us.NewUser(NewUser{
		"Bobby Tables",
		"bobby.tables@example.com",
		"correcthorsebatterystaple",
	})

	if err != nil {
		return NotifyInfo{}, User{}, nil, fmt.Errorf("Failed to create new user: %v", err)
	}

	ns := InitNotifyInfoService(db)

	n, err := ns.New(u.ID.String(), "p", "iPhone X", "1234456")

	if err != nil {
		return NotifyInfo{}, User{}, nil, fmt.Errorf("Failed to create new notify method: %v", err)
	}

	return n, u, db, nil
}

func TestSubscriptionService_GetSubsByStopTimeID(t *testing.T) {
	// TODO: need to modify this so we use different users and notifications ids for our test subs
	// It shows that our sql is working though... which is nice

	n, u, db, err := subSetup()

	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	ss := InitSubscriptionService(db)

	subs := []Subscription{}

	for i := 0; i < 5; i++ {
		ns := NewSubscription{
			TripID:          "df688c57-987c-4705-9e22-936342eb6e3f",
			StopTimeID:      "5cce0bca-d489-43d7-b3cb-48e0df054c8a",
			Days:            []Day{"Mon", "Tue", "Wed"},
			NotificationIDs: []string{n.ID},
			UserID:          u.ID.String(),
		}

		s, err := ss.new(ns)

		if err != nil {
			t.Fatalf("Failed to create new sub: %v", err)
		}

		subs = append(subs, s)
	}

	dbsubs, err := ss.GetSubsByStopTimeID("5cce0bca-d489-43d7-b3cb-48e0df054c8a")

	if err != nil {
		t.Fatalf("Failed to get subs from db: %v", err)
	}

	if !reflect.DeepEqual(subs, dbsubs) {
		t.Fatalf("Subs returned from GetSubsByStopTimeID not the same as expected")
	}

	// CleanUp
	for _, sub := range subs {
		_, err := db.Exec("DELETE FROM sub_notification WHERE sub_id = $1", sub.ID)

		if err != nil {
			t.Fatalf("Failed to delete created sub notification: %v\n", err)
		}

		_, err = db.Exec("DELETE FROM subscription WHERE sub_id = $1", sub.ID)

		if err != nil {
			t.Fatalf("Failed to delete created sub: %v\n", err)
		}
	}

	err = subCleanUp(Subscription{}, n, u, db)

	if err != nil {
		t.Fatalf("Failed to clean up: %v", err)
	}
}

func subCleanUp(s Subscription, ni NotifyInfo, u User, db *sql.DB) error {
	// We can send a blank Subscription struct if we have already cleaned up
	if s.ID != "" {
		_, err := db.Exec("DELETE FROM sub_notification WHERE sub_id = $1", s.ID)

		if err != nil {
			return fmt.Errorf("Failed to delete created subs: %v\n", err)
		}

		_, err = db.Exec("DELETE FROM subscription WHERE sub_id = $1", s.ID)

		if err != nil {
			return fmt.Errorf("Failed to delete created sub: %v\n", err)
		}
	}

	_, err := db.Exec("DELETE FROM notification WHERE notification_id = $1", ni.ID)

	if err != nil {
		return fmt.Errorf("Failed to delete created notification: %v\n", err)
	}
	_, err = db.Exec("DELETE FROM tokens WHERE user_id = $1", u.ID)

	if err != nil {
		return fmt.Errorf("Failed to delete created user token: %v\n", err)
	}

	_, err = db.Exec("DELETE FROM users WHERE user_id = $1", u.ID)

	if err != nil {
		return fmt.Errorf("Failed to delete created user: %v\n", err)
	}

	db.Close()

	return nil
}
