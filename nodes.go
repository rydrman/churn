package churn

import "fmt"

// StringNode is a graph component that outputs a single string value
type StringNode struct {
	BaseNode

	// OutValue is the output port of this node
	OutValue chan<- string
}

// PrintNode is a graph component that prints a string to stdout
type PrintNode struct {
	BaseNode
}

// InMessage is an input port which prints
// arbitrary message data to stdout
func (n *PrintNode) InMessage(msg interface{}) {

	fmt.Println(msg)

}

// IntNode is a graph component that outputs a single 64-bit integer
type IntNode struct {
	BaseNode

	// OutValue is the output port of this node
	OutValue chan<- int64
}

// FloatNode is a graph component that outputs a single 64-bit float
type FloatNode struct {
	BaseNode

	// OutValue is the output port of this node
	OutValue chan<- float64
}
