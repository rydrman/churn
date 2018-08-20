package churncore

// Subscription connects a sender to a compatible receiver
// and manages the transfer of messages between them
type Subscription struct {
	sender   *Sender
	receiver *Receiver
	onClose  func()
}

// Close ends this subscription
func (s *Subscription) Close() {
	if s.onClose != nil {
		s.onClose()
	}
}
