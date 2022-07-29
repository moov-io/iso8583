package prefix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Byte representations of EBCDIC (Code Page 1047) codes for digits
const (
	ebcdic1047_0 = 0xF0
	ebcdic1047_1 = 0xF1
	ebcdic1047_2 = 0xF2
	ebcdic1047_3 = 0xF3
	ebcdic1047_4 = 0xF4
	ebcdic1047_5 = 0xF5
	ebcdic1047_6 = 0xF6
	ebcdic1047_7 = 0xF7
	ebcdic1047_8 = 0xF8
	ebcdic1047_9 = 0xF9
)

// Some non-digit EBCDIC-encoded (Code Page 1047) characters
const (
	ebcdic1047l    = 0x93 // "l"
	ebcdic1047Plus = 0x4E // "+"
	ebcdic1047Dot  = 0x4B // "."
)

func TestEBCDIC1047PrefixersEncode(t *testing.T) {
	t.Parallel()
	for _, testCase := range []struct {
		prefixer       Prefixer
		maxLen         int
		dataLen        int
		expectedOutput []byte
	}{
		{
			prefixer:       EBCDIC1047.Fixed,
			maxLen:         64,
			dataLen:        64,
			expectedOutput: []byte{},
		},
		{
			prefixer:       EBCDIC1047.L,
			maxLen:         8,
			dataLen:        8,
			expectedOutput: []byte{ebcdic1047_8},
		},
		{
			prefixer:       EBCDIC1047.LL,
			maxLen:         14,
			dataLen:        9,
			expectedOutput: []byte{ebcdic1047_0, ebcdic1047_9},
		},
		{
			prefixer:       EBCDIC1047.LL,
			maxLen:         32,
			dataLen:        16,
			expectedOutput: []byte{ebcdic1047_1, ebcdic1047_6},
		},
		{
			prefixer:       EBCDIC1047.LLL,
			maxLen:         128,
			dataLen:        1,
			expectedOutput: []byte{ebcdic1047_0, ebcdic1047_0, ebcdic1047_1},
		},
		{
			prefixer:       EBCDIC1047.LLL,
			maxLen:         256,
			dataLen:        81,
			expectedOutput: []byte{ebcdic1047_0, ebcdic1047_8, ebcdic1047_1},
		},
		{
			prefixer:       EBCDIC1047.LLL,
			maxLen:         512,
			dataLen:        512,
			expectedOutput: []byte{ebcdic1047_5, ebcdic1047_1, ebcdic1047_2},
		},
		{
			prefixer:       EBCDIC1047.LLLL,
			maxLen:         1024,
			dataLen:        5,
			expectedOutput: []byte{ebcdic1047_0, ebcdic1047_0, ebcdic1047_0, ebcdic1047_5},
		},
		{
			prefixer:       EBCDIC1047.LLLL,
			maxLen:         2048,
			dataLen:        66,
			expectedOutput: []byte{ebcdic1047_0, ebcdic1047_0, ebcdic1047_6, ebcdic1047_6},
		},
		{
			prefixer:       EBCDIC1047.LLLL,
			maxLen:         4096,
			dataLen:        124,
			expectedOutput: []byte{ebcdic1047_0, ebcdic1047_1, ebcdic1047_2, ebcdic1047_4},
		},
		{
			prefixer:       EBCDIC1047.LLLL,
			maxLen:         8192,
			dataLen:        5432,
			expectedOutput: []byte{ebcdic1047_5, ebcdic1047_4, ebcdic1047_3, ebcdic1047_2},
		},
	} {
		encoded, err := testCase.prefixer.EncodeLength(testCase.maxLen, testCase.dataLen)
		require.NoError(t, err)
		require.Equal(t, testCase.expectedOutput, encoded)
	}
}

