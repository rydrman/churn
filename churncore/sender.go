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

// NewSender creates a new message sender using the given go channel.
// 'channel' must be receive-able, and the returned Sender will continually
// consume all channel values until the channel is closed
func NewSender(channel interface{}) (*Sender, error) {

	chanVal := reflect.ValueOf(channel)
	chanType := chanVal.Type()

	if chanType.Kind() != reflect.Chan {
		return nil, errors.Wrapf(errNotAChannel, "invalid type %T", channel)
	}

	if chanType.ChanDir() == reflect.SendDir {
		return nil, errSendOnly
	}

	s := &Sender{
		dataType: chanType.Elem(),
		channel:  chanVal,
		subs:     make(map[uuid.UUID]*Subscription),
	}

	go func() {
		alive := true
		for alive {
			alive = s.handleOne()
		}
	}()

	return s, nil

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
		sender:   s,
		receiver: r,
		onClose: func() {
			s.mutex.Lock()
			delete(s.subs, id)
			s.mutex.Unlock()
		},
	}

	s.mutex.Lock()
	s.subs[id] = subs
	s.mutex.Unlock()

	return subs, nil

}

// returns false if the underlying channel has been
// closed or became nil
func (s *Sender) handleOne() (wasHandled bool) {

	val, ok := s.channel.Recv()
	if !ok {
		return ok
	}
	for _, sub := range s.subs {
		sub.receiver.function.Call([]reflect.Value{val})
	}
	return true

}
