package field

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestHexBitmap(t *testing.T) {
	t.Run("Read only first bitmap", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		})

		// set bit: 10
		read, err := bitmap.Unpack([]byte("004000000000000000000000000000000000000000000000"))

		require.NoError(t, err)
		require.Equal(t, 16, read) // 16 is 8 bytes (one bitmap) encoded in hex

		require.True(t, bitmap.IsSet(10))
	})

	t.Run("Read two bitmaps", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		})

		// set bits: 1, 10, 70
		read, err := bitmap.Unpack([]byte("804000000000000004000000000000000000000000000000"))

		require.NoError(t, err)
		require.Equal(t, 32, read) // 32 is 16 bytes (two bitmaps) encoded in hex

		require.True(t, bitmap.IsSet(10))
		require.True(t, bitmap.IsSet(70))
	})

	t.Run("Read three bitmaps", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		})

		// set bits: 1, 10, 65, 140
		read, err := bitmap.Unpack([]byte("804000000000000080000000000000000010000000000000"))

		require.NoError(t, err)
		require.Equal(t, 48, read) // 48 is 24 bytes (three bitmaps) encoded in hex

		require.True(t, bitmap.IsSet(10))
		require.True(t, bitmap.IsSet(140))
	})

	t.Run("When not enough data to unpack", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		})

		_, err := bitmap.Unpack([]byte("4000"))

		require.Error(t, err)
		require.Contains(t, err.Error(), "for 1 bitmap: not enough data to read")
	})

	t.Run("When bit for secondary bitmap is set but not enough data to read", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		})

		// bits 2, 20, 65, 120 are set, but no data for third bitmap
		_, err := bitmap.Unpack([]byte("c0001000000000008000000000000100"))

		require.Error(t, err)
		require.Contains(t, err.Error(), "for 3 bitmap: not enough data to read")
	})

	t.Run("With primary bitmap only it returns signle bitmap length", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		})

		bitmap.Set(20) // first bitmap field

		data, err := bitmap.Pack()

		require.NoError(t, err)
		require.Len(t, data, 16) // 16 bytes is 8 bytes (one bitmap) encoded in hex
	})

	t.Run("With secondary bitmap it returns length of two bitmaps", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		})

		bitmap.Set(20) // first bitmap field
		bitmap.Set(70) // second bitmap field

		data, err := bitmap.Pack()

		require.NoError(t, err)
		require.Len(t, data, 32) // 32 bytes is 16 bytes (two bitmaps) encoded in hex
	})

	t.Run("With third bitmap it returns length of three bitmaps", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		})

		bitmap.Set(20)  // first bitmap field
		bitmap.Set(70)  // second bitmap field
		bitmap.Set(150) // third bitmap field

		data, err := bitmap.Pack()

		require.NoError(t, err)
		require.Len(t, data, 48) // 48 bytes is 24 bytes (three bitmaps) encoded in hex
	})
}

func TestBinaryBitmap(t *testing.T) {
	t.Run("With primary bitmap only it returns signle bitmap length", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.Fixed,
		})

		bitmap.Set(20) // first bitmap field

		data, err := bitmap.Pack()

		require.NoError(t, err)
		require.Len(t, data, 8)
	})

	t.Run("With secondary bitmap it returns length of two bitmaps", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.Fixed,
		})

		bitmap.Set(20) // first bitmap field
		bitmap.Set(70) // second bitmap field

		data, err := bitmap.Pack()

		require.NoError(t, err)
		require.Len(t, data, 16)
	})

	t.Run("With third bitmap it returns length of three bitmaps", func(t *testing.T) {
		bitmap := NewBitmap(&Spec{
			Description: "Bitmap",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.Fixed,
		})

		bitmap.Set(20)  // first bitmap field
		bitmap.Set(70)  // second bitmap field
		bitmap.Set(150) // third bitmap field

		data, err := bitmap.Pack()

		require.NoError(t, err)
		require.Len(t, data, 24)
	})
}

func TestBitmap_Unmarshal(t *testing.T) {
	spec := &Spec{
		Description: "Bitmap",
		Enc:         encoding.BytesToASCIIHex,
		Pref:        prefix.Hex.Fixed,
	}

	t.Run("Unmarshal gets bitmap into data parameter", func(t *testing.T) {
		bitmap := NewBitmap(spec)
		bitmap.Set(10) // set bit

		data := NewBitmap(nil)

		err := bitmap.Unmarshal(data)

		require.NoError(t, err)
		require.True(t, data.IsSet(10))
	})
}

func TestBitmap_SetData(t *testing.T) {
	spec := &Spec{
		Description: "Bitmap",
		Enc:         encoding.BytesToASCIIHex,
		Pref:        prefix.Hex.Fixed,
	}
	bitmapBytes := []byte("004000000000000000000000000000000000000000000000")

	t.Run("Nil data causes no side effects", func(t *testing.T) {
		bitmap := NewBitmap(spec)
		err := bitmap.SetData(nil)
		require.NoError(t, err)
		require.Equal(t, NewBitmap(spec), bitmap)
	})

	t.Run("non-Bitmap data type returns error", func(t *testing.T) {
		bitmap := NewBitmap(spec)

		str := &struct {
			a string
		}{"left"}

		err := bitmap.SetData(str)
		require.Error(t, err)
	})

	t.Run("Unpack sets the data field with the correct bitmap provided using SetData", func(t *testing.T) {
		bitmap := NewBitmap(spec)

		data := &Bitmap{}
		bitmap.SetData(data)
		// set bit: 10
		read, err := bitmap.Unpack(bitmapBytes)
		require.NoError(t, err)
		require.Equal(t, 16, read) // 16 is 8 bytes (one bitmap) encoded in hex

		bitmapBytes, err := bitmap.Bytes()
		require.NoError(t, err)
		dataBytes, err := data.Bytes()
		require.NoError(t, err)
		require.Equal(t, bitmapBytes, dataBytes)
	})

	t.Run("Pack returns bytes using the bitmap provided using SetData", func(t *testing.T) {
		bitmap := NewBitmap(spec)

		data := NewBitmap(nil)
		data.Set(20) // first bitmap field

		bitmap.SetData(data)

		packed, err := bitmap.Pack()
		require.NoError(t, err)
		require.Len(t, packed, 16) // 16 bytes is 8 bytes (one bitmap) encoded in hex
	})

	t.Run("SetBytes sets data to the data field", func(t *testing.T) {
		bitmap := NewBitmap(spec)

		data := &Bitmap{}
		bitmap.SetData(data)

		err := bitmap.SetBytes([]byte("a"))
		require.NoError(t, err)
		b, err := data.Bytes()
		require.NoError(t, err)
		require.Equal(t, []byte("a"), b)
	})
}

func TestBitmapNil(t *testing.T) {
	var str *Bitmap = nil

	bs, err := str.Bytes()
	require.NoError(t, err)
	require.Nil(t, bs)

	value, err := str.String()
	require.NoError(t, err)
	require.Equal(t, "", value)
}
