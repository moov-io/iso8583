package adapters

import (
	"strings"
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/specs"
	"github.com/stretchr/testify/require"
)

// spec87HexFixture is the canonical set of Spec87Hex field values. Each
// fixed field is a hex string of length 2*Spec.Length (the double-encode
// input form); F2 is Hex.LL variable so it is passed through unchanged.
type spec87HexFixture struct {
	MTI                  string `iso8583:"0"`
	PrimaryAccountNumber string `iso8583:"2"`
	ProcessingCode       string `iso8583:"3"`
	Amount               string `iso8583:"4"`
	TransmissionDateTime string `iso8583:"7"`
	STAN                 string `iso8583:"11"`
	ExpirationDate       string `iso8583:"14"`
}

func sampleFixture() *spec87HexFixture {
	return &spec87HexFixture{
		MTI:                  "00000100",
		PrimaryAccountNumber: "4242424242424242",
		ProcessingCode:       "000000000000",
		Amount:               "000000000000000000002599",
		TransmissionDateTime: "00000000001109120000",
		STAN:                 "000000000123",
		ExpirationDate:       "00002712",
	}
}

// packSpec87Hex builds a Spec87Hex message from data and returns the raw
// (asymmetric, pre-prep) wire bytes that msg.Pack produces.
func packSpec87Hex(t *testing.T, data any) []byte {
	t.Helper()
	msg := iso8583.NewMessage(specs.Spec87Hex)
	require.NoError(t, msg.Marshal(data))
	wire, err := msg.Pack()
	require.NoError(t, err)
	return wire
}

func TestPrepSpec87HexWire_RoundTrip(t *testing.T) {
	wire := packSpec87Hex(t, sampleFixture())

	prepped, err := PrepSpec87HexWire(specs.Spec87Hex, wire)
	require.NoError(t, err)
	// Fixed fields shrink 4*Length -> 2*Length, so the prepped wire must
	// be strictly shorter than what Pack emitted.
	require.Less(t, len(prepped), len(wire))

	// Prepped wire must now unpack without offset shift.
	msg := iso8583.NewMessage(specs.Spec87Hex)
	require.NoError(t, msg.Unpack(prepped))

	// Variable field (Hex.LL) passes through untouched.
	pan, err := msg.GetString(2)
	require.NoError(t, err)
	require.Equal(t, "4242424242424242", pan)

	// Render makes fixed fields human-readable again.
	RenderSpec87HexForDisplay(msg)
	mti, err := msg.GetMTI()
	require.NoError(t, err)
	require.Equal(t, "00000100", mti)
}

func TestPrepSpec87HexWire_Errors(t *testing.T) {
	t.Run("secondary bitmap rejected", func(t *testing.T) {
		// MTI (Hex.Fixed len 4 -> 16 wire chars) + bitmap (16 wire chars)
		// whose first decoded byte has the high bit set.
		wire := []byte(strings.Repeat("30", 8) + "8000000000000000")
		_, err := PrepSpec87HexWire(specs.Spec87Hex, wire)
		require.ErrorContains(t, err, "secondary bitmap is set")
	})

	t.Run("short wire", func(t *testing.T) {
		wire := packSpec87Hex(t, sampleFixture())
		_, err := PrepSpec87HexWire(specs.Spec87Hex, wire[:len(wire)-4])
		require.ErrorContains(t, err, "short wire")
	})
}

func TestRenderSpec87HexForDisplay(t *testing.T) {
	// MTI carries hex letters so we can assert uppercasing; PAN is a
	// variable field and must be left as-is.
	data := sampleFixture()
	data.MTI = "000001ab"

	wire := packSpec87Hex(t, data)
	prepped, err := PrepSpec87HexWire(specs.Spec87Hex, wire)
	require.NoError(t, err)

	msg := iso8583.NewMessage(specs.Spec87Hex)
	require.NoError(t, msg.Unpack(prepped))

	RenderSpec87HexForDisplay(msg)

	mti, err := msg.GetMTI()
	require.NoError(t, err)
	require.Equal(t, "000001AB", mti) // uppercased hex

	pan, err := msg.GetString(2)
	require.NoError(t, err)
	require.Equal(t, "4242424242424242", pan) // variable field untouched
}

func TestBitSet(t *testing.T) {
	// 0x80 = field 1 set; 0x01 = field 8 set.
	bitmap := []byte{0x80, 0x01}
	tests := []struct {
		n    int
		want bool
	}{
		{1, true},
		{2, false},
		{8, false},
		{16, true},
		{17, false}, // out of range -> false
	}
	for _, tt := range tests {
		require.Equalf(t, tt.want, bitSet(bitmap, tt.n), "field %d", tt.n)
	}
}
