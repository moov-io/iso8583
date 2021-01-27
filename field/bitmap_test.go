package field

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestHexBitmap(t *testing.T) {
	t.Run("when not enough data to unpack", func(t *testing.T) {
		field := NewBitmap(&Spec{
			Length:      16,
			Description: "Bitmap",
			Enc:         encoding.Hex,
			Pref:        prefix.Hex.Fixed,
		})
		_, err := field.Unpack([]byte("4000"))

		require.Error(t, err)
		require.Contains(t, err.Error(), "expected min data length is 16, but it is 4")
	})

	t.Run("decode primary bitmap", func(t *testing.T) {
		field := NewBitmap(&Spec{
			Length:      16,
			Description: "Bitmap",
			Enc:         encoding.Hex,
			Pref:        prefix.Hex.Fixed,
		})
		// bits 2, 20 are set
		length, err := field.Unpack([]byte("4000100000000000"))

		require.NoError(t, err)
		require.Equal(t, length, 16)
	})

	t.Run("decode primary and secondary bitmap", func(t *testing.T) {
		field := NewBitmap(&Spec{
			Length:      16,
			Description: "Bitmap",
			Enc:         encoding.Hex,
			Pref:        prefix.Hex.Fixed,
		})
		// bits 2, 20, 65,120 are set
		length, err := field.Unpack([]byte("c0001000000000008000000000000100"))

		require.NoError(t, err)
		require.Equal(t, length, 32)
	})

	t.Run("with primary bitmap only it returns only half of the full length", func(t *testing.T) {
		field := NewBitmap(&Spec{
			Length:      16,
			Description: "Bitmap",
			Enc:         encoding.Hex,
			Pref:        prefix.Hex.Fixed,
		})

		bitmap := field.(*Bitmap)

		// first 64 bits are referred to as the primary bit map
		bitmap.Set(2)
		bitmap.Set(20)

		data, err := bitmap.Pack()

		fmt.Println("str", string(data))

		require.NoError(t, err)
		require.Len(t, data, 16)
	})

	t.Run("with secondary bitmap only it returns only half of the full length", func(t *testing.T) {
		field := NewBitmap(&Spec{
			Length:      16,
			Description: "Bitmap",
			Enc:         encoding.Hex,
			Pref:        prefix.Hex.Fixed,
		})

		bitmap := field.(*Bitmap)

		// first 64 bits are referred to as the primary bit map
		bitmap.Set(2)
		bitmap.Set(20)
		bitmap.Set(65)
		bitmap.Set(120)

		data, err := bitmap.Pack()

		require.NoError(t, err)
		require.Len(t, data, 32)
	})
}

func TestBinaryBitmap(t *testing.T) {
	field := NewBitmap(&Spec{
		Length:      16,
		Description: "Bitmap",
		Enc:         encoding.Binary,
		Pref:        prefix.Binary.Fixed,
	})

	bitmap := field.(*Bitmap)

	t.Run("with primary bitmap only it returns only half of the full length", func(t *testing.T) {
		// first 64 bits are referred to as the primary bit map
		bitmap.Set(2)
		bitmap.Set(20)

		data, err := bitmap.Pack()

		require.NoError(t, err)
		require.Len(t, data, 8)
	})

	t.Run("with secondary bitmap only it returns only half of the full length", func(t *testing.T) {
		// bits 65 through 128 are referred to as the secondary bit map
		bitmap.Set(2)
		bitmap.Set(20)
		bitmap.Set(65)
		bitmap.Set(120)

		data, err := bitmap.Pack()

		require.NoError(t, err)
		require.Len(t, data, 16)
	})
}
