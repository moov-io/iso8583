package examples

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

//	The sample is for VSDC chip data usage

//		This field 55 VSDC chip data usage contains three subfields after the length subfield.
//		Positions:
//			1 2 3 4 ... 255
//	 Fields
//		 - Subfield 1: length Byte, a one-byte binary subfield  that contains the number of bytes in this field after the length subfield
//		 - Subfield 2: dataset ID, a one-byte binary identifier
//		 - Subfield 3: dataset length, 2-byte binary subfield that contains the total length of all TLV elements that follow.
//		 - Subfield 4:
//			Chip Card TLV data elements
//			Tag Length Value Tag Length Value
func TestICCField55(t *testing.T) {
	field55Spec := &field.Spec{
		Length:      255,
		Description: "Integrated Circuit Card (ICC) Data",
		Pref:        prefix.Binary.L,
		Tag: &field.TagSpec{
			// We have 1 byte length tag, that in the spec is seen as a Hex string
			// but will be encoded as a binary byte (ASCIIHexToBytes)
			// We sort the TLV tags by their hex values, but it's not important
			// Finally, if we have unknown tag, in order to skip it, we need to know
			// how long its value is. For this, we need to read the length tag.
			// To read the length tag, we need to know its length.
			// Setting PrefUnknownTLV to prefix.Binary.LL will read the length prefix
			// as a 2-byte binary value.
			Length:             1,
			Enc:                encoding.ASCIIHexToBytes,
			Sort:               sort.StringsByHex,
			SkipUnknownTLVTags: true,
			// if we have unknown TLV tags,
			PrefUnknownTLV: prefix.Binary.LL,
		},
		Subfields: map[string]field.Field{
			"01": field.NewComposite(&field.Spec{
				Length:      252,
				Description: "VSDC Data",
				Pref:        prefix.Binary.LL,
				Tag: &field.TagSpec{
					Enc:                encoding.BerTLVTag,
					Sort:               sort.StringsByHex,
					SkipUnknownTLVTags: true,
				},
				Subfields: map[string]field.Field{
					"9A": field.NewString(&field.Spec{
						Description: "Transaction Date",
						Enc:         encoding.Binary,
						Pref:        prefix.BerTLV,
					}),
					"9F02": field.NewString(&field.Spec{
						Description: "Amount, Authorized (Numeric)",
						Enc:         encoding.Binary,
						Pref:        prefix.BerTLV,
					}),
				},
			}),
		},
	}

	type VSDCData struct {
		TransactionDate string `iso8583:"9A"`
		Amount          string `iso8583:"9F02"`
	}

	type ICCData struct {
		VSDCData *VSDCData `iso8583:"01"`
	}

	filed55 := field.NewComposite(field55Spec)
	err := filed55.Marshal(&ICCData{
		VSDCData: &VSDCData{
			TransactionDate: "210720",
			Amount:          "000000000501",
		},
	})
	require.NoError(t, err)

	packed, err := filed55.Pack()
	require.NoError(t, err)

	require.Equal(t, "1A0100179A063231303732309F020C303030303030303030353031", fmt.Sprintf("%X", packed))
}
