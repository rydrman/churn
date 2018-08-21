package churn

import (
	"reflect"
	"testing"
)

type inputTester struct{ BaseNode }

func (*inputTester) InValue(int) {}
func (*inputTester) InNothing()  {}

func TestPortCatalog_catalogInPorts(t *testing.T) {

	n := new(inputTester)
	n.catalogInPorts(reflect.ValueOf(n))

	port := n.In("Value")
	if port == nil {
		t.Error("expected InValue(int) function to be cataloged")
	}

	port = n.In("Nothing")
	if port != nil {
		t.Error("expected function with no argument not to be cataloged")
	}

}

func TestPortCatalog_catalogOutPorts(t *testing.T) {

	n := &struct {
		BaseNode
		OutString chan<- string `desc:"valid port"`
		Out       chan<- string `desc:"should not consider prefix only"`
		Other     chan<- string `desc:"should not consider when prefix is missing"`
		Outher    chan<- string `desc:"should not consider when not camelcase"`
		OutAndIn  <-chan int    `desc:"should not consider recv channel"`
		OutOther  int           `desc:"not a port but starts with Out"`
	}{}

	n.catalogOutPorts(reflect.ValueOf(n))

	if len(n.Outs) > 1 {
		t.Errorf("expected only 1 port to be cataloged, got %d", len(n.Outs))
	}

	port := n.PortCatalog.Out("String")
	if port == nil {
		t.Fatal("expected OutValue cahnnel field to be cataloged as a port")
	}

	port = n.PortCatalog.Out("Other")
	if port != nil {
		t.Error("non channel field should not be seen as a port")
	}

	port = n.PortCatalog.Out("AndIn")
	if port != nil {
		t.Error("recv channel field should not be cataloged as a port", port)
	}

}
