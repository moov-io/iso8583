package iso8583

import (
	"reflect"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestMessageSpec_CreateMessageFields(t *testing.T) {
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
				Enc:         encoding.BytesToASCIIHex,
				Pref:        prefix.Hex.Fixed,
			}),
		},
	}

	fields := spec.CreateMessageFields()

	// test that derived fields have the same type as in the message spec
	require.True(t, reflect.TypeOf(fields[0]).Elem() == reflect.TypeOf(field.String{}))
	require.True(t, reflect.TypeOf(fields[1]).Elem() == reflect.TypeOf(field.Bitmap{}))

	// test that derived field have the same spec
	require.Equal(t, fields[0].Spec(), spec.Fields[0].Spec())
	require.Equal(t, fields[1].Spec(), spec.Fields[1].Spec())
}
