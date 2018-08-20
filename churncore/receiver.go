package churncore

import (
	"reflect"

	"github.com/pkg/errors"
)

var (
	errNotAFunction         = errors.New("receiver must be a function")
	errWrongNumberOfArgs    = errors.New("receiver func must take exactly one argument")
	errWrongNumberOfReturns = errors.New("receiver func may only return a single error value")
)

// Receiver represents a function that can handle messages of
// a specific go data type
type Receiver struct {
	function reflect.Value
	dataType reflect.Type
}

// NewReceiver creates a message receiver from the given function.
// 'handlerFunc' must take a single parameter of the desired message
// data type and may return nothing, or a single error, as required
func NewReceiver(handlerFunc interface{}) (*Receiver, error) {

	funcVal := reflect.ValueOf(handlerFunc)
	funcType := funcVal.Type()

	if funcType.Kind() != reflect.Func {
		return nil, errors.Wrapf(errNotAFunction, "invalid type %T", handlerFunc)
	}

	if funcType.NumIn() != 1 {
		return nil, errWrongNumberOfArgs
	}

	if funcType.NumOut() > 1 {
		return nil, errWrongNumberOfReturns
	}

	return &Receiver{
		dataType: funcType.In(0),
		function: funcVal,
	}, nil

}
