package field

import (
	"reflect"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestStringField(t *testing.T) {
	spec := &Spec{
		Length:      10,
		Description: "Field",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.Fixed,
		Pad:         padding.Left(' '),
	}
	str := NewString(spec)

	str.SetBytes([]byte("hello"))
	require.Equal(t, "hello", str.Value())

	packed, err := str.Pack()
	require.NoError(t, err)
	require.Equal(t, "     hello", string(packed))

	length, err := str.Unpack([]byte("     olleh"))
	require.NoError(t, err)
	require.Equal(t, 10, length)

	b, err := str.Bytes()
	require.NoError(t, err)
	require.Equal(t, "olleh", string(b))

	require.Equal(t, "olleh", str.Value())

	str = NewString(spec)
	str.Marshal(NewStringValue("hello"))
	packed, err = str.Pack()
	require.NoError(t, err)
	require.Equal(t, "     hello", string(packed))

	str = NewString(spec)
	length, err = str.Unpack([]byte("     olleh"))
	require.NoError(t, err)
	require.Equal(t, 10, length)
	require.Equal(t, "olleh", str.Value())

	str = NewString(spec)
	err = str.SetBytes([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, "hello", str.Value())

	str = NewString(spec)
	str.SetValue("hello")
	require.Equal(t, "hello", str.Value())
}

func TestStringNil(t *testing.T) {
	var str *String = nil

	bs, err := str.Bytes()
	require.NoError(t, err)
	require.Nil(t, bs)

	value, err := str.String()
	require.NoError(t, err)
	require.Equal(t, "", value)

	value = str.Value()
	require.Equal(t, "", value)
}

func TestStringPack(t *testing.T) {
	t.Run("returns error for zero value when fixed length and no padding specified", func(t *testing.T) {
		spec := &Spec{
			Length:      10,
			Description: "Field",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}
		str := NewString(spec)
		_, err := str.Pack()

		require.EqualError(t, err, "failed to encode length: field length: 0 should be fixed: 10")
	})
}

func TestStringFieldUnmarshal(t *testing.T) {
	stringField := NewStringValue("123456")

	vString := &String{}
	err := stringField.Unmarshal(vString)
	require.NoError(t, err)
	require.Equal(t, "123456", vString.Value())

	var s string
	err = stringField.Unmarshal(&s)
	require.NoError(t, err)
	require.Equal(t, "123456", s)

	var b int
	err = stringField.Unmarshal(&b)
	require.NoError(t, err)
	require.Equal(t, 123456, b)

	var i64 int64
	err = stringField.Unmarshal(&i64)
	require.NoError(t, err)
	require.Equal(t, int64(123456), i64)

	el := reflect.ValueOf(&i64).Elem()
	err = stringField.Unmarshal(el)
	require.NoError(t, err)
	require.Equal(t, int64(123456), el.Int())

	refStrValue := reflect.ValueOf(&s).Elem()
	err = stringField.Unmarshal(refStrValue)
	require.NoError(t, err)
	require.Equal(t, "123456", refStrValue.String())

	refIntValue := reflect.ValueOf(&b).Elem()
	err = stringField.Unmarshal(refIntValue)
	require.NoError(t, err)
	require.Equal(t, 123456, int(refIntValue.Int()))

	refStr := reflect.ValueOf(s)
	err = stringField.Unmarshal(refStr)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type string", err.Error())

	refStrPointer := reflect.ValueOf(&s)
	err = stringField.Unmarshal(refStrPointer)
	require.Error(t, err)
	require.Equal(t, "cannot set reflect.Value of type ptr", err.Error())

	err = stringField.Unmarshal(nil)
	require.Error(t, err)
	require.Equal(t, "unsupported type: expected *String, *string, or reflect.Value, got <nil>", err.Error())
}

func TestStringFieldMarshal(t *testing.T) {
	stringField := NewStringValue("")

	inputString := NewStringValue("123456")
	stringField.Marshal(inputString)
	require.Equal(t, "123456", stringField.Value())

	stringField.Marshal(&inputString)
	require.Equal(t, "123456", stringField.Value())

	inputStr := "123456"
	stringField.Marshal(inputStr)
	require.Equal(t, "123456", stringField.Value())

	stringField.Marshal(&inputStr)
	require.Equal(t, "123456", stringField.Value())

	inputInt := 123456
	stringField.Marshal(inputInt)
	require.Equal(t, "123456", stringField.Value())

	stringField.Marshal(&inputInt)
	require.Equal(t, "123456", stringField.Value())

	err := stringField.Marshal(nil)
	require.NoError(t, err)

	err = stringField.Marshal([]byte("123456"))
	require.Error(t, err)
	require.Equal(t, "data does not match required *String or (string, *string, int, *int) type", err.Error())
}

func TestStringJSONUnmarshal(t *testing.T) {
	input := []byte(`"4000"`)

	str := NewString(&Spec{
		Length:      4,
		Description: "Field",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.Fixed,
	})

	require.NoError(t, str.UnmarshalJSON(input))
	require.Equal(t, "4000", str.Value())
}

func TestStringJSONMarshal(t *testing.T) {
	str := NewStringValue("1000")
	marshalledJSON, err := str.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"1000"`, string(marshalledJSON))
}
