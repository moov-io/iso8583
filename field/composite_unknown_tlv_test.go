package field

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

// Test spec for BER-TLV with StoreUnknownTLVTags enabled
var storeUnknownTLVSpec = &Spec{
	Length:      999,
	Description: "ICC Data – EMV Having Multiple Tags with Unknown Tag Storage",
	Pref:        prefix.ASCII.LLL,
	Tag: &TagSpec{
		Enc:                 encoding.BerTLVTag,
		Sort:                sort.StringsByHex,
		SkipUnknownTLVTags:  true,
		StoreUnknownTLVTags: true,
	},
	Subfields: map[string]Field{
		"9A": NewHex(&Spec{
			Description: "Transaction Date",
			Enc:         encoding.Binary,
			Pref:        prefix.BerTLV,
		}),
		"9F02": NewHex(&Spec{
			Description: "Amount, Authorized (Numeric)",
			Enc:         encoding.Binary,
			Pref:        prefix.BerTLV,
		}),
	},
}

// Data struct that includes unknown TLV fields
type TLVDataWithUnknown struct {
	F9A   *Hex `index:"9A"`
	F9F02 *Hex `index:"9F02"`
	// Unknown fields - user defines these to capture specific unknown tags
	F9F36 *Binary `index:"9F36"`
	F9F37 *Binary `index:"9F37"`
}

func TestStoreUnknownTLVTags(t *testing.T) {
	t.Run("Unpack and Pack preserves unknown TLV tags", func(t *testing.T) {
		composite := NewComposite(storeUnknownTLVSpec)

		// Data contains:
		// - 9A (known): 3 bytes [0x21, 0x07, 0x20]
		// - 9F02 (known): 6 bytes [0x00, 0x00, 0x00, 0x00, 0x05, 0x01]
		// - 9F36 (unknown): 2 bytes [0x01, 0x57]
		// - 9F37 (unknown): 4 bytes [0x9b, 0xad, 0xbc, 0xab]
		inputData := []byte{
			0x30, 0x32, 0x36, // ASCII "026" - length prefix
			0x9a, 0x03, 0x21, 0x07, 0x20, // 9A: length 3, value 210720
			0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, // 9F02: length 6, value 000000000501
			0x9f, 0x36, 0x02, 0x01, 0x57, // 9F36: length 2, value 0157
			0x9f, 0x37, 0x04, 0x9b, 0xad, 0xbc, 0xab, // 9F37: length 4, value 9badbcab
		}

		read, err := composite.Unpack(inputData)
		require.NoError(t, err)
		require.Equal(t, len(inputData), read)

		// Verify all fields are stored (including unknown)
		subfields := composite.GetSubfields()
		require.Len(t, subfields, 4)
		require.Contains(t, subfields, "9A")
		require.Contains(t, subfields, "9F02")
		require.Contains(t, subfields, "9F36")
		require.Contains(t, subfields, "9F37")

		packed, err := composite.Pack()
		require.NoError(t, err)

		// We expect the packed data to match the original input
		require.Equal(t, inputData, packed)
	})

	t.Run("Unpack and Unmarshal unknown TLV tags to data struct", func(t *testing.T) {
		composite := NewComposite(storeUnknownTLVSpec)

		// Data contains:
		// - 9A (known): 3 bytes [0x21, 0x07, 0x20]
		// - 9F02 (known): 6 bytes [0x00, 0x00, 0x00, 0x00, 0x05, 0x01]
		// - 9F36 (unknown): 2 bytes [0x01, 0x57]
		// - 9F37 (unknown): 4 bytes [0x9b, 0xad, 0xbc, 0xab]
		inputData := []byte{
			0x30, 0x32, 0x36, // ASCII "026" - length prefix
			0x9a, 0x03, 0x21, 0x07, 0x20, // 9A: length 3, value 210720
			0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, // 9F02: length 6, value 000000000501
			0x9f, 0x36, 0x02, 0x01, 0x57, // 9F36: length 2, value 0157
			0x9f, 0x37, 0x04, 0x9b, 0xad, 0xbc, 0xab, // 9F37: length 4, value 9badbcab
		}

		_, err := composite.Unpack(inputData)
		require.NoError(t, err)

		// Unmarshal to data struct
		data := &TLVDataWithUnknown{}
		err = composite.Unmarshal(data)
		require.NoError(t, err)

		// Verify known fields
		require.Equal(t, "210720", data.F9A.Value())
		require.Equal(t, "000000000501", data.F9F02.Value())

		// Verify unknown fields were unmarshaled
		require.NotNil(t, data.F9F36)
		require.Equal(t, []byte{0x01, 0x57}, data.F9F36.Value())

		require.NotNil(t, data.F9F37)
		require.Equal(t, []byte{0x9b, 0xad, 0xbc, 0xab}, data.F9F37.Value())

		// Verify unknown field values using UnmarshalPath
		expected := []byte{0x01, 0x57}
		got := []byte{}

		err = composite.UnmarshalPath("9F36", &got)
		require.NoError(t, err)
		require.Equal(t, expected, got)

		expected = []byte{0x9b, 0xad, 0xbc, 0xab}
		got = []byte{}

		err = composite.UnmarshalPath("9F37", &got)
		require.NoError(t, err)
		require.Equal(t, expected, got)
	})

	t.Run("Unpack, modify unknown tags with Marshal, and Pack reflects changes", func(t *testing.T) {
		composite := NewComposite(storeUnknownTLVSpec)

		// Data contains:
		// - 9A (known): 3 bytes [0x21, 0x07, 0x20]
		// - 9F02 (known): 6 bytes [0x00, 0x00, 0x00, 0x00, 0x05, 0x01]
		// - 9F36 (unknown): 2 bytes [0x01, 0x57]
		// - 9F37 (unknown): 4 bytes [0x9b, 0xad, 0xbc, 0xab]
		inputData := []byte{
			0x30, 0x32, 0x36, // ASCII "026" - length prefix
			0x9a, 0x03, 0x21, 0x07, 0x20, // 9A: length 3, value 210720
			0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, // 9F02: length 6, value 000000000501
			0x9f, 0x36, 0x02, 0x01, 0x57, // 9F36: length 2, value 0157
			0x9f, 0x37, 0x04, 0x9b, 0xad, 0xbc, 0xab, // 9F37: length 4, value 9badbcab
		}

		_, err := composite.Unpack(inputData)
		require.NoError(t, err)

		// Create data struct with modified unknown field values
		data := &TLVDataWithUnknown{
			F9F36: NewBinaryValue([]byte{0xAA, 0xBB}),             // Changed from 0157 to AABB
			F9F37: NewBinaryValue([]byte{0x11, 0x22, 0x33, 0x44}), // Changed from 9badbcab to 11223344
		}

		// Marshal the modified data back to composite
		err = composite.Marshal(data)
		require.NoError(t, err)

		// Pack and verify the changes are reflected
		// Note: fields are sorted by hex value (sort.StringsByHex)
		packed, err := composite.Pack()
		require.NoError(t, err)

		// Expected output with modified unknown fields
		expectedData := []byte{
			0x30, 0x32, 0x36, // ASCII "026" - length prefix (same)
			0x9a, 0x03, 0x21, 0x07, 0x20, // 9A: length 3, value 210720 (unchanged)
			0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, // 9F02: length 6 (unchanged)
			0x9f, 0x36, 0x02, 0xAA, 0xBB, // 9F36: length 2, value AABB (modified)
			0x9f, 0x37, 0x04, 0x11, 0x22, 0x33, 0x44, // 9F37: length 4, value 11223344 (modified)
		}

		require.Equal(t, expectedData, packed)
	})
}

