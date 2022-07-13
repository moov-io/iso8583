package field

import (
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
	require.Equal(t, "hello", str.Value)

	packed, err := str.Pack()
	require.NoError(t, err)
	require.Equal(t, "     hello", string(packed))

	length, err := str.Unpack([]byte("     olleh"))
	require.NoError(t, err)
	require.Equal(t, 10, length)

	b, err := str.Bytes()
	require.NoError(t, err)
	require.Equal(t, "olleh", string(b))

	require.Equal(t, "olleh", str.Value)

	str = NewString(spec)
	str.SetData(NewStringValue("hello"))
	packed, err = str.Pack()
	require.NoError(t, err)
	require.Equal(t, "     hello", string(packed))

	str = NewString(spec)
	data := NewStringValue("")
	str.SetData(data)
	length, err = str.Unpack([]byte("     olleh"))
	require.NoError(t, err)
	require.Equal(t, 10, length)
	require.Equal(t, "olleh", data.Value)

	str = NewString(spec)
	data = &String{}
	str.SetData(data)
	err = str.SetBytes([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, "hello", data.Value)

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
	str := NewStringValue("hello")

	val := &String{}

	err := str.Unmarshal(val)

	require.NoError(t, err)
	require.Equal(t, "hello", val.Value)
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
	require.Equal(t, "4000", str.Value)
}

func TestStringJSONMarshal(t *testing.T) {
	str := NewStringValue("1000")
	marshalledJSON, err := str.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"1000"`, string(marshalledJSON))
}
