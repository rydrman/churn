package churn

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// OutPortCore is the core logic implementation for out ports
type OutPortCore struct {

	// the underlying data type that this port outputs
	dataType reflect.Type

	// the source channel of this port
	channel           reflect.Value
	channelBufferSize int

	// functions that accept a single parameter of the same type
	// as this port's driving channel
	receivers []reflect.Value
}

// NewOutPortCore initializes a port core for an output port
// driven by the given struct field. 'chanField' must be
// a field whose type is a send-only channel. 'bufferSize' is the
// size of buffer to use for the newly created channel
func NewOutPortCore(chanField reflect.StructField, bufferSize int) (*OutPortCore, error) {

	// must be a channel
	if chanField.Type.Kind() != reflect.Chan {
		return nil, errors.New("field type is not a channel")
	}

	// must be send-only
	if chanField.Type.ChanDir() != reflect.SendDir {
		return nil, errors.New("field type is not a send-only channel")
	}

	core := &OutPortCore{
		dataType: chanField.Type.Elem(),
	}

	// the chanField is expected to be unidirectional, but to
	// create a channel we need to define the bidirectional type
	chanType := reflect.ChanOf(reflect.BothDir, core.dataType)

	core.channel = reflect.MakeChan(chanType, bufferSize)
	go func() {
		alive := true
		for alive {
			alive = core.handleOne()
		}
	}()

	return core, nil

}

// returns false if the underlying channel has been
// closed or became nil
func (c *OutPortCore) handleOne() (wasHandled bool) {

	val, ok := c.channel.Recv()
	if !ok {
		return ok
	}
	for _, dest := range c.receivers {
		dest.Call([]reflect.Value{val})
	}
	return true

}

// AddReceiver registers the given port to receive new values
// coming from this out port
func (c *OutPortCore) AddReceiver(dest *Port) error {

	destCore, isInput := dest.PortCore.(*InPortCore)
	if !isInput {
		return errors.New("destination is not an in port")
	}

	if !c.dataType.AssignableTo(destCore.dataType) {
		return fmt.Errorf(
			"destination port %q does not support data of type %s (want %s)",
			dest.Name, c.dataType, destCore.dataType,
		)
	}

	c.receivers = append(c.receivers, destCore.function)
	return nil

}

func (c *OutPortCore) close() {

	c.channel.Close()

}
