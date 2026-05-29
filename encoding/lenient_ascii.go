package encoding

import (
	"fmt"
)

var (
	_            Encoder = (*lenientASCIIEncoder)(nil)
	LenientASCII         = &lenientASCIIEncoder{}
)

// LenientASCII is like ASCII but passes bytes > 0x7F through without raising
// an encoder error. The ISO 8583 spec defines text fields as 7-bit ASCII;
// the strict ASCII encoder enforces that contract. In practice, however,
// many real-world counterparties occasionally inject non-ASCII bytes into
// text fields (legacy peers using cp1252, copy-paste artifacts from Word or
// Excel, uninitialized buffers, or test harnesses that intentionally emit
// invalid payloads to verify the receiver's format-error path). When that
// happens, the strict ASCII encoder fails Unpack and the connection layer
// silently drops the whole message — preventing the application from
// replying with a proper response (e.g. ISO 8583 RC 30 Format Error) that
// echoes the original STAN/RRN.
//
// LenientASCII lets the message reach the application handler so it can
// inspect the offending field, decide how to respond, and reply with the
// echoed fields intact. Content validation (rejecting non-ASCII, control
// characters, or other domain-specific rules) is the caller's responsibility
// and should be performed in the handler, not the encoder.
//
// Use LenientASCII for fields where you want robust receive behavior over
// strict spec compliance. Continue using ASCII for fields where any
// non-ASCII byte should hard-fail at the encoder layer.
type lenientASCIIEncoder struct{}

func (e lenientASCIIEncoder) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

func (e lenientASCIIEncoder) Decode(data []byte, length int) ([]byte, int, error) {
	if length < 0 {
		return nil, 0, fmt.Errorf("invalid length: %d", length)
	}

	if length == 0 {
		return nil, 0, nil
	}

	if len(data) < length {
		return nil, 0, fmt.Errorf("not enough data to decode. expected len %d, got %d", length, len(data))
	}

	out := make([]byte, length)
	copy(out, data[:length])
	return out, length, nil
}
