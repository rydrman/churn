// Package churn defines a node-based, asynchronus computational engine
package churn

import (
	"path"
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const (
	inPortNamePrefix  = "In"
	outPortNamePrefix = "Out"
)

// Graph constructs and manages connections
// between execution nodes
type Graph struct {
	BaseComponent
	components     map[string]Component
	componentMutex sync.Mutex
}

// NewGraph initializes a new Graph instance
func NewGraph() *Graph {
	return &Graph{
		components: make(map[string]Component),
	}
}

// Add adds a node to this graph, where 'name' is expected to be unique
func (g *Graph) Add(name string, cmpt Component) error {

	g.componentMutex.Lock()
	defer g.componentMutex.Unlock()

	_, exists := g.components[name]
	if exists {
		return errors.Wrap(ErrNameTaken, name)
	}

	node, isNode := cmpt.(Node)
	if isNode {
		node.initialize(reflect.ValueOf(node))
	}

	g.components[name] = cmpt
	return nil

}

// Connect joins an out port on one node to the in port of another
func (g *Graph) Connect(sourcePortPath, destPortPath string) error {

	srcPort := g.GetOutPort(sourcePortPath)
	if srcPort == nil {
		return errors.Wrap(ErrPortNotExist, sourcePortPath)
	}
	destPort := g.GetInPort(destPortPath)
	if destPort == nil {
		return errors.Wrap(ErrPortNotExist, destPortPath)
	}

	return srcPort.core.AddReceiver(destPort)

}

// GetOutPort returns the out port specified by the given
// graph path, or nil if it does not exist
func (g *Graph) GetOutPort(portPath string) *Port {

	node := g.GetNode(portPath)
	if node == nil {
		return nil
	}

	_, _, portName := SplitGraphPath(portPath)
	return node.Ports().Out(portName)

}

// GetInPort returns the in port specified by the given
// graph path, or nil if it does not exist
func (g *Graph) GetInPort(portPath string) *Port {

	node := g.GetNode(portPath)
	if node == nil {
		return nil
	}

	_, _, portName := SplitGraphPath(portPath)
	return node.Ports().In(portName)

}

// GetNode returns the node identified in the given
// graph path. If it does not exist, or the discovered
// component is not a node, then nil is returned
func (g *Graph) GetNode(nodePath string) Node {

	cmpt := g.GetComponent(nodePath)
	node, _ := cmpt.(Node)
	return node

}

// GetComponent returns the component identified in the given
// graph path, or nil if such a component does not exist
func (g *Graph) GetComponent(cmptPath string) Component {

	location, name, _ := SplitGraphPath(cmptPath)
	if location == "." {
		return g.components[name]
	}

	panic("sub-graph lookup not supported: " + cmptPath) // TODO: handle sub-graphs

}

// SplitGraphPath splits a graph path into its three
// components given the shape of the path is:
//  graph/node.port
//
// if no graph is specified, a "." is returned.
func SplitGraphPath(graphPath string) (graph, node, port string) {

	graph, node = path.Split(graphPath)

	// clean and remove trailing slashes
	graph = path.Dir(graph)

	// the current node value only requires further
	// refinement if a port was specified
	if strings.Contains(node, ".") {
		parts := strings.Split(node, ".")
		node = strings.Join(parts[:len(parts)-1], ".")
		port = parts[len(parts)-1]
	}

	return

}
