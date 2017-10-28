package data

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"
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
			TripID:          "bc82fb93-8b60-40e8-83ed-ce520d6ed5a2",
			StopTimeID:      "a47cf7a1-c40f-4cee-84f5-b025e79c0935",
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

	dbsubs, err := ss.GetSubsByStopTimeID("a47cf7a1-c40f-4cee-84f5-b025e79c0935")

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

func TestSubscriptionService_GetSubsByTripID(t *testing.T) {
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
			TripID:          "bc82fb93-8b60-40e8-83ed-ce520d6ed5a2",
			StopTimeID:      "a47cf7a1-c40f-4cee-84f5-b025e79c0935",
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

	dbsubs, err := ss.GetSubsByTripID("bc82fb93-8b60-40e8-83ed-ce520d6ed5a2")

	if err != nil {
		t.Fatalf("Failed to get subs from db: %v", err)
	}

	if !reflect.DeepEqual(subs, dbsubs) {
		t.Fatalf("Subs returned from GetSubsByTripID not the same as expected")
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

func TestSubscriptionService_RecentlyNotified(t *testing.T) {
	n, u, db, err := subSetup()

	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	ss := InitSubscriptionService(db)

	ns1 := NewSubscription{
		TripID:          "bc82fb93-8b60-40e8-83ed-ce520d6ed5a2",
		StopTimeID:      "a47cf7a1-c40f-4cee-84f5-b025e79c0935",
		Days:            []Day{"Mon", "Tue", "Wed"},
		NotificationIDs: []string{n.ID},
		UserID:          u.ID.String(),
	}

	s1, err := ss.new(ns1)

	if err != nil {
		t.Fatalf("Failed to create new sub: %v", err)
	}

	ns2 := NewSubscription{
		TripID:          "0b8f67e8-78dc-4d77-bb57-546513e71430",
		StopTimeID:      "d20ff7e9-34d7-4d07-b86d-6a4ceb89daa3",
		Days:            []Day{"Mon", "Tue", "Wed"},
		NotificationIDs: []string{n.ID},
		UserID:          u.ID.String(),
	}

	s2, err := ss.new(ns2)

	_, err = ss.db.Exec("INSERT INTO notification_event (notification_event_id, sub_id, date_created) VALUES ($1, $2, $3)",
		"226cda79-957a-4e27-b41e-93f989b0ccf1",
		s1.ID,
		time.Now().Round(time.Second),
	)

	_, err = ss.db.Exec("INSERT INTO notification_event (notification_event_id, sub_id, date_created) VALUES ($1, $2, $3)",
		"6519996b-5e82-4b3e-aeec-a7a186c6b61b",
		s2.ID,
		time.Now().Round(time.Second).Add(-(time.Minute * time.Duration(120))),
	)

	tests := []struct {
		ID       string
		Expected bool
	}{
		{s1.ID, true},
		{s2.ID, false},
	}

	for _, test := range tests {
		used, err := ss.RecentlyNotified(test.ID)

		if err != nil {
			t.Fatalf("Failed to retrieve recently notified: %v", err)
		}

		if used != test.Expected {
			t.Fatalf("Recently notified doesn't match expected: Expected %v, got %v", used, test.Expected)
		}
	}

	//Clean up
	_, err = db.Exec("DELETE FROM sub_notification WHERE sub_id = $1", s1.ID)

	if err != nil {
		t.Fatalf("Failed to delete created sub notification: %v\n", err)
	}

	_, err = db.Exec("DELETE FROM notification_event WHERE sub_id = $1", s1.ID)

	if err != nil {
		t.Fatalf("Failed to delete created sub notification event: %v\n", err)
	}
	_, err = db.Exec("DELETE FROM subscription WHERE sub_id = $1", s1.ID)

	if err != nil {
		t.Fatalf("Failed to delete created sub: %v\n", err)
	}
	_, err = db.Exec("DELETE FROM sub_notification WHERE sub_id = $1", s2.ID)

	if err != nil {
		t.Fatalf("Failed to delete created sub notification: %v\n", err)
	}

	_, err = db.Exec("DELETE FROM notification_event WHERE sub_id = $1", s2.ID)

	if err != nil {
		t.Fatalf("Failed to delete created sub notification event: %v\n", err)
	}

	_, err = db.Exec("DELETE FROM subscription WHERE sub_id = $1", s2.ID)

	if err != nil {
		t.Fatalf("Failed to delete created sub: %v\n", err)
	}

	err = subCleanUp(Subscription{}, n, u, db)

	if err != nil {
		t.Fatalf("Failed to clean up: %v", err)
	}
}

func TestSubscriptionService_Notified(t *testing.T) {
	n, u, db, err := subSetup()

	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	ss := InitSubscriptionService(db)

	ns := NewSubscription{
		TripID:          "bc82fb93-8b60-40e8-83ed-ce520d6ed5a2",
		StopTimeID:      "a47cf7a1-c40f-4cee-84f5-b025e79c0935",
		Days:            []Day{"Mon", "Tue", "Wed"},
		NotificationIDs: []string{n.ID},
		UserID:          u.ID.String(),
	}

	s, err := ss.new(ns)

	if err != nil {
		t.Fatalf("Failed to create new sub: %v", err)
	}

	err = ss.Notified(s)

	if err != nil {
		t.Fatalf("Failed to mark subscription as notified: %v", err)
	}

	row := ss.db.QueryRow("SELECT COUNT(*) from notification_event WHERE sub_id = $1", s.ID)

	var count int

	err = row.Scan(&count)

	if err != nil {
		t.Fatalf("Failed to get notification event count from db: %v", err)
	}

	if count != 1 {
		t.Fatalf("Number of notification events for subscriptions unexpected. Expected 1 got %v", count)
	}

	//Clean up

	_, err = db.Exec("DELETE FROM notification_event WHERE sub_id = $1", s.ID)

	if err != nil {
		t.Fatalf("Failed to delete created sub notification event: %v\n", err)
	}

	err = subCleanUp(s, n, u, db)

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
