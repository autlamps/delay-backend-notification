package push

import (
	"flag"
	"os"
	"testing"
)

var firebaseKey string
var device string

func init() {
	flag.StringVar(&firebaseKey, "FIREBASE_KEY", "", "firebase url")
	flag.StringVar(&device, "DEVICE", "", "device token to send notification to -- wont actually be notified")
	flag.Parse()

	if firebaseKey == "" {
		os.Getenv("FIREBASE_KEY")
	}

	if device == "" {
		os.Getenv("DEVICE")
	}
}

func TestFirebase_Send(t *testing.T) {
	fb := Firebase{serverKey: firebaseKey, testing: true}

	td := struct{ Foo string }{"bar"}

	err := fb.Send(device, "Hello", "Please work", td)

	if err != nil {
		t.Fatalf("Failed to push test to firebase: %v", err)
	}
}
