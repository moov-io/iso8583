package field

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

// TestCompositeField_MultipleOccurrences verifies that when the same data field
// appears multiple times in a composite field, none of the data is lost.
//
// Enhancement: Handling Multiple Occurrences of the Same Data Fields
// When decoding, duplicates are preserved by adding numeric suffixes:
// - First occurrence: "1a"
// - Second occurrence: "1a_1"
// - Third occurrence: "1a_2"
// And so on...
//
// During encoding, these suffixes are removed and all occurrences are packed
// in order using the original field name, ensuring data is not lost.
func TestCompositeField_MultipleOccurrences(t *testing.T) {
	// Create a spec with TLV-style subfields that can repeat
	spec := &Spec{
		Length:      100,
		Description: "Field 104 Composite with TLV",
		Pref:        prefix.ASCII.LL,
		Tag: &TagSpec{
			Length: 2,
			Enc:    encoding.ASCII,
			Pad:    padding.Left('0'),
			Sort:   sort.StringsByInt,
		},
		Subfields: map[string]Field{
			"1a": NewString(&Spec{
				Length:      10,
				Description: "Subfield 1a",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
		},
	}

	// Test Case 1: Pack and then Unpack with single occurrence
	t.Run("single_occurrence", func(t *testing.T) {
		composite := NewComposite(spec)

		type CompositeData struct {
			SubField1a *String `index:"1a"`
		}

		inputData := &CompositeData{
			SubField1a: NewStringValue("HELLO"),
		}

		err := composite.Marshal(inputData)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		// Unpack the data
		newComposite := NewComposite(spec)
		_, err = newComposite.Unpack(packed)
		require.NoError(t, err)

		subfields := newComposite.GetSubfields()
		require.Len(t, subfields, 1, "Should have 1 subfield")
		require.Contains(t, subfields, "1a", "Should have subfield 1a")

		value1a, ok := subfields["1a"].(*String)
		require.True(t, ok, "Subfield 1a should be a String")
		require.Equal(t, "HELLO", value1a.Value(), "Subfield 1a value should be HELLO")
	})

	// Test Case 2: Verify data handling with tag-based composite
	t.Run("multiple_tags_no_data_loss", func(t *testing.T) {
		composite := NewComposite(spec)

		type CompositeData struct {
			SubField1a *String `index:"1a"`
		}

		// Pack data with subfield 1a
		data := &CompositeData{
			SubField1a: NewStringValue("DATATEST"),
		}

		err := composite.Marshal(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		// Verify that data was packed
		require.NotEmpty(t, packed, "Packed data should not be empty")

		// When unpacking, verify data is restored correctly
		unpackComposite := NewComposite(spec)
		_, err = unpackComposite.Unpack(packed)
		require.NoError(t, err)

		subfields := unpackComposite.GetSubfields()
		require.Len(t, subfields, 1, "Should have 1 subfield after unpacking")

		value1a, ok := subfields["1a"].(*String)
		require.True(t, ok, "Subfield 1a should be a String")
		require.Equal(t, "DATATEST", value1a.Value(), "Unpacked data should match original")
	})

	// Test Case 3: Verify no data loss through pack/unpack cycle
	t.Run("pack_unpack_consistency", func(t *testing.T) {
		testData := "ORIGINAL"

		// Pack phase
		composite1 := NewComposite(spec)
		type CompositeData struct {
			SubField1a *String `index:"1a"`
		}

		err := composite1.Marshal(&CompositeData{
			SubField1a: NewStringValue(testData),
		})
		require.NoError(t, err)

		packed, err := composite1.Pack()
		require.NoError(t, err)

		// Unpack phase
		composite2 := NewComposite(spec)
		_, err = composite2.Unpack(packed)
		require.NoError(t, err)

		// Re-pack phase to verify it cycles correctly
		repacked, err := composite2.Pack()
		require.NoError(t, err)

		// Original packed data should match repacked data
		require.Equal(t, packed, repacked, "Repacked data should match original packed data")

		// Verify final unpacking still has correct data
		composite3 := NewComposite(spec)
		_, err = composite3.Unpack(repacked)
		require.NoError(t, err)

		subfields := composite3.GetSubfields()
		value1a, ok := subfields["1a"].(*String)
		require.True(t, ok, "Subfield 1a should be a String")
		require.Equal(t, testData, value1a.Value(), "Final unpacked data should match original")
	})

	// Test Case 4: BUG REPRO - Same tag appearing multiple times (e.g., "1a" appears 3 times)
	// Without the fix, only the last occurrence is kept, others are overwritten
	// With the fix, all occurrences are preserved with numeric suffixes: "1a", "1a_1", "1a_2"
	t.Run("same_tag_multiple_occurrences_bug_repro", func(t *testing.T) {
		// Manually construct raw bytes with same tag appearing 3 times
		// Format: LL (length prefix as ASCII decimal) + [Tag(2 chars) + LL(2 chars) + Data]...

		// Each field: "1a" (tag) + "03" (length=3 in LL format) + "ABC" (data) = 7 bytes
		// Three occurrences total length: 7+7+7 = 21 bytes content
		// LL prefix should be "21" in ASCII

		// Create raw composite data with 3 occurrences of tag "1a"
		rawData := []byte(
			"21" + // LL: Total length is 21 bytes
				"1a03ABC" + // First occurrence of tag "1a" with LL length "03" and data "ABC"
				"1a03DEF" + // Second occurrence of tag "1a" with LL length "03" and data "DEF"
				"1a03GHI", // Third occurrence of tag "1a" with LL length "03" and data "GHI"
		)

		unpackComposite := NewComposite(spec)
		_, err := unpackComposite.Unpack(rawData)
		require.NoError(t, err)

		subfields := unpackComposite.GetSubfields()

		// BUG VERIFICATION:
		// With the bug: Only "1a" would exist, containing "GHI" (last value overwrites previous)
		// With the fix: "1a", "1a_1", "1a_2" should all exist with different values

		// Check if first occurrence exists
		value1a, ok := subfields["1a"].(*String)
		require.True(t, ok, "First occurrence field '1a' should exist")
		require.Equal(t, "ABC", value1a.Value(), "First occurrence should have value 'ABC'")

		// Check if second occurrence exists with suffix
		value1a_1, ok := subfields["1a_1"].(*String)
		require.True(t, ok, "Second occurrence field '1a_1' should exist (BUG if this fails: data loss!)")
		require.Equal(t, "DEF", value1a_1.Value(), "Second occurrence should have value 'DEF'")

		// Check if third occurrence exists with suffix
		value1a_2, ok := subfields["1a_2"].(*String)
		require.True(t, ok, "Third occurrence field '1a_2' should exist (BUG if this fails: data loss!)")
		require.Equal(t, "GHI", value1a_2.Value(), "Third occurrence should have value 'GHI'")

		// Verify we have exactly 3 fields
		require.Len(t, subfields, 3, "Should have 3 subfields (one original + two with suffixes)")
	})
}
