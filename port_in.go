package churn

import (
	"reflect"

	"github.com/pkg/errors"
)

// InPortCore is the core logic implementation for in ports
type InPortCore struct {

	// the function to be called in response to values from
	// upstream connections
	function reflect.Value

	// the underlying data type that this port accepts
	dataType reflect.Type
}

// NewInPortCore creates a new input port from the given function value
//
// 'function' should take a single parameter of the expected data type
// for this port. During graph execution, any other parameters will
// be passed as their zero value, and any return values will be
// silently discarded
func NewInPortCore(function reflect.Value) *InPortCore {

	function.Type()
	return &InPortCore{
		dataType: function.Type().In(0),
		function: function,
	}

}

// AddReceiver always returns an error as receivers are not
// supported on in ports
func (c *InPortCore) AddReceiver(*Port) error {
	return errors.New("AddReceiver called on input port")
}

// PortSlice provides helper methods for working with
// slices of ports
type PortSlice []*Port

// FindByName returns the first port in this slice with the given
// name, or nil if no port has such a name
func (p PortSlice) FindByName(name string) *Port {

	for _, port := range p {
		if name == port.Name {
			return port
		}
	}
	return nil

}
