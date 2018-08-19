package churn

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func Example() {

	var (
		graph   = NewGraph()
		strNode = new(StringNode)
	)

	err := graph.Add("MessageSource", strNode)
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to add message source"))
	}

	err = graph.Add("Printer", new(PrintNode))
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to add printer"))
	}

	err = graph.Connect("MessageSource.Value", "Printer.Message")
	if err != nil {
		fmt.Println(errors.Wrap(err, "failed to connect"))
	}

	strNode.OutValue <- "Hello, World!"

	// let the network propagate
	time.Sleep(1 * time.Millisecond)

	// Output:
	// Hello, World!

}

func TestGraph_Connect(t *testing.T) {

	g := NewGraph()
	g.Add("NodeOut", new(StringNode))
	g.Add("NodeIn", new(PrintNode))

	err := g.Connect("Unknown.Value", "NodeIn.Message")
	if !IsPortNotExist(err) {
		t.Errorf("expected NodeNotExist when connecting from unknown node, got: %s", err)
	}

	err = g.Connect("NodeOut.Value", "Unkown.Message")
	if !IsPortNotExist(err) {
		t.Errorf("expected NodeNotExist when connecting into unknown node, got: %s", err)
	}

}

func TestGraph_Add(t *testing.T) {

	g := NewGraph()
	err := g.Add("MyNode", new(StringNode))
	if err != nil {
		t.Error(errors.Wrap(err, "unexpected error adding string node"))
	}
	err = g.Add("MyNode", new(StringNode))
	if !IsNameTaken(err) {
		t.Errorf("expected ErrNameTaken when adding duplicate node, got %v", err)
	}

}

func TestSplitGraphPath_FullPath(t *testing.T) {

	loc, name, port := SplitGraphPath("loc/name.port")
	if loc != "loc" || name != "name" || port != "port" {
		t.Errorf("expected (loc, name, port), got: (%s, %s, %s)", loc, name, port)
	}

}

func TestSplitGraphPath_PartialPaths(t *testing.T) {

	loc, name, port := SplitGraphPath("name.port")
	if loc != "." || name != "name" || port != "port" {
		t.Errorf("expected (., name, port), got: (%s, %s, %s)", loc, name, port)
	}

	loc, name, port = SplitGraphPath("name")
	if loc != "." || name != "name" || port != "" {
		t.Errorf("expected (., name, ), got: (%s, %s, %s)", loc, name, port)
	}

	loc, name, port = SplitGraphPath(".port")
	if loc != "." || name != "" || port != "port" {
		t.Errorf("expected (., , port), got: (%s, %s, %s)", loc, name, port)
	}

}
