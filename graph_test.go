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

	messageSource := graph.SafeAdd("MessageSource", strNode)
	printer := graph.SafeAdd("Printer", new(PrintNode))

	err := graph.Connect(messageSource+".Value", printer+".Message")
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

func TestGraph_SafeAdd(t *testing.T) {

	graph := NewGraph()
	node0 := new(StringNode)
	node1 := new(StringNode)
	node2 := new(StringNode)

	desired := "Node"
	node0Name := graph.SafeAdd(desired, node0)
	node1Name := graph.SafeAdd(desired, node1)
	node2Name := graph.SafeAdd(desired, node2)

	if node0Name != "Node" {
		t.Errorf(
			"expected first node name not to be changed from %q, got %q",
			desired, node0Name,
		)
	}

	if node1Name != "Node1" {
		t.Errorf(
			"expected second node name to be changed to %q, got %q",
			"Node1", node1Name,
		)
	}

	if node2Name != "Node2" {
		t.Errorf(
			"expected first node name to be changed to %q, got %q",
			"Node2", node2Name,
		)
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
