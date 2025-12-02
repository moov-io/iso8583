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
[![Twitter](https://img.shields.io/twitter/follow/moov?style=social)](https://twitter.com/moov?lang=en)

# moov-io/iso8583

Moov's mission is to give developers an easy way to create and integrate bank processing into their own software products. Our open source projects are each focused on solving a single responsibility in financial services and designed around performance, scalability, and ease of use.

ISO8583 implements an ISO 8583 message reader and writer in Go. ISO 8583 is an international standard for card-originated financial transaction messages that defines both message format and communication flow. It's used by major card networks around the globe including Visa, Mastercard, and Verve. The standard supports card purchases, withdrawals, deposits, refunds, reversals, balance inquiries, inter-account transfers, administrative messages, secure key exchanges, and more.

## Table of contents

- [Project status](#project-status)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Go version support policy](#go-version-support-policy)
- [How to](#how-to)
	- [Defining Message Specifications](#defining-message-specifications)
    - [Working with ISO 8583 Messages](#working-with-iso-8583-messages)
    - [Setting Message Data](#setting-message-data)
    - [Getting Message Data](#getting-message-data)
	- [Inspecting message fields](#inspecting-message-fields)
	- [JSON Encoding and Decoding](#json-encoding-and-decoding)
- [ISO8583 CLI](#cli)
- [Learn more](#learn-more)
- [Getting help](#getting-help)
- [Contributing](#contributing)
- [Related projects](#related-projects)

## Project status

Moov ISO8583 is a Go package that's been **thoroughly tested and trusted in the real world**. The project has proven its reliability and robustness in real-world, high-stakes scenarios. Please let us know if you encounter any missing feature/bugs/unclear documentation by opening up [an issue](https://github.com/moov-io/iso8583/issues/new) or asking on our [#iso8583 Community Slack channel](https://moov-io.slack.com/archives/C014UT7C3ST).
. Thanks!

## Installation

```
go get github.com/moov-io/iso8583
```

## Quick Start

The following example demonstrates how to:
- Define message structure using Go types
- Pack a message for transmission
- Unpack and parse a received message

```go
package main

import (
	"fmt"
	"os"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/examples"
)

// Define types for the message fields
type Authorization struct {
	MTI                  string               `iso8583:"0"`  // Message Type Indicator
	PrimaryAccountNumber string               `iso8583:"2"`  // PAN
	ProcessingCode       string               `iso8583:"3"`  // Processing code
	Amount               int64                `iso8583:"4"`  // Transaction amount
	STAN                 string               `iso8583:"11"` // System Trace Audit Number
	ExpirationDate       string               `iso8583:"14"` // YYMM
	AcceptorInformation  *AcceptorInformation `iso8583:"43"` // Merchant details
}

type AcceptorInformation struct {
	Name    string `iso8583:"1"`
	City    string `iso8583:"2"`
	Country string `iso8583:"3"`
}

func main() {
	// Pack the message
	msg := iso8583.NewMessage(examples.Spec)

	authData := &Authorization{
		MTI:                  "0100",
		PrimaryAccountNumber: "4242424242424242",
		ProcessingCode:       "000000",
		Amount:               2599,
		ExpirationDate:       "2201",
		AcceptorInformation: &AcceptorInformation{
			Name:    "Merchant Name",
			City:    "Denver",
			Country: "US",
		},
	}

	// Set the field values
	err := msg.Marshal(authData)
	if err != nil {
		panic(err)
	}

	// Pack the message
	packed, err := msg.Pack()
	if err != nil {
		panic(err)
	}

	// send packed message to the server
	// ...

	// Unpack the message
	msg = iso8583.NewMessage(examples.Spec)
	err = msg.Unpack(packed)
	if err != nil {
		panic(err)
	}

	// get individual field values
	var amount int64
	err = msg.UnmarshalPath("4", &amount)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Amount: %d\n", amount)

	// get value of composite subfield
	var acceptorName string
	err = msg.UnmarshalPath("43.1", &acceptorName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Acceptor Name: %s\n", acceptorName)

	// Get the field values into data structure
	authData = &Authorization{}
	err = msg.Unmarshal(authData)
	if err != nil {
		panic(err)
	}

	// Print the entire message
	iso8583.Describe(msg, os.Stdout)
}
```

## Go version support policy

### Always up-to-date, never left behind

While we strive to embrace the latest language enhancements, we also appreciate the need for a certain degree of backward compatibility. We understand that not everyone can update to the latest version immediately. Our philosophy is to move forward and embrace the new, but without leaving anyone immediately behind.

#### Which versions do we support now?

As of today, we are supporting the following versions as referenced in the [setup-go action step](https://github.com/actions/setup-go#using-stableoldstable-aliases):

* `stable` (which points to the current Go version)
* `oldstable` (which points to the previous Go version)

The [setup-go](https://github.com/actions/setup-go) action automatically manages versioning, allowing us to always stay aligned with the latest and preceding Go releases.

#### What does this mean for you?

Whenever a new version of Go is released, we will update our systems and ensure that our project remains fully compatible with it. At the same time, we will continue to support the previous version. However, once a new version is released, the 'previous previous' version will no longer be officially supported. 

#### Continuous integration

To ensure our promise of support for these versions, we've configured our GitHub CI actions to test our code with both the current and previous versions of Go. This means you can feel confident that the project will work as expected if you're using either of these versions.

## How to

### Defining Message Specifications

Most ISO 8583 implementations use confidential specifications that vary between payment systems, so you'll likely need to create your own specification. We provide example specifications in [/specs](./specs) directory that you can use as a starting point.

#### Core Concepts

The package maps ISO 8583 concepts to the following types:

- **MessageSpec** - Defines the complete message format with fields
- **field.Spec** - Defines field's structure and behavior
- **field.Field** - Represents an ISO 8583 data element with its value storing and handling logic:
  - `field.String` - For alphanumeric fields
  - `field.Numeric` - For numeric fields
  - `field.Binary` - For binary data fields 
  - `field.Composite` - For structured data like TLV/BER-TLV fields or fields with positional subfields

Each field specification consists of these elements:

| Element          | Notes                                                                                                                                                                                                                       | Example                    |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------|
| `Length`         | Maximum length of field (bytes, characters or digits), for both fixed and variable lengths.                                                                                                                                 | `10`                       |
| `Description`    | Describes what the data field holds.                                                                                                                                                                                        | `"Primary Account Number"` |
| `Enc`            | Sets the encoding type (`ASCII`, `Binary`, `BCD`, `LBCD`, `EBCDIC`).                                                                                                                                                           | `encoding.ASCII`           |
| `Pref`           | Sets the encoding (`ASCII`, `Binary`, `BCD`, `EBCDIC`) of the field length and its type as fixed or variable (`Fixed`, `L`, `LL`, `LLL`, `LLLL`). The number of 'L's corresponds to the number of digits in a variable length. | `prefix.ASCII.Fixed`       |
| `Pad` (optional) | Sets padding direction and type.                                                                                                                                                                                            | `padding.Left('0')`        |

Note: While some ISO 8583 specifications do not have field 0 and field 1, we use them for MTI and Bitmap as they are technically regular fields. We use `String` field for MTI and `Bitmap` field for the bitmap.

For more advanced examples including handling of BER-TLV data, positional subfields, and various encoding types, see:
- [message_test.go](message_test.go) - Complex message specifications and field types
- [field/composite_test.go](field/composite_test.go) - Working with composite fields and subfields

### Working with ISO 8583 Messages

The package provides two key operations for working with ISO 8583 messages:

#### Building Messages

To build a message:
1. Set message data using Go structs or individual field operations
2. Pack the message into bytes using `Pack()`
3. Send the bytes over the network

```
Set Data → Pack → Network →
```

#### Processing Messages

When receiving a message:
1. Unpack the received bytes using `Unpack()`
2. Get message data using Go structs or individual field operations
3. Process the data in your application

```
→ Network → Unpack → Get Data
```

### Setting Message Data

After defining your specification, you can set message data in two ways: working with individual fields or using Go structs. While individual field access is available, using structs provides a cleaner approach.

#### Using Go Structs (Recommended)

Define a struct that maps your business data to ISO 8583 fields using native Go types and `iso8583` tags:

```go
type Authorization struct {
    MTI                  string    `iso8583:"0"`  // Message Type Indicator
    PrimaryAccountNumber string    `iso8583:"2"`  // PAN
    ProcessingCode       string    `iso8583:"3"`  // Processing code
    Amount               int64     `iso8583:"4"`  // Transaction amount
    STAN                 string    `iso8583:"11"` // System Trace Audit Number
    LocalTime            string    `iso8583:"12"` // HHmmss
    LocalDate            string    `iso8583:"13"` // MMDD
    MerchantType         string    `iso8583:"18"` // Merchant category code
    AcceptorInfo         *Acceptor `iso8583:"43"` // Merchant details
}

type Acceptor struct {
    Name    string `iso8583:"1"`
    City    string `iso8583:"2"`
    Country string `iso8583:"3"`
}
```

Then create and populate your message:

```go
// Create new message
msg := iso8583.NewMessage(Spec)

// Prepare transaction data
auth := &Authorization{
    MTI:                  "0100",
    PrimaryAccountNumber: "4242424242424242",
    ProcessingCode:       "000000",
    Amount:               9999,
    STAN:                 "000123",
    LocalTime:            "152059",
    LocalDate:            "0205",
    MerchantType:         "5411",
    AcceptorInfo: &Acceptor{
        Name:    "ACME Store",
        City:    "New York",
        Country: "US",
    },
}

// Marshal data into message
err := msg.Marshal(auth)
if err != nil {
    panic(err)
}

// Pack for transmission
data, err := msg.Pack()
if err != nil {
    panic(err)
}
// data is ready to be sent
```

If you want empty values to be included in the message, you can use the `keepzero` option in the `iso8583` field tag:

```go
type Authorization struct {
    // ...
    AdditionalData       string    `iso8583:"48,keepzero"` // Additional data
}
```

For such fields, the field bit will be set in the bitmap, but the value will be empty. You should set padding for such fields to ensure the field length is correct.

#### Working with Individual Fields

<details>
<summary>Click to show individual field operations</summary>

For simple operations, you can set fields directly:

```go
msg := iso8583.NewMessage(spec)

// Set MTI
msg.MTI("0100")

// Set field values as strings
err := msg.Field(2, "4242424242424242")
err = msg.Field(3, "000000")
err = msg.Field(4, "9999")

// For binary fields
err = msg.BinaryField(52, []byte{0x1A, 0x2B, 0x3C, 0x4D})

// Pack for transmission
data, err := msg.Pack()
```

Note: Individual field operations are limited to string or []byte values. The underlying field handles type conversion.
</details>

#### Legacy: Using Field Types

<details>
<summary>Click to show legacy field type usage</summary>

In previous versions, it was common to use package-specific field types. While still supported, we recommend using native Go types instead:

```go
type Authorization struct {
    MTI         *field.String `iso8583:"0"`
    PAN         *field.String `iso8583:"2"`
    Amount      *field.Numeric `iso8583:"4"`
    LocalTime   *field.String `iso8583:"12"`
}

auth := &Authorization{
    MTI:       field.NewStringValue("0100"),
    PAN:       field.NewStringValue("4242424242424242"),
    Amount:    field.NewNumericValue(9999),
    LocalTime: field.NewStringValue("152059"),
}

msg.Marshal(auth)
```
</details>

### Getting Message Data

When you receive a packed ISO 8583 message, you can unpack and access its data in two ways: using Go structs or accessing individual fields. Using structs provides a cleaner approach to working with message data.

#### Using Go Structs (Recommended)

Define a struct matching your expected message format using native Go types:

```go
type Authorization struct {
    MTI                  string    `iso8583:"0"`  // Message Type Indicator
    PrimaryAccountNumber string    `iso8583:"2"`  // PAN
    ProcessingCode       string    `iso8583:"3"`  // Processing code
    Amount              int64     `iso8583:"4"`  // Transaction amount
    STAN                string    `iso8583:"11"` // System Trace Audit Number
    LocalTime           string    `iso8583:"12"` // HHmmss
    LocalDate           string    `iso8583:"13"` // MMDD
    MerchantType        string    `iso8583:"18"` // Merchant category code
    AcceptorInfo        *Acceptor `iso8583:"43"` // Merchant details
}

type Acceptor struct {
    Name    string `iso8583:"1"`
    City    string `iso8583:"2"`
    Country string `iso8583:"3"`
}
```

Then unpack and parse your message:

```go
// Create message with appropriate spec
msg := iso8583.NewMessage(specs.Spec87ASCII)

// Unpack received bytes
err := msg.Unpack(receivedData)
if err != nil {
    panic(err)
}

// Parse into struct
var auth Authorization
err = msg.Unmarshal(&auth)
if err != nil {
    panic(err)
}

// Work with parsed data
fmt.Printf("Transaction amount: %d\n", auth.Amount)
fmt.Printf("Merchant: %s, %s\n", auth.AcceptorInfo.Name, auth.AcceptorInfo.City)

// Print full message contents (with sensitive data masked)
iso8583.Describe(msg, os.Stdout)
```

#### Working with Individual Fields

<details>
<summary>Click to show individual field operations</summary>

For quick access to specific fields:

```go
msg := iso8583.NewMessage(spec)
err := msg.Unpack(receivedData)

// Get MTI
mti, err := msg.GetMTI()

// Get field values as strings
pan, err := msg.GetString(2)
proc, err := msg.GetString(3)
amount, err := msg.GetString(4)

// Get binary field values
pinData, err := msg.GetBytes(52)
```

Note: Individual field access is limited to string or []byte values.
</details>

#### Legacy: Using Field Types

<details>
<summary>Click to show legacy field type usage</summary>

Previously, it was common to use package-specific field types. While still supported, we recommend using native Go types instead:

```go
type LegacyAuthorization struct {
    MTI         *field.String  `iso8583:"0"`
    PAN         *field.String  `iso8583:"2"`
    Amount      *field.Numeric `iso8583:"4"`
    LocalTime   *field.String  `iso8583:"12"`
}

var auth LegacyAuthorization
msg.Unmarshal(&auth)

fmt.Println(auth.MTI.Value())
fmt.Println(auth.Amount.Value())
```
</details>

### Inspecting Message Fields

There is a `Describe` function in the package that displays all message fields
in a human-readable way. Here is an example of how you can print message fields
with their values to STDOUT:

```go
// print message to os.Stdout
iso8583.Describe(message, os.Stdout)
```

and it will produce the following output:

```
MTI........................................: 0100
Bitmap.....................................: 000000000000000000000000000000000000000000000000
Bitmap bits................................: 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
F000 Message Type Indicator................: 0100
F002 Primary Account Number................: 4242****4242
F003 Processing Code.......................: 123456
F004 Transaction Amount....................: 100
F020 PAN Extended Country Code.............: 4242****4242
F035 Track 2 Data..........................: 4000****0506=2512111123400001230
F036 Track 3 Data..........................: 011234****3445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**
F045 Track 1 Data..........................: B4815****1896^YATES/EUGENE L^^^356858      00998000000
F052 PIN Data..............................: 12****78
F055 ICC Data – EMV Having Multiple Tags...: ICC  ... Tags
```

by default, we apply `iso8583.DefaultFilters` to mask the values of the fields
with sensitive data. You can define your filter functions and redact specific
fields like this:

```go
filterAll = func(in string, data field.Field) string {
	runesInString := utf8.RuneCountInString(in)

	return strings.Repeat("*", runesInString)
}

// filter only value of the field 2
iso8583.Describe(message, os.Stdout, filterAll(2, filterAll))

// outputs:
// F002 Primary Account Number................: ************
```

If you want to view unfiltered values, you can use no-op filters `iso8583.DoNotFilterFields` that we defined:

```go
// display unfiltered field values
iso8583.Describe(message, os.Stdout, DoNotFilterFields()...)
```

### JSON Encoding and Decoding

You can serialize message into JSON format:

```go
jsonMessage, err := json.Marshal(message)
```

it will produce the following JSON (bitmap is not included, as it's only used to unpack message from the binary representation):

```json
{
   "0":"0100",
   "2":"4242424242424242",
   "3":123456,
   "4":"100"
}
```

Also, you can unmarshal JSON into `iso8583.Message`:

```go
input := `{"0":"0100","2":"4242424242424242","4":"100"}`

message := NewMessage(spec)
if err := json.Unmarshal([]byte(input), message); err != nil {
    // handle err
}

// access indidual fields or using struct
```

### Sending and Receiving Messages

While this package handles message formatting and parsing, for network operations we recommend using our companion package [moov-io/iso8583-connection](https://github.com/moov-io/iso8583-connection). It provides robust client/server communication with features like:

- Message sending and receiving
- Request/response matching
- Connection management
- Support for both acquiring and issuing services
- Testing utilities

#### Network Headers

All messages between the client/server (ISO host and endpoint) have a message length header. It can be a 4 bytes ASCII or 2 bytes BCD encoded length or any other custom header format. We provide a `network.Header` interface to simplify the reading and writing of the network header.

Following network headers implementations are available in [network](./network) package:

* Binary2Bytes - message length encoded in 2 bytes, e.g, {0x00 0x73} for 115
  bytes of the message
* ASCII4Bytes - message length encoded in 4 bytes ASCII, e.g., 0115 for 115
  bytes of the message
* BCD2Bytes - message length encoded in 2 bytes BCD, e.g, {0x01, 0x15} for 115
  bytes of the message
* VMLH (Visa Message Length Header) - message length encoded in 2 bytes + 2 reserved bytes

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
read, err := io.ReadFull(conn, buf)
if err != nil {
	// handle error
}
if reqLen != header.Length() {
	// handle error
}

message := iso8583.NewMessage(specs.Spec87ASCII)
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
MTI........................................: 0100
Bitmap.....................................: 000000000000000000000000000000000000000000000000
Bitmap bits................................: 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
F000 Message Type Indicator................: 0100
F002 Primary Account Number................: 4242****4242
F003 Processing Code.......................: 123456
F004 Transaction Amount....................: 100
F020 PAN Extended Country Code.............: 4242****4242
F035 Track 2 Data..........................: 4000****0506=2512111123400001230
F036 Track 3 Data..........................: 011234****3445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**
F045 Track 1 Data..........................: B4815****1896^YATES/EUGENE L^^^356858      00998000000
F052 PIN Data..............................: 12****78
F055 ICC Data – EMV Having Multiple Tags...: ICC  ... Tags
```

You can specify which of the built-in specs to use to the describe message via
the `spec` flag:

```
➜ ./bin/iso8583 describe -spec spec87ascii msg.bin
```

You can also define your spec in JSON format and describe message using the spec file with `spec-file` flag:

```
➜ ./bin/iso8583 describe -spec-file ./examples/specs/spec87ascii.json msg.bin
```

Please, check the example of the JSON spec file [spec87ascii.json](./examples/specs/spec87ascii.json).


## Learn more

- [How to Define Composite Fields](./docs/composite-fields.md)
- [Intro to ISO 8583](./docs/intro.md)
- [Message Type Indicator](./docs/mti.md)
- [Bitmaps](./docs/bitmap.md)
- [How Tos](./docs/howtos.md)
- [Data Fields](./docs/data-elements.md)
- [Mastering ISO 8583 messages with Golang](https://alovak.com/2024/08/15/mastering-iso-8583-messages-with-golang/)
- [Mastering ISO 8583 Message Networking with Golang](https://alovak.com/2024/08/27/mastering-iso-8583-message-networking-with-golang/)
- [ISO 8583 Terms and Definitions](https://www.iso.org/obp/ui/#iso:std:iso:8583:-1:ed-1:v1:en)

## Getting help

 channel | info
 ------- | -------
[Project Documentation](https://github.com/moov-io/iso8583/tree/master/docs) | Our project documentation available online.
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/iso8583/issues/new) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel (`#iso8583`) to have an interactive discussion about the development of the project.

## Contributing

**While [Spec87ASCII](./specs/spec87ascii.go) is appropriate for most users, we hope to see improvements and variations of this specification for different systems by the community. Please do not hesitate to contribute issues, questions, or PRs to cover new use cases. Tests are also appreciated if possible!**

Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) to get started! Check out our [issues for first time contributors](https://github.com/moov-io/iso8583/contribute) for something to help out with.

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.18 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/iso8583/releases/latest) as well. We highly recommend you use a tagged release for production.

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
