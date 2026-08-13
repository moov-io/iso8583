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
	asciiJson, err := ExportJSON(Spec87ASCII)
	require.NoError(t, err)

	asciiSpec, err := ImportJSON(asciiJson)
	require.NoError(t, err)

	require.Exactly(t, Spec87ASCII, asciiSpec)

	hexJson, err := ExportJSON(Spec87Hex)
	require.NoError(t, err)

	hexSpec, err := ImportJSON(hexJson)
	require.NoError(t, err)

	require.Exactly(t, Spec87Hex.Name, hexSpec.Name)
}

func TestImportingJSONWithTrack2Spec(t *testing.T) {
	track2Json, err := os.ReadFile("../examples/fields/track2.json")
	require.NoError(t, err)

	track2Spec, err := ImportJSON(track2Json)
	require.NoError(t, err)
	require.Exactly(t, Spec87Track2, track2Spec)
}

func TestExampleJSONSpec(t *testing.T) {
	asciiJson, err := os.ReadFile("../examples/specs/spec87ascii.json")
	require.NoError(t, err)

	asciiSpec, err := ImportJSON(asciiJson)
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

	specJSON, err := ExportJSON(testSpec)
	require.NoError(t, err)
	importedSpec, err := ImportJSON(specJSON)
	require.NoError(t, err)
	reexportedJSON, err := ExportJSON(testSpec)
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

	spec, err := ImportJSON(specJSON)
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

	exportedJSON, err := ExportJSON(spec)
	require.NoError(t, err)

	spec, err = ImportJSON(exportedJSON)
	require.NoError(t, err)

	compositeRestored = field.NewComposite(spec.Fields[1].Spec())
	err = compositeRestored.Marshal(&data)
	require.NoError(t, err)

	packed, err = compositeRestored.Pack()
	require.NoError(t, err)
	require.Equal(t, "22004000000000000011 456", string(packed))
}

func TestExportImportTagSpecTLVFields(t *testing.T) {
	spec := &iso8583.MessageSpec{
		Name: "TLV Spec",
		Fields: map[int]field.Field{
			1: field.NewComposite(&field.Spec{
				Length:      999,
				Description: "TLV field with unknown tag handling",
				Pref:        prefix.ASCII.LLL,
				Tag: &field.TagSpec{
					Length:              2,
					Enc:                 encoding.ASCII,
					Pad:                 padding.Left('0'),
					Sort:                sort.StringsByInt,
					SkipUnknownTLVTags:  true,
					StoreUnknownTLVTags: true,
					PrefUnknownTLV:      prefix.ASCII.LL,
				},
				Subfields: map[string]field.Field{
					"1": field.NewString(&field.Spec{
						Length:      10,
						Description: "Sub 1",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.LL,
					}),
				},
			}),
		},
	}

	specJSON, err := ExportJSON(spec)
	require.NoError(t, err)

	importedSpec, err := ImportJSON(specJSON)
	require.NoError(t, err)

	// Verify the TLV-specific tag fields round-tripped correctly
	importedTag := importedSpec.Fields[1].Spec().Tag
	require.True(t, importedTag.SkipUnknownTLVTags)
	require.True(t, importedTag.StoreUnknownTLVTags)
	require.NotNil(t, importedTag.PrefUnknownTLV)
	require.Equal(t, "ASCII.LL", importedTag.PrefUnknownTLV.Inspect())

	// Verify JSON re-export is stable
	reexportedJSON, err := ExportJSON(importedSpec)
	require.NoError(t, err)
	require.Equal(t, specJSON, reexportedJSON)
}

func TestBuilderYAML(t *testing.T) {
	asciiYAML, err := ExportYAML(Spec87ASCII)
	require.NoError(t, err)

	asciiSpec, err := ImportYAML(asciiYAML)
	require.NoError(t, err)

	require.Exactly(t, Spec87ASCII, asciiSpec)

	hexYAML, err := ExportYAML(Spec87Hex)
	require.NoError(t, err)

	hexSpec, err := ImportYAML(hexYAML)
	require.NoError(t, err)

	require.Exactly(t, Spec87Hex.Name, hexSpec.Name)
}

