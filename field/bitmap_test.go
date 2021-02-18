package field

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/utils"
	"github.com/stretchr/testify/require"
)

func TestHexBitmap(t *testing.T) {
	// test
	// when there is only fields for the first bitmap
	// then the legth of the bitmap should be 16 (hex)
	// when there are fields for the second bitmap
	// then the length of the bitmap should be 32
	// when there are fields for the third bitmap
	// then the length of the bitmap should be 48

	// b1.Set(10)
	// 004000000000000000000000000000000000000000000000

	// b2.Set(1) //second bitmap presents
	// b2.Set(10)
	// b2.Set(70)
	// 804000000000000004000000000000000000000000000000

	// b3.Set(1)  //second bitmap presents
	// b3.Set(65) //third bitmap presents
	// b3.Set(10)
	// b3.Set(140)
	// 804000000000000080000000000000000010000000000000
	t.Run("Read only first bitmap", func(t *testing.T) {
		field := &Bitmap{
			spec: &Spec{
				Description: "Bitmap",
				Enc:         encoding.Hex,
				Pref:        prefix.Hex.Fixed,
			},
			bitmap: utils.NewBitmap(192),
		}

		read, err := field.Unpack([]byte("004000000000000000000000000000000000000000000000"))

		require.NoError(t, err)
		require.Equal(t, 16, read)

		require.True(t, field.IsSet(10))
	})

	t.Run("Read two bitmaps", func(t *testing.T) {
		field := &Bitmap{
			spec: &Spec{
				Description: "Bitmap",
				Enc:         encoding.Hex,
				Pref:        prefix.Hex.Fixed,
			},
			bitmap: utils.NewBitmap(192),
		}

		read, err := field.Unpack([]byte("804000000000000004000000000000000000000000000000"))

		require.NoError(t, err)
		require.Equal(t, 32, read)

		require.True(t, field.IsSet(10))
		require.True(t, field.IsSet(70))
	})

	t.Run("Read three bitmaps", func(t *testing.T) {
		field := &Bitmap{
			spec: &Spec{
				Length:      16,
				Description: "Bitmap",
				Enc:         encoding.Hex,
				Pref:        prefix.Hex.Fixed,
			},
			bitmap: utils.NewBitmap(192),
		}

		read, err := field.Unpack([]byte("804000000000000080000000000000000010000000000000"))

		require.NoError(t, err)
		require.Equal(t, 48, read)

		require.True(t, field.IsSet(10))
		require.True(t, field.IsSet(140))
	})

	t.Run("when not enough data to unpack", func(t *testing.T) {
		field := NewBitmap(&Spec{
			Length:      16,
			Description: "Bitmap",
			Enc:         encoding.Hex,
			Pref:        prefix.Hex.Fixed,
		})
		_, err := field.Unpack([]byte("4000"))

		require.Error(t, err)
		require.Contains(t, err.Error(), "not enough data to read 1 bitmap")
	})

	t.Run("when bit for secondary bitmap is set but not enough data to read", func(t *testing.T) {
		field := NewBitmap(&Spec{
			Length:      16,
			Description: "Bitmap",
			Enc:         encoding.Hex,
			Pref:        prefix.Hex.Fixed,
		})
		// bits 2, 20, 65, 120 are set, but no data for third bitmap
		_, err := field.Unpack([]byte("c0001000000000008000000000000100"))

		// error
		require.Error(t, err)
		require.Contains(t, err.Error(), "not enough data to read 3 bitmap")
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
