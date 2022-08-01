package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// EBCDIC Code Page 1047 Encodings, taken from https://www.ibm.com/docs/en/personal-communications/5.9?topic=pages-1047103-latin-open-systems
var ebcdic1047CharacterEncodings = map[string]byte{
	"â":  0x42,
	"ä":  0x43,
	"à":  0x44,
	"á":  0x45,
	"ã":  0x46,
	"å":  0x47,
	"ç":  0x48,
	"ñ":  0x49,
	"¢":  0x4A,
	".":  0x4B,
	"<":  0x4C,
	"(":  0x4D,
	"+":  0x4E,
	"|":  0x4F,
	"&":  0x50,
	"é":  0x51,
	"ê":  0x52,
	"ë":  0x53,
	"è":  0x54,
	"í":  0x55,
	"î":  0x56,
	"ï":  0x57,
	"ì":  0x58,
	"ß":  0x59,
	"!":  0x5A,
	"$":  0x5B,
	"*":  0x5C,
	")":  0x5D,
	";":  0x5E,
	"^":  0x5F,
	"-":  0x60,
	"/":  0x61,
	"Â":  0x62,
	"Ä":  0x63,
	"À":  0x64,
	"Á":  0x65,
	"Ã":  0x66,
	"Å":  0x67,
	"Ç":  0x68,
	"Ñ":  0x69,
	"¦":  0x6A,
	",":  0x6B,
	"%":  0x6C,
	"_":  0x6D,
	">":  0x6E,
	"?":  0x6F,
	"ø":  0x70,
	"É":  0x71,
	"Ê":  0x72,
	"Ë":  0x73,
	"È":  0x74,
	"Í":  0x75,
	"Î":  0x76,
	"Ï":  0x77,
	"Ì":  0x78,
	"`":  0x79,
	":":  0x7A,
	"#":  0x7B,
	"@":  0x7C,
	"'":  0x7D,
	"=":  0x7E,
	"\"": 0x7F,
	"Ø":  0x80,
	"a":  0x81,
	"b":  0x82,
	"c":  0x83,
	"d":  0x84,
	"e":  0x85,
	"f":  0x86,
	"g":  0x87,
	"h":  0x88,
	"i":  0x89,
	"«":  0x8A,
	"»":  0x8B,
	"ð":  0x8C,
	"ý":  0x8D,
	"þ":  0x8E,
	"±":  0x8F,
	"°":  0x90,
	"j":  0x91,
	"k":  0x92,
	"l":  0x93,
	"m":  0x94,
	"n":  0x95,
	"o":  0x96,
	"p":  0x97,
	"q":  0x98,
	"r":  0x99,
	"ª":  0x9A,
	"º":  0x9B,
	"æ":  0x9C,
	"¸":  0x9D,
	"Æ":  0x9E,
	"¤":  0x9F,
	"µ":  0xA0,
	"~":  0xA1,
	"s":  0xA2,
	"t":  0xA3,
	"u":  0xA4,
	"v":  0xA5,
	"w":  0xA6,
	"x":  0xA7,
	"y":  0xA8,
	"z":  0xA9,
	"¡":  0xAA,
	"¿":  0xAB,
	"Ð":  0xAC,
	"[":  0xAD,
	"Þ":  0xAE,
	"®":  0xAF,
	"¬":  0xB0,
	"£":  0xB1,
	"¥":  0xB2,
	"·":  0xB3,
	"©":  0xB4,
	"§":  0xB5,
	"¶":  0xB6,
	"¼":  0xB7,
	"½":  0xB8,
	"¾":  0xB9,
	"Ý":  0xBA,
	"¨":  0xBB,
	"¯":  0xBC,
	"]":  0xBD,
	"´":  0xBE,
	"×":  0xBF,
	"{":  0xC0,
	"A":  0xC1,
	"B":  0xC2,
	"C":  0xC3,
	"D":  0xC4,
	"E":  0xC5,
	"F":  0xC6,
	"G":  0xC7,
	"H":  0xC8,
	"I":  0xC9,
	"ô":  0xCB,
	"ö":  0xCC,
	"ò":  0xCD,
	"ó":  0xCE,
	"õ":  0xCF,
	"}":  0xD0,
	"J":  0xD1,
	"K":  0xD2,
	"L":  0xD3,
	"M":  0xD4,
	"N":  0xD5,
	"O":  0xD6,
	"P":  0xD7,
	"Q":  0xD8,
	"R":  0xD9,
	"¹":  0xDA,
	"û":  0xDB,
	"ü":  0xDC,
	"ù":  0xDD,
	"ú":  0xDE,
	"ÿ":  0xDF,
	"\\": 0xE0,
	"÷":  0xE1,
	"S":  0xE2,
	"T":  0xE3,
	"U":  0xE4,
	"V":  0xE5,
	"W":  0xE6,
	"X":  0xE7,
	"Y":  0xE8,
	"Z":  0xE9,
	"²":  0xEA,
	"Ô":  0xEB,
	"Ö":  0xEC,
	"Ò":  0xED,
	"Ó":  0xEE,
	"Õ":  0xEF,
	"0":  0xF0,
	"1":  0xF1,
	"2":  0xF2,
	"3":  0xF3,
	"4":  0xF4,
	"5":  0xF5,
	"6":  0xF6,
	"7":  0xF7,
	"8":  0xF8,
	"9":  0xF9,
	"³":  0xFA,
	"Û":  0xFB,
	"Ü":  0xFC,
	"Ù":  0xFD,
	"Ú":  0xFE,
}