func TestExampleYAMLSpec(t *testing.T) {
	asciiYAML, err := os.ReadFile("../examples/specs/spec87ascii.yaml")
	require.NoError(t, err)

	asciiSpec, err := ImportYAML(asciiYAML)
	require.NoError(t, err)
	require.Exactly(t, Spec87ASCII, asciiSpec)
}

func TestSpecWithCompositeFieldsYAML(t *testing.T) {
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

	specYAML, err := ExportYAML(testSpec)
	require.NoError(t, err)
	importedSpec, err := ImportYAML(specYAML)
	require.NoError(t, err)
	reexportedYAML, err := ExportYAML(testSpec)
	require.NoError(t, err)
	require.Equal(t, specYAML, reexportedYAML)

	// We can't compare sort functions for equality, so nil them out to check the rest
	testSpec.Fields[1].Spec().Tag.Sort = nil
	importedSpec.Fields[1].Spec().Tag.Sort = nil
	testSpec.Fields[1].Spec().Subfields["1"].Spec().Tag.Sort = nil
	importedSpec.Fields[1].Spec().Subfields["1"].Spec().Tag.Sort = nil

	require.Exactly(t, testSpec, importedSpec)
}

func TestExportImportTagSpecTLVFieldsYAML(t *testing.T) {
	spec := &iso8583.MessageSpec{
		Name: "TLV Spec",
		Fields: map[int]field.Field{
			1: field.NewComposite(&field.Spec{
				Length:      999,
				Description: "TLV field with unknown tag handling",
				Pref:        prefix.ASCII.LLL,
				Tag: &field.TagSpec{
					Length:              2,
					Enc:                 encoding.ASCII,
					Pad:                 padding.Left('0'),
					Sort:                sort.StringsByInt,
					SkipUnknownTLVTags:  true,
					StoreUnknownTLVTags: true,
					PrefUnknownTLV:      prefix.ASCII.LL,
				},
				Subfields: map[string]field.Field{
					"1": field.NewString(&field.Spec{
						Length:      10,
						Description: "Sub 1",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.LL,
					}),
				},
			}),
		},
	}

	specYAML, err := ExportYAML(spec)
	require.NoError(t, err)

	importedSpec, err := ImportYAML(specYAML)
	require.NoError(t, err)

	// Verify the TLV-specific tag fields round-tripped correctly
	importedTag := importedSpec.Fields[1].Spec().Tag
	require.True(t, importedTag.SkipUnknownTLVTags)
	require.True(t, importedTag.StoreUnknownTLVTags)
	require.NotNil(t, importedTag.PrefUnknownTLV)
	require.Equal(t, "ASCII.LL", importedTag.PrefUnknownTLV.Inspect())

	// Verify YAML re-export is stable
	reexportedYAML, err := ExportYAML(importedSpec)
	require.NoError(t, err)
	require.Equal(t, specYAML, reexportedYAML)
}

func TestCrossFormatRoundTrip(t *testing.T) {
	// JSON export → import → YAML export → import, verify equivalence
	jsonData, err := ExportJSON(Spec87ASCII)
	require.NoError(t, err)

	specFromJSON, err := ImportJSON(jsonData)
	require.NoError(t, err)

	yamlData, err := ExportYAML(specFromJSON)
	require.NoError(t, err)

	specFromYAML, err := ImportYAML(yamlData)
	require.NoError(t, err)

	require.Exactly(t, Spec87ASCII, specFromYAML)
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

	specJSON, err := ExportJSON(spec)
	require.NoError(t, err)

	spec, err = ImportJSON(specJSON)
	require.NoError(t, err)
}

