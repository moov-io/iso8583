// Package adapters provides optional wire-format preprocessors and
// display post-processors for specs whose built-in pack/unpack is
// asymmetric (for example specs.Spec87Hex, where fixed-length string
// fields are double-encoded on the wire).
//
// Design shape:
//
//   - WirePrepper transforms raw wire bytes before Unpack.
//   - DisplayRenderer post-processes an unpacked message before
//     Describe.
//   - Adapter combines both. Package-ships adapters live in registry.go
//     as a static map keyed by *iso8583.MessageSpec.
//   - UnpackWith reads the static map; no setup, no mutation, no
//     registry argument to thread through call sites.
//
// Typical use:
//
//	msg, err := adapters.UnpackWith(spec, wire)
//
// If spec has an entry in the static map, the registered adapter runs
// around the normal Unpack. If not, UnpackWith reduces to
// iso8583.NewMessage + Unpack with no behavioral change.
package adapters

import (
	"fmt"

	"github.com/moov-io/iso8583"
)

// WirePrepper preprocesses raw wire bytes before Unpack. Returning an
// error aborts the unpack.
type WirePrepper interface {
	PrepWire(spec *iso8583.MessageSpec, wire []byte) ([]byte, error)
}

// DisplayRenderer post-processes an unpacked message so Describe (or
// any string-rendering caller) produces human-readable output.
// Implementations mutate msg in place.
type DisplayRenderer interface {
	RenderForDisplay(msg *iso8583.Message) error
}

// Adapter is the combined interface. The static registry in
// registry.go stores values of this type, so every built-in adapter
// must implement both methods. An adapter that logically needs only
// one side can stub the other with a no-op.
type Adapter interface {
	WirePrepper
	DisplayRenderer
}

// UnpackWith applies the adapter registered for spec (if any) around
// a normal Unpack. Looks the adapter up in the static registry
// defined in registry.go.
//
// The returned message may be non-nil even when err is non-nil (for
// example when Unpack fails after preprocessing) so callers can still
// describe partial results.
func UnpackWith(spec *iso8583.MessageSpec, wire []byte) (*iso8583.Message, error) {
	adapter := adapterFor[spec]

	if adapter != nil {
		prepped, err := adapter.PrepWire(spec, wire)
		if err != nil {
			return nil, fmt.Errorf("prep wire: %w", err)
		}
		wire = prepped
	}

	msg := iso8583.NewMessage(spec)
	if err := msg.Unpack(wire); err != nil {
		return msg, fmt.Errorf("unpacking message: %w", err)
	}

	if adapter != nil {
		if err := adapter.RenderForDisplay(msg); err != nil {
			return msg, fmt.Errorf("render for display: %w", err)
		}
	}

	return msg, nil
}
