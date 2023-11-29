package field

import (
	"reflect"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestHexNumericFieldWithLengthAndLeftPadding(t *testing.T) {
	spec := &field.Spec{
		Length:      6,
		Description: "Field",
		Enc:         encoding.Binary,
		Pref:        prefix.BerTLV,
		Pad:         padding.Left(0),
	}

	t.Run("odd number of digits", func(t *testing.T) {
		hexNumeric := NewHexNumeric(spec)

		// set bytes
		err := hexNumeric.SetBytes([]byte{0x01, 0x23})
		require.NoError(t, err)
		require.Equal(t, int64(123), hexNumeric.Value())

		// pack
		packedBytes, err := hexNumeric.Pack()
		require.NoError(t, err)
		require.Equal(t, []byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x01, 0x23}, packedBytes)

		// unpack
		length, err := hexNumeric.Unpack([]byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x01, 0x23})
		require.NoError(t, err)
		require.Equal(t, 7, length)

		// get bytes
		bytes, err := hexNumeric.Bytes()
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x23}, bytes)

		// string
		str, err := hexNumeric.String()
		require.NoError(t, err)
		require.Equal(t, "123", str)

		// value
		require.Equal(t, int64(123), hexNumeric.Value())

		// set value
		hexNumeric.SetValue(int64(456))
		require.Equal(t, int64(456), hexNumeric.Value())
	})

	t.Run("even number of digits", func(t *testing.T) {
		hexNumeric := NewHexNumeric(spec)

		// set bytes
		err := hexNumeric.SetBytes([]byte{0x41, 0x23})
		require.NoError(t, err)
		require.Equal(t, int64(4123), hexNumeric.Value())

		// pack
		packedBytes, err := hexNumeric.Pack()
		require.NoError(t, err)
		require.Equal(t, []byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x41, 0x23}, packedBytes)

		// unpack
		length, err := hexNumeric.Unpack([]byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x41, 0x23})
		require.NoError(t, err)
		require.Equal(t, 7, length)

		// get bytes
		bytes, err := hexNumeric.Bytes()
		require.NoError(t, err)
		require.Equal(t, []byte{0x41, 0x23}, bytes)

		// string
		str, err := hexNumeric.String()
		require.NoError(t, err)
		require.Equal(t, "4123", str)

		// value
		require.Equal(t, int64(4123), hexNumeric.Value())

		// set value
		hexNumeric.SetValue(int64(4567))
		require.Equal(t, int64(4567), hexNumeric.Value())
	})

	t.Run("empty", func(t *testing.T) {
		hexNumeric := NewHexNumeric(spec)

		// set bytes
		err := hexNumeric.SetBytes([]byte{})
		require.NoError(t, err)
		require.Equal(t, int64(0), hexNumeric.Value())

		// pack
		packedBytes, err := hexNumeric.Pack()
		require.NoError(t, err)
		require.Equal(t, []byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, packedBytes)

		// unpack
		length, err := hexNumeric.Unpack([]byte{0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		require.NoError(t, err)
		require.Equal(t, 7, length)

		// get bytes
		bytes, err := hexNumeric.Bytes()
		require.NoError(t, err)
		require.Equal(t, []byte{0x00}, bytes)

		// string
		str, err := hexNumeric.String()
		require.NoError(t, err)
		require.Equal(t, "0", str)

		// value
		require.Equal(t, int64(0), hexNumeric.Value())
	})
}

func TestHexNumericFieldNoLength(t *testing.T) {
	spec := &field.Spec{
		Description: "Field",
		Enc:         encoding.Binary,
		Pref:        prefix.BerTLV,
		Pad:         padding.Left(0),
	}

	t.Run("odd number of digits", func(t *testing.T) {
		hexNumeric := NewHexNumeric(spec)

		// set bytes
		err := hexNumeric.SetBytes([]byte{0x00, 0x01, 0x23})
		require.NoError(t, err)
		require.Equal(t, int64(123), hexNumeric.Value())

		// pack
		packedBytes, err := hexNumeric.Pack()
		require.NoError(t, err)
		require.Equal(t, []byte{0x02, 0x01, 0x23}, packedBytes)

		// unpack
		length, err := hexNumeric.Unpack([]byte{0x02, 0x01, 0x23})
		require.NoError(t, err)
		require.Equal(t, 3, length)

		// get bytes
		bytes, err := hexNumeric.Bytes()
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x23}, bytes)

		// string
		str, err := hexNumeric.String()
		require.NoError(t, err)
		require.Equal(t, "123", str)

		// value
		require.Equal(t, int64(123), hexNumeric.Value())

		// set value
		hexNumeric.SetValue(int64(456))
		require.Equal(t, int64(456), hexNumeric.Value())
	})

	t.Run("even number of digits", func(t *testing.T) {
		hexNumeric := NewHexNumeric(spec)

		// set bytes
		err := hexNumeric.SetBytes([]byte{0x41, 0x23})
		require.NoError(t, err)
		require.Equal(t, int64(4123), hexNumeric.Value())

		// pack
		packedBytes, err := hexNumeric.Pack()
		require.NoError(t, err)
		require.Equal(t, []byte{0x02, 0x41, 0x23}, packedBytes)

		// unpack
		length, err := hexNumeric.Unpack([]byte{0x02, 0x41, 0x23})
		require.NoError(t, err)
		require.Equal(t, 3, length)

		// get bytes
		bytes, err := hexNumeric.Bytes()
		require.NoError(t, err)
		require.Equal(t, []byte{0x41, 0x23}, bytes)

		// string
		str, err := hexNumeric.String()
		require.NoError(t, err)
		require.Equal(t, "4123", str)

		// value
		require.Equal(t, int64(4123), hexNumeric.Value())

		// set value
		hexNumeric.SetValue(int64(4567))
		require.Equal(t, int64(4567), hexNumeric.Value())

	})

	t.Run("empty", func(t *testing.T) {
		hexNumeric := NewHexNumeric(spec)

		// set bytes
		err := hexNumeric.SetBytes([]byte{})
		require.NoError(t, err)
		require.Equal(t, int64(0), hexNumeric.Value())

		// pack
		packedBytes, err := hexNumeric.Pack()
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x00}, packedBytes)

		// unpack
		length, err := hexNumeric.Unpack([]byte{0x01, 0x00})
		require.NoError(t, err)
		require.Equal(t, 2, length)

		// get bytes
		bytes, err := hexNumeric.Bytes()
		require.NoError(t, err)
		require.Equal(t, []byte{0x00}, bytes)

		// string
		str, err := hexNumeric.String()
		require.NoError(t, err)
		require.Equal(t, "0", str)

		// value
		require.Equal(t, int64(0), hexNumeric.Value())
	})
}

