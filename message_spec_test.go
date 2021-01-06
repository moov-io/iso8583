package iso8583

import (
	"reflect"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestMessageSpec(t *testing.T) {
	// spec := &MessageSpec{
	// 	Fields: map[int]field.Field{
	// 		0: field.NewStringField(&field.Spec{
	// 			Length:      4,
	// 			Description: "Message Type Indicator",
	// 			Enc:         encoding.ASCII,
	// 			Pref:        prefix.ASCII.Fixed,
	// 		}),
	// 		1: field.NewBitmapField(&field.Spec{
	// 			Length:      16,
	// 			Description: "Bitmap",
	// 			Enc:         encoding.Hex,
	// 			Pref:        prefix.Hex.Fixed,
	// 		}),
	// 		2: field.NewStringField(&field.Spec{
	// 			Length:      19,
	// 			Description: "Primary Account Number",
	// 			Enc:         encoding.ASCII,
	// 			Pref:        prefix.ASCII.LL,
	// 		}),
	// 	},
	// }

	// fmt.Println(spec)
}

func TestMessageSpec_CreateMessageFields(t *testing.T) {
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
		},
	}

	fields := spec.CreateMessageFields()

	// test that derived fields have the same type as in the message spec
	require.True(t, reflect.TypeOf(fields[0]).Elem() == reflect.TypeOf(field.StringField{}))
	require.True(t, reflect.TypeOf(fields[1]).Elem() == reflect.TypeOf(field.BitmapField{}))

	// test that derived field have the same spec
	require.Equal(t, fields[0].Spec(), spec.Fields[0].Spec())
	require.Equal(t, fields[1].Spec(), spec.Fields[1].Spec())
}
