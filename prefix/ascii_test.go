package prefix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsciiVarPrefixer_EncodeLengthDigitsValidation(t *testing.T) {
	pref := asciiVarPrefixer{
		Digits: 2,
	}

	_, err := pref.EncodeLength(999, 123)

	require.Contains(t, err.Error(), "number of digits in length: 123 exceeds: 2")
}

func TestAsciiVarPrefixer_EncodeLengthMaxLengthValidation(t *testing.T) {
	pref := asciiVarPrefixer{
		Digits: 2,
	}

	_, err := pref.EncodeLength(20, 22)

	require.Contains(t, err.Error(), "field length: 22 is larger than maximum: 20")
}

func TestAsciiVarPrefixer_DecodeLengthMaxLengthValidation(t *testing.T) {
	pref := asciiVarPrefixer{
		Digits: 3,
	}

	_, _, err := pref.DecodeLength(20, []byte("22"))

	require.Contains(t, err.Error(), "not enough data length: 2 to read: 3 byte digits")
}

func TestAsciiVarPrefixer_LHelpers(t *testing.T) {
	tests := []struct {
		pref   Prefixer
		digits int
		maxLen int
		in     int
		out    []byte
	}{
		{ASCII.L, 1, 5, 3, []byte("3")},
		{ASCII.LL, 2, 20, 2, []byte("02")},
		{ASCII.LL, 2, 20, 12, []byte("12")},
		{ASCII.LLL, 3, 340, 2, []byte("002")},
		{ASCII.LLL, 3, 340, 200, []byte("200")},
		{ASCII.LLLL, 4, 9999, 1234, []byte("1234")},
	}

	// test encoding
	for _, tt := range tests {
		t.Run(tt.pref.Inspect()+"_EncodeLength", func(t *testing.T) {
			got, err := tt.pref.EncodeLength(tt.maxLen, tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.out, got)
		})
	}

	// test decoding
	for _, tt := range tests {
		t.Run(tt.pref.Inspect()+"_DecodeLength", func(t *testing.T) {
			got, read, err := tt.pref.DecodeLength(tt.maxLen, tt.out)
			require.NoError(t, err)
			require.Equal(t, tt.in, got)
			require.Equal(t, tt.digits, read)
		})
	}
}

func TestAsciiFixedPrefixer(t *testing.T) {
	pref := asciiFixedPrefixer{}

	// Fixed prefixer returns empty byte slice as
	// size is not encoded into field
	data, err := pref.EncodeLength(8, 8)

	require.NoError(t, err)
	require.Equal(t, 0, len(data))

	// Fixed prefixer returns configured len
	// rather than read it from data
	dataLen, read, err := pref.DecodeLength(8, []byte("data"))

	require.NoError(t, err)
	require.Equal(t, 8, dataLen)
	require.Equal(t, 0, read)
}

func TestAsciiFixedPrefixer_EncodeLengthValidation(t *testing.T) {
	pref := asciiFixedPrefixer{}

	_, err := pref.EncodeLength(8, 12)

	require.Contains(t, err.Error(), "field length: 12 should be fixed: 8")
}
