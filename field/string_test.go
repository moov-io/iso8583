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
	str := NewStringValue("123456")

	val1 := &String{}
	err := str.Unmarshal(val1)
	require.NoError(t, err)
	require.Equal(t, "123456", val1.Value())

	var val2 string
	err = str.Unmarshal(&val2)
	require.NoError(t, err)
	require.Equal(t, "123456", val2)

	var val3 int
	err = str.Unmarshal(&val3)
	require.NoError(t, err)
	require.Equal(t, 123456, val3)

	val4 := reflect.ValueOf(&val2).Elem()
	err = str.Unmarshal(val4)
	require.NoError(t, err)
	require.Equal(t, "123456", val4.String())

	val5 := reflect.ValueOf(&val3).Elem()
	err = str.Unmarshal(val5)
	require.NoError(t, err)
	require.Equal(t, 123456, int(val5.Int()))

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
	require.Equal(t, "unsupported type: expected *String, *string, or reflect.Value, got <nil>", err.Error())
}

func TestStringFieldMarshal(t *testing.T) {
	str := NewStringValue("")
	vString := NewStringValue("123456")
	str.Marshal(vString)
	require.Equal(t, "123456", str.Value())

	str.Marshal(&vString)
	require.Equal(t, "123456", str.Value())

	vstring := "123456"
	str.Marshal(vstring)
	require.Equal(t, "123456", str.Value())

	str.Marshal(&vstring)
	require.Equal(t, "123456", str.Value())

	vint := 123456
	str.Marshal(vint)
	require.Equal(t, "123456", str.Value())

	str.Marshal(&vint)
	require.Equal(t, "123456", str.Value())

	err := str.Marshal(nil)
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
