package notification

import (
	"database/sql"
	"fmt"

	"encoding/json"
	"log"

	"errors"

	"github.com/autlamps/delay-backend-notification/data"
	"github.com/autlamps/delay-backend-notification/input"
	"github.com/autlamps/delay-backend-notification/notify"
	"github.com/autlamps/delay-backend-notification/static"

	"time"

	"sync"

	"github.com/autlamps/delay-backend-notification/push"
	_ "github.com/lib/pq"
)

const LOOK_AHEAD = 5

// Conf stores our initial string connection values before being turned into services
type Conf struct {
	DBURL         string
	MQURL         string
	FirebaseToken string
}

// Env stores our services to be used
type Env struct {
	Notification  notify.Notifier
	Trips         static.TripStore
	StopTimes     static.StopTimeStore
	Routes        static.RouteStore
	Subscriptions data.SubscriptionStore
	NotifyInfo    data.NotifyInfoStore
	Push          push.Pusher
	wg            sync.WaitGroup
}

// EnvFromConf
func EnvFromConf(c Conf) (Env, error) {
	db, err := sql.Open("postgres", c.DBURL)

	if err != nil {
		return Env{}, fmt.Errorf("Failed to open db connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return Env{}, fmt.Errorf("Failed to ping db connection: %v", err)
	}

	n, err := notify.InitService(c.MQURL)

	if err != nil {
		return Env{}, err
	}

	return Env{
		Notification:  n,
		Trips:         static.TripServiceInit(db),
		StopTimes:     static.StopTimeServiceInit(db),
		Routes:        static.RouteServiceInit(db),
		Subscriptions: data.InitSubscriptionService(db),
		NotifyInfo:    data.InitNotifyInfoService(db),
		Push:          push.InitFirebase(c.FirebaseToken),
	}, nil
}

// Start contains our main loop
func (e *Env) Start(ec <-chan bool) {
	nc, err := e.Notification.Receive()

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case _ = <-ec:
			log.Printf("Exit recieved. Waiting on goroutines to finish.")
			e.wg.Wait()
			return // Exit if we receive value on exit channel and all our goroutines are finished
		case nm := <-nc: // Receive from collection service
			b := nm.Body
			n := input.Notification{}

			err := json.Unmarshal(b, &n)

			if err != nil {
				log.Printf("notification - Start: failed to unmarshal json of body: Body %v, err %v", b, err)
				continue
			}

			stopTimeIndex, err := findStopTimeIndex(n.StopTimes, n.StopTimeID)

			if err != nil {
				log.Printf("notification - Start: notification id not in returned stoptimes: Expected %v, got %v", n.StopTimeID, n.StopTimes)
				continue
			}

			// Get all stoptimes from our current stoptime until the end of the trip
			stsToNotify := n.StopTimes[stopTimeIndex:]

			for i, st := range stsToNotify {
				// We only want to look ahead and notify for stoptimes up to LOOK_AHEAD times
				if i >= LOOK_AHEAD {
					break
				}

				go e.processStopTime(st, n)
			}

			nm.Ack(true)
		}
	}
}

// processStopTime obtains all subscriptions for a stoptime id then dispatches goroutines to notify them
func (e *Env) processStopTime(st static.StopTime, n input.Notification) {
	e.wg.Add(1) // Tell our wait group that we're doing work and to not stop on us

	subs, err := e.Subscriptions.GetSubsByStopTimeID(st.ID)

	if err != nil {
		log.Printf("notification - processStopTime: failed to get subscriptions for stoptime id: %v. Error: %v", n.StopTimeID, err)
		return
	}

	// If nobody is subscribed then skip over this stoptime
	if len(subs) < 1 {
		return
	}

	eta := st.Arrival.Add(time.Second * time.Duration(n.Delay))

	for _, s := range subs {
		go e.processSubscription(s, eta, st, n)
	}

	e.wg.Done() // Tell our wait group that we're done
}

// processSubscription gets all notification methods of a single subscription and notifies them
func (e *Env) processSubscription(s data.Subscription, eta time.Time, st static.StopTime, n input.Notification) {
	e.wg.Add(1) // Tell our wait group that we're doing work and to not stop on us

	recentlyNotified, err := e.Subscriptions.RecentlyNotified(s.ID)

	if err != nil {
		log.Printf("notification - processSubscriptions: err retrieving whether or not sub was recently notified: %v", err)
		return
	}

	// If user was recently notified about this subscription then skip it and move onto the next
	if recentlyNotified {
		return
	}

	// If user doesn't wish to receive notifications for this subscription for today then don't notify them
	if !s.SubscribedForToday() {
		return
	}

	var lateEarly string

	if n.Delay > 0 {
		lateEarly = "Late"
	} else {
		lateEarly = "Early"
	}

	title := fmt.Sprintf("%v is Running %v", n.Route.LongName, lateEarly)
	message := fmt.Sprintf("%v expected to arrive at %v by %v", n.Route.ShortName, st.StopInfo.Code, eta.Format("15:04"))

	nData := struct {
		StopTime static.StopTime `json:"stop_time"`
		Route    static.Route    `json:"route"`
		Delay    int             `json:"delay"`
		Eta      string          `json:"eta"`
	}{
		StopTime: st,
		Route:    n.Route,
		Delay:    n.Delay,
		Eta:      eta.Format("15:04"),
	}

	for _, nid := range s.NotificationIDs {
		ntfy, err := e.NotifyInfo.Get(nid)

		if err != nil {
			log.Printf("notification - processSubscriptions: failed to get notification method: %v", err)
			return
		}

		switch ntfy.Type {
		case data.PUSH:
			err := e.Push.Send(ntfy.Value, title, message, nData)

			if err != nil {
				log.Printf("notification - processSubscription: failed to call Push.Send: %v", err)
			}
		case data.TXT: // not implemented... clearly
			fmt.Printf("txt notification: %v", eta)
		case data.EMAIL: // not implemented... clearly
			fmt.Printf("email notification: %v", eta)
		}
	}

	err = e.Subscriptions.Notified(s)

	if err != nil {
		fmt.Printf("notification - processSubscriptions: failed to set sub as notified: %v", err)
		return
	}

	e.wg.Done() // Tell our wait group that we're done
}

var ErrIDNotInSlice = errors.New("notification - given id not in slice")

// findStopTimeIndex takes a "haystack" of stoptimes and a "needle" stoptime id and returns the index of that needle
func findStopTimeIndex(h []static.StopTime, n string) (int, error) {
	for i, s := range h {
		if s.ID == n {
			return i, nil
		}
	}

	return -1, ErrIDNotInSlice
}
