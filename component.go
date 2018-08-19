package churn

// Component is an element that can exist within the graph
type Component interface {
	component()
}

// BaseComponent contains the core component logic and
// should be embeded in all component definitions
type BaseComponent struct{}

func (*BaseComponent) component() {}
