package kubernetes

import (
	_ "embed"
	"fmt"
)

//go:embed schema/schema.json
var schema string

func (m *Mixin) PrintSchema() error {
	fmt.Fprintf(m.Out, schema)
	return nil
}
