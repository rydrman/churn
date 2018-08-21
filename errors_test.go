package churn

import (
	"testing"

	"github.com/pkg/errors"
)

func TestPanicIfError(t *testing.T) {

	var err error
	panicIfError(err) // should not panic

	err = errors.New("error")
	defer func() {
		actual := recover()
		if actual != err {
			t.Errorf("expected a panic with %q, got %v", err, actual)
		}
	}()
	panicIfError(err)

}
