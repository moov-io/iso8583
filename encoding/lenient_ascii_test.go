package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLenientASCII(t *testing.T) {
	enc := &lenientASCIIEncoder{}

	t.Run("Decode plain ASCII passes through unchanged", func(t *testing.T) {
		res, read, err := enc.Decode([]byte("hello"), 5)

		require.NoError(t, err)
		require.Equal(t, []byte("hello"), res)
		require.Equal(t, 5, read)
	})

	t.Run("Decode passes bytes > 0x7F through without error", func(t *testing.T) {
		// Latin-1 / cp1252 punctuation that strict ASCII rejects:
		//   0xBB = » (right-pointing double angle quotation mark)
		//   0xF5 = õ (Latin small letter o with tilde)
		//   0x9F = Ÿ (Latin capital letter Y with diaeresis, cp1252)
		//   0x80 = € (euro sign, cp1252)
		data := []byte{'a', 0xBB, 'b', 0xF5, 'c', 0x9F, 'd', 0x80, 'e'}

		res, read, err := enc.Decode(data, len(data))

		require.NoError(t, err)
		require.Equal(t, data, res)
		require.Equal(t, len(data), read)
	})

	t.Run("Decode passes UTF-8 multibyte through without error", func(t *testing.T) {
		// "hello, 世界!" — Chinese characters encoded as UTF-8 multibyte
		// sequences; every continuation byte has bit 7 set.
		data := []byte("hello, 世界!")

		res, read, err := enc.Decode(data, len(data))

		require.NoError(t, err)
		require.Equal(t, data, res)
		require.Equal(t, len(data), read)
	})

	t.Run("Decode reads only length bytes", func(t *testing.T) {
		res, read, err := enc.Decode([]byte("hello world"), 5)

		require.NoError(t, err)
		require.Equal(t, []byte("hello"), res)
		require.Equal(t, 5, read)
	})

	t.Run("Decode preserves control characters", func(t *testing.T) {
		data := []byte{0x00, 0x01, 0x1F, 0x7F}

		res, read, err := enc.Decode(data, len(data))

		require.NoError(t, err)
		require.Equal(t, data, res)
		require.Equal(t, len(data), read)
	})

	t.Run("Decode errors when length is negative", func(t *testing.T) {
		_, _, err := enc.Decode([]byte("hello"), -1)

		require.Error(t, err)
		require.EqualError(t, err, "invalid length: -1")
	})

	t.Run("Decode errors when data shorter than length", func(t *testing.T) {
		_, _, err := enc.Decode([]byte("hello"), 6)

		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 6, got 5")
	})

	t.Run("Decode errors when input is nil and length > 0", func(t *testing.T) {
		_, _, err := enc.Decode(nil, 6)

		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 6, got 0")
	})

	t.Run("Decode returns empty slice when length is zero", func(t *testing.T) {
		res, read, err := enc.Decode([]byte("hello"), 0)

		require.NoError(t, err)
		require.Empty(t, res)
		require.Equal(t, 0, read)
	})

	t.Run("Encode plain ASCII passes through unchanged", func(t *testing.T) {
		res, err := enc.Encode([]byte("hello"))

		require.NoError(t, err)
		require.Equal(t, []byte("hello"), res)
	})

	t.Run("Encode passes bytes > 0x7F through without error", func(t *testing.T) {
		data := []byte{'a', 0xBB, 'b', 0xF5, 'c', 0x9F, 'd', 0x80, 'e'}

		res, err := enc.Encode(data)

		require.NoError(t, err)
		require.Equal(t, data, res)
	})

	t.Run("Encode passes UTF-8 multibyte through without error", func(t *testing.T) {
		data := []byte("hello, 世界!")

		res, err := enc.Encode(data)

		require.NoError(t, err)
		require.Equal(t, data, res)
	})

	t.Run("Encode returns empty slice when input is empty", func(t *testing.T) {
		res, err := enc.Encode([]byte{})

		require.NoError(t, err)
		require.Empty(t, res)
	})

	t.Run("Encode returns a copy independent of the input", func(t *testing.T) {
		// Mutating the input after Encode must not affect the returned slice.
		data := []byte("hello")
		res, err := enc.Encode(data)
		require.NoError(t, err)

		data[0] = 'x'
		require.Equal(t, []byte("hello"), res)
	})

	t.Run("Decode returns a copy independent of the input", func(t *testing.T) {
		data := []byte("hello")
		res, _, err := enc.Decode(data, 5)
		require.NoError(t, err)

		data[0] = 'x'
		require.Equal(t, []byte("hello"), res)
	})
}

func TestLenientASCIISingleton(t *testing.T) {
	// LenientASCII must satisfy the Encoder interface and be safe to use as
	// a shared, immutable singleton (consistent with ASCII, Binary, etc.).
	var _ Encoder = LenientASCII

	a, err := LenientASCII.Encode([]byte("abc"))
	require.NoError(t, err)
	require.Equal(t, []byte("abc"), a)

	b, _, err := LenientASCII.Decode([]byte("abc"), 3)
	require.NoError(t, err)
	require.Equal(t, []byte("abc"), b)
}

func FuzzDecodeLenientASCII(f *testing.F) {
	enc := &lenientASCIIEncoder{}

	f.Fuzz(func(t *testing.T, data []byte, length int) {
		enc.Decode(data, length)
	})
}

func FuzzEncodeLenientASCII(f *testing.F) {
	enc := &lenientASCIIEncoder{}

	f.Fuzz(func(t *testing.T, data []byte) {
		enc.Encode(data)
	})
}
