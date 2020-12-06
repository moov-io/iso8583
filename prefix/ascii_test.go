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

	require.Contains(t, err.Error(), "Number of digits exceeds: 2")
}

func TestAsciiVarPrefixer_EncodeLengthMaxLengthValidation(t *testing.T) {
	pref := asciiVarPrefixer{
		Digits: 2,
	}

	_, err := pref.EncodeLength(20, 22)

	require.Contains(t, err.Error(), "Field length: 22 is larger than maximum: 20")
}

func TestAsciiVarPrefixer_DecodeLengthMaxLengthValidation(t *testing.T) {
	pref := asciiVarPrefixer{
		Digits: 3,
	}

	_, err := pref.DecodeLength(20, []byte("22"))

	require.Contains(t, err.Error(), "Not enought data length: 2 to read: 3 byte digits")
}

func TestAsciiVarPrefixer_LHelpers(t *testing.T) {
	tests := []struct {
		pref   Prefixer
		maxLen int
		in     int
		out    []byte
	}{
		{ASCII.L, 5, 3, []byte("3")},
		{ASCII.LL, 20, 12, []byte("12")},
		{ASCII.LLL, 340, 200, []byte("200")},
		{ASCII.LLLL, 9999, 1234, []byte("1234")},
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
			got, err := tt.pref.DecodeLength(tt.maxLen, tt.out)
			require.NoError(t, err)
			require.Equal(t, tt.in, got)
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
	dataLen, err := pref.DecodeLength(8, []byte("data"))

	require.NoError(t, err)
	require.Equal(t, 8, dataLen)
}

func TestAsciiFixedPrefixer_EncodeLengthValidation(t *testing.T) {
	pref := asciiFixedPrefixer{}

	_, err := pref.EncodeLength(8, 12)

	require.Contains(t, err.Error(), "Field length: 12 should be fixed: 8")
}
