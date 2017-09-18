// Input contains the messages received by the notification service
package input

type Notification struct {
	TripID     string
	StopTimeID string
	Delay      int
	Lat        float64
	Lon        float64
}
