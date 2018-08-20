package churncore

import (
	"reflect"
	"sync"

	"github.com/satori/go.uuid"

	"github.com/pkg/errors"
)

var (
	errNotAChannel          = errors.New("not a channel")
	errSendOnly             = errors.New("cannot be a send-only channel")
	errIncompatibleReceiver = errors.New("incompatible receiver type")
)

// Sender produces messages that can be handled by receivers
type Sender struct {
	dataType reflect.Type
	channel  reflect.Value
	subs     map[uuid.UUID]*Subscription

	// mutex is held anytime a value is being sent or
	// the subscription set is being modified
	mutex sync.Mutex
}

// NewSender creates a new source for messages of the given go data type
func NewSender(channel interface{}) (*Sender, error) {

	chanVal := reflect.ValueOf(channel)
	chanType := chanVal.Type()

	if chanType.Kind() != reflect.Chan {
		return nil, errors.Wrapf(errNotAChannel, "invalid type %T", channel)
	}

	if chanType.ChanDir() == reflect.SendDir {
		return nil, errSendOnly
	}

	return &Sender{
		dataType: chanType.Elem(),
		channel:  chanVal,
		subs:     make(map[uuid.UUID]*Subscription),
	}, nil

}

// Subscribe creates a subscription from this sender to the
// given receiver, which will cause the receiver's underlying
// function to be called for every sent value until the
// subscription is closed
func (s *Sender) Subscribe(r *Receiver) (*Subscription, error) {

	if !s.dataType.AssignableTo(r.dataType) {
		return nil, errors.Wrapf(
			errIncompatibleReceiver,
			"cannot assign [%s] -> [%s]", s.dataType, r.dataType,
		)
	}

	id := uuid.NewV4()
	subs := &Subscription{
		onClose: func() {
			s.mutex.Lock()
			s.subs[id] = nil
			s.mutex.Unlock()
		},
	}

	s.mutex.Lock()
	s.subs[id] = subs
	s.mutex.Unlock()

	return subs, nil

}
