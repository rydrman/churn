package churncore

import (
	"testing"

	"github.com/pkg/errors"
)

func TestNewSender(t *testing.T) {

	_, err := NewSender("string")
	if errors.Cause(err) != errNotAChannel {
		t.Errorf("expected a non-channel to return a relevant error, got: %s", err)
	}

	_, err = NewSender(make(chan<- int))
	if errors.Cause(err) != errSendOnly {
		t.Errorf("expected a send-only channel to return relevant error, got: %s", err)
	}

	_, err = NewSender(make(chan int))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

}

func TestSender_Subscribe_Valid(t *testing.T) {

	ch := make(chan string)
	sender, err := NewSender(ch)
	if err != nil {
		t.Fatal(err)
	}

	receiver, err := NewReceiver(func(string) {})
	if err != nil {
		t.Fatal(err)
	}

	subs, err := sender.Subscribe(receiver)
	if err != nil {
		t.Fatal(err)
	}

	subs.Close()

}

func TestSender_Subscribe_Invalid(t *testing.T) {

	ch := make(chan string)
	sender, err := NewSender(ch)
	if err != nil {
		t.Fatal(err)
	}

	receiver, err := NewReceiver(func(int) {})
	if err != nil {
		t.Fatal(err)
	}

	_, err = sender.Subscribe(receiver)
	if errors.Cause(err) != errIncompatibleReceiver {
		t.Errorf("expected relevant error subscribing with incompatible type, got: %s", err)
	}

}