func TestEBCDIC1047PrefixersEncodeErrors(t *testing.T) {
	t.Parallel()

	t.Run("data longer than maximum allowed length", func(t *testing.T) {
		t.Parallel()
		for _, testCase := range []struct {
			prefixer      Prefixer
			maxLen        int
			dataLen       int
			expectedError string
		}{
			{
				prefixer:      EBCDIC1047.L,
				maxLen:        8,
				dataLen:       9,
				expectedError: "field length [9] is larger than maximum [8]",
			},
			{
				prefixer:      EBCDIC1047.LL,
				maxLen:        52,
				dataLen:       73,
				expectedError: "field length [73] is larger than maximum [52]",
			},
			{
				prefixer:      EBCDIC1047.LLL,
				maxLen:        512,
				dataLen:       999,
				expectedError: "field length [999] is larger than maximum [512]",
			},
			{
				prefixer:      EBCDIC1047.LLLL,
				maxLen:        1024,
				dataLen:       2048,
				expectedError: "field length [2048] is larger than maximum [1024]",
			},
		} {
			encoded, err := testCase.prefixer.EncodeLength(testCase.maxLen, testCase.dataLen)
			require.Nil(t, encoded)
			require.EqualError(t, err, testCase.expectedError)
		}
	})

	t.Run("length has too many digits", func(t *testing.T) {
		// N.B. this error case should never be reached, as the maxLen of any field should be
		// within the bounds set on the number of digits required to represent this number. This is
		// only possible in incorrectly defined schemes.
		t.Parallel()
		for _, testCase := range []struct {
			prefixer      Prefixer
			maxLen        int
			dataLen       int
			expectedError string
		}{
			{
				prefixer:      EBCDIC1047.L,
				maxLen:        52,
				dataLen:       10,
				expectedError: "number of digits in data [10] exceeds its maximum indicator [1]",
			},
			{
				prefixer:      EBCDIC1047.LL,
				maxLen:        101,
				dataLen:       100,
				expectedError: "number of digits in data [100] exceeds its maximum indicator [2]",
			},
			{
				prefixer:      EBCDIC1047.LLL,
				maxLen:        1333,
				dataLen:       1001,
				expectedError: "number of digits in data [1001] exceeds its maximum indicator [3]",
			},
			{
				prefixer:      EBCDIC1047.LLLL,
				maxLen:        11111,
				dataLen:       10908,
				expectedError: "number of digits in data [10908] exceeds its maximum indicator [4]",
			},
		} {
			encoded, err := testCase.prefixer.EncodeLength(testCase.maxLen, testCase.dataLen)
			require.Nil(t, encoded)
			require.EqualError(t, err, testCase.expectedError)
		}
	})

	t.Run("fixed length error", func(t *testing.T) {
		t.Parallel()
		encoded, err := EBCDIC1047.Fixed.EncodeLength(128, 127)
		require.Nil(t, encoded)
		require.EqualError(t, err, "field length [127] should be fixed [128]")
	})
}

func TestEBCDIC1047PrefixersDecode(t *testing.T) {
	t.Parallel()
	for _, testCase := range []struct {
		prefixer       Prefixer
		maxLen         int
		data           []byte
		expectedOutput int
		expectedRead   int
	}{
		{
			prefixer:       EBCDIC1047.Fixed,
			maxLen:         64,
			data:           []byte{},
			expectedOutput: 64,
			expectedRead:   0,
		},
		{
			prefixer:       EBCDIC1047.L,
			maxLen:         8,
			data:           []byte{ebcdic1047_7},
			expectedOutput: 7,
			expectedRead:   1,
		},
		{
			prefixer:       EBCDIC1047.LL,
			maxLen:         14,
			data:           []byte{ebcdic1047_0, ebcdic1047_1},
			expectedOutput: 1,
			expectedRead:   2,
		},
		{
			prefixer:       EBCDIC1047.LL,
			maxLen:         32,
			data:           []byte{ebcdic1047_2, ebcdic1047_8},
			expectedOutput: 28,
			expectedRead:   2,
		},
		{
			prefixer:       EBCDIC1047.LLL,
			maxLen:         128,
			data:           []byte{ebcdic1047_0, ebcdic1047_0, ebcdic1047_7},
			expectedOutput: 7,
			expectedRead:   3,
		},
		{
			prefixer:       EBCDIC1047.LLL,
			maxLen:         128,
			data:           []byte{ebcdic1047_0, ebcdic1047_6, ebcdic1047_2},
			expectedOutput: 62,
			expectedRead:   3,
		},
		{
			prefixer:       EBCDIC1047.LLL,
			maxLen:         512,
			data:           []byte{ebcdic1047_2, ebcdic1047_0, ebcdic1047_2},
			expectedOutput: 202,
			expectedRead:   3,
		},
		{
			prefixer:       EBCDIC1047.LLLL,
			maxLen:         1024,
			data:           []byte{ebcdic1047_0, ebcdic1047_0, ebcdic1047_0, ebcdic1047_1},
			expectedOutput: 1,
			expectedRead:   4,
		},
		{
			prefixer:       EBCDIC1047.LLLL,
			maxLen:         2048,
			data:           []byte{ebcdic1047_0, ebcdic1047_0, ebcdic1047_9, ebcdic1047_1},
			expectedOutput: 91,
			expectedRead:   4,
		},
		{
			prefixer:       EBCDIC1047.LLLL,
			maxLen:         4096,
			data:           []byte{ebcdic1047_0, ebcdic1047_7, ebcdic1047_7, ebcdic1047_7},
			expectedOutput: 777,
			expectedRead:   4,
		},
		{
			prefixer:       EBCDIC1047.LLLL,
			maxLen:         8192,
			data:           []byte{ebcdic1047_1, ebcdic1047_9, ebcdic1047_9, ebcdic1047_3},
			expectedOutput: 1993,
			expectedRead:   4,
		},
	} {
		length, read, err := testCase.prefixer.DecodeLength(testCase.maxLen, testCase.data)
		require.NoError(t, err)
		require.Equal(t, testCase.expectedOutput, length)
		require.Equal(t, testCase.expectedRead, read)
	}
}

