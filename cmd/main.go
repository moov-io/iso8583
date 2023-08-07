package main

import (
	"fmt"
	"os"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
)

var spec *iso8583.MessageSpec = &iso8583.MessageSpec{
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		}),
		2: field.NewString(&field.Spec{
			Length:      19,
			Description: "Primary Account Number",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		3: field.NewComposite(&field.Spec{
			Length:      99,
			Description: "Processing Code",
			Pref:        prefix.TLV.LL,
			Tag: &field.TagSpec{
				Sort: sort.StringsByInt,
			},
			Subfields: map[string]field.Field{
				"1": field.NewString(&field.Spec{
					Length:      2,
					Description: "Transaction Type",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"2": field.NewString(&field.Spec{
					Length:      2,
					Description: "From Account",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"3": field.NewString(&field.Spec{
					Length:      99,
					Description: "To Account",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.LL,
				}),
			},
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

type ISO struct {
	BIT_0 *field.String `index:"0"`
	BIT_2 *field.String `index:"2"`
	BIT_4 *field.String `index:"4"`
}

func main() {
	ascii := "010070000000000000002164242424242424242326112234305560004000000000100"
	isomessage := iso8583.NewMessage(spec)
	isomessage.Unpack([]byte(ascii))

	iso8583.Describe(isomessage, os.Stdout)

	data := &ISO{}
	err := isomessage.Unmarshal(data)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	fmt.Println("BIT_0: ", data.BIT_0.Value())
	fmt.Println("BIT_2: ", data.BIT_2.Value())
	fmt.Println("BIT_4: ", data.BIT_4.Value())
}
