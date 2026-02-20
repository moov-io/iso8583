package iso8583

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

// unknownTagsSpec is a minimal MessageSpec with field 55 as a BER-TLV
// composite that skips unknown tags.
var unknownTagsSpec = &MessageSpec{
	Name: "Test Spec – Unknown Tags",
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Length:      8,
			Description: "Bitmap",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.Fixed,
		}),
		// Field 55: ICC / EMV data — BER-TLV with skip-unknown enabled
		55: field.NewComposite(&field.Spec{
			Length:      999,
			Description: "ICC Data",
			Pref:        prefix.ASCII.LLL,
			Tag: &field.TagSpec{
				Enc:                encoding.BerTLVTag,
				Sort:               sort.StringsByHex,
				SkipUnknownTLVTags: true,
				AuditUnknownTLVTags: true,
			},
			Subfields: map[string]field.Field{
				"9A": field.NewHex(&field.Spec{
					Description: "Transaction Date",
					Enc:         encoding.Binary,
					Pref:        prefix.BerTLV,
				}),
				"9F02": field.NewHex(&field.Spec{
					Description: "Amount, Authorized",
					Enc:         encoding.Binary,
					Pref:        prefix.BerTLV,
				}),
			},
		}),
	},
}

// nestedUnknownTagsSpec adds a second composite field (field 56) whose
// subfields are themselves composites, so we can verify the callback chain
// builds paths like "56.01.9F36".
var nestedUnknownTagsSpec = &MessageSpec{
	Name: "Test Spec – Nested Unknown Tags",
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Length:      8,
			Description: "Bitmap",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.Fixed,
		}),
		// Field 56: outer composite whose subfield "01" is itself a TLV composite
		56: field.NewComposite(&field.Spec{
			Length:      999,
			Description: "Nested Composite",
			Pref:        prefix.ASCII.LLL,
			Tag: &field.TagSpec{
				Length:              2,
				Enc:                 encoding.ASCII,
				Sort:                sort.StringsByInt,
				AuditUnknownTLVTags: true,
			},
			Subfields: map[string]field.Field{
				"01": field.NewComposite(&field.Spec{
					Length:      999,
					Description: "Inner TLV",
					Pref:        prefix.ASCII.LLL,
					Tag: &field.TagSpec{
						Enc:                encoding.BerTLVTag,
						Sort:               sort.StringsByHex,
						SkipUnknownTLVTags: true,
						AuditUnknownTLVTags: true,
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
		}),
	},
}

func TestMessageUnknownTags(t *testing.T) {
	t.Run("returns full paths for unknown TLV tags in a top-level composite", func(t *testing.T) {
		// Bit 55 is in byte index 6 of the 8-byte bitmap (0-indexed).
		// bit 55 → byte 6, bit offset 6 from MSB → 0b00000010 = 0x02
		rawMsg := []byte(
			"0200", // MTI
		)
		rawMsg = append(rawMsg,
			// Bitmap: only bit 55 set
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,
		)
		// Field 55 content: length "026" + TLV bytes
		rawMsg = append(rawMsg, []byte("026")...)
		rawMsg = append(rawMsg,
			0x9a, 0x03, 0x21, 0x07, 0x20,                         // 9A  (known)
			0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, // 9F02 (known)
			0x9f, 0x36, 0x02, 0x01, 0x57,                         // 9F36 (unknown)
			0x9f, 0x37, 0x04, 0x9b, 0xad, 0xbc, 0xab,             // 9F37 (unknown)
		)

		msg := NewMessage(unknownTagsSpec)
		require.NoError(t, msg.Unpack(rawMsg))

		unknown := msg.UnknownTags()
		require.ElementsMatch(t, []string{"55.9F36", "55.9F37"}, unknown)
	})

	t.Run("returns empty slice when no unknown tags are present", func(t *testing.T) {
		rawMsg := []byte("0200")
		rawMsg = append(rawMsg, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00)
		rawMsg = append(rawMsg, []byte("014")...)
		rawMsg = append(rawMsg,
			0x9a, 0x03, 0x21, 0x07, 0x20,
			0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01,
		)

		msg := NewMessage(unknownTagsSpec)
		require.NoError(t, msg.Unpack(rawMsg))

		require.Empty(t, msg.UnknownTags())
	})

	t.Run("resets unknown tags on re-unpack", func(t *testing.T) {
		withUnknown := []byte("0200")
		withUnknown = append(withUnknown, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00)
		withUnknown = append(withUnknown, []byte("026")...)
		withUnknown = append(withUnknown,
			0x9a, 0x03, 0x21, 0x07, 0x20,
			0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01,
			0x9f, 0x36, 0x02, 0x01, 0x57,
			0x9f, 0x37, 0x04, 0x9b, 0xad, 0xbc, 0xab,
		)

		msg := NewMessage(unknownTagsSpec)
		require.NoError(t, msg.Unpack(withUnknown))
		require.Len(t, msg.UnknownTags(), 2)

		// Re-unpack with only known tags — unknown tags list must be cleared
		onlyKnown := []byte("0200")
		onlyKnown = append(onlyKnown, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00)
		onlyKnown = append(onlyKnown, []byte("014")...)
		onlyKnown = append(onlyKnown,
			0x9a, 0x03, 0x21, 0x07, 0x20,
			0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01,
		)
		require.NoError(t, msg.Unpack(onlyKnown))
		require.Empty(t, msg.UnknownTags())
	})

	t.Run("builds full path for unknown tags in nested composites", func(t *testing.T) {
		// Bit 56 → byte index 6 of the 8-byte bitmap, LSB → 0x01.
		//
		// Field 56 layout (outer composite, positional tag "01"):
		//   tag "01" (2 ASCII bytes) + inner composite content
		//   inner composite: LLL prefix + TLV: 9A (known) + 9F36 (unknown)
		innerTLV := []byte{
			0x9a, 0x03, 0x21, 0x07, 0x20, // 9A  (known)   — 5 bytes
			0x9f, 0x36, 0x02, 0x01, 0x57, // 9F36 (unknown) — 5 bytes
		}
		innerField := append([]byte(fmt.Sprintf("%03d", len(innerTLV))), innerTLV...)

		outerContent := append([]byte("01"), innerField...)
		outerField := append([]byte(fmt.Sprintf("%03d", len(outerContent))), outerContent...)

		rawMsg := []byte("0200")
		rawMsg = append(rawMsg, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00)
		rawMsg = append(rawMsg, outerField...)

		msg := NewMessage(nestedUnknownTagsSpec)
		require.NoError(t, msg.Unpack(rawMsg))

		unknown := msg.UnknownTags()
		require.ElementsMatch(t, []string{"56.01.9F36"}, unknown)
	})
}
