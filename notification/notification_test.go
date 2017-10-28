package notification

import (
	"testing"

	"github.com/autlamps/delay-backend-notification/static"
)

func TestFindStopTimeIndex(t *testing.T) {

	data := []static.StopTime{
		{ID: "3d72f2f8-ba68-41f4-988f-1deec9e64db0"},
		{ID: "121e1f20-5c84-415e-9f40-db04b3c3627a"},
		{ID: "a6d1d7f6-64e1-417e-b54e-4b9c071d4e53"},
	}

	tests := []struct {
		ID          string
		ExpectedInt int
		ExpectedErr error
	}{
		{"3d72f2f8-ba68-41f4-988f-1deec9e64db0", 0, nil},
		{"deba62ca-4188-49e6-bfba-9b15f883ef52", -1, ErrIDNotInSlice},
	}

	for _, test := range tests {

		loc, err := FindStopTimeIndex(data, test.ID)

		if loc != test.ExpectedInt {
			t.Fatalf("Location not correct. Expected %v, got %v", test.ExpectedInt, loc)
		}

		if err != test.ExpectedErr {
			t.Fatalf("Error not correct. Expected %v, got %v", test.ExpectedErr, err)
		}
	}
}