func TestStoreUnknownTLVTagsDisabled(t *testing.T) {
	t.Run("Unknown tags are not stored when StoreUnknownTLVTags is false", func(t *testing.T) {
		// Create spec with SkipUnknownTLVTags=true but StoreUnknownTLVTags=false
		spec := &Spec{
			Length:      999,
			Description: "TLV without storage",
			Pref:        prefix.ASCII.LLL,
			Tag: &TagSpec{
				Enc:                 encoding.BerTLVTag,
				Sort:                sort.StringsByHex,
				SkipUnknownTLVTags:  true,
				StoreUnknownTLVTags: false, // explicitly false
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

		// Data contains:
		// - 9A (unknown): 3 bytes [0x21, 0x07, 0x20]
		// - 9F02 (known): 6 bytes [0x00, 0x00, 0x00, 0x00, 0x05, 0x01]
		// - 9F36 (unknown): 2 bytes [0x01, 0x57]
		// - 9F37 (unknown): 4 bytes [0x9b, 0xad, 0xbc, 0xab]
		inputData := []byte{
			0x30, 0x32, 0x36, // ASCII "026" - length prefix
			0x9a, 0x03, 0x21, 0x07, 0x20, // 9A: length 3, value 210720
			0x9f, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x05, 0x01, // 9F02: length 6, value 000000000501
			0x9f, 0x36, 0x02, 0x01, 0x57, // 9F36: length 2, value 0157
			0x9f, 0x37, 0x04, 0x9b, 0xad, 0xbc, 0xab, // 9F37: length 4, value 9badbcab
		}

		_, err := composite.Unpack(inputData)
		require.NoError(t, err)

		// Verify only known field is stored
		subfields := composite.GetSubfields()
		require.Len(t, subfields, 1)
		require.Contains(t, subfields, "9F02")
		require.NotContains(t, subfields, "9A")
		require.NotContains(t, subfields, "9F36")
		require.NotContains(t, subfields, "9F37")
	})
}
