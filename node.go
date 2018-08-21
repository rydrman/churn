package churn

// Node is a graph component that can participate in the
// graph execution by exposing any number of input and output
// ports
type Node interface {

	// In returns the in port on this node with the given name,
	// or nil if no port exists with that name
	In(name string) *Port
	// Out returns the out port on this node with the given name,
	// or nil if no port exists with that name
	Out(name string) *Port

	// Init is called when this node is being initialized in a graph
	Init()

	setupBaseNode(node Node)
}

// BaseNode contains the core node logic that must
// be embeded into all node definitions
type BaseNode struct {
	BaseComponent
	PortCatalog
}

// Init can be overridden for custom node initialization
func (n *BaseNode) Init() {}

func (n *BaseNode) setupBaseNode(node Node) {

	n.PortCatalog = *CatalogPorts(node)

}
