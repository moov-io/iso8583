package prefix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBerTLVPrefixer(t *testing.T) {
	tests := []struct {
		desc     string
		numBytes int
		length   int
		data     []byte
	}{
		{"single byte prefix", 1, 126, []byte{0b01111110}},
		{"two byte prefix", 2, 131, []byte{0b10000001, 0b10000011}},
		{"three byte prefix", 3, 65039, []byte{0b10000010, 0b11111110, 0b00001111}},
	}

	// test encoding
	for _, tt := range tests {
		t.Run(tt.desc+"_EncodeLength", func(t *testing.T) {
			got, err := BerTLV.EncodeLength(0, tt.length)
			require.NoError(t, err)
			require.Equal(t, tt.data, got)
		})
	}

	// test decoding
	for _, tt := range tests {
		t.Run(tt.desc+"_DecodeLength", func(t *testing.T) {
			got, read, err := BerTLV.DecodeLength(0, tt.data)
			require.NoError(t, err)
			require.Equal(t, tt.length, got)
			require.Equal(t, tt.numBytes, read)
		})
	}
}

func TestBerTLVPrefixer_DecodeReturnsErrOnIncorrectInput(t *testing.T) {
	t.Run("if initial byte larger than number of subsequent bytes", func(t *testing.T) {
		// First byte indicates that there should be 3 additional bytes to follow
		// However, there are only 2.
		_, _, err := BerTLV.DecodeLength(0, []byte{0b10000011, 0b11111110, 0b00001111})
		require.EqualError(t, err, "failed to read long form TLV length: unexpected EOF")
	})

	t.Run("if input is empty", func(t *testing.T) {
		_, _, err := BerTLV.DecodeLength(0, []byte{})
		require.EqualError(t, err, "failed to decode TLV length: EOF")
	})
}
