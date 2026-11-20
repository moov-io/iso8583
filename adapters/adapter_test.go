package adapters

import (
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/specs"
	"github.com/stretchr/testify/require"
)

// TestUnpackWith_Spec87Hex exercises the full pipeline through the real
// static registry: a registered spec gets its wire prepped and its fixed
// fields re-rendered for display.
func TestUnpackWith_Spec87Hex(t *testing.T) {
	wire := packSpec87Hex(t, sampleFixture())

	msg, err := UnpackWith(specs.Spec87Hex, wire)
	require.NoError(t, err)
	require.NotNil(t, msg)

	mti, err := msg.GetMTI()
	require.NoError(t, err)
	require.Equal(t, "00000100", mti)

	pan, err := msg.GetString(2)
	require.NoError(t, err)
	require.Equal(t, "4242424242424242", pan)
}

// TestUnpackWith_NoAdapterPassthrough: a spec absent from the registry
// must behave exactly like a plain NewMessage + Unpack.
func TestUnpackWith_NoAdapterPassthrough(t *testing.T) {
	require.Nil(t, adapterFor[specs.Spec87ASCII], "precondition: Spec87ASCII must not be registered")

	type ascii struct {
		MTI string `iso8583:"0"`
		PAN string `iso8583:"2"`
	}
	src := iso8583.NewMessage(specs.Spec87ASCII)
	require.NoError(t, src.Marshal(&ascii{MTI: "0100", PAN: "4242424242424242"}))
	wire, err := src.Pack()
	require.NoError(t, err)

	// Reference: plain unpack.
	want := iso8583.NewMessage(specs.Spec87ASCII)
	require.NoError(t, want.Unpack(wire))

	got, err := UnpackWith(specs.Spec87ASCII, wire)
	require.NoError(t, err)

	wantMTI, _ := want.GetMTI()
	gotMTI, _ := got.GetMTI()
	require.Equal(t, wantMTI, gotMTI)

	wantPAN, _ := want.GetString(2)
	gotPAN, _ := got.GetString(2)
	require.Equal(t, wantPAN, gotPAN)
}

// TestUnpackWith_PrepError: a registered spec fed malformed wire must
// surface the prep failure wrapped with the "prep wire" context.
func TestUnpackWith_PrepError(t *testing.T) {
	// Valid MTI segment + bitmap flagging a secondary bitmap the adapter
	// does not support.
	wire := []byte("30303030303030308000000000000000")
	_, err := UnpackWith(specs.Spec87Hex, wire)
	require.ErrorContains(t, err, "prep wire")
}

// TestRegistry guards the static registry contract: every entry is a
// usable Adapter, and the known spec is present.
func TestRegistry(t *testing.T) {
	require.NotEmpty(t, adapterFor)
	for spec, a := range adapterFor {
		require.NotNilf(t, a, "adapter for %q is nil", spec.Name)
	}
	require.NotNil(t, adapterFor[specs.Spec87Hex])
}
