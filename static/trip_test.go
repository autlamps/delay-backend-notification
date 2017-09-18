package static

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestTripService_GetTripByGTFSID(t *testing.T) {
	db, err := sql.Open("postgres", dburl)
	defer db.Close()

	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	ts := TripServiceInit(db)

	tests := []struct {
		id       string
		expected Trip
	}{
		{"13223097957-20170829094406_v57.13",
			Trip{ID: "236a9757-6aeb-49e9-ad7c-c56d660a10ed",
				RouteID:   "3dd07748-5586-4edf-a714-f7ba4999af7e",
				ServiceID: "b49995c1-22f4-40d6-8996-245d0f7b4f89",
				GTFSID:    "13223097957-20170829094406_v57.13",
				Headsign:  "City Centre"}},
	}

	for _, test := range tests {

		trip, err := ts.GetTripByGTFSID(test.id)

		if err != nil {
			t.Errorf("Failed to retrieve trip. %v", err)
		}

		if !trip.IsEqual(test.expected) {
			t.Errorf("Failed to retrieve correct trip. Expected id %v, got %v", test.expected.ID, trip.ID)
		}

	}
}

func TestTrip_IsEqual(t *testing.T) {

	t1 := Trip{ID: "f5efea1b-60e7-4239-a6e7-47d14a858399",
		RouteID:   "0f6775a1-ce65-4e6e-82ba-f94e599a1a57",
		ServiceID: "069423ce-9866-4796-b6ad-8eeb8ad87f2a",
		GTFSID:    "1080081195-20170807091914_v56.25",
		Headsign:  "City Centre"}

	t2 := Trip{ID: "f5efea1b-60e7-4239-a6e7-47d14a858399",
		RouteID:   "0f6775a1-ce65-4e6e-82ba-f94e599a1a57",
		ServiceID: "069423ce-9866-4796-b6ad-8eeb8ad87f2a",
		GTFSID:    "1080081195-20170807091914_v56.25",
		Headsign:  "City Centre"}

	if !t1.IsEqual(t2) {
		t.Errorf("Identical trips not equal :(")
	}
}
