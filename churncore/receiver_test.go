package churncore

import (
	"testing"

	"github.com/pkg/errors"
)

func TestNewReceiver(t *testing.T) {

	_, err := NewReceiver("string")
	if errors.Cause(err) != errNotAFunction {
		t.Errorf("expected non-function receiver to give relevant error, got: %s", err)
	}

	_, err = NewReceiver(func() {})
	if errors.Cause(err) != errWrongNumberOfArgs {
		t.Errorf("expected function with no arguments to give relevant error, got: %s", err)
	}

	_, err = NewReceiver(func(int) (a, b error) { return })
	if errors.Cause(err) != errWrongNumberOfReturns {
		t.Errorf("expected function with 2 return values to give relevant error, got: %s", err)

	}

	_, err = NewReceiver(func(int) error { return nil })
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	_, err = NewReceiver(func(int) {})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

}
