package churn

import (
	"reflect"
	"strings"
)

// Node is a graph component that can participate in the
// graph execution by exposing any number of input and output
// ports
type Node interface {
	Ports() PortCatalog
	initialize(self reflect.Value)
}

// BaseNode contains the core node logic that must
// be embeded into all node definitions
type BaseNode struct {
	BaseComponent
	PortCatalog
}

func (n *BaseNode) initialize(self reflect.Value) {

	n.PortCatalog = PortCatalog{}
	n.catalogInPorts(self)
	n.catalogOutPorts(self)

}

func (n *BaseNode) catalogInPorts(self reflect.Value) {

	selfType := self.Type()
	for i := 0; i < selfType.NumMethod(); i++ {

		meth := selfType.Method(i)
		name := meth.Name
		if !strings.HasPrefix(name, inPortNamePrefix) {
			continue
		}

		n.Ins = append(n.Ins, &Port{
			Name: strings.TrimPrefix(name, inPortNamePrefix),
			core: NewInPortCore(self.Method(i)),
		})

	}

}

func (n *BaseNode) catalogOutPorts(self reflect.Value) {

	selfType := self.Type()
	selfKind := selfType.Kind()
	for selfKind == reflect.Ptr || selfKind == reflect.Interface {
		self = self.Elem()
		selfType = self.Type()
		selfKind = selfType.Kind()
	}

	for i := 0; i < selfType.NumField(); i++ {

		field := selfType.Field(i)
		if field.Anonymous {
			// TODO: support fields from embedded structs
			continue
		}
		if !strings.HasPrefix(field.Name, outPortNamePrefix) {
			continue
		}

		// a returned error signifies that this field did not
		// meet the final requirements, so we move on
		core, err := NewOutPortCore(field)
		if err != nil {
			continue
		}

		name := strings.TrimPrefix(field.Name, outPortNamePrefix)
		n.Outs = append(n.Outs, &Port{
			Name: name,
			core: core,
		})

		// NOTE: any previous channel value is blindly
		// replaces and not closed
		self.FieldByName(field.Name).Set(core.channel)

	}

}

// PortCatalog describes a collection of in and out ports
type PortCatalog struct {
	Ins  PortSlice
	Outs PortSlice
}

// In returns the first in port found in this catalog
// with the given name, or nil
func (c PortCatalog) In(name string) *Port {
	return c.Ins.FindByName(name)
}

// Out returns the first in port found in this catalog
// with the given name, or nil
func (c PortCatalog) Out(name string) *Port {
	return c.Outs.FindByName(name)
}

func (c PortCatalog) copy() PortCatalog {

	dupe := PortCatalog{
		Ins:  make([]*Port, len(c.Ins)),
		Outs: make([]*Port, len(c.Outs)),
	}
	copy(dupe.Ins, c.Ins)
	copy(dupe.Outs, c.Outs)
	return dupe

}
