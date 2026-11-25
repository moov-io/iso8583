// Spec87Hex adapter. Spec87Hex pairs field.NewString with
// encoding.BytesToASCIIHex, a combination the encoder's own
// documentation warns against: pack writes the value double-hex-encoded
// (4 * spec.Length wire bytes for a fixed field), but unpack only
// consumes 2 * spec.Length. Variable fields round-trip correctly on
// their own.
//
// PrepSpec87HexWire walks the wire in bitmap order, shrinks each
// fixed-field value with one hex.Decode pass, and leaves variable
// fields untouched. RenderSpec87HexForDisplay re-encodes the shrunk
// fields so Describe prints the hex-string form the caller passed to
// Marshal. Spec87HexAdapter bundles both as a single Adapter.

package adapters

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/field"
)

// PrepSpec87HexWire adjusts a Spec87Hex pack output so it can be
// unpacked without offset shift.
//
// The spec must be structurally compatible with specs.Spec87Hex: fields
// implemented as *field.String with encoding.BytesToASCIIHex, using
// Hex.Fixed or Hex.L/LL/LLL prefixes. Other combinations return an
// error naming the offending field.
//
// Precondition: wire must be the double-hex-encoded output of this
// library's Spec87Hex Pack. This shrink is the inverse of that Pack;
// running it on wire a peer packed in the single-hex-encoded form will
// consume the wrong byte count and corrupt (or fail) the parse.
func PrepSpec87HexWire(spec *iso8583.MessageSpec, wire []byte) ([]byte, error) {
	var out bytes.Buffer
	pos := 0

	mti, ok := spec.Fields[0]
	if !ok {
		return nil, fmt.Errorf("field 0 not defined in spec")
	}
	consumed, err := handleField(mti, wire[pos:], &out)
	if err != nil {
		return nil, fmt.Errorf("field 0: %w", err)
	}
	pos += consumed

	bitmapField, ok := spec.Fields[1]
	if !ok {
		return nil, fmt.Errorf("field 1 not defined in spec")
	}
	bmpWireStart := pos
	consumed, err = handleField(bitmapField, wire[pos:], &out)
	if err != nil {
		return nil, fmt.Errorf("field 1: %w", err)
	}
	pos += consumed

	// Bitmap sits in out and wire in its on-wire (hex-encoded) form.
	// Decode once to get the raw bitmap bytes for walking set fields.
	bitmap, err := hex.DecodeString(string(wire[bmpWireStart : bmpWireStart+consumed]))
	if err != nil {
		return nil, fmt.Errorf("decoding bitmap: %w", err)
	}

	maxField := 64
	if len(bitmap) > 0 && bitmap[0]&0x80 != 0 {
		return nil, fmt.Errorf("secondary bitmap is set; adapter does not yet support fields 65-128")
	}

	for num := 2; num <= maxField; num++ {
		if !bitSet(bitmap, num) {
			continue
		}
		f, ok := spec.Fields[num]
		if !ok {
			return nil, fmt.Errorf("field %d set in bitmap but not defined in spec", num)
		}
		consumed, err := handleField(f, wire[pos:], &out)
		if err != nil {
			return nil, fmt.Errorf("field %d: %w", num, err)
		}
		pos += consumed
	}

	if pos != len(wire) {
		return nil, fmt.Errorf("%d unexpected trailing bytes after last field", len(wire)-pos)
	}
	return out.Bytes(), nil
}

