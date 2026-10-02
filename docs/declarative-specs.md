# Defining a Spec as a Document

You can write a `MessageSpec` in Go. You can also write it as a YAML or JSON
document and load it at runtime. The two forms give the same spec. This guide
is about the document form. The other guides use Go.

This guide is not about [JSON Encoding and
Decoding](../README.md#json-encoding-and-decoding). That section serializes a
*message*, which is the field values. This guide is about the *spec*, which is
the definition that the library uses to parse messages.

## Why a document

A spec in Go is part of the binary. To change a dialect, you must build and
deploy the program again. Other tools cannot read the spec. For example, a
validator, a documentation generator, or an operations team cannot see what the
switch expects.

A spec in a document has its own version history. When a card scheme sends a
bulletin, you can compare the old and the new document. The same binary can
load the new document without a new build.

## Reading and writing

```go
import "github.com/moov-io/iso8583/specs"

spec, err := specs.ImportYAML(raw)   // or ImportJSON
raw, err := specs.ExportYAML(spec)   // or ExportJSON
```

To make your first document, use `ExportYAML`:

1. Build the spec in Go, as you do now.
2. Export it with `ExportYAML`.
3. Keep the file.

## What a document looks like

```yaml
name: Example
fields:
    "0":
        type: String
        length: 4
        description: Message Type Indicator
        enc: ASCII
        prefix: ASCII.Fixed
    "1":
        type: Bitmap
        description: Bitmap
        enc: HexToASCII
        prefix: Hex.Fixed
    "2":
        type: String
        length: 19
        description: Primary Account Number
        enc: ASCII
        prefix: ASCII.LL
    "4":
        type: Numeric
        length: 12
        description: Amount, Transaction
        enc: ASCII
        prefix: ASCII.Fixed
        padding:
            type: Left
            pad: "0"
```

Field keys are strings, because JSON object keys must be strings. The values of
`type`, `enc`, `prefix` and `padding.type` are names. The importer finds the
Go value for each name. The lists in [What the importer
accepts](#what-the-importer-accepts) give all the names.

The library includes a complete document:
[`examples/specs/spec87ascii.yaml`](../examples/specs/spec87ascii.yaml). The
same spec is also in JSON:
[`spec87ascii.json`](../examples/specs/spec87ascii.json).

## Composite fields

A composite has `subfields`. It also has a `tag` block or a `bitmap`, but not
both. If a composite has neither, the library cannot build it, and the import
fails. For the meaning of each shape, see the [composite fields
guide](composite-fields.md).

### Positional subfields

```yaml
    "3":
        type: Composite
        length: 6
        description: Processing Code
        prefix: ASCII.Fixed
        tag:
            sort: StringsByInt
        subfields:
            "1":
                type: String
                length: 2
                description: Transaction Type
                enc: ASCII
                prefix: ASCII.Fixed
            "2":
                type: String
                length: 2
                description: From Account
                enc: ASCII
                prefix: ASCII.Fixed
            "3":
                type: String
                length: 2
                description: To Account
                enc: ASCII
                prefix: ASCII.Fixed
```

The `tag` block can also have `length`, `enc`, `padding`, `skipUnknownTLVTags`,
`storeUnknownTLVTags` and `prefUnknownTLV`. These keys match the fields of
`field.TagSpec`. The `padding` key in the `tag` block has the same `{type, pad}`
form as the `padding` key of a field. In Go, the `TagSpec` field is one `Pad`
value.

### TLV

A BER-TLV field has the same shape, but the tag is encoded in the data. It is
not a position. This is DE 55 in Go:

```go
55: field.NewComposite(&field.Spec{
	Length:      255,
	Description: "ICC Data",
	Pref:        prefix.ASCII.LLL,
	Tag: &field.TagSpec{
		Enc:                encoding.BerTLVTag,
		Sort:               sort.StringsByHex,
		SkipUnknownTLVTags: true,
	},
	Subfields: map[string]field.Field{
		"9F02": field.NewNumeric(&field.Spec{
			Description: "Amount, Authorised",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.BerTLV,
		}),
		"9F1E": field.NewString(&field.Spec{
			Description: "Interface Device Serial Number",
			Enc:         encoding.ASCII,
			Pref:        prefix.BerTLV,
		}),
	},
}),
```

This is the `ExportYAML` output for that spec:

```yaml
    "55":
        type: Composite
        length: 255
        description: ICC Data
        prefix: ASCII.LLL
        tag:
            enc: BerTLVTag
            sort: StringsByHex
            skipUnknownTLVTags: true
        subfields:
            9F1E:
                type: String
                description: Interface Device Serial Number
                enc: ASCII
                prefix: BerTLV
            9F02:
                type: Numeric
                description: Amount, Authorised
                enc: HexToASCII
                prefix: BerTLV
```

The subfields of one composite can have different encodings. In this example,
the serial number is ASCII text, and the amount is hex digits.

Be careful with these two points:

- Go and the document use different names for one encoding. Go uses
  `BytesToASCIIHex`. The document uses `HexToASCII`. The lists below give the
  names for the document, not the Go identifiers.
- The order of the subfields in the file is not the packing order. The `sort`
  key sets the packing order. For this reason, a `tag` block must have `sort`.
  If it does not, the import fails.

A field can also have two more keys, when they apply:

- `bitmap`: a nested field spec, for a composite that has a bitmap.
- `disableAutoExpand`: for a bitmap field.

## What the importer accepts

The importer rejects an unknown `type`, `enc`, `prefix` or `tag.sort`. The
error message gives the field that has the value.

The importer does not reject an unknown `padding.type` on a field, or an unknown
`tag.enc`. It ignores the value, and the spec has no padder or no tag encoding.
Make sure that these two values are in the lists.

**`type`**

`Binary`, `Bitmap`, `Composite`, `Hex`, `Numeric`, `String`, `Track1`, `Track2`,
`Track3`

**`enc`**

`ASCII`, `ASCIIToHex`, `BCD`, `BerTLVTag`, `Binary`, `EBCDIC`, `EBCDIC1047`,
`HexToASCII`, `LBCD`

**`prefix`**

`{ASCII|BCD|Binary|EBCDIC|Hex}.{Fixed|L|LL|LLL|LLLL}`, plus `BerTLV` and
`None.Fixed`

**`padding.type`**

`Left`, `Right`, `None`

**`tag.sort`**

`StringsByInt`, `StringsByHex`

A document gives the name of a sort function. It cannot contain the function.
Thus a document can only use a sort function that the importer knows. If a
dialect needs a different order, register the sort function, or keep the spec
in Go.

## Round trip

When this package exports a document, you can import the document again. The
result is an equivalent spec, with all composites and subfields. To make sure
that this is true for your spec, add a test to your repository:

1. Export your spec.
2. Import the document again.
3. Pack the same message with the two specs, and compare the results.

```go
raw, err := specs.ExportYAML(mySpec)
require.NoError(t, err)

back, err := specs.ImportYAML(raw)
require.NoError(t, err)
```

## What a document cannot carry

A document cannot contain a Go value. It can only contain a name. Thus a
document cannot contain these items:

- A custom field type.
- A custom encoder.
- A sort function that is not `StringsByInt` or `StringsByHex`.

If your spec uses one of these items, keep the spec in Go.

This limit also applies if you add behavior to a spec. The behavior stays after
an export and an import only if a document can give it as a name or as data.
