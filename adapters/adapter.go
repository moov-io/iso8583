// Package adapters applies optional per-spec fixups around Unpack for
// specs whose pack/unpack is asymmetric, such as specs.Spec87Hex.
//
// A WirePrepper rewrites raw wire bytes before Unpack; a DisplayRenderer
// post-processes the unpacked message for readable output. UnpackWithAdapter
// runs both around a normal Unpack, or a plain Unpack when the adapter is nil.
//
// PrepWire is the inverse of a specific Pack, so apply an adapter only to wire
// this library produced for the matching spec. For wire of unknown provenance,
// pass a nil adapter (plain Unpack).
package adapters

import (
	"fmt"

	"github.com/moov-io/iso8583"
)

// WirePrepper rewrites raw wire bytes before Unpack.
type WirePrepper interface {
	PrepWire(spec *iso8583.MessageSpec, wire []byte) ([]byte, error)
}

// DisplayRenderer post-processes an unpacked message for readable output.
type DisplayRenderer interface {
	RenderForDisplay(msg *iso8583.Message) error
}

// Adapter combines both sides of the fixup.
type Adapter interface {
	WirePrepper
	DisplayRenderer
}

// UnpackWithAdapter unpacks wire against spec. A non-nil adapter preps the
// wire before Unpack and renders the message after; a nil adapter is a plain
// Unpack. The returned message may be non-nil alongside an error so callers
// can still describe partial results.
func UnpackWithAdapter(spec *iso8583.MessageSpec, wire []byte, adapter Adapter) (*iso8583.Message, error) {
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
