package churn

import (
	"reflect"
	"testing"
)

type intIONode struct {
	BaseNode
}

func (intIONode) InValue(int64) {}

func TestOutPortCore_AddReceiver_InvalidType(t *testing.T) {

	graph := NewGraph()
	strNode := new(StringNode)
	intNode := new(intIONode)
	graph.Add("String", strNode)
	graph.Add("Int", intNode)

	strNode.Ports().Out("Value").Core().AddReceiver(intNode.Ports().In("Value"))

}

func TestOutPortCore_AddReceiver_ReceiverIn(t *testing.T) {

	graph := NewGraph()
	strNode := new(StringNode)
	graph.Add("Source", strNode)

	valPort := strNode.Ports().Out("Value")

	err := valPort.Core().AddReceiver(valPort)
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
