package field

import (
	"bytes"
	"strings"
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

	t.Run("ReadFrom reads data from the reader", func(t *testing.T) {
		str := NewString(spec)

		length, err := str.ReadFrom(strings.NewReader("     olleh"))

		require.NoError(t, err)
		require.Equal(t, "olleh", str.Value)
		require.Equal(t, 10, length)
	})

	t.Run("WritesTo writes data to the writer", func(t *testing.T) {
		str := NewString(spec)
		str.Value = "hello"

		var buf bytes.Buffer

		length, err := str.WriteTo(&buf)

		require.NoError(t, err)
		require.Equal(t, "     hello", buf.String())
		require.Equal(t, 10, length)
	})
}