func TestHexNumericPackErrors(t *testing.T) {
	t.Run("returns error for zero value when fixed length and no padding specified", func(t *testing.T) {
		spec := &field.Spec{
			Length:      10,
			Description: "Field",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}
		hexNumeric := NewHexNumeric(spec)
		_, err := hexNumeric.Pack()

		// zero value for HexNumeric is 0, so we have default field length 1
		require.ErrorContains(t, err, "failed to encode length")
	})
}

func TestHexNumericFieldUnmarshal(t *testing.T) {
	hexNumericField := NewHexNumericValue(123456)

	vHexNumeric := &HexNumeric{}
	err := hexNumericField.Unmarshal(vHexNumeric)
	require.NoError(t, err)
	require.Equal(t, int64(123456), vHexNumeric.Value())

	var s string
	err = hexNumericField.Unmarshal(&s)
	require.NoError(t, err)
	require.Equal(t, "123456", s)

	var n int64
	err = hexNumericField.Unmarshal(&n)
	require.NoError(t, err)
	require.Equal(t, int64(123456), n)

	refStrValue := reflect.ValueOf(&s).Elem()
	err = hexNumericField.Unmarshal(refStrValue)
	require.NoError(t, err)
	require.Equal(t, "123456", refStrValue.String())

	refIntValue := reflect.ValueOf(&n).Elem()
	err = hexNumericField.Unmarshal(refIntValue)
	require.NoError(t, err)
	require.Equal(t, int64(123456), int64(refIntValue.Int()))

	refStr := reflect.ValueOf(s)
	err = hexNumericField.Unmarshal(refStr)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type string", err.Error())

	refStrPointer := reflect.ValueOf(&s)
	err = hexNumericField.Unmarshal(refStrPointer)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type ptr", err.Error())

	err = hexNumericField.Unmarshal(nil)
	require.Error(t, err)
	require.Equal(t, "unsupported type: expected *HexNumeric, *string, *int, or reflect.Value, got <nil>", err.Error())
}

func TestHexNumericFieldMarshal(t *testing.T) {
	hexNumericField := NewHexNumericValue(0)

	inputHexNumeric := NewHexNumericValue(123456)
	err := hexNumericField.Marshal(inputHexNumeric)
	require.NoError(t, err)
	require.Equal(t, int64(123456), hexNumericField.Value())

	inputStr := "123456"
	err = hexNumericField.Marshal(inputStr)
	require.NoError(t, err)
	require.Equal(t, int64(123456), hexNumericField.Value())

	err = hexNumericField.Marshal(&inputStr)
	require.NoError(t, err)
	require.Equal(t, int64(123456), hexNumericField.Value())

	var inputInt64 int64 = 123456
	err = hexNumericField.Marshal(inputInt64)
	require.NoError(t, err)
	require.Equal(t, int64(123456), hexNumericField.Value())

	err = hexNumericField.Marshal(&inputInt64)
	require.NoError(t, err)
	require.Equal(t, int64(123456), hexNumericField.Value())

	err = hexNumericField.Marshal(nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), hexNumericField.Value())

	err = hexNumericField.Marshal([]byte("123456"))
	require.Error(t, err)
	require.Equal(t, "data does not match required *HexNumeric or (int, *int, string, *string) type", err.Error())
}

func TestHexNumericJSONMarshal(t *testing.T) {
	hexNumeric := NewHexNumericValue(4321)
	marshalledJSON, err := hexNumeric.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "4321", string(marshalledJSON))
}

func TestHexNumericJSONUnmarshal(t *testing.T) {
	input := []byte(`4321`)

	hexNumeric := NewHexNumeric(&field.Spec{
		Length:      4,
		Description: "Field",
		Enc:         encoding.Binary,
		Pref:        prefix.BerTLV,
	})

	require.NoError(t, hexNumeric.UnmarshalJSON(input))
	require.Equal(t, int64(4321), hexNumeric.Value())
}
