# moov-io/iso8583

[![GoDoc](https://godoc.org/github.com/moov-io/iso8583?status.svg)](https://godoc.org/github.com/moov-io/iso8583)
[![Build Status](https://github.com/moov-io/iso8583/workflows/Go/badge.svg)](https://github.com/moov-io/iso8583/actions)
[![Coverage Status](https://codecov.io/gh/moov-io/iso8583/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/iso8583)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/iso8583)](https://goreportcard.com/report/github.com/moov-io/iso8583)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/iso8583/master/LICENSE)

Package `github.com/moov-io/iso8583` implements a file reader and writer written in Go decorated with a HTTP API for creating, parsing, and validating financial transaction card originated interchange messaging.

## Project Status

Moov ISO8583 is under active development.

Please star the project if you are interested in its progress. If you have layers above ISO 8583 to simplify tasks or found bugs we would appreciate an issue or pull request. Thanks!


## Usage

Length encode and MTI encode types:

* bcd - BCD encoding of field length (only for Ll* and Lll* fields)
* ascii - ASCII encoding of field length (only for Ll* and Lll* fields)


Encode types:

* bcd - BCD encoding
* rbcd - BCD encoding with "right-aligned" value with odd length (for ex. "643" as [6 67] == "0643"), only for Numeric, Llnumeric and Lllnumeric fields
* ascii - ASCII encoding

### Example

```go
package main

import (
	"fmt"

	"github.com/moov-io/iso8583"
)

type Data struct {
	No   *iso8583.Numeric      `field:"3" length:"6" encode:"bcd"` // bcd value encoding
	Oper *iso8583.Numeric      `field:"26" length:"2" encode:"ascii"` // ascii value encoding
	Ret  *iso8583.Alphanumeric `field:"39" length:"2"`
	Sn   *iso8583.Llvar        `field:"45" length:"23" encode:"bcd,ascii"` // bcd length encoding, ascii value encoding
	Info *iso8583.Lllvar       `field:"46" length:"42" encode:"bcd,ascii"`
	Mac  *iso8583.Binary       `field:"64" length:"8"`
}

func main() {
	data := &Data{
		No:   iso8583.NewNumeric("001111"),
		Oper: iso8583.NewNumeric("22"),
		Ret:  iso8583.NewAlphanumeric("ok"),
		Sn:   iso8583.NewLlvar([]byte("abc001")),
		Info: iso8583.NewLllvar([]byte("你好 golang!")),
		Mac:  iso8583.NewBinary([]byte("a1s2d3f4")),
	}
	msg := iso8583.NewMessage("0800", data)
	msg.MtiEncode = iso8583.BCD
	b, err := msg.Bytes()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("% x\n", b)
}
```

Then you will get:

```
08 00 20 00 00 40 02 0c 00 01 00 11 11 32 32 6f 6b 06 61 62 63 30 30 31 00 14 e4 bd a0 e5 a5 bd 20 67 6f 6c 61 6e 67 21 61 31 73 32 64 33 66 34
```
## Guides

- A layman's guide to understanding the ISO8583 Financial Transaction Message [PDF](http://www.lytsing.org/downloads/iso8583.pdf)
- Wikipedia [ISO 8583](https://en.wikipedia.org/wiki/ISO_8583)
- Code Project - [Introduction to ISO 8583](https://www.codeproject.com/Articles/100084/Introduction-to-ISO)
- ISO 8583-1:2003 [Purchased specification](https://www.iso.org/standard/31628.html)

## Supported and Tested Platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows
- Rasberry Pi

## Contributing

Yes please! Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) to get started! Checkout our [issues for first time contributors](https://github.com/moov-io/iso8583/contribute) for something to help out with.

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) and uses Go 1.14 or higher. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/iso8583/releases/latest) as well. We highly recommend you use a tagged release for production.

### Releasing

To make a release of iso8583 simply open a pull request with `CHANGELOG.md` and `version.go` updated with the next version number and details. You'll also need to push the tag (i.e. `git push origin v1.0.0`) to origin in order for CI to make the release.


## License

Apache License 2.0 See [LICENSE](LICENSE) for details.

Work is derived from [@ideazxy](https://github.com/ideazxy/iso8583), [@zmwilliam](https://github.com/zmwilliam/iso8583), [@rdingwall](https://github.com/rdingwall/iso8583)