// Export names field types with reflection, so any field type in the library
// can be written to a spec document, while import resolves them through
// FieldConstructor. This pins the round trip for every field type so the two
// directions cannot drift apart.
func TestExportImportAllFieldTypes(t *testing.T) {
	spec := &iso8583.MessageSpec{
		Name: "all field types",
		Fields: map[int]field.Field{
			0: field.NewString(&field.Spec{
				Length:      4,
				Description: "Message Type Indicator",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			1: field.NewBitmap(&field.Spec{
				Length:      16,
				Description: "Bitmap",
				Enc:         encoding.BytesToASCIIHex,
				Pref:        prefix.Hex.Fixed,
			}),
			2: field.NewNumeric(&field.Spec{
				Length:      19,
				Description: "Primary Account Number",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			3: field.NewBinary(&field.Spec{
				Length:      8,
				Description: "Binary field",
				Enc:         encoding.Binary,
				Pref:        prefix.Binary.Fixed,
			}),
			52: field.NewHex(&field.Spec{
				Length:      8,
				Description: "PIN Data",
				Enc:         encoding.Binary,
				Pref:        prefix.Binary.Fixed,
			}),
			35: field.NewTrack2(&field.Spec{
				Length:      37,
				Description: "Track 2 Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			36: field.NewTrack3(&field.Spec{
				Length:      104,
				Description: "Track 3 Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LLL,
			}),
			45: field.NewTrack1(&field.Spec{
				Length:      76,
				Description: "Track 1 Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
		},
	}

	jsonData, err := ExportJSON(spec)
	require.NoError(t, err)

	fromJSON, err := ImportJSON(jsonData)
	require.NoError(t, err)
	require.Exactly(t, spec, fromJSON)

	yamlData, err := ExportYAML(spec)
	require.NoError(t, err)

	fromYAML, err := ImportYAML(yamlData)
	require.NoError(t, err)
	require.Exactly(t, spec, fromYAML)
}

// emvICCDataSpecYAML is a BER-TLV composite as it appears in a real spec
// document. It exercises the composite tag encoding and Hex subfields
// together, which is how EMV data is described.
const emvICCDataSpecYAML = `name: EMV ICC Data
fields:
    "55":
        type: Composite
        length: 999
        description: ICC Data, EMV having multiple tags
        prefix: ASCII.LLL
        tag:
            enc: BerTLVTag
            sort: StringsByHex
        subfields:
            "82":
                type: Hex
                description: Application Interchange Profile
                enc: Binary
                prefix: BerTLV
            "95":
                type: Hex
                description: Terminal Verification Results
                enc: Binary
                prefix: BerTLV
            "9A":
                type: Hex
                description: Transaction Date
                enc: Binary
                prefix: BerTLV
            "9F02":
                type: Hex
                description: Amount, Authorized
                enc: Binary
                prefix: BerTLV
            "9F26":
                type: Hex
                description: Application Cryptogram
                enc: Binary
                prefix: BerTLV
            "9F36":
                type: Hex
                description: Application Transaction Counter
                enc: Binary
                prefix: BerTLV
`

func TestImportExportEMVCompositeSpec(t *testing.T) {
	spec, err := ImportYAML([]byte(emvICCDataSpecYAML))
	require.NoError(t, err)

	iccSpec := spec.Fields[55].Spec()
	require.Equal(t, encoding.BerTLVTag, iccSpec.Tag.Enc)
	require.Len(t, iccSpec.Subfields, 6)
	require.Equal(t, "Application Cryptogram", iccSpec.Subfields["9F26"].Spec().Description)

	exported, err := ExportYAML(spec)
	require.NoError(t, err)

	// Exporting is stable: importing the exported document and exporting it
	// again produces the same bytes.
	reimported, err := ImportYAML(exported)
	require.NoError(t, err)
	reexported, err := ExportYAML(reimported)
	require.NoError(t, err)
	require.Equal(t, string(exported), string(reexported))
}

// EBCDIC1047 is resolvable by name on import and by struct name on export.
func TestExportImportEBCDIC1047Encoding(t *testing.T) {
	spec := &iso8583.MessageSpec{
		Name: "ebcdic 1047",
		Fields: map[int]field.Field{
			2: field.NewString(&field.Spec{
				Length:      10,
				Description: "EBCDIC 1047 field",
				Enc:         encoding.EBCDIC1047,
				Pref:        prefix.EBCDIC.Fixed,
			}),
		},
	}

	jsonData, err := ExportJSON(spec)
	require.NoError(t, err)

	fromJSON, err := ImportJSON(jsonData)
	require.NoError(t, err)
	require.Exactly(t, spec, fromJSON)
}
