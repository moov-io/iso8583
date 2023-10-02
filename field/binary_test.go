package field

import (
	"reflect"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestBinaryField(t *testing.T) {
	spec := &Spec{
		Length:      10,
		Description: "Field",
		Enc:         encoding.Binary,
		Pref:        prefix.Binary.Fixed,
	}

	in := []byte("1234567890")

	t.Run("Pack returns binary data", func(t *testing.T) {
		bin := NewBinary(spec)
		bin.SetBytes(in)

		packed, err := bin.Pack()

		require.NoError(t, err)
		require.Equal(t, in, packed)
	})

	t.Run("String returns binary data encoded in HEX", func(t *testing.T) {
		bin := NewBinary(spec)
		bin.value = in

		str, err := bin.String()

		require.NoError(t, err)
		require.Equal(t, "31323334353637383930", str)
	})

	t.Run("Unpack returns binary data", func(t *testing.T) {
		bin := NewBinary(spec)

		n, err := bin.Unpack(in)

		require.NoError(t, err)
		require.Equal(t, len(in), n)
		require.Equal(t, in, bin.value)
	})

	t.Run("Marshal sets data to the field", func(t *testing.T) {
		bin := NewBinary(spec)
		bin.Marshal(NewBinaryValue(in))

		packed, err := bin.Pack()

		require.NoError(t, err)
		require.Equal(t, in, packed)
	})

	t.Run("Unmarshal gets data from the field", func(t *testing.T) {
		bin := NewBinaryValue([]byte{1, 2, 3})
		val := &Binary{}

		err := bin.Unmarshal(val)

		require.NoError(t, err)
		require.Equal(t, []byte{1, 2, 3}, val.value)
	})

	t.Run("SetBytes sets data to the data field", func(t *testing.T) {
		bin := NewBinary(spec)
		data := &Binary{}
		bin.Marshal(data)

		err := bin.SetBytes(in)
		require.NoError(t, err)

		require.Equal(t, in, bin.value)
	})

	// SetValue sets data to the data field
	t.Run("SetValue sets data to the data field", func(t *testing.T) {
		bin := NewBinary(spec)
		bin.SetValue(in)

		require.Equal(t, in, bin.Value())
	})

	t.Run("Unpack sets data to data value", func(t *testing.T) {
		bin := NewBinary(spec)
		data := NewBinaryValue([]byte{})
		bin.Marshal(data)

		n, err := bin.Unpack(in)

		require.NoError(t, err)
		require.Equal(t, len(in), n)
		require.Equal(t, in, bin.value)
	})

	t.Run("UnmarshalJSON unquotes input before handling it", func(t *testing.T) {
		input := []byte(`"500000000000000000000000000000000000000000000000"`)

		bin := NewBinary(spec)
		require.NoError(t, bin.UnmarshalJSON(input))

		str, err := bin.String()
		require.NoError(t, err)

		require.Equal(t, `500000000000000000000000000000000000000000000000`, str)
	})

	t.Run("MarshalJSON returns string hex repr of binary field", func(t *testing.T) {
		bin := NewBinaryValue([]byte{0xAB})
		marshalledJSON, err := bin.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `"AB"`, string(marshalledJSON))
	})

	t.Run("returns error for zero value when fixed length and no padding specified", func(t *testing.T) {
		bin := NewBinary(spec)
		_, err := bin.Pack()

		require.EqualError(t, err, "failed to encode length: field length: 0 should be fixed: 10")
	})
}

func TestBinaryNil(t *testing.T) {
	var str *Binary = nil

	bs, err := str.Bytes()
	require.NoError(t, err)
	require.Nil(t, bs)

	value, err := str.String()
	require.NoError(t, err)
	require.Equal(t, "", value)

	bs = str.Value()
	require.Nil(t, bs)
}

func TestBinaryFieldUnmarshal(t *testing.T) {
	testValue := []byte{0x12, 0x34, 0x56}
	str := NewBinaryValue(testValue)

	val1 := &Binary{}
	err := str.Unmarshal(val1)
	require.NoError(t, err)
	require.Equal(t, testValue, val1.Value())

	var val2 string
	err = str.Unmarshal(&val2)
	require.NoError(t, err)
	require.Equal(t, "123456", val2)

	var val3 []byte
	err = str.Unmarshal(&val3)
	require.NoError(t, err)
	require.Equal(t, testValue, val3)

	val4 := reflect.ValueOf(&val2).Elem()
	err = str.Unmarshal(val4)
	require.NoError(t, err)
	require.Equal(t, "123456", val4.String())

	val5 := reflect.ValueOf(&val3).Elem()
	err = str.Unmarshal(val5)
	require.NoError(t, err)
	require.Equal(t, testValue, val5.Bytes())

	val6 := reflect.ValueOf(val2)
	err = str.Unmarshal(val6)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type string", err.Error())

	val7 := reflect.ValueOf(&val2)
	err = str.Unmarshal(val7)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type ptr", err.Error())

	err = str.Unmarshal(nil)
	require.Error(t, err)
	require.Equal(t, "unsupported type: expected *Binary, *string, *[]byte, or reflect.Value, got <nil>", err.Error())
}

func TestBinaryFieldMarshal(t *testing.T) {
	testValue := []byte{0x12, 0x34, 0x56}
	str := NewBinaryValue(nil)

	vstring := "123456"
	err := str.Marshal(vstring)
	require.NoError(t, err)
	require.Equal(t, testValue, str.Value())

	err = str.Marshal(&vstring)
	require.NoError(t, err)
	require.Equal(t, testValue, str.Value())

	err = str.Marshal(testValue)
	require.NoError(t, err)
	require.Equal(t, testValue, str.Value())

	err = str.Marshal(&testValue)
	require.NoError(t, err)
	require.Equal(t, testValue, str.Value())

	err = str.Marshal(nil)
	require.Error(t, err)
	require.Equal(t, "data does not match required *Binary or (string, *string, []byte, *[]byte) type", err.Error())
}
