package field

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

func TestSkipUnknownTLVTagsMalformedLength(t *testing.T) {
	spec := &Spec{
		Length:      999,
		Description: "TLV skip-only",
		Pref:        prefix.ASCII.LLL,
		Tag: &TagSpec{
			Enc:                 encoding.BerTLVTag,
			Sort:                sort.StringsByHex,
			SkipUnknownTLVTags:  true,
			StoreUnknownTLVTags: false,
		},
		Subfields: map[string]Field{
			"9F02": NewHex(&Spec{
				Description: "Amount, Authorized (Numeric)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
		},
	}

	composite := NewComposite(spec)

	// 8A is an unknown tag. Its BER-TLV length uses the long form with 8
	// length bytes encoding 2^63, which decodes to a negative int.
	inputData := []byte{
		0x30, 0x31, 0x30, // ASCII "010" length prefix
		0x8A,                                           // unknown tag (1 byte)
		0x88,                                           // long form: 8 length bytes follow
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 2^63
	}

	_, err := composite.Unpack(inputData)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not enough data to unpack unknown TLV tag")
}
