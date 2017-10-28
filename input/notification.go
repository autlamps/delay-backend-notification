// Input contains the messages received by the notification service
package input

import "github.com/autlamps/delay-backend-notification/static"

type Notification struct {
	Cancelled  bool
	TripID     string
	StopTimeID string
	Delay      int
	Lat        float64
	Lon        float64
	Route      static.Route
	Trip       static.Trip
	StopTimes  []static.StopTime
}
