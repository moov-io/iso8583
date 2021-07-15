[![Moov Banner Logo](https://user-images.githubusercontent.com/20115216/104214617-885b3c80-53ec-11eb-8ce0-9fc745fb5bfc.png)](https://github.com/moov-io)

<p align="center">
  <a href="https://github.com/moov-io/iso8583/tree/master/docs">Project Documentation</a>
  ·
  <a href="https://slack.moov.io/">Community</a>
  ·
  <a href="https://moov.io/blog/">Blog</a>
  <br>
  <br>
</p>

[![GoDoc](https://godoc.org/github.com/moov-io/iso8583?status.svg)](https://godoc.org/github.com/moov-io/iso8583)
[![Build Status](https://github.com/moov-io/iso8583/workflows/Go/badge.svg)](https://github.com/moov-io/iso8583/actions)
[![Coverage Status](https://codecov.io/gh/moov-io/iso8583/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/iso8583)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/iso8583)](https://goreportcard.com/report/github.com/moov-io/iso8583)
[![Repo Size](https://img.shields.io/github/languages/code-size/moov-io/iso8583?label=project%20size)](https://github.com/moov-io/iso8583)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/iso8583/master/LICENSE)
[![Slack Channel](https://slack.moov.io/badge.svg?bg=e01563&fgColor=fffff)](https://slack.moov.io/)
[![GitHub Stars](https://img.shields.io/github/stars/moov-io/iso8583)](https://github.com/moov-io/iso8583)
[![Twitter](https://img.shields.io/twitter/follow/moov_io?style=social)](https://twitter.com/moov_io?lang=en)

# moov-io/iso8583

Moov's mission is to give developers an easy way to create and integrate bank processing into their own software products. Our open source projects are each focused on solving a single responsibility in financial services and designed around performance, scalability, and ease of use.

ISO8583 implements an ISO 8583 message reader and writer in Go. ISO 8583 is an international standard for card-originated financial transaction messages that defines both message format and communication flow. It's used by major card networks around the globe including Visa, Mastercard, and Verve. The standard supports card purchases, withdrawals, deposits, refunds, reversals, balance inquiries, inter-account transfers, administrative messages, secure key exchanges, and more.

## Table of contents

- [Project status](#project-status)
- [Go module](#go-library)
	- [Define specification](#define-your-specification)
	- [Build message](#build-and-pack-the-message)
	- [Parse message](#parse-the-message-and-access-the-data)
	- [JSON encoding](#json-encoding)
	- [Network header](#network-header)
- [ISO8583 CLI](#cli)
- [Learn about ISO 8583](#learn-about-iso-8583)
- [Getting help](#getting-help)
- [Contributing](#contributing)
- [Related projects](#related-projects)

## Project status

Moov ISO8583 currently offers a Go package with plans for an API in the near future. Please star the project if you are interested in its progress. The project supports generating and parsing ISO8583 messages. Feedback on this early version of the project is appreciated and vital to its success. Please let us know if you encounter any bugs/unclear documentation or have feature suggestions by opening up an issue. Thanks!

## Go library

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) and Go v1.14 or higher. See [Golang's install instructions](https://golang.org/doc/install) for help in setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/iso8583/releases/latest) as well. We highly recommend you use a tagged release for production.

### Installation

```
go get github.com/moov-io/iso8583
```

### Define your specification

Our default specification [(spec87.go](./spec87.go)) is suitable for the majority of use cases. Simply instantiate a new message using Spec87:
```
isomessage := iso8583.NewMessage(iso8583.Spec87)
```
If this spec does not meet your needs, we encourage you to modify it or create your own using the information below. 

First, you need to define the format of the message fields that are described in your ISO8583 specification. Each data field has a type and its own spec. You can create a `NewBitmap`, `NewString`, or `NewNumeric` field. Each individual field spec consists of a few elements:

| Element          | Notes                                                                                                                                                                                                                       | Example                    |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------|
| `Length`         | Maximum length of field (bytes, characters or digits), for both fixed and variable lengths.                                                                                                                                 | `10`                       |
| `Description`    | Describes what the data field holds.                                                                                                                                                                                        | `"Primary Account Number"` |
| `Enc`            | Sets the encoding type (`ASCII`, `Hex`, `Binary`, `BCD`, `LBCD`).                                                                                                                                                           | `encoding.ASCII`           |
| `Pref`           | Sets the encoding (`ASCII`, `Hex`, `Binary`, `BCD`) of the field length and its type as fixed or variable (`Fixed`, `L`, `LL`, `LLL`, `LLLL`). The number of 'L's corresponds to the number of digits in a variable length. | `prefix.ASCII.Fixed`       |
| `Pad` (optional) | Sets padding direction and type.                                                                                                                                                                                            | `padding.Left('0')`        |

While some ISO8583 specifications do not have field 0 and field 1, we use them for MTI and Bitmap. Because technically speaking, they are just regular fields. We use field specs to describe MTI and Bitmap too. We currently use the `String` field for MTI, while we have a separate `Bitmap` field for the bitmap.

The following example creates a full specification with three individual fields (excluding MTI and Bitmap):

```go
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
			Enc:         encoding.Hex,
			Pref:        prefix.Hex.Fixed,
		}),

		// Message fields:
		2: field.NewString(&field.Spec{
			Length:      19,
			Description: "Primary Account Number",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		3: field.NewNumeric(&field.Spec{
			Length:      6,
			Description: "Processing Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		4: field.NewString(&field.Spec{
			Length:      12,
			Description: "Transaction Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
	},
}
```

### Build and pack the message

After the specification is defined, you can build a message. Having a binary representation of your message that's packed according to the provided spec lets you send it directly to a payment system! There are two ways to generate a message: typed or untyped.

Notice in the examples below, you do not need to set the bitmap value manually, as it is automatically generated for you during packing.

#### Untyped fields

If you don't care about types, and strings are fine for you, then you can easily set message fields like this:

```go
// create message with defined spec
message := NewMessage(spec)

// sets message type indicator at field 0
message.MTI("0100")

// set all message fields you need as strings
message.Field(2, "4242424242424242")
message.Field(3, "123456")
message.Field(4, "100")

// generate binary representation of the message into rawMessage
rawMessage, err := message.Pack()

// now you can send rawMessage over the wire
```

The binary message created will be `"01007000000000000000164242424242424242123456000000000100"`.

#### Typed fields

In many cases, you may want to work with specific types (numbers, strings, time, etc.) when setting message fields.

First, you need to define a struct that corresponds to the spec field types. Here is an example:

```go
// use the same types from message specification
type ISO87Data struct {
	F2 *field.StringField
	F3 *field.NumericField
	F4 *field.StringField
}

message := NewMessage(spec)
message.MTI("0100")

// now, pass data with fields into the message 
err := message.SetData(&ISO87Data{
	F2: field.NewStringValue("4242424242424242"),
	F3: field.NewNumericValue(123456),
	F4: field.NewStringValue("100"),
})

// pack the message and send it to your provider
rawMessage, err := message.Pack()
```

### Parse the message and access the data

When you have a binary (packed) message and you know the specification it follows, you can unpack it and access the data. Again, you have two options for data access: untyped and typed.

#### Untyped fields

With this approach you can easily access fields as strings:

```go
message := NewMessage(spec)
message.Unpack(rawMessage)

message.GetMTI() // MTI: 0100
message.GetString(2) // Card number: 4242424242424242
message.GetString(3) // Processing code: 123456
message.GetString(4) // Transaction amount: 100
```

#### Typed fields

To get typed field values just pass an empty struct for the data you want to access:

```go
// use the same types from message specification
type ISO87Data struct {
	F2 *field.StringField
	F3 *field.NumericField
	F4 *field.StringField
}

message := NewMessage(spec)
message.SetData(&ISO87Data{})

// let's unpack binary message
err := message.Unpack(rawMessage)

// to get access to typed data we have to get Data from the message and convert it into our ISO87Data type
data := message.Data().(*ISO87Data)

// now you have typed values
data.F2.Value // is a string "4242424242424242"
data.F3.Value // is an int 123456
data.F4.Value // is a string "100"
```

For complete code samples please check [./message_test.go](./message_test.go).

### JSON encoding

You can serialize message into JSON format:

```go
message := iso8583.NewMessage(spec)
message.MTI("0100")
message.Field(2, "4242424242424242")
message.Field(3, "123456")
message.Field(4, "100")

jsonMessage, err := json.Marshal(message)
```

it will produce following JSON:

```json
{
   "0":"0100",
   "1":"700000000000000000000000000000000000000000000000",
   "2":"4242424242424242",
   "3":123456,
   "4":"100"
}
```

### Network Header

All messages between the client/server (ISO host and endpoint) have a message
length header. It can be a 4 bytes ASCII or 2 bytes BCD encoded length. We
provide a `network.Header` interface to simplify the reading and writing of the
network header.

Following network headers are supported:

* ASCII4Bytes - message length encoded in 4 bytes ASCII, e.g., 0115 for 115
  bytes of the message
* BCD2Bytes - message length encoded in 2 bytes BCD, e.g, {0x01, x15} for 115
  bytes of the message

You can read network header from the network connection like this:

```go
header := network.NewBCD2BytesHeader()
_, err := header.ReadFrom(conn)
if err != nil {
	// handle error
}

// Make a buffer to hold message
buf := make([]byte, header.Length())
// Read the incoming message into the buffer.
read, err := conn.Read(buf)
if err != nil {
	// handle error
}
if reqLen != header.Length() {
	// handle error
}

message := iso8583.NewMessage(iso8583.Spec87)
message.Unpack(buf)
```

Here is an example of how to write network header into network connection:

```go
header := network.NewBCD2BytesHeader()
packed, err := message.Pack()
if err != nil {
	// handle error
}
header.SetLength(len(packed))
_, err = header.WriteTo(conn)
if err != nil {
	// handle error
}
n, err := conn.Write(packed)
if err != nil {
	// handle error
}
```

## CLI

CLI suports following command:

* `display` to display ISO8583 message in a human-readable format

### Installation

`iso8583` CLI is available as downloadable binaries from the [releases page](https://github.com/moov-io/iso8583/releases/latest) for MacOS, Windows and Linux.

Here is an example how to install MacOS version:

```
wget -O ./iso8583 https://github.com/moov-io/iso8583/releases/download/v0.4.6/iso8583_0.4.6_darwin_amd64 && chmod +x ./iso8583
```

Now you can run CLI:

```
➜ ./iso8583
Work seamlessly with ISO 8583 from the command line.

Usage:
  iso8583 <command> [flags]

Available commands:
  describe: display ISO 8583 file in a human-readable format
```


### Display

To display ISO8583 message in a human-readable format

Example:

```
➜ ./bin/iso8583 describe msg.bin
ISO 8583 Message:
MTI.............................: 0100
Bitmap..........................: A2380000008000000400000000000000
Bitmap bits.....................: 10100010 10100010 10100010 10100010 10100010 10100010 10100010 10100010
001 Processing Code.............: 000001
011 System Trace Audit Number...: 005835
```

## Learn about ISO 8583

- [Intro to ISO 8583](./docs/intro.md)
- [Message Type Indicator](./docs/mti.md)
- [Bitmaps](./docs/bitmap.md)
- [Data Fields](./docs/data-elements.md)
- [ISO 8583 Terms and Definitions](https://www.iso.org/obp/ui/#iso:std:iso:8583:-1:ed-1:v1:en)

## Getting help

 channel | info
 ------- | -------
[Project Documentation](https://github.com/moov-io/iso8583/tree/master/docs) | Our project documentation available online.
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/iso8583/issues/new) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel (`#iso8583`) to have an interactive discussion about the development of the project.

## Contributing

**While [Spec87](./spec87.go) is appropriate for most users, we hope to see improvements and variations of this specification for different systems by the community. Please do not hesitate to contribute issues, questions, or PRs to cover new use cases. Tests are also appreciated if possible!**

Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) to get started! Check out our [issues for first time contributors](https://github.com/moov-io/iso8583/contribute) for something to help out with.

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) and uses Go v1.14 or higher. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/iso8583/releases/latest) as well. We highly recommend you use a tagged release for production.

## Related projects
As part of Moov's initiative to offer open source fintech infrastructure, we have a large collection of active projects you may find useful:

- [Moov ACH](https://github.com/moov-io/ach) provides ACH file generation and parsing, supporting all Standard Entry Codes for the primary method of money movement throughout the United States.

- [Moov Watchman](https://github.com/moov-io/watchman) offers search functions over numerous trade sanction lists from the United States and European Union.

- [Moov Fed](https://github.com/moov-io/fed) implements utility services for searching the United States Federal Reserve System such as ABA routing numbers, financial institution name lookup, and FedACH and Fedwire routing information.

- [Moov Wire](https://github.com/moov-io/wire) implements an interface to write files for the Fedwire Funds Service, a real-time gross settlement funds transfer system operated by the United States Federal Reserve Banks.

- [Moov ImageCashLetter](https://github.com/moov-io/imagecashletter) implements Image Cash Letter (ICL) files used for Check21, X.9 or check truncation files for exchange and remote deposit in the U.S.

- [Moov Metro2](https://github.com/moov-io/metro2) provides a way to easily read, create, and validate Metro 2 format, which is used for consumer credit history reporting by the United States credit bureaus.


## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
