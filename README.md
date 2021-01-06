# alovak/iso8583

Package `github.com/alovak/iso8583` implements ISO8583 standard in GO.

...

## Getting Started

### Install

```
go get github.com/alovak/iso8583
```

### Define your spec

First, you need to define the format of the message fields that are described in your ISO8583 specification.
Here is how you can do this:

```go
spec := &MessageSpec{
	Fields: map[int]field.Field{
		0: field.NewStringField(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		1: field.NewBitmapField(&field.Spec{
			Length:      16,
			Description: "Bitmap",
			Enc:         encoding.Hex,
			Pref:        prefix.Hex.Fixed,
		}),
		2: field.NewStringField(&field.Spec{
			Length:      19,
			Description: "Primary Account Number",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		3: field.NewNumericField(&field.Spec{
			Length:      6,
			Description: "Processing Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		4: field.NewStringField(&field.Spec{
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

When the specification is defined it's time to build the message. There are two ways to do this: typed and untyped.

#### Untyped fields

If you don't worry about types and strings are ok for you, then you can easily set message fields like this:


```go
// create message with defined spec
message := iso8583.NewMessage(spec87)

message.MTI("0100")

// set all message fields you need as strings
message.Field(2, "4242424242424242")
message.Field(3, "123456")
message.Field(4, "000000000100")

// get binary representation of the message into rawMessage
rawMessage, err := message.Pack()

// now you can send rawMessage over the wire
```

#### Typed fields

In many cases, you may want to work with types: numbers, strings, time, etc. We have got that covered!

First, you need to define the struct that corresponds to the spec field types. Here an example:

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


Having a binary representation of your message that is packed according to the provided spec lets you send it directly to the payment system!

### Parse the message and access the data

When you have a binary (packed) message and you know the specification it follows, you can unpack it and access the data. Again, you have two options for data access: untyped and typed.

#### Untyped fields

With this approach you can access fields as strings no sweat:

```go
message := iso8583.NewMessage(spec)
message.Unpack(binaryData)

message.GetMTI() // MTI: 0100
message.GetString(2) // Card number: 4242424242424242
message.GetString(3) // Processing code: 123456
message.GetString(4) // Transaction amount: 100

// ...
```

#### Typed fields

To get typed field values just pass empty struct for the data you want to access:

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
err := message.Unpack(binaryData)

// to get access to typed data we have to get Data from the message
// and convert it into our ISO87Data type
data := message.Data().(*ISO87Data)

// now you have typed values
data.F2.Value // is a string "4242424242424242"
data.F3.Value // is an int 123456
data.F4.Value // is a string "100"
```

For real code samples please check [./message_test](./message_test).

# License

Apache License 2.0 See [LICENSE](LICENSE) for details.
