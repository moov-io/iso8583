package prefixer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// prefixer.NONE

// func TestPrefixer(t *testing.T) {
// 	require.Equal(t, []byte(nil), NONE.Encode(12))
// 	require.Equal(t, -1, NONE.Decode([]byte("16")))

// 	require.Equal(t, []byte("12"), ASCII.Encode(12))
// 	require.Equal(t, 16, ASCII.Decode([]byte("16")))
// }

func TestAsciiPrefixer_EncodeLengthDigitsCheck(t *testing.T) {
	pref := asciiPrefixer{
		MaxLen: 999,
		Digits: 2,
	}

	_, err := pref.EncodeLength(123)

	require.Contains(t, err.Error(), "Number of digits exceeds: 2")
}

func TestAsciiPrefixer_EncodeLengthMaxLength(t *testing.T) {
	pref := asciiPrefixer{
		MaxLen: 20,
		Digits: 2,
	}

	_, err := pref.EncodeLength(22)

	require.Contains(t, err.Error(), "Provided length: 22 is larger than maximum: 20")
}

func TestAsciiPrefixer_DecodeLengthMaxLength(t *testing.T) {
	pref := asciiPrefixer{
		MaxLen: 20,
		Digits: 3,
	}

	_, err := pref.DecodeLength([]byte("22"))

	require.Contains(t, err.Error(), "Not enought data length: 2 to read: 3 byte digits")
}

func TestAsciiPrefixer_VarLength(t *testing.T) {
	tests := []struct {
		pref Prefixer
		in   int
		out  []byte
	}{
		{ASCII.L(5), 3, []byte("3")},
		{ASCII.LL(20), 12, []byte("12")},
		{ASCII.LLL(340), 200, []byte("200")},
		{ASCII.LLLL(9999), 1234, []byte("1234")},
	}

	// test encoding
	for _, tt := range tests {
		t.Run(tt.pref.Inspect()+"_EncodeLength", func(t *testing.T) {
			got, err := tt.pref.EncodeLength(tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.out, got)
		})
	}

	// test decoding
	for _, tt := range tests {
		t.Run(tt.pref.Inspect()+"_DecodeLength", func(t *testing.T) {
			got, err := tt.pref.DecodeLength(tt.out)
			require.NoError(t, err)
			require.Equal(t, tt.in, got)
		})
	}
}
