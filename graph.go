// Package churn defines a node-based, asynchronus computational engine
package churn

import (
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/rydrman/churn/churncore"

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

	channelBufferSize int
}

// NewGraph initializes a new Graph instance
func NewGraph(options ...GraphOption) *Graph {

	g := &Graph{
		components: make(map[string]Component),
	}
	for _, option := range options {
		option.Apply(g)
	}
	return g

}

// SafeAdd adds the given component to this graph. If a component with
// 'desiredName' already exists, then an integer will be appended to
// the end of the name such that if becomes unique. The returned
// string is the final accepted name of the added node.
func (g *Graph) SafeAdd(desiredName string, cmpt Component) (name string) {

	name = desiredName
	err := g.Add(desiredName, cmpt)
	for count := 1; IsNameTaken(err); count++ {
		name = desiredName + strconv.Itoa(count)
		err = g.Add(name, cmpt)
	}
	return

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
		node.setupBaseNode(node)
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

	sender := srcPort.core.(*churncore.Sender)
	receiver := destPort.core.(*churncore.Receiver)

	_, err := sender.Subscribe(receiver)
	// TODO: store the subscription for cleanup
	return err

}

// GetOutPort returns the out port specified by the given
// graph path, or nil if it does not exist
func (g *Graph) GetOutPort(portPath string) *Port {

	node := g.GetNode(portPath)
	if node == nil {
		return nil
	}

	_, _, portName := SplitGraphPath(portPath)
	return node.Out(portName)

}

// GetInPort returns the in port specified by the given
// graph path, or nil if it does not exist
func (g *Graph) GetInPort(portPath string) *Port {

	node := g.GetNode(portPath)
	if node == nil {
		return nil
	}

	_, _, portName := SplitGraphPath(portPath)
	return node.In(portName)

}

// GetNode returns the node identified in the given
// graph path. If it does not exist, or the discovered
// component is not a node, then nil is returned
func (g *Graph) GetNode(nodePath string) Node {

	cmpt := g.GetComponent(nodePath)
	node, _ := cmpt.(Node)
	return node

}

// GetSubGraph returns the sub-graph identified in the given
// graph path. If it does not exist, or the discovered
// component is not a sub-graph, then nil is returned
func (g *Graph) GetSubGraph(graphPath string) *Graph {

	cmpt := g.GetComponent(graphPath)
	graph, _ := cmpt.(*Graph)
	return graph

}

// GetComponent returns the component identified in the given
// graph path, or nil if such a component does not exist
func (g *Graph) GetComponent(cmptPath string) Component {

	location, name, port := SplitGraphPath(cmptPath)
	if location == "." {
		return g.components[name]
	}

	parts := strings.Split(location, "/")
	subGraphName := parts[0]
	subGraphPath := BuildGraphPath(path.Join(parts[1:]...), name, port)
	subGraph := g.GetSubGraph(subGraphName)
	if subGraph == nil {
		return nil
	}
	return subGraph.GetComponent(subGraphPath)

}

// Close ends all node execution and tears down the graph node network
func (g *Graph) Close() {

	for _, cmpt := range g.components {
		cmpt.close()
	}
	return

}

// BuildGraphPath cleans and construct a valid graph path string
// from the given components. Any parameter may be an empty string
// to omit that portion of the path, although relative or partial
// paths may only be valid in some contexts
func BuildGraphPath(location, component, port string) string {

	p := location
	if location != "" {
		p = path.Clean(location)
	}
	p = path.Join(p, component)
	if port != "" {
		p = p + "." + port
	}
	return p

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
