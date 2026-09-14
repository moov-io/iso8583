# Defining a Spec as a Document

A `MessageSpec` can be written in Go, and it can be written as YAML or JSON and
loaded at runtime. Both produce the same spec. This guide covers the second
form. The rest of the documentation builds specs in Go.

Not to be confused with [JSON Encoding and
Decoding](../README.md#json-encoding-and-decoding): that serialises a *message*,
its field values. This is about the *spec*, the definition those messages are
parsed against.

## Why a document

A spec in Go lives inside the binary. Changing a dialect means a build and a
deploy, and nothing outside this program can read the spec: not a validator, not
a documentation generator, not the operations team looking at what the switch
expects.

A spec as a document can be versioned on its own, diffed when a scheme bulletin
arrives, and loaded by a build that never changes.

## Reading and writing

```go
import "github.com/moov-io/iso8583/specs"

spec, err := specs.ImportYAML(raw)   // or ImportJSON
raw, err := specs.ExportYAML(spec)   // or ExportJSON
```

`ExportYAML` on a spec built in Go is the fastest way to get a first document:
build it once the way you already do, export it, and keep the file.

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

Field keys are strings because JSON object keys are. The `type`, `enc`, `prefix`
and `padding.type` values are names the importer resolves. The lists below give
every name it accepts.

A complete one ships with the library:
[`examples/specs/spec87ascii.yaml`](../examples/specs/spec87ascii.yaml), and the
same spec as
[`spec87ascii.json`](../examples/specs/spec87ascii.json).

## Composite fields

A composite carries its `subfields` and either a `tag` block or a `bitmap`,
exactly one of the two: a composite that declares neither cannot be built, and
the import fails. See the [composite fields guide](composite-fields.md) for what
each shape means.

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

The `tag` block also takes `length`, `enc`, `padding`, `skipUnknownTLVTags`,
`storeUnknownTLVTags` and `prefUnknownTLV`, matching `field.TagSpec`. Its
`padding` is the same `{type, pad}` block a field takes, even though the Go
field it fills is a single `Pad`.

### TLV

A BER-TLV field is the same shape with the tag encoded rather than positional.
Written in Go, DE 55 looks like this:

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

`ExportYAML` on that spec gives:

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

Subfields of one composite can be encoded differently, as they are here: the
serial number is ASCII text, the amount is hex digits.

Two things are easy to trip on. Go calls the encoding `BytesToASCIIHex`, and the
document calls it `HexToASCII`. The lists below give the document's names, not
the Go identifiers. And subfield order in the file is not the packing order.
`sort` decides that, which is why it has to be there.

Two more keys exist on a field and appear where they apply: `bitmap`, a nested
field spec for a composite that carries one, and `disableAutoExpand`, on a
bitmap field.

## What the importer accepts

The importer rejects anything outside these lists as it reads the document, and
names the field it came from.

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

A document names a sort rather than writing it out, so a spec can only use one
the importer knows. A dialect that needs its own ordering has to register it, or
stay in Go.

## Round trip

A document exported by this package imports back into an equivalent spec,
composites and their subfields included. That is worth a test in your own
repository: export the spec you use, import it back, and pack the same message
with both.

```go
raw, err := specs.ExportYAML(mySpec)
require.NoError(t, err)

back, err := specs.ImportYAML(raw)
require.NoError(t, err)
```

## What a document cannot carry

Everything in a spec that is a Go value rather than a name. A custom field
type, a custom encoder, a sort function that is not one of the two above: these
have no representation in the document, and a spec that uses them stays in Go.

This is also the boundary worth knowing about before designing anything that
puts behaviour in a spec. Behaviour survives the round trip when it can be
named or written down, and not otherwise.
