// Static registry of built-in adapters, kept in its own file so adding
// a new adapter is a one-line change in a well-known place.
//
// Lookup semantics: a spec absent from the map returns the zero value
// (nil), which UnpackWith treats as "no adapter — plain Unpack path."
// Specs without an asymmetric pack/unpack simply have no entry here;
// Spec87ASCII, for example, is not listed.
//
// To add support for a new spec:
//
//  1. Create a new file in this package (e.g. spec2003.go) defining
//     an adapter type that implements the Adapter interface
//     (both WirePrepper and DisplayRenderer, or one with a no-op for
//     the other).
//  2. Add a single line below pairing the spec with the adapter value.
//
// Downstream callers who want to use a custom adapter not listed here
// can still build their own pipeline by calling the exported free
// functions (PrepSpec87HexWire, RenderSpec87HexForDisplay, ...) or
// implementing their own WirePrepper/DisplayRenderer and invoking them
// around a plain iso8583.NewMessage + Unpack.

package adapters

import (
	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/specs"
)

// adapterFor maps each supported spec to its adapter. A missing entry
// means "no adapter — use plain Unpack."
var adapterFor = map[*iso8583.MessageSpec]Adapter{
	specs.Spec87Hex: &Spec87HexAdapter{},
}
