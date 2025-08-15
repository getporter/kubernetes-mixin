package kubernetes

import (
	_ "embed"
	"fmt"
)

//go:embed schema/schema.json
var schema string

func (m *Mixin) PrintSchema() error {
	fmt.Fprint(m.Out, schema)
	return nil
}