// some randomly-chosen phrases with some interesting characters
var knownEncodings = []struct {
	Phrase   string
	Encoding []byte
}{
	{
		Phrase:   "hello, world!",
		Encoding: []byte{0x88, 0x85, 0x93, 0x93, 0x96, 0x6B, 0x40, 0xA6, 0x96, 0x99, 0x93, 0x84, 0x5A},
	},
	{
		Phrase:   "¿Cómo estás?",
		Encoding: []byte{0xAB, 0xC3, 0xCE, 0x94, 0x96, 0x40, 0x85, 0xA2, 0xA3, 0x45, 0xA2, 0x6F},
	},
	{
		Phrase:   "Ágætis byrjun",
		Encoding: []byte{0x65, 0x87, 0x9C, 0xA3, 0x89, 0xA2, 0x40, 0x82, 0xA8, 0x99, 0x91, 0xA4, 0x95},
	},
}

func TestEBCDIC1047SingleCharacterEncode(t *testing.T) {
	t.Parallel()
	for character, expectedEncoding := range ebcdic1047CharacterEncodings {
		encoding, err := EBCDIC1047.Encode([]byte(character))
		require.NoError(t, err)
		require.Len(t, encoding, 1)
		require.Equal(t, expectedEncoding, encoding[0])
	}
}

func TestEBCDIC1047SingleCharacterDecode(t *testing.T) {
	t.Parallel()
	for expectedCharacter, byteChar := range ebcdic1047CharacterEncodings {
		decoding, length, err := EBCDIC1047.Decode([]byte{byteChar}, 1)
		require.NoError(t, err)
		require.Equal(t, 1, length)
		require.Equal(t, expectedCharacter, string(decoding))
	}
}

func TestEBCDIC1047Encode(t *testing.T) {
	t.Parallel()
	for _, testCase := range knownEncodings {
		encoding, err := EBCDIC1047.Encode([]byte(testCase.Phrase))
		require.NoError(t, err)
		require.Equal(t, testCase.Encoding, encoding)
	}
}

func TestEBCDIC1047Decode(t *testing.T) {
	t.Parallel()

	t.Run("errors on invalid data length", func(t *testing.T) {
		t.Parallel()

		decoding, length, err := EBCDIC1047.Decode([]byte("test"), 5)
		require.Nil(t, decoding)
		require.Zero(t, length)
		require.EqualError(t, err, "not enough data to decode. expected len 5, got 4")
	})

	t.Run("decode whole string", func(t *testing.T) {
		t.Parallel()
		for _, testCase := range knownEncodings {
			decoding, length, err := EBCDIC1047.Decode(testCase.Encoding, len(testCase.Encoding))
			require.NoError(t, err)
			require.Equal(t, len(testCase.Encoding), length)
			require.Equal(t, testCase.Phrase, string(decoding))
		}
	})

	t.Run("decode partial string", func(t *testing.T) {
		t.Parallel()
		messageToDecode := []byte{0xF6, 0x40, 0xBF, 0x40, 0xF9, 0x40, 0x7E, 0x40, 0xF4, 0xF2} // 6 × 9 = 42
		lengthToDecode := 7
		expectedPartialMessage := "6 × 9 =" // first 7 characters
		decoding, length, err := EBCDIC1047.Decode(messageToDecode, lengthToDecode)
		require.NoError(t, err)
		require.Equal(t, lengthToDecode, length)
		require.Equal(t, expectedPartialMessage, string(decoding))
	})
}
