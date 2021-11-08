package kubernetes

import (
	_ "embed"
	"fmt"
)

//go:embed schema/schema.json
var schema string

func (m *Mixin) PrintSchema() error {
	schema := m.GetSchema()
	fmt.Fprintf(m.Out, schema)
	return nil
}

func (m *Mixin) GetSchema() string {
	return schema
}
