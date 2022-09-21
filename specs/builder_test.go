package specs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
)

func TestBuilder(t *testing.T) {

	asciiJson, err := Builder.ExportJSON(Spec87ASCII)
	require.NoError(t, err)

	asciiSpec, err := Builder.ImportJSON(asciiJson)
	require.NoError(t, err)

	require.Exactly(t, Spec87ASCII, asciiSpec)

	hexJson, err := Builder.ExportJSON(Spec87Hex)
	require.NoError(t, err)

	hexSpec, err := Builder.ImportJSON(hexJson)
	require.NoError(t, err)

	require.Exactly(t, Spec87Hex.Name, hexSpec.Name)

}

func TestExampleJSONSpec(t *testing.T) {
	asciiJson, err := os.ReadFile("../examples/specs/spec87ascii.json")
	require.NoError(t, err)

	asciiSpec, err := Builder.ImportJSON(asciiJson)
	require.NoError(t, err)
	require.Exactly(t, Spec87ASCII, asciiSpec)
}

func TestSpecWithCompositeFields(t *testing.T) {
	testSpec := &iso8583.MessageSpec{
		Name: "ISO 8583 v1987 ASCII",
		Fields: map[int]field.Field{
			1: field.NewComposite(&field.Spec{
				Length:      3,
				Description: "example with a tag with encoding",
				Pref:        prefix.EBCDIC.LLL,
				Tag: &field.TagSpec{
					Length: 2,
					Enc:    encoding.EBCDIC,
					Pad:    padding.Left('0'),
					Sort:   sort.StringsByInt,
				},
				Subfields: map[string]field.Field{
					"1": field.NewComposite(&field.Spec{
						Length:      7,
						Description: "example subfield with nested subfields",
						Pref:        prefix.EBCDIC.Fixed,
						Tag: &field.TagSpec{
							Sort: sort.StringsByInt,
						},
						Subfields: map[string]field.Field{
							"1": field.NewString(&field.Spec{
								Length:      4,
								Description: "Date",
								Enc:         encoding.EBCDIC,
								Pref:        prefix.EBCDIC.Fixed,
							}),
							"2": field.NewString(&field.Spec{
								Length:      3,
								Description: "Data",
								Enc:         encoding.EBCDIC,
								Pref:        prefix.EBCDIC.Fixed,
							}),
						},
					}),
					"2": field.NewNumeric(&field.Spec{
						Length:      5,
						Description: "num field",
						Enc:         encoding.EBCDIC,
						Pref:        prefix.EBCDIC.Fixed,
						Pad:         padding.Left('0'),
					}),
				}}),
			30: field.NewNumeric(&field.Spec{
				Length:      5,
				Description: "field key that is not next number",
				Enc:         encoding.EBCDIC,
				Pref:        prefix.EBCDIC.Fixed,
				Pad:         padding.Left('0'),
			}),
		},
	}

	specJSON, err := Builder.ExportJSON(testSpec)
	require.NoError(t, err)
	importedSpec, err := Builder.ImportJSON(specJSON)
	require.NoError(t, err)
	reexportedJSON, err := Builder.ExportJSON(testSpec)
	require.NoError(t, err)
	require.Equal(t, specJSON, reexportedJSON)

	// We can't compare sort functions for equality, so nil them out to check the rest
	testSpec.Fields[1].Spec().Tag.Sort = nil
	importedSpec.Fields[1].Spec().Tag.Sort = nil
	testSpec.Fields[1].Spec().Subfields["1"].Spec().Tag.Sort = nil
	importedSpec.Fields[1].Spec().Subfields["1"].Spec().Tag.Sort = nil

	require.Exactly(t, testSpec, importedSpec)
}
