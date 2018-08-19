package churn

// GraphOption is a type that applies one or more initialization
// options to a new graph instance
type GraphOption interface {
	Apply(*Graph)
}

// OptionFunc is a function that can be given as a graph option
type OptionFunc func(*Graph)

// Apply calls the underlying option function for g
func (f OptionFunc) Apply(g *Graph) { f(g) }

// ChannelBufferSize sets the buffer size for out port channels
// created in a graph network
func ChannelBufferSize(size int) GraphOption {
	return OptionFunc(func(g *Graph) {
		g.channelBufferSize = size
	})
}
