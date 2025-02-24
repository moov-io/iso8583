# Guide to Composite Fields

## Overview

Composite fields in ISO8583 messages can represent complex data structures like TLV (Tag-Length-Value) or structured fields with subfields. This guide covers how to define and work with composite fields.

## Field Packing and Unpacking Flow

Composite fields follow a general structure of `[tag][length][value]` and are processed recursively:

### Packing Process
1. When packing a message, we build the binary data field by field
2. When a field is a composite, we recursively pack its subfields
3. For each subfield, we add its tag (if applicable), length (if applicable), and value
4. The packed values are concatenated according to field specifications
5. The final composite field includes its own length prefix

### Unpacking Process
1. When unpacking a message, we read and parse fields in sequence
2. For composite fields, we first read the field's length
3. We then recursively unpack each subfield within that length
4. Each subfield parser reads the tag, length, and value
5. The data is stored in the field structure for access

This recursive approach allows handling complex nested structures.

## Basic Composite Field Structure

A composite field consists of:
1. A base field specification (`Spec`)
2. Either a `Tag` or `Bitmap` definition for organizing subfields
3. Subfield definitions


## Common Field Types

### 1. Positional Subfields

Positional subfields represent a fixed format where multiple data elements are concatenated together without any explicit tags or separators in the binary data. Each element has a predefined position and length. This format is common in older financial message formats where space efficiency was critical.

In positional subfields, you'll typically see specifications like "Type (bytes 1-2), Name (bytes 3-20)", "Account Number (bytes 21-30)", etc. The binary data is just the raw values placed next to each other.

```go
spec := &field.Spec{
    Length:      30,             // Total length of all subfields combined
    Description: "Structured Field",
    Pref:        prefix.ASCII.Fixed,
    Tag: &field.TagSpec{
        Sort: sort.StringsByInt, // Only sorting is needed, no encoding
    },
    Subfields: map[string]field.Field{
        "01": field.NewString(&field.Spec{ // First subfield
            Length:      2,
            Description: "Type",
            Enc:         encoding.ASCII,
            Pref:        prefix.ASCII.Fixed,
        }),
        "02": field.NewString(&field.Spec{ // Second subfield
            Length:      18,
            Description: "Name",
            Enc:         encoding.ASCII,
            Pref:        prefix.ASCII.Fixed,
        }),
        "03": field.NewString(&field.Spec{ // Third subfield
            Length:      10,
            Description: "Account Number",
            Enc:    encoding.ASCII,
            Pref:   prefix.ASCII.Fixed,
        }),
    },
}
```

When working with positional subfields:
- This represents a set of values with specific lengths concatenated together, variable length is supported
- No tags appear in the data itself
- Keys in the map (like "01", "02", "03") are just identifiers for programmatic access
- The `Sort` function determines the order in which subfields are packed/unpacked
- The field's total `Length` must equal the sum of all subfield lengths in case of fixed-length subfields, or maximum allowed length
- It's your responsibility to order subfields correctly by assigning appropriate key values and using the right Sort function

### 2. TLV Fields (Tag-Length-Value)

In TLV (Tag-Length-Value) structures, each data element has an explicit tag to identify it, followed by a length indicator and then the value itself.

In TLV formats, both the tag length and the value length encoding are explicitly defined in the specification. This differs from BER-TLV where the encoding rules are standardized and more dynamic.

```go
spec := &field.Spec{
    Length:      999,            // Maximum allowable length
    Description: "TLV Data Field",
    Pref:        prefix.ASCII.LLL,
    Tag: &field.TagSpec{
        Length: 2,              // Tag length in bytes - explicitly defined
        Enc:    encoding.ASCII, // Tag encoding
        Sort:   sort.StringsByInt,
    },
    Subfields: map[string]field.Field{
        "01": field.NewString(&field.Spec{
            Length:      10,
            Description: "Name",
            Enc:         encoding.ASCII,
            Pref:        prefix.ASCII.LL,  // Length encoding for this value
        }),
        "02": field.NewString(&field.Spec{
            Length:      20,
            Description: "Address",
            Enc:         encoding.ASCII,
            Pref:        prefix.ASCII.LL,  // Length encoding for this value
        }),
    },
}
```

Key characteristics of TLV Fields:
- Both tag length and value length encoding are explicitly defined in the specification
- Tags with the length are present in the packed binary data
- This differs from BER-TLV where tag length and value length are dynamically encoded in the data itself

### 3. BER-TLV Fields

Unlike basic TLV formats, BER-TLV has standardized rules for encoding tag and length fields. The tag and length encoding are inherent to the format rather than explicitly defined in the specification.

```go
spec := &field.Spec{
    Length:      999,           // Maximum length, not actual
    Description: "EMV Data",
    Pref:        prefix.ASCII.LLL,
    Tag: &field.TagSpec{
        Enc:  encoding.BerTLVTag, // No Length needed - BerTLV handles tag length dynamically
        Sort: sort.StringsByHex,
    },
    Subfields: map[string]field.Field{
        "9F02": field.NewHex(&field.Spec{
            Description: "Amount, Authorized",
            Enc:         encoding.Binary,
            Pref:        prefix.BerTLV, // BerTLV handles value length encoding dynamically
        }),
        // Other EMV tags...
    },
}
```

