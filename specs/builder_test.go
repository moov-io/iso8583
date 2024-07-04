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
				},
			}),
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

func TestSpecWithCompositeBitmapedFields(t *testing.T) {
	specJSON := []byte(`
{
	"name": "TEST Spec",
	"fields": {
		"1": {
			"type": "Composite",
			"length": 255,
			"description": "Private use field",
			"prefix": "ASCII.LL",
			"bitmap": {
					"type": "Bitmap",
					"length": 8,
					"description": "Bitmap",
					"enc": "HexToASCII",
					"prefix": "Hex.Fixed",
					"disableAutoExpand": true
			},
			"subfields": {
				"1": {
					"type": "String",
					"length": 2,
					"description": "Cardholder certificate Serial Number",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				},
				"2": {
					"type": "String",
					"length": 2,
					"description": "Merchant certificate Serial Number",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				},
				"3": {
					"type": "String",
					"length": 2,
					"description": "Transaction ID",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				},
				"4": {
					"type": "String",
					"length": 20,
					"description": "CAVV",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				},
				"5": {
					"type": "String",
					"length": 20,
					"description": "CAVV",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				},
				"6": {
					"type": "String",
					"length": 2,
					"description": "Cardholder certificate Serial Number",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				},
				"7": {
					"type": "String",
					"length": 2,
					"description": "Merchant certificate Serial Number",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				},
				"8": {
					"type": "String",
					"length": 2,
					"description": "Transaction ID",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				},
				"9": {
					"type": "String",
					"length": 20,
					"description": "CAVV",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				},
				"10": {
					"type": "String",
					"length": 6,
					"description": "CVV2",
					"enc": "ASCII",
					"prefix": "ASCII.Fixed"
				}
			}
		}
	}
}`)

	spec, err := Builder.ImportJSON(specJSON)
	require.NoError(t, err)

	data := struct {
		F1  *field.String
		F2  *field.String
		F3  *field.String
		F4  *field.String
		F5  *field.String
		F6  *field.String
		F7  *field.String
		F8  *field.String
		F9  *field.String
		F10 *field.String
	}{
		F10: field.NewStringValue("11 456"),
	}

	compositeRestored := field.NewComposite(spec.Fields[1].Spec())
	err = compositeRestored.Marshal(&data)
	require.NoError(t, err)

	packed, err := compositeRestored.Pack()
	require.NoError(t, err)
	require.Equal(t, "22004000000000000011 456", string(packed))

	exportedJSON, err := Builder.ExportJSON(spec)
	require.NoError(t, err)

	spec, err = Builder.ImportJSON(exportedJSON)
	require.NoError(t, err)

	compositeRestored = field.NewComposite(spec.Fields[1].Spec())
	err = compositeRestored.Marshal(&data)
	require.NoError(t, err)

	packed, err = compositeRestored.Pack()
	require.NoError(t, err)
	require.Equal(t, "22004000000000000011 456", string(packed))
}

func TestExportImportWithNonePrefixField(t *testing.T) {
	spec := &iso8583.MessageSpec{
		Fields: map[int]field.Field{
			3: field.NewComposite(&field.Spec{
				Description: "Processing code",
				Pref:        prefix.None.Fixed,
				Tag: &field.TagSpec{
					Enc:  encoding.ASCII,
					Sort: sort.StringsByInt,
				},
				Subfields: map[string]field.Field{
					"1": field.NewNumeric(&field.Spec{
						Length:      3,
						Description: "Transaction code",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
				},
			}),
		},
	}

	specJSON, err := Builder.ExportJSON(spec)
	require.NoError(t, err)

	spec, err = Builder.ImportJSON(specJSON)
	require.NoError(t, err)
}
