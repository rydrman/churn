package churn

import (
	"reflect"
	"testing"
)

type intInNode struct {
	BaseNode
}

func (intInNode) InValue(int64) {}

func TestOutPortCore_AddReceiver_InvalidType(t *testing.T) {

	graph := NewGraph()
	strNode := new(StringNode)
	intNode := new(intInNode)
	graph.Add("String", strNode)
	graph.Add("Int", intNode)

	strNode.Out("Value").AddReceiver(intNode.In("Value"))

}

func TestOutPortCore_AddReceiver_ReceiverIn(t *testing.T) {

	graph := NewGraph()
	strNode := new(StringNode)
	graph.Add("Source", strNode)

	valPort := strNode.Out("Value")

	err := valPort.AddReceiver(valPort)
	if err == nil {
		t.Errorf("expected an error adding input port as a receiver")
	}

}

func TestOutPortCore_handleOne(t *testing.T) {

	channel := make(chan int)
	close(channel)
	core := &OutPortCore{
		channel: reflect.ValueOf(channel),
	}
	wasHandled := core.handleOne()
	if wasHandled {
		t.Error("expected nil channel value to not handle successfully")
	}

}
