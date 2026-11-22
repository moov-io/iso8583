package adapters

import (
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/specs"
	"github.com/stretchr/testify/require"
)

// Spec87Hex wire, prepped and rendered through the adapter, round-trips.
func TestUnpackWithAdapter_Spec87Hex(t *testing.T) {
	wire := packSpec87Hex(t, sampleFixture())

	msg, err := UnpackWithAdapter(specs.Spec87Hex, wire, &Spec87HexAdapter{})
	require.NoError(t, err)

	mti, err := msg.GetMTI()
	require.NoError(t, err)
	require.Equal(t, "00000100", mti)

	pan, err := msg.GetString(2)
	require.NoError(t, err)
	require.Equal(t, "4242424242424242", pan)
}

// A nil adapter is a plain Unpack: no prep, no render.
func TestUnpackWithAdapter_NilIsPlainUnpack(t *testing.T) {
	type ascii struct {
		MTI string `iso8583:"0"`
		PAN string `iso8583:"2"`
	}
	src := iso8583.NewMessage(specs.Spec87ASCII)
	require.NoError(t, src.Marshal(&ascii{MTI: "0100", PAN: "4242424242424242"}))
	wire, err := src.Pack()
	require.NoError(t, err)

	msg, err := UnpackWithAdapter(specs.Spec87ASCII, wire, nil)
	require.NoError(t, err)

	mti, _ := msg.GetMTI()
	require.Equal(t, "0100", mti)
}

// The adapter is bound by argument, not by spec pointer, so a spec rebuilt
// into a different pointer still round-trips.
func TestUnpackWithAdapter_RebuiltSpec(t *testing.T) {
	wire := packSpec87Hex(t, sampleFixture())

	rebuilt := *specs.Spec87Hex
	msg, err := UnpackWithAdapter(&rebuilt, wire, &Spec87HexAdapter{})
	require.NoError(t, err)

	mti, err := msg.GetMTI()
	require.NoError(t, err)
	require.Equal(t, "00000100", mti)
}

// AdapterForSpec resolves a built-in spec and errors on an unregistered one.
func TestAdapterForSpec(t *testing.T) {
	a, err := AdapterForSpec(specs.Spec87Hex)
	require.NoError(t, err)
	require.NotNil(t, a)

	_, err = AdapterForSpec(specs.Spec87ASCII)
	require.Error(t, err)
}

// Variable-length fields (Hex.LL and Hex.LLL) pass through the prep unchanged
// and round-trip across different lengths.
func TestUnpackWithAdapter_VariableLengthFields(t *testing.T) {
	type msg struct {
		MTI    string `iso8583:"0"`
		PAN    string `iso8583:"2"`  // Hex.LL
		Track3 string `iso8583:"36"` // Hex.LLL
	}

	cases := []struct {
		name   string
		pan    string
		track3 string
	}{
		{"short", "123", "4567"},
		{"longer", "1234567890123456", "0123456789abcdef0123"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := iso8583.NewMessage(specs.Spec87Hex)
			require.NoError(t, src.Marshal(&msg{MTI: "00000100", PAN: tc.pan, Track3: tc.track3}))
			wire, err := src.Pack()
			require.NoError(t, err)

			got, err := UnpackWithAdapter(specs.Spec87Hex, wire, &Spec87HexAdapter{})
			require.NoError(t, err)

			pan, err := got.GetString(2)
			require.NoError(t, err)
			require.Equal(t, tc.pan, pan)

			track3, err := got.GetString(36)
			require.NoError(t, err)
			require.Equal(t, tc.track3, track3)
		})
	}
}

// Trailing bytes after the last field are reported, not silently dropped.
func TestUnpackWithAdapter_TrailingBytes(t *testing.T) {
	wire := packSpec87Hex(t, sampleFixture())
	wire = append(wire, '0', '0')

	_, err := UnpackWithAdapter(specs.Spec87Hex, wire, &Spec87HexAdapter{})
	require.ErrorContains(t, err, "trailing bytes")
}

// A prep failure is wrapped with "prep wire" context.
func TestUnpackWithAdapter_PrepError(t *testing.T) {
	// Valid MTI segment + bitmap flagging an unsupported secondary bitmap.
	wire := []byte("30303030303030308000000000000000")

	_, err := UnpackWithAdapter(specs.Spec87Hex, wire, &Spec87HexAdapter{})
	require.ErrorContains(t, err, "prep wire")
}
