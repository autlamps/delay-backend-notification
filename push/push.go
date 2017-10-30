// Push implements firebase push notifications
package push

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrInvalidJSONGiven = errors.New("push - Firebase: Invalid JSON sent to firebase")
var ErrInvalidFirebaseAuth = errors.New("push - Firebase: invalid auth token given")

type Pusher interface {
	Send(to, t, b string, d interface{}) error
}

type Firebase struct {
	serverKey string
	testing   bool
}

// InitFirebase returns a firebase service we can use to send push notification
func InitFirebase(key string) *Firebase {
	return &Firebase{serverKey: key, testing: false}
}

const FIREBASE_URL = "https://fcm.googleapis.com/fcm/send"

// Send sends a push notification to the user via firebase
func (f *Firebase) Send(to, t, b string, d interface{}) error {
	n := struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		Sound string `json:"sound"`
	}{
		Title: t,
		Body:  b,
		Sound: "default",
	}

	r := struct {
		To           string      `json:"to"`
		Notification interface{} `json:"notification"`
		TimeToLive   int         `json:"time_to_live"`
		Data         interface{} `json:"data"`
		DryRun       bool        `json:"dry_run"`
	}{
		To:           to,
		Notification: n,
		Data:         d,
		TimeToLive:   150,
		DryRun:       f.testing,
	}

	njsn, err := json.Marshal(r)

	if err != nil {
		return fmt.Errorf("push - Send: Failed to marshal JSON notification: %v", err)
	}

	c := http.Client{Timeout: time.Second * time.Duration(20)}

	req, _ := http.NewRequest("POST", FIREBASE_URL, bytes.NewBuffer(njsn))
	req.Header.Add("Authorization", fmt.Sprintf("key=%v", f.serverKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)

	if err != nil {
		return fmt.Errorf("push - Send: failed to call firebase: %v", err)
	}

	if resp.StatusCode == 400 {
		return ErrInvalidJSONGiven
	}

	if resp.StatusCode == 401 {
		return ErrInvalidFirebaseAuth
	}

	if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
		return fmt.Errorf("push - Send: firebase error code: %v", resp.StatusCode)
	}

	return nil
}
