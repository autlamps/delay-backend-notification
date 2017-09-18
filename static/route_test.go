package static

import (
	"database/sql"
	"testing"
)

func TestRouteService_GetRouteByID(t *testing.T) {
	db, err := sql.Open("postgres", dburl)
	defer db.Close()

	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	rs := RouteServiceInit(db)

	tests := []struct {
		id       string
		expected Route
	}{
		{"42d70439-930a-4149-b690-37010dbf36ea",
			Route{"42d70439-930a-4149-b690-37010dbf36ea",
				"12014-20170829094406_v57.13",
				"cb90705c-32b2-4c27-a808-a0d519548b8c",
				"120",
				"Akoranga to Henderson"}},
	}

	for _, test := range tests {
		route, err := rs.GetRouteByID(test.id)

		if err != nil {
			t.Errorf("Failed to retrieve route. %v", err)
		}

		if !route.IsEqual(test.expected) {
			t.Errorf("Failed to retrieve correct trip. Expected id %v, got %v", test.expected.ID, route.ID)
		}
	}
}

func TestRoute_IsEqual(t *testing.T) {

	r1 := Route{"3ad6312a-9a56-4bd5-9b89-4c0b9687db95",
		"30009-20170724124507_v56.18",
		"09e076cf-e471-453e-adcc-3da322502160",
		"120",
		"Akoranga to Henderson"}

	r2 := Route{"3ad6312a-9a56-4bd5-9b89-4c0b9687db95",
		"30009-20170724124507_v56.18",
		"09e076cf-e471-453e-adcc-3da322502160",
		"120",
		"Akoranga to Henderson"}

	if !r1.IsEqual(r2) {
		t.Errorf("Identical routes not equal :(")
	}
}
