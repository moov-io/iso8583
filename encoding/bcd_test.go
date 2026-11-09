package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yerden/go-util/bcd"
)

func TestBCD(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		res, read, err := BCD.Decode([]byte{0x12, 0x34}, 4)

		require.NoError(t, err)
		require.Equal(t, []byte("1234"), res)
		require.Equal(t, 2, read)

		res, read, err = BCD.Decode([]byte{0x01, 0x23}, 3)

		require.NoError(t, err)
		require.Equal(t, []byte("123"), res)
		require.Equal(t, 2, read)

		res, read, err = BCD.Decode([]byte{0x12, 0x30}, 3)

		require.NoError(t, err)
		require.Equal(t, []byte("230"), res)
		require.Equal(t, 2, read)

		res, read, err = BCD.Decode([]byte{0x21, 0x43, 0x55}, 4)

		require.NoError(t, err)
		require.Equal(t, []byte("2143"), res)
		require.Equal(t, 2, read)

		res, read, err = BCD.Decode([]byte{0x21, 0x43, 0xff}, 4)

		require.NoError(t, err)
		require.Equal(t, []byte("2143"), res)
		require.Equal(t, 2, read)

		_, _, err = BCD.Decode([]byte{0x21, 0x43}, 6)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 3, got 2")

		_, _, err = BCD.Decode(nil, 6)
		require.Error(t, err)
		require.EqualError(t, err, "not enough data to decode. expected len 3, got 0")

		_, _, err = BCD.Decode([]byte{0xAB, 0xCD}, 4)
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform BCD decoding")
		require.ErrorIs(t, err, bcd.ErrBadBCD)
	})

	t.Run("Encode", func(t *testing.T) {
		res, err := BCD.Encode([]byte("0110"))
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x10}, res)

		// right justified by default
		res, err = BCD.Encode([]byte("123"))
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x23}, res)

		_, err = BCD.Encode([]byte("abc"))
		require.Error(t, err)
		require.EqualError(t, err, "failed to perform BCD encoding")
		require.ErrorIs(t, err, bcd.ErrBadInput)
	})
}

func BenchmarkBCDEncode(b *testing.B) {
	src := []byte("1234567890123456")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := BCD.Encode(src); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBCDecode(b *testing.B) {
	src := []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, _, err := BCD.Decode(src, len(src)*2); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBCDRoundTrip(b *testing.B) {
	src := []byte("1234567890123456")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encoded, err := BCD.Encode(src)
		if err != nil {
			b.Fatal(err)
		}
		decoded, _, err := BCD.Decode(encoded, len(src))
		if err != nil {
			b.Fatal(err)
		}
		if !bytes.Equal(decoded, src) {
			b.Fatal("BCD round trip mismatch")
		}
	}
}

func FuzzDecodeBCD(f *testing.F) {
	enc := &bcdEncoder{}

	f.Fuzz(func(t *testing.T, data []byte, length int) {
		enc.Decode(data, length)
	})
}
