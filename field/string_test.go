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

func TestStringFieldZeroLength(t *testing.T) {
	str := NewStringValue("")
	str.SetSpec(&Spec{
		Length:      10,
		Description: "Field",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.Fixed,
	})

	packed, err := str.Pack()
	require.NoError(t, err)
	require.Equal(t, []byte{}, packed)
	require.Equal(t, "", string(packed))
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
