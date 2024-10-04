package field

import (
	"errors"
	"reflect"
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

	t.Run("packing and unpacking with variable length", func(t *testing.T) {
		spec := &Spec{
			Length:      5, // 5 bytes, 10 hex chars
			Description: "Field",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.LL,
		}

		f := NewHexValue("AABBCCDDEE")
		f.SetSpec(spec)

		packed, err := f.Pack()

		require.NoError(t, err)
		require.Equal(t, []byte{0x00, 0x05, 0xaa, 0xbb, 0xcc, 0xdd, 0xee}, packed)

		f = NewHex(spec)
		read, err := f.Unpack(packed)

		require.NoError(t, err)
		require.Equal(t, 7, read)
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

func TestHexFieldUnmarshal(t *testing.T) {
	testValue := []byte{0x12, 0x34, 0x56}
	hexField := NewHexValue("123456")

	vHex := &Hex{}
	err := hexField.Unmarshal(vHex)
	require.NoError(t, err)
	require.Equal(t, "123456", vHex.Value())
	buf, _ := vHex.Bytes()
	require.Equal(t, testValue, buf)

	var s string
	err = hexField.Unmarshal(&s)
	require.NoError(t, err)
	require.Equal(t, "123456", s)

	var b []byte
	err = hexField.Unmarshal(&b)
	require.NoError(t, err)
	require.Equal(t, testValue, b)

	refStrValue := reflect.ValueOf(&s).Elem()
	err = hexField.Unmarshal(refStrValue)
	require.NoError(t, err)
	require.Equal(t, "123456", refStrValue.String())

	refBytesValue := reflect.ValueOf(&b).Elem()
	err = hexField.Unmarshal(refBytesValue)
	require.NoError(t, err)
	require.Equal(t, testValue, refBytesValue.Bytes())

	refStr := reflect.ValueOf(s)
	err = hexField.Unmarshal(refStr)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type string", err.Error())

	refStrPointer := reflect.ValueOf(&s)
	err = hexField.Unmarshal(refStrPointer)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type ptr", err.Error())

	err = hexField.Unmarshal(nil)
	require.Error(t, err)
	require.Equal(t, "unsupported type: expected *Hex, *string, *[]byte, or reflect.Value, got <nil>", err.Error())
}

func TestHexFieldMarshal(t *testing.T) {
	testValue := []byte{0x12, 0x34, 0x56}
	hexField := NewHexValue("")

	inputStr := "123456"
	err := hexField.Marshal(inputStr)
	require.NoError(t, err)
	require.Equal(t, "123456", hexField.Value())
	buf, _ := hexField.Bytes()
	require.Equal(t, testValue, buf)

	err = hexField.Marshal(&inputStr)
	require.NoError(t, err)
	require.Equal(t, "123456", hexField.Value())
	buf, _ = hexField.Bytes()
	require.Equal(t, testValue, buf)

	err = hexField.Marshal(testValue)
	require.NoError(t, err)
	require.Equal(t, "123456", hexField.Value())
	buf, _ = hexField.Bytes()
	require.Equal(t, testValue, buf)

	err = hexField.Marshal(&testValue)
	require.NoError(t, err)
	require.Equal(t, "123456", hexField.Value())
	buf, _ = hexField.Bytes()
	require.Equal(t, testValue, buf)

	err = hexField.Marshal(nil)
	require.NoError(t, err)

	err = hexField.Marshal(123456)
	require.Error(t, err)
	require.Equal(t, "data does not match required *Hex or (string, *string, []byte, *[]byte) type", err.Error())
}
