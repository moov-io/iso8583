package iso8583

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

func TestUnknownTags(t *testing.T) {
	spec := &MessageSpec{
		Fields: map[int]field.Field{
			0: field.NewString(&field.Spec{
				Length:      4,
				Description: "Message Type Indicator",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			1: field.NewBitmap(&field.Spec{
				Description: "Bitmap",
				Enc:         encoding.BytesToASCIIHex,
				Pref:        prefix.Hex.Fixed,
			}),
			3: field.NewComposite(&field.Spec{
				Length:      999,
				Description: "Processing Code",
				Pref:        prefix.ASCII.LLL,
				Tag: &field.TagSpec{
					Sort: sort.StringsByInt,
				},
				Subfields: map[string]field.Field{
					"1": field.NewString(&field.Spec{
						Length:      2,
						Description: "Transaction Type",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
					"2": field.NewComposite(&field.Spec{
						Length:      100,
						Description: "Nested TLV",
						Pref:        prefix.ASCII.LLL,
						Tag: &field.TagSpec{
							Enc:                 encoding.BerTLVTag,
							Sort:                sort.StringsByHex,
							SkipUnknownTLVTags:  true,
							StoreUnknownTLVTags: true,
						},
						Subfields: map[string]field.Field{
							"9A": field.NewHex(&field.Spec{
								Description: "Transaction Date",
								Enc:         encoding.Binary,
								Pref:        prefix.BerTLV,
							}),
							"9F02": field.NewHex(&field.Spec{
								Description: "Amount, Authorized (Numeric)",
								Enc:         encoding.Binary,
								Pref:        prefix.BerTLV,
							}),
						},
					}),
				},
			}),
		},
	}

	t.Run("returns unknown tags from nested composites", func(t *testing.T) {
		msg := NewMessage(spec)

		// Build data with known + unknown TLV tags inside field 3.2
		data := []byte("01002000000000000000")
		// LLL for composite field 3
		data = append(data, []byte("031")...)
		// Subfield 1 value
		data = append(data, []byte("00")...)
		// LLL for TLV field 3.2
		data = append(data, []byte("026")...)
		// Known tags
		data = append(data, 0x9a, 0x03, 0x21, 0x07, 0x20)                         // 9A
		data = append(data, 0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01) // 9F02
		// Unknown tags
		data = append(data, 0x9f, 0x36, 0x02, 0x01, 0x57)                   // 9F36
		data = append(data, 0x9f, 0x37, 0x04, 0x9b, 0xad, 0xbc, 0xab)       // 9F37

		err := msg.Unpack(data)
		require.NoError(t, err)

		unknownTags := UnknownTags(msg)
		require.Len(t, unknownTags, 2)
		require.Contains(t, unknownTags, "3.2.9F36")
		require.Contains(t, unknownTags, "3.2.9F37")

		// verify we can inspect the field data
		f9F36Bytes, err := unknownTags["3.2.9F36"].Bytes()
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x57}, f9F36Bytes)
	})

	t.Run("returns empty map when no unknown tags", func(t *testing.T) {
		msg := NewMessage(spec)

		// Build data with only known tags
		data := []byte("01002000000000000000")
		// LLL for composite field 3
		data = append(data, []byte("019")...)
		// Subfield 1 value
		data = append(data, []byte("00")...)
		// LLL for TLV field 3.2
		data = append(data, []byte("014")...)
		// Only known tags
		data = append(data, 0x9a, 0x03, 0x21, 0x07, 0x20)                         // 9A
		data = append(data, 0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01) // 9F02

		err := msg.Unpack(data)
		require.NoError(t, err)

		unknownTags := UnknownTags(msg)
		require.Empty(t, unknownTags)
	})

	t.Run("returns empty when StoreUnknownTLVTags is disabled", func(t *testing.T) {
		noStoreSpec := &MessageSpec{
			Fields: map[int]field.Field{
				0: field.NewString(&field.Spec{
					Length:      4,
					Description: "Message Type Indicator",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				1: field.NewBitmap(&field.Spec{
					Description: "Bitmap",
					Enc:         encoding.BytesToASCIIHex,
					Pref:        prefix.Hex.Fixed,
				}),
				3: field.NewComposite(&field.Spec{
					Length:      999,
					Description: "Processing Code",
					Pref:        prefix.ASCII.LLL,
					Tag: &field.TagSpec{
						Enc:                 encoding.BerTLVTag,
						Sort:                sort.StringsByHex,
						SkipUnknownTLVTags:  true,
						StoreUnknownTLVTags: false,
					},
					Subfields: map[string]field.Field{
						"9A": field.NewHex(&field.Spec{
							Description: "Transaction Date",
							Enc:         encoding.Binary,
							Pref:        prefix.BerTLV,
						}),
					},
				}),
			},
		}

		msg := NewMessage(noStoreSpec)

		data := []byte("01002000000000000000")
		// LLL for field 3
		data = append(data, []byte("017")...)
		// Known tag
		data = append(data, 0x9a, 0x03, 0x21, 0x07, 0x20) // 9A
		// Unknown tags - skipped but not stored
		data = append(data, 0x9f, 0x36, 0x02, 0x01, 0x57)             // 9F36
		data = append(data, 0x9f, 0x37, 0x04, 0x9b, 0xad, 0xbc, 0xab) // 9F37

		err := msg.Unpack(data)
		require.NoError(t, err)

		// Unknown tags are not stored, so UnknownTags returns empty
		unknownTags := UnknownTags(msg)
		require.Empty(t, unknownTags)
	})
}
