# moov-io/iso8583

[![GoDoc](https://godoc.org/github.com/moov-io/iso8583?status.svg)](https://godoc.org/github.com/moov-io/iso8583)
[![Build Status](https://github.com/moov-io/iso8583/workflows/Go/badge.svg)](https://github.com/moov-io/iso8583/actions)
[![Coverage Status](https://codecov.io/gh/moov-io/iso8583/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/iso8583)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/iso8583)](https://goreportcard.com/report/github.com/moov-io/iso8583)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/iso8583/master/LICENSE)

Package `github.com/moov-io/iso8583` implements a message reader and writer written in Go decorated with a HTTP API for creating, parsing, and validating financial transaction card originated interchange messaging.

Docs: [API Endpoints](https://moov-io.github.io/iso8583/api/)

## Getting Started

### Docker

We publish a [public Docker image `moov/iso8583`](https://hub.docker.com/r/moov/iso8583/tags) on Docker Hub with tagged release of the package. No configuration is required to serve on `:8080`.


Start the Docker image:
```
docker run -p 8080:8080 moov/iso8583:latest
```

Upload a file and validate it
```
curl -XPOST --form "input=@./test/testdata/iso_reversal_message_advice.dat" http://localhost:8080/validator
```
```
{"status":"valid file"}
```
with specification file
```
curl -XPOST --form "input=@./test/testdata/iso_reversal_message_advice.dat" --form "spec=@./test/testdata/specification_ver_1987.json" http://localhost:8080/validator
```
```
{"status":"valid file"}
```

Convert a message between formats
```
curl -XPOST --form "file=@./test/testdata/iso_reversal_message_advice.dat" --form "format=json" http://localhost:8080/convert
```
```
{
	"mti": "0420",
	"bitmap": "0111001000110000000000001000000100000000000000000000000000000000",
	"elements": {
		"2": "",
		"3": "180000",
		"4": "000000030000",
		"7": "0109080646",
		"11": "100331",
		"12": "001120",
		"25": "33",
		"32": "00011122233"
	}
}
```
```
curl -XPOST --form "file=@./test/testdata/iso_reversal_message_advice.dat" --form "format=xml" http://localhost:8080/convert
```
```
<isoMessage>
	<MTI>0420</MTI>
	<Bitmap>0111001000110000000000001000000100000000000000000000000000000000</Bitmap>
	<DataElements>
		<Element Number="2"></Element>
		<Element Number="3">180000</Element>
		<Element Number="4">000000030000</Element>
		<Element Number="7">0109080646</Element>
		<Element Number="11">100331</Element>
		<Element Number="12">001120</Element>
		<Element Number="25">33</Element>
		<Element Number="32">00011122233</Element>
	</DataElements>
</isoMessage>
```

### Go Library

There is a Go library which can read and write iso8583 message. We write unit tests and fuzz the code to help ensure our code is production ready for everyone. Iso8583 uses [Go Modules](https://github.com/golang/go/wiki/Modules) to manage dependencies and suggests Go 1.14 or greater.

To clone our code and verify our tests on your system run:

```
$ git clone git@github.com:moov-io/iso8583.git
$ cd iso8583

$ go test ./...
ok      github.com/moov-io/iso8583      0.015s
ok      github.com/moov-io/iso8583/cmd/iso8583  21.908s
?       github.com/moov-io/iso8583/pkg/client   [no test files]
ok      github.com/moov-io/iso8583/pkg/lib      0.137s
ok      github.com/moov-io/iso8583/pkg/server   5.901s
ok      github.com/moov-io/iso8583/pkg/utils    0.028s
```

## Formats and configuration file
### message formats
Iso8583 have supported 3 message types: iso8583, json, xml.
Iso8583 specification defines a message format, but don't define json and xml format.
Iso8583 package have defined json and xml format of message, specification file (configuration file) that use to define message structure.

Json format:
```
{
	"mti": "0420",
	"bitmap": "0111001000110000000000001000000100000000000000000000000000000000",
	"elements": {
		"2": "",
		"3": "180000",
		...
	}
}
```
Bitmap is binary string and isn't hex string, value type of elements is string.

XML format:
```
<isoMessage>
	<MTI>0420</MTI>
	<Bitmap>0111001000110000000000001000000100000000000000000000000000000000</Bitmap>
	<DataElements>
		<Element Number="2"></Element>
		<Element Number="3">180000</Element>
        ...
	</DataElements>
</isoMessage>
```
Bitmap is binary string and isn't hex string.

### message specification file (configuration file)
The first digit of the MTI indicates the iso8583 version in which the message is encoded.
Iso8583 message structure are difference between version.
User can manage iso8583 message with several versions of iso8583 using the specification file feature (configuration file)
message specification file supported json format only.

```
{
	"elements": {
		"1": {
			"Describe": "b 64",
			"Description": "Second Bitmap"
		},
		"2": {
			"Describe": "n..19",
			"Description": "Primary account number (PAN)"
		},
        ...
    },
	"encoding": {
		"mti_enc": "CHAR",
		"bmp_enc": "HEX",
		"len_enc": "CHAR",
		"num_enc": "CHAR",
		"chr_enc": "ASCII",
		"bin_enc": "HEX",
		"trk_enc": "EBCDIC"
	},
    "message_types": {
        "0100" : {
            "mandatory_hex_mask": "72300000000000000000000000000000",
            "optional_hex_mask": "000C0661A9C000000000000000000000",
        }
        ...
    },
```
Data element is specify using describe that indicate data type and length.
In many case describes are attributes of message data element in iso8583 specification document.
Encoding define encoding/decoding type about any part of message.
Available types are "CHAR", "HEX", "EBCDIC", "ASCII", "BCD", "RBCD".
Message Types define mandatory fields and optional fields of message using hex string.

## Commands

iso8583 has command line interface to manage iso8583 messages and to lunch web service.

```
iso8583 --help

Usage:
   [command]

Available Commands:
  convert     Convert iso8583 message format
  help        Help about any command
  print       Print iso8583 message
  validator   Validate iso8583 message
  web         Launches web server

Flags:
  -h, --help           help for this command
      --input string   iso8583 message (the message types are iso8583 raw message, xml, json. default is $PWD/iso8583_message.dat)
      --spec string    specification file (default is $PWD/iso8583_specification.json)

Use " [command] --help" for more information about a command.
```

Each interaction that the library supports is exposed in a command-line option:

 Command | Info
 ------- | -------
`convert` | The convert command allows users to convert from a iso8583 message to another message format. Result will create a iso8583 message.
`print` | The print command allows users to print a iso8583 message with special file format (json, xml, iso8583).
`validator` | The validator command allows users to validate a iso8583 message.
`web` | The web command will launch a web server with endpoints to manage iso8583 messages.

### message convert

```
iso8583 convert --help

Usage:
   convert [output] [flags]

Flags:
      --format string   format of iso8583 message(required) (default "iso8583")
  -h, --help            help for convert

Global Flags:
      --input string   iso8583 message (the message types are iso8583 raw message, xml, json. default is $PWD/iso8583_message.dat)
      --spec string    specification file (default is $PWD/iso8583_specification.json)
```

The output parameter is the full path name to convert new iso8583 message.
The format parameter is supported 3 types that are "json", "xml" and  "iso8583".
The input parameter is source iso8583 message, supported "json", "xml" and  "iso8583".
The spec parameter is specification file.

example:
```
iso8583 convert output/output.json --input testdata/iso_reversal_message_advice.dat --format json
```

### message print

```
iso8583 print --help

Usage:
   print [flags]

Flags:
      --format string   print format (default "iso8583")
  -h, --help            help for print

Global Flags:
      --input string   iso8583 message (the message types are iso8583 raw message, xml, json. default is $PWD/iso8583_message.dat)
      --spec string    specification file (default is $PWD/iso8583_specification.json)
```

The format parameter is supported 3 types that are "json", "xml" and  "iso8583".
The input parameter is source iso8583 message, supported "json", "xml" and  "iso8583".
The spec parameter is specification file.

example:
```
iso8583 print --input testdata/iso_reversal_message_advice.dat --format json
{
	"mti": "0420",
	"bitmap": "0111001000110000000000001000000100000000000000000000000000000000",
	"elements": {
		"2": "",
		"3": "180000",
		"4": "000000030000",
		"7": "0109080646",
		"11": "100331",
		"12": "001120",
		"25": "33",
		"32": "00011122233"
	}
}
```

### message validate

```
iso8583 validator --help

Usage:
   validator [flags]

Flags:
  -h, --help   help for validator

Global Flags:
      --input string   iso8583 message (the message types are iso8583 raw message, xml, json. default is $PWD/iso8583_message.dat)
      --spec string    specification file (default is $PWD/iso8583_specification.json)
```

The input parameter is source iso8583 message, supported "json", "xml" and  "iso8583".

example:
```
iso8583 validator --input testdata/iso_reversal_message_advice.dat
```

### web server

```
iso8583 web --help

Usage:
   web [flags]

Flags:
  -h, --help   help for web
  -t, --test   test server

Global Flags:
      --input string   iso8583 message (the message types are iso8583 raw message, xml, json. default is $PWD/iso8583_message.dat)
      --spec string    specification file (default is $PWD/iso8583_specification.json)
```

The port parameter is port number of web service.

example:
```
iso8583 web
```

Web server have some endpoints to manage iso8583 messages

Method | Endpoint | Content-Type | Info
 ------- | ------- | ------- | -------
 `POST` | `/convert` | multipart/form-data | convert iso8583 messages. will download new file.
 `GET` | `/health` | text/plain | check web server.
 `POST` | `/print` | multipart/form-data | print iso8583 messages.
 `POST` | `/validator` | multipart/form-data | validate iso8583 messages.

web page example to use iso8583 web server:

```
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Test ISO8583 APIs</title>
</head>
<body>
<h1>Upload single file with fields</h1>

<form action="http://localhost:8080/convert" method="post" enctype="multipart/form-data">
    Format: <input type="text" name="format"><br>
    Input File: <input type="file" name="input"><br><br>
    Input File: <input type="file" name="spec"><br><br>
    <input type="submit" value="Submit">
</form>
</body>
</html>
```

## Docker

You can run the [moov/iso8583 Docker image](https://hub.docker.com/r/moov/iso8583) which defaults to starting the HTTP server.

```
docker run -p 8080:8080 moov/iso8583:latest
```

## Getting Help

 channel | info
 ------- | -------
  Google Group [moov-users](https://groups.google.com/forum/#!forum/moov-users)| The Moov users Google group is for contributors other people contributing to the Moov project. You can join them without a google account by sending an email to [moov-users+subscribe@googlegroups.com](mailto:moov-users+subscribe@googlegroups.com). After receiving the join-request message, you can simply reply to that to confirm the subscription.
Twitter [@moov_io](https://twitter.com/moov_io)	| You can follow Moov.IO's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/iso8583/issues) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel (`#iso8583`) to have an interactive discussion about the development of the project.

## Supported and Tested Platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## Contributing

Yes please! Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](https://github.com/moov-io/ach/blob/master/CODE_OF_CONDUCT.md) to get started! [Checkout our issues](https://github.com/moov-io/iso8583/issues) for something to help out with.

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) and uses Go 1.14 or higher. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/iso8583/releases/latest) as well. We highly recommend you use a tagged release for production.

## License

Apache License 2.0 See [LICENSE](LICENSE) for details.
