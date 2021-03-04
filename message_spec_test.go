package iso8583

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
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
				Enc:         encoding.Hex,
				Pref:        prefix.Hex.Fixed,
			}),
			2: field.NewNumeric(&field.Spec{
				Length:      12,
				Description: "Field 4",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			3: field.NewString(&field.Spec{
				Length:      999,
				Identifier:  "Additional data National",
				Description: "Additional data (National)",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LLL,
			}),
		},
	}

	jsonMessage := `{
	"000": {
		"type": "String",
		"length": 4,
		"enc": "ASCII",
		"pref": "ASCII.Fixed",
		"description": "Message Type Indicator"
	},
	"001": {
		"type": "Bitmap",
		"length": 0,
		"enc": "Hex",
		"pref": "Hex.Fixed",
		"description": "Bitmap"
	},
	"002": {
		"type": "Numeric",
		"length": 12,
		"enc": "ASCII",
		"pref": "ASCII.Fixed",
		"description": "Field 4",
		"pad": {
			"type": "Left",
			"padder": "0"
		}
	},
	"003": {
		"type": "String",
		"length": 999,
		"enc": "ASCII",
		"pref": "ASCII.LLL",
		"description": "Additional data (National)",
		"identifier": "Additional data National"
	}
}`

	xmlMessage := `<Specification>
	<F000>
		<Type>String</Type>
		<Length>4</Length>
		<Enc>ASCII</Enc>
		<Pref>ASCII.Fixed</Pref>
		<Description>Message Type Indicator</Description>
	</F000>
	<F001>
		<Type>Bitmap</Type>
		<Length>0</Length>
		<Enc>Hex</Enc>
		<Pref>Hex.Fixed</Pref>
		<Description>Bitmap</Description>
	</F001>
	<F002>
		<Type>Numeric</Type>
		<Length>12</Length>
		<Enc>ASCII</Enc>
		<Pref>ASCII.Fixed</Pref>
		<Description>Field 4</Description>
		<Pad>
			<Type>Left</Type>
			<Padder>0</Padder>
		</Pad>
	</F002>
	<F003>
		<Type>String</Type>
		<Length>999</Length>
		<Enc>ASCII</Enc>
		<Pref>ASCII.LLL</Pref>
		<Description>Additional data (National)</Description>
		<Identifier>Additional data National</Identifier>
	</F003>
</Specification>`

	t.Run("Test Fields", func(t *testing.T) {
		fields := spec.CreateMessageFields()

		// test that derived fields have the same type as in the message spec
		require.True(t, reflect.TypeOf(fields[0]).Elem() == reflect.TypeOf(field.String{}))
		require.True(t, reflect.TypeOf(fields[1]).Elem() == reflect.TypeOf(field.Bitmap{}))
		require.True(t, reflect.TypeOf(fields[2]).Elem() == reflect.TypeOf(field.Numeric{}))
		require.True(t, reflect.TypeOf(fields[3]).Elem() == reflect.TypeOf(field.String{}))

		// test that derived field have the same spec
		require.Equal(t, fields[0].Spec(), spec.Fields[0].Spec())
		require.Equal(t, fields[1].Spec(), spec.Fields[1].Spec())
		require.Equal(t, fields[2].Spec(), spec.Fields[2].Spec())
		require.Equal(t, fields[3].Spec(), spec.Fields[3].Spec())
	})

	t.Run("Json Marshal Test", func(t *testing.T) {
		jsonBuf, err := json.MarshalIndent(spec, "", "\t")

		require.Equal(t, nil, err)
		require.Equal(t, jsonMessage, string(jsonBuf))
	})

	t.Run("Json Unmarshal Test", func(t *testing.T) {
		spec1 := &MessageSpec{}
		err := json.Unmarshal([]byte(jsonMessage), spec1)

		require.Equal(t, nil, err)

		// test that derived fields have the same type as in the message spec
		require.True(t, reflect.TypeOf(spec1.Fields[0]).Elem() == reflect.TypeOf(field.String{}))
		require.True(t, reflect.TypeOf(spec1.Fields[1]).Elem() == reflect.TypeOf(field.Bitmap{}))
		require.True(t, reflect.TypeOf(spec1.Fields[2]).Elem() == reflect.TypeOf(field.Numeric{}))
		require.True(t, reflect.TypeOf(spec1.Fields[3]).Elem() == reflect.TypeOf(field.String{}))

		// test that derived field have the same spec
		require.Equal(t, spec1.Fields[0].Spec(), spec.Fields[0].Spec())
		require.Equal(t, spec1.Fields[1].Spec(), spec.Fields[1].Spec())
		require.Equal(t, spec1.Fields[2].Spec(), spec.Fields[2].Spec())
		require.Equal(t, spec1.Fields[3].Spec(), spec.Fields[3].Spec())
	})

	t.Run("Xml Marshal Test", func(t *testing.T) {
		xmlBuf, err := xml.MarshalIndent(spec, "", "\t")

		require.Equal(t, nil, err)
		require.Equal(t, xmlMessage, string(xmlBuf))
	})

	t.Run("Xml Unmarshal Test", func(t *testing.T) {
		spec1 := &MessageSpec{}
		err := xml.Unmarshal([]byte(xmlMessage), spec1)

		require.Equal(t, nil, err)

		// test that derived fields have the same type as in the message spec
		require.True(t, reflect.TypeOf(spec1.Fields[0]).Elem() == reflect.TypeOf(field.String{}))
		require.True(t, reflect.TypeOf(spec1.Fields[1]).Elem() == reflect.TypeOf(field.Bitmap{}))
		require.True(t, reflect.TypeOf(spec1.Fields[2]).Elem() == reflect.TypeOf(field.Numeric{}))
		require.True(t, reflect.TypeOf(spec1.Fields[3]).Elem() == reflect.TypeOf(field.String{}))

		// test that derived field have the same spec
		require.Equal(t, spec1.Fields[0].Spec(), spec.Fields[0].Spec())
		require.Equal(t, spec1.Fields[1].Spec(), spec.Fields[1].Spec())
		require.Equal(t, spec1.Fields[2].Spec(), spec.Fields[2].Spec())
		require.Equal(t, spec1.Fields[3].Spec(), spec.Fields[3].Spec())
	})
}
