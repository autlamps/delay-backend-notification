package static

import (
	"testing"

	"time"

	"database/sql"

	_ "github.com/lib/pq"
)

func TestStopTimeService_GetStopsByTripID(t *testing.T) {
	db, err := sql.Open("postgres", dburl)
	defer db.Close()

	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	sts := StopTimeServiceInit(db)

	// Note: find short trips to test by running: SELECT count(trip_id), trip_id FROM stop_times GROUP BY trip_Id ORDER BY count(trip_id) ASC;
	tests := []struct {
		in       string
		expected []StopTime
	}{
		{"d0d6dece-e75d-4e1a-8bbc-7eefd6a07758", StopTimeArray{
			{"847af420-5527-49b1-85f0-9a18fb143d87",
				"d0d6dece-e75d-4e1a-8bbc-7eefd6a07758",
				time.Date(0, 1, 1, 14, 45, 0, 0, time.UTC),
				time.Date(0, 1, 1, 14, 45, 0, 0, time.UTC),
				1,
				Stop{
					"67cca212-476a-4e60-bf25-8bc5318c3d23",
					"9670",
					"Devonport Ferry Terminal",
					-36.83317,
					174.7954,
				}},
			{"532789a3-1abb-4d48-886e-24212998212a",
				"d0d6dece-e75d-4e1a-8bbc-7eefd6a07758",
				time.Date(0, 1, 1, 14, 55, 0, 0, time.UTC),
				time.Date(0, 1, 1, 14, 55, 0, 0, time.UTC),
				2,
				Stop{
					"c8e15bd6-0066-4840-8c37-5f8b1bcb0f6c",
					"9600",
					"Downtown Ferry Terminal Pier 1",
					-36.84243,
					174.76708,
				}},
		}},
	}

	for _, test := range tests {
		var sta StopTimeArray

		sta, err := sts.GetStopTimesByTripID(test.in)

		if err != nil {
			t.Errorf("Failed to retrieve stoptimes %v", err)
		}

		if !sta.IsEqual(test.expected) {
			t.Errorf("Stoptimes not equal. \nE %v, \nG %v", test.expected, sta)
		}
	}
}

func TestStop_IsEqual(t *testing.T) {
	s1 := Stop{
		"c8e15bd6-0066-4840-8c37-5f8b1bcb0f6c",
		"9600",
		"Downtown Ferry Terminal Pier 1",
		-36.84243,
		174.76708,
	}

	s2 := Stop{
		"c8e15bd6-0066-4840-8c37-5f8b1bcb0f6c",
		"9600",
		"Downtown Ferry Terminal Pier 1",
		-36.84243,
		174.76708,
	}

	if !s1.IsEqual(s2) {
		t.Errorf("Identical stops not equal :(")
	}
}

func TestStopTime_IsEqual(t *testing.T) {
	st1 := StopTime{"847af420-5527-49b1-85f0-9a18fb143d87",
		"d0d6dece-e75d-4e1a-8bbc-7eefd6a07758",
		time.Date(0, 0, 0, 06, 45, 0, 0, time.UTC),
		time.Date(0, 0, 0, 06, 55, 0, 0, time.UTC),
		1,
		Stop{
			"c8d731d1-2e1c-4ca2-ac26-fead1695320d",
			"9600",
			"Downtown Ferry Terminal Pier 1",
			-36.84243,
			174.76708,
		}}

	st2 := StopTime{"847af420-5527-49b1-85f0-9a18fb143d87",
		"d0d6dece-e75d-4e1a-8bbc-7eefd6a07758",
		time.Date(0, 0, 0, 06, 45, 0, 0, time.UTC),
		time.Date(0, 0, 0, 06, 55, 0, 0, time.UTC),
		1,
		Stop{
			"c8d731d1-2e1c-4ca2-ac26-fead1695320d",
			"9600",
			"Downtown Ferry Terminal Pier 1",
			-36.84243,
			174.76708,
		}}

	if !st1.IsEqual(st2) {
		t.Errorf("Identical stoptimes not equal :(")
	}
}

func TestStopTimes_IsEqual(t *testing.T) {
	st1 := StopTimeArray{
		{"847af420-5527-49b1-85f0-9a18fb143d87",
			"d0d6dece-e75d-4e1a-8bbc-7eefd6a07758",
			time.Date(0, 0, 0, 06, 45, 0, 0, time.UTC),
			time.Date(0, 0, 0, 06, 55, 0, 0, time.UTC),
			1,
			Stop{
				"c8d731d1-2e1c-4ca2-ac26-fead1695320d",
				"9600",
				"Downtown Ferry Terminal Pier 1",
				-36.84243,
				174.76708,
			}},
		{"8e89b74f-7cb2-4c33-b9a4-cf937a30ecb1",
			"d0d6dece-e75d-4e1a-8bbc-7eefd6a07758",
			time.Date(0, 0, 0, 06, 55, 0, 0, time.UTC),
			time.Date(0, 0, 0, 06, 55, 0, 0, time.UTC),
			2,
			Stop{
				"295a2c04-41f8-46f7-9904-020a51b92955",
				"9670",
				"Devonport Ferry Terminal",
				-36.83317,
				174.7954,
			}},
	}

	st2 := StopTimeArray{
		{"847af420-5527-49b1-85f0-9a18fb143d87",
			"d0d6dece-e75d-4e1a-8bbc-7eefd6a07758",
			time.Date(0, 0, 0, 06, 45, 0, 0, time.UTC),
			time.Date(0, 0, 0, 06, 55, 0, 0, time.UTC),
			1,
			Stop{
				"c8d731d1-2e1c-4ca2-ac26-fead1695320d",
				"9600",
				"Downtown Ferry Terminal Pier 1",
				-36.84243,
				174.76708,
			}},
		{"8e89b74f-7cb2-4c33-b9a4-cf937a30ecb1",
			"d0d6dece-e75d-4e1a-8bbc-7eefd6a07758",
			time.Date(0, 0, 0, 06, 55, 0, 0, time.UTC),
			time.Date(0, 0, 0, 06, 55, 0, 0, time.UTC),
			2,
			Stop{
				"295a2c04-41f8-46f7-9904-020a51b92955",
				"9670",
				"Devonport Ferry Terminal",
				-36.83317,
				174.7954,
			}},
	}

	if !st1.IsEqual(st2) {
		t.Errorf("Identical StopTimeArray not equal :(")
	}
}