func TestEBCDIC1047PrefixersDecodeErrors(t *testing.T) {
	t.Parallel()

	t.Run("insufficient data to read length", func(t *testing.T) {
		t.Parallel()
		for _, testCase := range []struct {
			prefixer      Prefixer
			maxLen        int
			data          []byte
			expectedError string
		}{
			{
				prefixer:      EBCDIC1047.L,
				maxLen:        8,
				data:          []byte{},
				expectedError: "not enough data length [0] to read [1] byte digits",
			},
			{
				prefixer:      EBCDIC1047.LL,
				maxLen:        16,
				data:          []byte{ebcdic1047_8},
				expectedError: "not enough data length [1] to read [2] byte digits",
			},
			{
				prefixer:      EBCDIC1047.LLL,
				maxLen:        128,
				data:          []byte{ebcdic1047_0, ebcdic1047_0},
				expectedError: "not enough data length [2] to read [3] byte digits",
			},
			{
				prefixer:      EBCDIC1047.LLLL,
				maxLen:        8,
				data:          []byte{ebcdic1047_0, ebcdic1047_0, ebcdic1047_9},
				expectedError: "not enough data length [3] to read [4] byte digits",
			},
		} {
			length, read, err := testCase.prefixer.DecodeLength(testCase.maxLen, testCase.data)
			require.Zero(t, length)
			require.Zero(t, read)
			require.EqualError(t, err, testCase.expectedError)
		}
	})

	t.Run("data length to large", func(t *testing.T) {
		t.Parallel()
		for _, testCase := range []struct {
			prefixer      Prefixer
			maxLen        int
			data          []byte
			expectedError string
		}{
			{
				prefixer:      EBCDIC1047.L,
				maxLen:        8,
				data:          []byte{ebcdic1047_9},
				expectedError: "data length [9] is larger than maximum [8]",
			},
			{
				prefixer:      EBCDIC1047.LL,
				maxLen:        16,
				data:          []byte{ebcdic1047_2, ebcdic1047_0},
				expectedError: "data length [20] is larger than maximum [16]",
			},
			{
				prefixer:      EBCDIC1047.LLL,
				maxLen:        128,
				data:          []byte{ebcdic1047_1, ebcdic1047_9, ebcdic1047_4},
				expectedError: "data length [194] is larger than maximum [128]",
			},
			{
				prefixer:      EBCDIC1047.LLLL,
				maxLen:        8000,
				data:          []byte{ebcdic1047_8, ebcdic1047_0, ebcdic1047_9, ebcdic1047_2},
				expectedError: "data length [8092] is larger than maximum [8000]",
			},
		} {
			length, read, err := testCase.prefixer.DecodeLength(testCase.maxLen, testCase.data)
			require.Zero(t, length)
			require.Zero(t, read)
			require.EqualError(t, err, testCase.expectedError)
		}
	})

	t.Run("non-numeric length", func(t *testing.T) {
		t.Parallel()
		for _, testCase := range []struct {
			prefixer      Prefixer
			maxLen        int
			data          []byte
			expectedError string
		}{
			{
				prefixer:      EBCDIC1047.L,
				maxLen:        9,
				data:          []byte{ebcdic1047Plus},
				expectedError: "length [+] is not a valid integer length field",
			},
			{
				prefixer:      EBCDIC1047.LL,
				maxLen:        13,
				data:          []byte{ebcdic1047l, ebcdic1047_3},
				expectedError: "length [l3] is not a valid integer length field",
			},
			{
				prefixer:      EBCDIC1047.LLL,
				maxLen:        128,
				data:          []byte{ebcdic1047_1, ebcdic1047Dot, ebcdic1047_0},
				expectedError: "length [1.0] is not a valid integer length field",
			},
			{
				prefixer:      EBCDIC1047.LLLL,
				maxLen:        9999,
				data:          []byte{ebcdic1047l, ebcdic1047l, ebcdic1047l, ebcdic1047l},
				expectedError: "length [llll] is not a valid integer length field",
			},
		} {
			length, read, err := testCase.prefixer.DecodeLength(testCase.maxLen, testCase.data)
			require.Zero(t, length)
			require.Zero(t, read)
			require.EqualError(t, err, testCase.expectedError)
		}
	})
}

func TestEBCDIC1047PrefixersInspect(t *testing.T) {
	t.Parallel()
	require.Equal(t, "EBCDIC.Fixed", EBCDIC1047.Fixed.Inspect())
	require.Equal(t, "EBCDIC.L", EBCDIC1047.L.Inspect())
	require.Equal(t, "EBCDIC.LL", EBCDIC1047.LL.Inspect())
	require.Equal(t, "EBCDIC.LLL", EBCDIC1047.LLL.Inspect())
	require.Equal(t, "EBCDIC.LLLL", EBCDIC1047.LLLL.Inspect())
}
