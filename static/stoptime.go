package static

import (
	"database/sql"
	"time"
)

// StopTime represents a specific stop on a trip. Also includes embedded Stop info
type StopTime struct {
	ID           string
	TripID       string
	Arrival      time.Time
	Departure    time.Time
	StopSequence int
	StopInfo     Stop
}

// Stop represents a physical stop. Embedded into StopTime instead of being its own service for ease of use
type Stop struct {
	ID   string
	Code string
	Name string
	Lat  float64
	Lon  float64
}

// StopTimeArray is simply a slice of StopTime
type StopTimeArray []StopTime

// StopTimeStore defines the methods that a concrete StopTimeService should implement
type StopTimeStore interface {
	GetStopTimesByTripID(tripID string) (StopTimeArray, error)
	getStopByID(id string) (Stop, error)
}

// StopTimeService implements StopTimeStore in PSQL
type StopTimeService struct {
	db *sql.DB
}

// StopTimeServiceInit initializes a new StopTimeService
func StopTimeServiceInit(db *sql.DB) *StopTimeService {
	return &StopTimeService{db: db}
}

// GetStopTimesByTripID returns all stops of the given trip
func (sts *StopTimeService) GetStopTimesByTripID(tripID string) (StopTimeArray, error) {
	var sta StopTimeArray

	rows, err := sts.db.Query("SELECT stoptime_id, trip_id, arrival_time, departure_time, stop_id, "+
		"stop_sequence from stop_times WHERE trip_id = $1 ORDER BY stop_sequence ASC", tripID)

	if err != nil {
		return sta, err
	}

	for rows.Next() {
		st := StopTime{}
		var stopID string

		if err := rows.Scan(&st.ID, &st.TripID, &st.Arrival, &st.Departure, &stopID, &st.StopSequence); err != nil {
			return sta, err // TODO: decide what to do here. Do we inject logger and log it? Stop execution?
		}

		st.StopInfo, err = sts.getStopByID(stopID)

		if err != nil {
			return sta, err
		}

		sta = append(sta, st)
	}

	return sta, nil
}

// getStopByID returns a single stop
func (sts *StopTimeService) getStopByID(id string) (Stop, error) {
	s := Stop{}

	row := sts.db.QueryRow("SELECT stop_id, stop_code, stop_name, stop_lat, stop_lon FROM stops WHERE stop_id = $1", id)

	err := row.Scan(&s.ID, &s.Code, &s.Name, &s.Lat, &s.Lon)

	if err != nil {
		return s, err
	}

	return s, nil
}

// IsEqual returns true if the given StopTime is equal to the StopTime the method is run on
func (st StopTime) IsEqual(x StopTime) bool {

	if st.ID != x.ID {
		return false
	}
	if st.TripID != x.TripID {
		return false
	}
	if !st.Arrival.Equal(x.Arrival) {
		return false
	}
	if !st.Departure.Equal(x.Departure) {
		return false
	}
	if st.StopSequence != x.StopSequence {
		return false
	}
	if !st.StopInfo.IsEqual(x.StopInfo) {
		return false
	}

	return true
}

// IsEqual returns true if the given Stop is equal to the Stop the method is run on
func (s Stop) IsEqual(x Stop) bool {

	if s.ID != x.ID {
		return false
	}
	if s.Lon != x.Lon {
		return false
	}
	if s.Lat != x.Lat {
		return false
	}
	if s.Name != x.Name {
		return false
	}

	return true
}

// IsEqual returns true if the given StopTimeArray is equal to the StopTimeArray the method is run on
func (st StopTimeArray) IsEqual(x StopTimeArray) bool {

	if len(st) != len(x) {
		return false
	}

	for i := 0; i < len(st); i++ {
		if !st[i].IsEqual(x[i]) {
			return false
		}
	}

	return true
}
