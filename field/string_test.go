package field

import (
	"testing"

	"github.com/franizus/iso8583/encoding"
	"github.com/franizus/iso8583/padding"
	"github.com/franizus/iso8583/prefix"
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
}
