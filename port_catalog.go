package churn

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/rydrman/churn/churncore"
)

// PortCatalog describes a collection of in and out ports
type PortCatalog struct {
	Ins  PortSlice
	Outs PortSlice
}

// CatalogPorts builds a record for all ports detected on the given node
func CatalogPorts(node Node) *PortCatalog {

	catalog := new(PortCatalog)
	nodeVal := reflect.ValueOf(node)
	catalog.catalogInPorts(nodeVal)
	catalog.catalogOutPorts(nodeVal)
	return catalog

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

func (c *PortCatalog) catalogInPorts(node reflect.Value) {

	nodeType := node.Type()
	for i := 0; i < nodeType.NumMethod(); i++ {

		meth := nodeType.Method(i)
		if !strings.HasPrefix(meth.Name, inPortNamePrefix) {
			continue
		}

		name := strings.TrimPrefix(meth.Name, inPortNamePrefix)

		// just the prefix alone is not enough of a name
		if name == "" {
			continue
		}

		// prefix must be followed by an uppercase letter
		// to be properly camel-cased
		if !unicode.IsUpper(rune(name[0])) {
			continue
		}

		core, err := churncore.NewReceiver(node.Method(i).Interface())
		if err != nil {
			// a returned error simply means that the method signature
			// was not in fact valid as an input port
			continue
		}

		c.Ins = append(c.Ins, &Port{
			Name: name,
			core: core,
		})

	}

}

func (c *PortCatalog) catalogOutPorts(node reflect.Value) {

	nodeType := node.Type()
	nodeKind := nodeType.Kind()
	for nodeKind == reflect.Ptr || nodeKind == reflect.Interface {
		node = node.Elem()
		nodeType = node.Type()
		nodeKind = nodeType.Kind()
	}

	for i := 0; i < nodeType.NumField(); i++ {

		field := nodeType.Field(i)
		if field.Anonymous {
			// TODO: support fields from embedded structs
			continue
		}
		if !strings.HasPrefix(field.Name, outPortNamePrefix) {
			continue
		}

		name := strings.TrimPrefix(field.Name, outPortNamePrefix)

		// just the prefix alone is not enough of a name
		if name == "" {
			continue
		}

		// prefix must be followed by an uppercase letter
		// to be properly camel-cased
		if !unicode.IsUpper(rune(name[0])) {
			continue
		}

		if field.Type.Kind() != reflect.Chan ||
			field.Type.ChanDir() == reflect.RecvDir {
			continue
		}

		// NOTE: any previous channel value is blindly
		// replaced and not closed
		ch := reflect.MakeChan(
			reflect.ChanOf(reflect.BothDir, field.Type.Elem()), 0,
		)
		node.Field(i).Set(ch)

		// a returned error signifies that this field did not
		// meet the final requirements, so we move on
		core, err := churncore.NewSender(ch.Interface())
		panicIfError(err) // should never happend

		c.Outs = append(c.Outs, &Port{
			Name: name,
			core: core,
		})

	}

}
