package field

import (
	"errors"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/utils"
	"github.com/stretchr/testify/require"
)

func TestHexField(t *testing.T) {
	spec := &Spec{
		Length:      5, // 5 bytes, 10 hex chars
		Description: "Field",
		Enc:         encoding.Binary,
		Pref:        prefix.ASCII.Fixed,
	}

	t.Run("packing", func(t *testing.T) {
		f := NewHexValue("AABBCCDDEE")
		f.SetSpec(spec)

		packed, err := f.Pack()

		require.NoError(t, err)
		require.Equal(t, []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee}, packed)
	})

	t.Run("unpacking", func(t *testing.T) {
		f := NewHex(spec)
		read, err := f.Unpack([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee})

		require.NoError(t, err)
		require.Equal(t, 5, read)
		require.Equal(t, "AABBCCDDEE", f.Value())
	})

	t.Run("marshaling", func(t *testing.T) {
		f := NewHexValue("AABBCCDDEE")
		f2 := &Hex{}

		f2.Marshal(f)

		require.Equal(t, f.Value(), f2.Value())
	})

	t.Run("unmarshaling", func(t *testing.T) {
		f := NewHexValue("AABBCCDDEE")
		f2 := &Hex{}

		f.Unmarshal(f2)

		require.Equal(t, f.Value(), f2.Value())
	})

	t.Run("JSON marshaling/unmarshaling", func(t *testing.T) {
		// when marshaling, we should get the hex string, not base64
		f := NewHexValue("AABBCCDDEE")
		f.SetSpec(spec)

		b, err := f.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, "\"AABBCCDDEE\"", string(b))

		var f2 Hex
		err = f2.UnmarshalJSON(b)
		require.NoError(t, err)
		require.Equal(t, f.Value(), f2.Value())
	})

	t.Run("methods", func(t *testing.T) {
		f := NewHex(spec)
		f.SetBytes([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee})

		require.Equal(t, "AABBCCDDEE", f.Value())

		str, err := f.String()
		require.NoError(t, err)
		require.Equal(t, "AABBCCDDEE", str)

		b, err := f.Bytes()
		require.NoError(t, err)
		require.Equal(t, []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee}, b)

		// SetValue
		f.SetValue("EEBBCCDDEE")
		require.Equal(t, "EEBBCCDDEE", f.Value())
	})

	t.Run("errors", func(t *testing.T) {
		f := NewHex(spec)
		f.SetBytes([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee})

		// invalid length
		f.SetValue("AABBCCDDE")

		_, err := f.Bytes()
		require.EqualError(t, err, "encoding/hex: odd length hex string")

		// invalid hex
		f.SetValue("AABBCCDDEG")
		_, err = f.Bytes()
		require.EqualError(t, err, "encoding/hex: invalid byte: U+0047 'G'")

		_, err = f.Pack()
		require.EqualError(t, err, "converting hex field into bytes")

		var e *utils.SafeError
		require.True(t, errors.As(err, &e))
		require.Equal(t, "converting hex field into bytes: encoding/hex: invalid byte: U+0047 'G'", e.UnsafeError())
	})
}

func TestHexNil(t *testing.T) {
	var f *Hex = nil

	bs, err := f.Bytes()
	require.NoError(t, err)
	require.Nil(t, bs)

	value, err := f.String()
	require.NoError(t, err)
	require.Equal(t, "", value)

	value = f.Value()
	require.Equal(t, "", value)
}

func TestHexPack(t *testing.T) {
	t.Run("returns error for zero value when fixed length and no padding specified", func(t *testing.T) {
		spec := &Spec{
			Length:      10,
			Description: "Field",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}
		str := NewHex(spec)
		_, err := str.Pack()

		require.EqualError(t, err, "failed to encode length: field length: 0 should be fixed: 10")
	})
}
