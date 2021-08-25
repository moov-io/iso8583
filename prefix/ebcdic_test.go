package prefix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEBCDICVarPrefixer_EncodeLengthDigitsValidation(t *testing.T) {
	_, err := EBCDIC.LL.EncodeLength(999, 123)

	require.Contains(t, err.Error(), "number of digits in length: 123 exceeds: 2")
}

func TestEBCDICVarPrefixer_EncodeLengthMaxLengthValidation(t *testing.T) {
	_, err := EBCDIC.LL.EncodeLength(20, 22)

	require.Contains(t, err.Error(), "field length: 22 is larger than maximum: 20")
}

func TestEBCDICVarPrefixer_DecodeLengthMaxLengthValidation(t *testing.T) {
	_, _, err := EBCDIC.LLL.DecodeLength(20, []byte{0x22})

	require.Contains(t, err.Error(), "length mismatch: want to read 3 bytes, get only 1")
}

func TestEBCDICVarPrefixer_LHelpers(t *testing.T) {
	tests := []struct {
		pref      Prefixer
		bytesRead int
		maxLen    int
		in        int
		out       []byte
	}{
		{EBCDIC.L, 1, 5, 3, []byte{0xf3}},
		{EBCDIC.LL, 2, 20, 2, []byte{0xf0, 0xf2}},
		{EBCDIC.LL, 2, 20, 12, []byte{0xf1, 0xf2}},
		{EBCDIC.LLL, 3, 340, 2, []byte{0xf0, 0xf0, 0xf2}},
		{EBCDIC.LLL, 3, 340, 200, []byte{0xf2, 0xf0, 0xf0}},
		{EBCDIC.LLLL, 4, 9999, 1234, []byte{0xf1, 0xf2, 0xf3, 0xf4}},
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

func TestEBCDICFixedPrefixer(t *testing.T) {
	pref := ebcdicFixedPrefixer{}

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

func TestEBCDICFixedPrefixer_EncodeLengthValidation(t *testing.T) {
	pref := ebcdicFixedPrefixer{}

	_, err := pref.EncodeLength(8, 12)

	require.Contains(t, err.Error(), "field length: 12 should be fixed: 8")
}
