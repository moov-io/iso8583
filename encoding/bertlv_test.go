package encoding

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBerTLVTag(t *testing.T) {
	tests := []struct {
		desc     string
		numBytes int
		hexTag   []byte
		asciiTag []byte
	}{
		{"PAN (single byte tag)", 1, []byte{0x5A}, []byte("5A")},
		{"CVM List (single byte tag)", 1, []byte{0x8E}, []byte("8E")},
		{"Acquirer Identifier (two byte tag)", 2, []byte{0x5F, 0x2A}, []byte("5F2A")},
		{"BIC (two byte tag)", 2, []byte{0x5F, 0x54}, []byte("5F54")},
		{"Authorized Amount", 2, []byte{0x9F, 0x02}, []byte("9F02")},
		{"ATC Register (two byte tag)", 2, []byte{0x9F, 0x13}, []byte("9F13")},
		{"Imaginary three byte tag", 3, []byte{0x9F, 0xA8, 0x13}, []byte("9FA813")},
	}

	for _, tt := range tests {
		t.Run(tt.desc+"_Decode", func(t *testing.T) {
			asciiTag, read, err := BerTLVTag.Decode(tt.hexTag, 0)
			require.NoError(t, err)
			require.Equal(t, tt.asciiTag, asciiTag)
			require.Equal(t, tt.numBytes, read)
		})
	}

	for _, tt := range tests {
		t.Run(tt.desc+"_Encode", func(t *testing.T) {
			hexTag, err := BerTLVTag.Encode(tt.asciiTag)
			require.NoError(t, err)
			require.Equal(t, tt.hexTag, hexTag)
		})
	}
}

func TestBerTLVTag_DecodeOnInvalidInput(t *testing.T) {
	t.Run("when bytes are nil", func(t *testing.T) {
		_, _, err := BerTLVTag.Decode(nil, 0)
		require.EqualError(t, err, "failed to read byte")
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("when bytes are empty", func(t *testing.T) {
		_, _, err := BerTLVTag.Decode([]byte{}, 0)
		require.EqualError(t, err, "failed to read byte")
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("when bits 5-1 of first byte set but 2nd byte does not exist", func(t *testing.T) {
		_, _, err := BerTLVTag.Decode([]byte{0x5F}, 0)
		require.EqualError(t, err, "failed to decode TLV tag")
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("when MSB of 2nd byte set but 3nd byte does not exist", func(t *testing.T) {
		_, _, err := BerTLVTag.Decode([]byte{0x5F, 0xA8}, 0)
		require.EqualError(t, err, "failed to decode TLV tag")
		require.ErrorIs(t, err, io.EOF)
	})
}
