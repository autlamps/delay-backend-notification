package static

import "database/sql"

// Route represents a route stored in database
type Route struct {
	ID        string
	GTFSID    string
	AgencyID  string
	ShortName string
	LongName  string
}

// RouteStore defines methods that a concrete implementation should implement
type RouteStore interface {
	GetRouteByID(id string) (Route, error)
}

// RouteService is a psql implementation of RouteStore
type RouteService struct {
	db *sql.DB
}

// RouteServiceInit initializes a RouteService
func RouteServiceInit(db *sql.DB) *RouteService {
	return &RouteService{db: db}
}

// GetRouteByID returns a a route with the given id or an error
func (rs *RouteService) GetRouteByID(id string) (Route, error) {
	r := Route{}

	row := rs.db.QueryRow("SELECT route_id, gtfs_route_id, agency_id, route_short_name, route_long_name FROM routes where route_id = $1", id)
	err := row.Scan(&r.ID, &r.GTFSID, &r.AgencyID, &r.ShortName, &r.LongName)

	if err != nil {
		return r, err
	}

	return r, nil
}

// IsEqual returns true if the given route is equal to the route this method is being called on
func (r Route) IsEqual(x Route) bool {

	if r.ID != x.ID {
		return false
	}

	if r.GTFSID != x.GTFSID {
		return false
	}

	if r.AgencyID != x.AgencyID {
		return false
	}

	if r.ShortName != x.ShortName {
		return false
	}

	if r.LongName != x.LongName {
		return false
	}

	return true
}
