package field

import (
	"reflect"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestNumericField(t *testing.T) {
	spec := &Spec{
		Length:      10,
		Description: "Field",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.Fixed,
		Pad:         padding.Left(' '),
	}
	numeric := NewNumeric(spec)

	numeric.SetBytes([]byte("100"))
	require.Equal(t, int64(100), numeric.Value())

	packed, err := numeric.Pack()
	require.NoError(t, err)
	require.Equal(t, "       100", string(packed))

	length, err := numeric.Unpack([]byte("      9876"))
	require.NoError(t, err)
	require.Equal(t, 10, length)

	b, err := numeric.Bytes()
	require.NoError(t, err)
	require.Equal(t, "9876", string(b))

	require.Equal(t, int64(9876), numeric.Value())

	numeric = NewNumeric(spec)
	numeric.Marshal(NewNumericValue(9876))
	packed, err = numeric.Pack()
	require.NoError(t, err)
	require.Equal(t, "      9876", string(packed))

	numeric = NewNumeric(spec)
	data := NewNumericValue(0)
	numeric.Marshal(data)
	length, err = numeric.Unpack([]byte("      9876"))
	require.NoError(t, err)
	require.Equal(t, 10, length)
	require.Equal(t, int64(9876), numeric.Value())

	numeric = NewNumeric(spec)
	numeric.SetValue(9876)

	require.Equal(t, int64(9876), numeric.Value())
}

func TestNumericNil(t *testing.T) {
	var str *Numeric = nil

	bs, err := str.Bytes()
	require.NoError(t, err)
	require.Nil(t, bs)

	value, err := str.String()
	require.NoError(t, err)
	require.Equal(t, "", value)

	n := str.Value()
	require.Equal(t, int64(0), n)
}

func TestNumericPack(t *testing.T) {
	t.Run("returns error for zero value when fixed length and no padding specified", func(t *testing.T) {
		spec := &Spec{
			Length:      10,
			Description: "Field",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}
		numeric := NewNumeric(spec)
		_, err := numeric.Pack()

		// zero value for Numeric is 0, so we have default field length 1
		require.EqualError(t, err, "failed to encode length: data length: 1 should be fixed: 10")
		require.True(t, prefix.IsLengthError(err), "error should be a length error")
	})
}

func TestNumericFieldUnmarshal(t *testing.T) {
	numericField := NewNumericValue(123456)

	vNumeric := &Numeric{}
	err := numericField.Unmarshal(vNumeric)
	require.NoError(t, err)
	require.Equal(t, int64(123456), vNumeric.Value())

	var s string
	err = numericField.Unmarshal(&s)
	require.NoError(t, err)
	require.Equal(t, "123456", s)

	var b int64
	err = numericField.Unmarshal(&b)
	require.NoError(t, err)
	require.Equal(t, int64(123456), b)

	refStrValue := reflect.ValueOf(&s).Elem()
	err = numericField.Unmarshal(refStrValue)
	require.NoError(t, err)
	require.Equal(t, "123456", refStrValue.String())

	refIntValue := reflect.ValueOf(&b).Elem()
	err = numericField.Unmarshal(refIntValue)
	require.NoError(t, err)
	require.Equal(t, int64(123456), int64(refIntValue.Int()))

	refStr := reflect.ValueOf(s)
	err = numericField.Unmarshal(refStr)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type string", err.Error())

	refStrPointer := reflect.ValueOf(&s)
	err = numericField.Unmarshal(refStrPointer)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type ptr", err.Error())

	err = numericField.Unmarshal(nil)
	require.Error(t, err)
	require.Equal(t, "unsupported type: expected *Numeric, *int, or reflect.Value, got <nil>", err.Error())
}

func TestNumericFieldMarshal(t *testing.T) {
	numericField := NewNumericValue(0)

	inputNumeric := NewNumericValue(123456)
	numericField.Marshal(inputNumeric)
	require.Equal(t, int64(123456), numericField.Value())

	numericField.Marshal(&inputNumeric)
	require.Equal(t, int64(123456), numericField.Value())

	inputStr := "123456"
	numericField.Marshal(inputStr)
	require.Equal(t, int64(123456), numericField.Value())

	numericField.Marshal(&inputStr)
	require.Equal(t, int64(123456), numericField.Value())

	var inputInt64 int64 = 123456
	numericField.Marshal(inputInt64)
	require.Equal(t, int64(123456), numericField.Value())

	numericField.Marshal(&inputInt64)
	require.Equal(t, int64(123456), numericField.Value())

	err := numericField.Marshal(nil)
	require.NoError(t, err)

	err = numericField.Marshal([]byte("123456"))
	require.Error(t, err)
	require.Equal(t, "data does not match require *Numeric or (int64, *int64, string, *string) type", err.Error())
}

func TestNumericFieldWithNotANumber(t *testing.T) {
	numeric := NewNumeric(&Spec{
		Length:      10,
		Description: "Field",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.Fixed,
		Pad:         padding.Left(' '),
	})

	err := numeric.SetBytes([]byte("hello"))
	require.Error(t, err)
	require.EqualError(t, err, "failed to convert into number")
	require.Equal(t, int64(0), numeric.Value())

	packed, err := numeric.Pack()
	require.NoError(t, err)
	require.Equal(t, "         0", string(packed))

	_, err = numeric.Unpack([]byte("hhhhhhhhhh"))
	require.Error(t, err)
	require.EqualError(t, err, "failed to set bytes: failed to convert into number")
}

func TestNumericFieldZeroLeftPaddedZero(t *testing.T) {
	numeric := NewNumeric(&Spec{
		Length:      4,
		Description: "Field",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.Fixed,
		Pad:         padding.Left('0'),
	})

	numeric.SetBytes([]byte("0"))
	require.Equal(t, int64(0), numeric.Value())

	packed, err := numeric.Pack()

	require.NoError(t, err)
	require.Equal(t, "0000", string(packed))

	length, err := numeric.Unpack([]byte("0000"))

	require.NoError(t, err)
	require.Equal(t, 4, length)

	bs, err := numeric.Bytes()
	require.NoError(t, err)
	require.Equal(t, "0", string(bs))

	require.Equal(t, int64(0), numeric.Value())
}

func TestNumericSetBytesSetsDataOntoDataStruct(t *testing.T) {
	numeric := NewNumeric(&Spec{
		Length:      1,
		Description: "Field",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.Fixed,
	})

	data := &Numeric{}
	err := numeric.Marshal(data)
	require.NoError(t, err)

	err = numeric.SetBytes([]byte("9"))
	require.NoError(t, err)

	require.Equal(t, int64(9), numeric.Value())
}

func TestNumericJSONMarshal(t *testing.T) {
	numeric := NewNumericValue(1)
	marshalledJSON, err := numeric.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "1", string(marshalledJSON))
}

func TestNumericJSONUnmarshal(t *testing.T) {
	input := []byte(`4000`)

	numeric := NewNumeric(&Spec{
		Length:      4,
		Description: "Field",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.Fixed,
	})

	require.NoError(t, numeric.UnmarshalJSON(input))
	require.Equal(t, int64(4000), numeric.Value())
}

func TestNumericJSONMarshalInt64(t *testing.T) {
	numeric := NewNumericValue(9223372036854775807)
	marshalledJSON, err := numeric.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "9223372036854775807", string(marshalledJSON))
}

func TestNumericJSONUnmarshalInt64(t *testing.T) {
	input := []byte(`9223372036854775807`)

	numeric := NewNumeric(&Spec{
		Length:      4,
		Description: "Field",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.Fixed,
	})

	require.NoError(t, numeric.UnmarshalJSON(input))
	require.Equal(t, int64(9223372036854775807), numeric.Value())
}
