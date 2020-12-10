# Go API client for client

Package github.com/moov-io/iso8583 implements a file reader and writer written in Go decorated with a HTTP API for creating, parsing, and validating financial transaction card originated interchange messaging. User can seed iso8583's specification file as json file, iso8583 message with several formsts (mesage binary, json, xml)
 | Input      | Output     |
 |------------|------------|
 | JSON       | JSON       |
 | XML        | XML        |
 | MESSAGE    | MESSAGE    |
 

## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 0.0.1
- Package version: 1.0.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/oauth2
go get golang.org/x/net/context
go get github.com/antihax/optional
```

Put the package under your project folder and add the following in import:

```golang
import "./client"
```

## Documentation for API Endpoints

All URIs are relative to *https://local.moov.io:8208*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*Iso8583MessageApi* | [**Convert**](docs/Iso8583MessageApi.md#convert) | **Post** /convert | Convert iso8583 message
*Iso8583MessageApi* | [**Health**](docs/Iso8583MessageApi.md#health) | **Get** /health | health iso8583 service
*Iso8583MessageApi* | [**Print**](docs/Iso8583MessageApi.md#print) | **Post** /print | Print iso8583 message with specific format
*Iso8583MessageApi* | [**Validator**](docs/Iso8583MessageApi.md#validator) | **Post** /validator | Validate iso8583 message


## Documentation For Models

 - [Attribute](docs/Attribute.md)
 - [AttributeItem](docs/AttributeItem.md)
 - [Element](docs/Element.md)
 - [Encoding](docs/Encoding.md)
 - [IsoMessage](docs/IsoMessage.md)
 - [Specification](docs/Specification.md)


## Documentation For Authorization



## GatewayAuth

- **Type**: HTTP basic authentication

Example

```golang
auth := context.WithValue(context.Background(), sw.ContextBasicAuth, sw.BasicAuth{
    UserName: "username",
    Password: "password",
})
r, err := client.Service.Operation(auth, args)
```



## Author