With BER-TLV fields:
- Do not set explicit tag length, as BER-TLV encoding dynamically determines tag length
- The field's Length attribute only indicates the maximum allowed length
- Tag and value lengths are encoded within the data itself according to BER-TLV rules
- Length encoding for values is handled by prefix.BerTLV
- This differs from regular TLV where tag length and value length encoding are explicitly defined

> **Note:** To work with BerTLV subfields, we suggest using the `encoding.Binary` for the entire field and then use [moov-io/bertlv](https://github.com/moov-io/bertlv) to parse the subfields.


### 4. Nested Composite Fields with Data Set IDs

Some composite fields use a hierarchical approach with "Data Set IDs" that organize related data elements. Each Data Set has its own ID, length, and a value, which is a collection of TLV (or BerTLV) fields.


```go
spec := &field.Spec{
    Length:      255,
    Description: "Extended Transaction Data",
    Pref:        prefix.Binary.L,
    Tag: &field.TagSpec{
        Length: 1,                // Data Set ID length
        Enc:    encoding.ASCIIHexToBytes,
        Sort:   sort.StringsByHex,

        // Configure handling of unknown Data Set IDs
        SkipUnknownTLVTags: true,
        // Data Set length prefix (2 bytes)
        // It matches the Pref of the inner subfield
        PrefUnknownTLV:     prefix.Binary.LL, 
    },
    Subfields: map[string]field.Field{
        "56": field.NewComposite(&field.Spec{ // Data Set ID "56"
            Length:      1535,
            Description: "Merchant Information Data",
            Pref:        prefix.Binary.LL, // Data Set length prefix
            Tag: &field.TagSpec{
                Enc:                encoding.BerTLVTag,
                Sort:               sort.StringsByHex,
                SkipUnknownTLVTags: true,
            },
            Subfields: map[string]field.Field{
                "9F": field.NewString(&field.Spec{
                    Length:      11,
                    Description: "Merchant Identifier",
                    Enc:         encoding.EBCDIC,
                    Pref:        prefix.BerTLV,
                }),
                "80": field.NewString(&field.Spec{
                    Length:      15,
                    Description: "Terminal Identifier",
                    Enc:         encoding.EBCDIC,
                    Pref:        prefix.BerTLV,
                }),
            },
        }),
    },
}
```

This structure represents:
1. The outer field with Data Set IDs as its tags
2. Each Data Set ID (e.g., "56") mapped to its own composite field
3. Within each Data Set, a BER-TLV structure with tags like "01", "02", etc. for subfields. But it can be any other composite field.

> **Important:** When implementing fields with Data Set IDs where values are subfields, build the structure in a nested way using Data Set IDs as tags (not as values).

### Handling Unknown Tags and Data Sets

The library can be configured to skip unknown tags or Data Set IDs during parsing:

```go
spec := &field.Spec{
    // ... other settings ...
    Tag: &field.TagSpec{
        SkipUnknownTLVTags: true,            // Enable skipping unknown tags
        PrefUnknownTLV:     prefix.ASCII.L,  // Format for unknown tag length
    },
}
```

When the parser encounters an unknown tag or Data Set ID:

1. It checks if `SkipUnknownTLVTags` is enabled
2. If enabled, it uses `PrefUnknownTLV` to determine how to read the length, unless it's a BER-TLV field
3. It then reads the specified length of data
4. It skips that data and continues parsing subsequent fields

This is crucial when:
- The specification may evolve with new tags or you don't want to handle all tags
- Working with Data Set IDs where some IDs might be unknown to your implementation

> **Note:** Using [moov-io/bertlv](https://github.com/moov-io/bertlv) handles cases like unknown tags, tags with the same name, etc.


## Example: Complete Implementation with Data Set IDs

Here's a complete example for a field with Data Set IDs and nested BER-TLV data:

```go
// Field 125 specification
field125Spec := &field.Spec{
    Length:      255,
    Description: "Extended Transaction Data",
    Pref:        prefix.Binary.L,
    Tag: &field.TagSpec{
        Length: 1,                         // Data Set ID length in bytes
        Enc:    encoding.ASCIIHexToBytes, // Data Set ID below is ASCII hex and will be converted to bytes (1 byte)
        Sort:   sort.StringsByHex,
        // Handle unknown Data Set IDs
        SkipUnknownTLVTags: true,
        PrefUnknownTLV:     prefix.Binary.LL,  // Data Set length format
    },
    Subfields: map[string]field.Field{
        "56": field.NewComposite(&field.Spec{ // Data Set ID
            Length:      1535,
            Description: "Merchant Information Data",
            Pref:        prefix.Binary.LL,  // Length of the entire Data Set value
            Tag: &field.TagSpec{
                Enc:                encoding.BerTLVTag,
                Sort:               sort.StringsByHex,
                SkipUnknownTLVTags: true,
            },
            Subfields: map[string]field.Field{
                "01": field.NewString(&field.Spec{
                    Length:      11,
                    Description: "Merchant Identifier",
                    Enc:         encoding.EBCDIC,
                    Pref:        prefix.BerTLV,
                }),
                "02": field.NewString(&field.Spec{
                    Length:      15,
                    Description: "Terminal Identifier",
                    Enc:         encoding.EBCDIC,
                    Pref:        prefix.BerTLV,
                }),
                "81": field.NewString(&field.Spec{
                    Length:      25,
                    Description: "Merchant Legal Name",
                    Enc:         encoding.EBCDIC,
                    Pref:        prefix.BerTLV,
                }),
                "82": field.NewString(&field.Spec{
                    Length:      25,
                    Description: "Merchant DBA Name",
                    Enc:         encoding.EBCDIC,
                    Pref:        prefix.BerTLV,
                }),
            },
        }),
        
        // Data Set ID "65" - Terminal Information
        "65": field.NewComposite(&field.Spec{
            Length:      1535,
            Description: "Terminal Information Data",
            Pref:        prefix.Binary.LL,
            Tag: &field.TagSpec{
                Enc:                encoding.BerTLVTag,
                Sort:               sort.StringsByHex,
                SkipUnknownTLVTags: true,
            },
            Subfields: map[string]field.Field{
                "03": field.NewHex(&field.Spec{
                    Description: "Terminal Capabilities",
                    Enc:         encoding.Binary,
                    Pref:        prefix.BerTLV,
                }),
                "04": field.NewHex(&field.Spec{
                    Description: "Terminal Type",
                    Enc:         encoding.Binary,
                    Pref:        prefix.BerTLV,
                }),
                "87": field.NewString(&field.Spec{
                    Length:      8,
                    Description: "Terminal Serial Number",
                    Enc:         encoding.EBCDIC,
                    Pref:        prefix.BerTLV,
                }),
            },
        }),
    },
}

// Corresponding data struct definitions that can be used to marshal/unmarshal

type Field125Data struct {
    MerchantData  *MerchantDataSet  `index:"56"`
    TerminalData  *TerminalDataSet  `index:"65"`
}

type MerchantDataSet struct {
    MerchantId    *field.String `index:"01"`
    TerminalId    *field.String `index:"02"`
    LegalName     *field.String `index:"81"`
    DBAName       *field.String `index:"82"`
}

type TerminalDataSet struct {
    Capabilities  *field.Hex    `index:"03"`
    Type          *field.Hex    `index:"04"`
    SerialNumber  *field.String `index:"87"`
}
```

## Best Practices for Testing Composite Fields

A good approach is to test composite field specs individually before integrating them into a larger message spec. This helps isolate any specification issues:

```go
type TestData struct {
	MerchantData *MerchantData `index:"56"`
}

type MerchantData struct {
	MerchantID *field.String `index:"01"`
}

func TestDataSetCompositeField(t *testing.T) {
	// Define the specification for a composite field with Data Set IDs
	spec := &field.Spec{
		Length:      255,
		Description: "Test Data Set Field",
		Pref:        prefix.Binary.L,
		Tag: &field.TagSpec{
			Length:             1,
			Enc:                encoding.ASCIIHexToBytes,
			Sort:               sort.StringsByHex,
			SkipUnknownTLVTags: true,
			PrefUnknownTLV:     prefix.Binary.LL,
		},
		Subfields: map[string]field.Field{
			"56": field.NewComposite(&field.Spec{
				Length:      255,
				Description: "Merchant Data",
				Pref:        prefix.Binary.LL,
				Tag: &field.TagSpec{
					Enc:  encoding.BerTLVTag,
					Sort: sort.StringsByHex,
				},
				Subfields: map[string]field.Field{
					"01": field.NewString(&field.Spec{
						Length:      11,
						Description: "Merchant ID",
						Enc:         encoding.EBCDIC,
						Pref:        prefix.BerTLV,
					}),
				},
			}),
		},
	}

	// Create a new composite field with our spec
	composite := field.NewComposite(spec)

	data := &TestData{
		MerchantData: &MerchantData{
			MerchantID: field.NewStringValue("12345ABCDE"),
		},
	}

	// Marshal the data into the field
	err := composite.Marshal(data)
	require.NoError(t, err)

	// Pack the field to binary
	packed, err := composite.Pack()
	require.NoError(t, err)

	// Optional: Inspect the packed data for debugging
	t.Logf("Packed data (hex): %X", packed)

	// Create a new field for unpacking
	unpackedField := field.NewComposite(spec)

	// Unpack the binary data
	read, err := unpackedField.Unpack(packed)
	require.NoError(t, err)
	require.Equal(t, len(packed), read, "Should read all bytes")

	// Unmarshal into a new struct to verify data
	unpacked := &TestData{}

	err = unpackedField.Unmarshal(unpacked)
	require.NoError(t, err)

	// Verify the data matches
	require.NotNil(t, unpacked.MerchantData)
	require.NotNil(t, unpacked.MerchantData.MerchantID)
	require.Equal(t, "12345ABCDE", unpacked.MerchantData.MerchantID.Value())
}
```
