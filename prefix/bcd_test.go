package prefix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBCDVarPrefixer_EncodeLengthDigitsValidation(t *testing.T) {
	_, err := BCD.LL.EncodeLength(999, 123)

	require.Contains(t, err.Error(), "number of digits in length: 123 exceeds: 2")
}

func TestBCDVarPrefixer_EncodeLengthMaxLengthValidation(t *testing.T) {
	_, err := BCD.LL.EncodeLength(20, 22)

	require.Contains(t, err.Error(), "field length: 22 is larger than maximum: 20")
}

func TestBCDVarPrefixer_DecodeLengthMaxLengthValidation(t *testing.T) {
	_, _, err := BCD.LLL.DecodeLength(20, []byte{0x22})

	require.Contains(t, err.Error(), "length mismatch: want to read 2 bytes, get only 1")
}

func TestBCDVarPrefixer_LHelpers(t *testing.T) {
	tests := []struct {
		pref      Prefixer
		bytesRead int
		maxLen    int
		in        int
		out       []byte
	}{
		{BCD.L, 1, 5, 3, []byte{0x03}},
		{BCD.LL, 1, 20, 2, []byte{0x02}},
		{BCD.LL, 1, 20, 12, []byte{0x12}},
		{BCD.LLL, 2, 340, 2, []byte{0x00, 0x02}},
		{BCD.LLL, 2, 340, 200, []byte{0x02, 0x00}},
		{BCD.LLLL, 2, 9999, 1234, []byte{0x12, 0x34}},
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
			require.Equal(t, tt.bytesRead, read)
		})
	}
}

func TestBCDFixedPrefixer(t *testing.T) {
	pref := bcdFixedPrefixer{}

	// Fixed prefixer returns empty byte slice as
	// size is not encoded into field
	data, err := pref.EncodeLength(8, 8)

	require.NoError(t, err)
	require.Equal(t, 0, len(data))

	// Fixed prefixer returns configured len
	// rather than read it from data
	dataLen, read, err := pref.DecodeLength(8, []byte("1234"))

	require.NoError(t, err)
	require.Equal(t, 8, dataLen)
	require.Equal(t, 0, read)
}

func TestBCDFixedPrefixer_EncodeLengthValidation(t *testing.T) {
	pref := bcdFixedPrefixer{}

	_, err := pref.EncodeLength(8, 12)

	require.Contains(t, err.Error(), "field length: 12 should be fixed: 8")
}
