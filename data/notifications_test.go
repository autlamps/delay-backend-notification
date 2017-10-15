package data

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func notifySetup() (User, *sql.DB, error) {
	db, err := sql.Open("postgres", dburl)

	if err != nil {
		return User{}, nil, fmt.Errorf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return User{}, nil, fmt.Errorf("Failed to ping db: %v", err)
	}

	us := InitUserService(db)

	u, err := us.NewUser(NewUser{
		"Bobby Tables",
		"bobby.tables@example.com",
		"correcthorsebatterystaple",
	})

	if err != nil {
		return User{}, nil, fmt.Errorf("Failed to create new user: %v", err)
	}

	return u, db, nil
}

func notifyCleanup(u User, ni NotifyInfo, db *sql.DB) error {
	// If we've already cleaned up the ni dont notifyCleanup
	if ni.ID != "" {
		_, err := db.Exec("DELETE FROM notification WHERE notification_id = $1", ni.ID)

		if err != nil {
			return fmt.Errorf("Failed to delete created notification: %v\n", err)
		}
	}

	_, err := db.Exec("DELETE FROM tokens WHERE user_id = $1", u.ID)

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

func TestNotifyInfoService_New(t *testing.T) {
	u, db, err := notifySetup()

	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	nis := InitNotifyInfoService(db)

	ni, err := nis.New(u.ID.String(), PUSH, "iPhone X", "123467876543214")

	if err != nil {
		t.Fatalf("Failed to insert new notification info into db: %v", err)
	}

	dbni := NotifyInfo{}

	row := db.QueryRow("SELECT notification_id, user_id, type, name, value, date_created FROM notification WHERE notification_id = $1", ni.ID)

	err = row.Scan(&dbni.ID, &dbni.UserID, &dbni.Type, &dbni.Name, &dbni.Value, &dbni.Created)

	if err != nil {
		t.Fatalf("Failed to retrieve notify info from db: %v", err)
	}

	dbni.Created = dbni.Created.In(time.Local)

	if !reflect.DeepEqual(ni, dbni) {
		t.Fatalf("Retrieved not the same as saved. Expected %v, got %v", ni, dbni)
	}

	// Cleanup
	err = notifyCleanup(u, ni, db)

	if err != nil {
		t.Fatalf("Failed to notifyCleanup: %v", err)
	}
}

func TestNotifyInfoService_Delete(t *testing.T) {
	u, db, err := notifySetup()

	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	nis := InitNotifyInfoService(db)

	ni, err := nis.New(u.ID.String(), PUSH, "iPhone X", "123467876543214")

	if err != nil {
		t.Fatalf("Failed to insert new notification info into db: %v", err)
	}

	err = nis.Delete(ni.ID)

	if err != nil {
		t.Fatalf("Err returned from delete func: %v", err)
	}

	dbni := NotifyInfo{}

	row := db.QueryRow("SELECT notification_id, user_id, type, name, value, date_created FROM notification WHERE notification_id = $1", ni.ID)

	err = row.Scan(&dbni.ID, &dbni.UserID, &dbni.Type, &dbni.Name, &dbni.Value, &dbni.Created)

	if err != sql.ErrNoRows {
		t.Fatal("Notification not deleted")
	}

	err = notifyCleanup(u, ni, db)

	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}
}

func TestNotifyInfoService_GetAll(t *testing.T) {
	u, db, err := notifySetup()

	if err != nil {
		t.Fatalf("Failed to notifySetup test: %v", err)
	}

	nis := InitNotifyInfoService(db)
	notifications := []NotifyInfo{}

	for i := 0; i < 5; i++ {
		ni, err := nis.New(u.ID.String(), PUSH, fmt.Sprintf("iPhone %v", i), fmt.Sprintf("%v", i))

		if err != nil {
			t.Fatalf("Failed to create notificaton: %v", err)
		}

		notifications = append(notifications, ni)
	}

	dbnotifications, err := nis.GetAll(u.ID.String())

	if err != nil {
		t.Fatalf("Failed to get all notifications: %v", err)
	}

	if len(notifications) != len(notifications) {
		t.Fatalf("Number of notifications not the same: Expected %v, got %v", len(notifications), len(dbnotifications))
	}

	//Clean up
	for _, ni := range notifications {
		err := nis.Delete(ni.ID)

		if err != nil {
			t.Fatalf("Failed to delete: %v", ni)
		}
	}

	err = notifyCleanup(u, NotifyInfo{}, db)

	if err != nil {
		t.Fatalf("Failed to run notifyCleanup: %v", err)
	}
}
