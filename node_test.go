package churn

import (
	"reflect"
	"testing"
)

func TestBaseNode_catalogInPorts(t *testing.T) {

	n := new(PrintNode)
	n.catalogInPorts(reflect.ValueOf(n))

	port := n.In("Message")
	if port == nil {
		t.Fatal("expected printer Message port to be cataloged")
	}

}

func TestBaseNode_catalogOutPorts(t *testing.T) {

	n := &struct {
		BaseNode
		OutString chan<- string `desc:"valid port"`
		Other     chan<- string `desc:"should not be consider when prefix is missing"`
		OutAndIn  chan string   `desc:"should not consider recv channel"`
		OutOther  int           `desc:"not a port but starts with Out"`
	}{}

	n.catalogOutPorts(reflect.ValueOf(n), 0)

	if len(n.Outs) > 1 {
		t.Errorf("expected only 1 port to be cataloged, got %d", len(n.Outs))
	}

	port := n.Out("String")
	if port == nil {
		t.Fatal("expected OutValue cahnnel field to be cataloged as a port")
	}

	port = n.Out("Other")
	if port != nil {
		t.Error("non channel field should not be ned as a port")
	}

	port = n.Out("AndIn")
	if port != nil {
		t.Error("recv channel field should not be cataloged as a port")
	}

}