func handleField(f field.Field, wire []byte, out *bytes.Buffer) (int, error) {
	fs := f.Spec()

	// Bitmap pack/unpack is symmetric (both move 2 * bitmapLength wire
	// bytes). No shrink needed — pass the on-wire bytes through so
	// Unpack reads them exactly as Pack wrote them.
	if _, isBitmap := f.(*field.Bitmap); isBitmap {
		length := fs.Length
		if length == 0 {
			length = defaultBitmapLength
		}
		onWire := 2 * length
		if onWire > len(wire) {
			return 0, fmt.Errorf("short wire: want %d bytes, have %d", onWire, len(wire))
		}
		out.Write(wire[:onWire])
		return onWire, nil
	}

	switch fs.Pref.Inspect() {
	case "Hex.Fixed":
		// String + Hex.Fixed is the asymmetric case: Pack wrote
		// 4 * Length wire bytes, Unpack expects 2 * Length. Shrink
		// with one hex.Decode pass.
		onWire := 4 * fs.Length
		if onWire > len(wire) {
			return 0, fmt.Errorf("short wire: want %d bytes, have %d", onWire, len(wire))
		}
		dec, err := hex.DecodeString(string(wire[:onWire]))
		if err != nil {
			return 0, fmt.Errorf("hex decode: %w", err)
		}
		out.Write(dec)
		return onWire, nil

	case "Hex.L", "Hex.LL", "Hex.LLL":
		valueLen, prefBytes, err := fs.Pref.DecodeLength(fs.Length, wire)
		if err != nil {
			return 0, fmt.Errorf("decode length: %w", err)
		}
		onWire := prefBytes + 2*valueLen
		if onWire > len(wire) {
			return 0, fmt.Errorf("short wire: want %d bytes, have %d", onWire, len(wire))
		}
		out.Write(wire[:onWire])
		return onWire, nil
	}

	return 0, fmt.Errorf("unsupported prefix %s", fs.Pref.Inspect())
}

// defaultBitmapLength mirrors field.defaultBitmapLength, which is
// unexported. ISO 8583 primary bitmap is 8 bytes.
const defaultBitmapLength = 8

// RenderSpec87HexForDisplay walks the message's String fields so
// Describe prints the original hex-string form the caller passed to
// Marshal instead of the raw bytes that Spec87Hex stores internally.
//
// Rule applied per field:
//
//   - Hex.Fixed string fields: the adapter shrank these from 4*Length to
//     2*Length wire bytes before Unpack, so stored bytes are the semantic
//     hex-decoded form. Hex-encode them for display.
//   - Hex.L/LL/LLL string fields: adapter passed these through unchanged,
//     so stored bytes equal the user's input bytes. Display as-is.
//
// Call this after a successful Unpack of a Spec87Hex message. It mutates
// the message in place.
func RenderSpec87HexForDisplay(msg *iso8583.Message) {
	for _, f := range msg.GetFields() {
		sf, ok := f.(*field.String)
		if !ok {
			continue
		}
		if sf.Spec().Pref.Inspect() != "Hex.Fixed" {
			continue
		}
		raw, err := sf.Bytes()
		if err != nil || len(raw) == 0 {
			continue
		}
		sf.SetValue(strings.ToUpper(hex.EncodeToString(raw)))
	}
}

// bitSet reports whether the bit for 1-based field number n is set in
// an ISO 8583 bitmap (bits are MSB-first within each byte).
func bitSet(bitmap []byte, n int) bool {
	idx := (n - 1) / 8
	bit := 7 - ((n - 1) % 8)
	if idx >= len(bitmap) {
		return false
	}
	return bitmap[idx]&(1<<bit) != 0
}

// Spec87HexAdapter adapts specs.Spec87Hex by preprocessing the wire
// before Unpack and re-rendering fixed-field values after Unpack for
// human-readable display. It holds no state; its methods delegate to
// PrepSpec87HexWire and RenderSpec87HexForDisplay, which are also usable
// directly.
type Spec87HexAdapter struct{}

// Compile-time assertion that Spec87HexAdapter satisfies the combined
// Adapter interface.
var _ Adapter = (*Spec87HexAdapter)(nil)

// PrepWire delegates to the free function PrepSpec87HexWire.
func (*Spec87HexAdapter) PrepWire(spec *iso8583.MessageSpec, wire []byte) ([]byte, error) {
	return PrepSpec87HexWire(spec, wire)
}

// RenderForDisplay delegates to the free function
// RenderSpec87HexForDisplay. It never returns an error today but the
// signature matches DisplayRenderer so future error paths can propagate.
func (*Spec87HexAdapter) RenderForDisplay(msg *iso8583.Message) error {
	RenderSpec87HexForDisplay(msg)
	return nil
}
