package adapters

import (
	"fmt"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/specs"
)

// adapterFor binds each built-in spec that needs a wire-format fixup to its
// adapter. Add a new built-in adapter here in one line.
var adapterFor = map[*iso8583.MessageSpec]Adapter{
	specs.Spec87Hex: &Spec87HexAdapter{},
}

// AdapterForSpec returns the adapter registered for spec, or an error if
// none is registered. Lookup is by pointer identity, so spec must be a
// built-in singleton; a spec rebuilt via ImportJSON/ImportYAML is a different
// pointer. Pass the adapter directly to UnpackWithAdapter when you hold such a spec.
func AdapterForSpec(spec *iso8583.MessageSpec) (Adapter, error) {
	adapter, ok := adapterFor[spec]
	if !ok {
		return nil, fmt.Errorf("no adapter registered for spec %q", spec.Name)
	}
	return adapter, nil
}
