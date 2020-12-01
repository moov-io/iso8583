# moov-io/iso8583

Package `github.com/moov-io/iso8583` implements ISO8583 standard in GO.

...

## Getting Started

### Install

```
go get github.com/moov-io/iso8583
```

### Define your spec

First, you need to define format of the the message fields that are described in your ISO8583 specification.
Here is how you can do this:

```go
spec87 := &spec.MessageSpec{
    Fields: map[int]spec.Packer{
        0: spec.NewField("Message Type Indicator", encoding.ASCII, prefixer.ASCII.Fixed(4)),

        // Bitmap, 16 bytes, fixed
        1: spec.Bitmap("Bitmap", encoding.Hex, prefixer.Hex.Fixed(16)),

        // LLVAR19
        2: spec.NewField("Primary Account Number", encoding.ASCII, prefixer.ASCII.LL(19)),

        // 6 bytes, fixed
        3: spec.NewField("Processing Code", encoding.ASCII, prefixer.ASCII.Fixed(6)),

        // 12 bytes, fixed
        4: spec.NewField("Transaction Amount", encoding.ASCII, prefixer.ASCII.Fixed(12)),
	
	// ...
    },
}
```

### Build and pack the message

When specification is defined it's time to build the message:

```go
message := iso8583.NewMessage(spec87)
message.Field(0, "0100")
message.Field(2, "4242424242424242")
message.Field(3, "123456")
message.Field(4, "000000000100")


binaryData, err := message.Pack()

// ...
```

Having a binary representation of your message that's packed according to the provided spec lets you send it directly to the payment system!

### Parse the message

When you have a binary (packed) message and you know the specification it follows, you can unpack it and access the data easily:

```go
message := iso8583.NewMessage(spec87)
message.Unpack(binaryData)

message.GetString(0) // MTI: 0100
message.GetString(2) // Card number: 4242424242424242
// ...

```

# License

Apache License 2.0 See [LICENSE](LICENSE) for details.
