package iso8583

import (
	"reflect"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestStructWithStringType(t *testing.T) {
	spec := &MessageSpec{
		Fields: map[int]field.Field{
			0: field.NewString(&field.Spec{
				Length:      4,
				Description: "Message Type Indicator",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			1: field.NewBitmap(&field.Spec{
				Length:      16,
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
			3: field.NewNumeric(&field.Spec{
				Length:      6,
				Description: "Processing Code",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
		},
	}

	t.Run("pack", func(t *testing.T) {
		panInt := 4242424242424242
		panStr := "4242424242424242"

		tests := []struct {
			name                 string
			input                interface{}
			expectedPackedString string
		}{
			// Tests for string type
			{
				name: "struct with string type and value set",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber string `index:"2"`
				}{
					MTI:                  "0110",
					PrimaryAccountNumber: panStr,
				},
				expectedPackedString: "011040000000000000000000000000000000164242424242424242",
			},
			{
				name: "struct with string type and no value",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber string `index:"2"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011000000000000000000000000000000000",
			},
			{
				name: "struct with string type, no value and keepzero tag - length prefix is set to 0 and no value is following",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber string `index:"2,keepzero"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "01104000000000000000000000000000000000",
			},

			// Tests for *string type
			{
				name: "struct with *string type and value set",
				input: struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *string `index:"2"`
				}{
					MTI:                  "0110",
					PrimaryAccountNumber: &panStr,
				},
				expectedPackedString: "011040000000000000000000000000000000164242424242424242",
			},
			{
				name: "struct with *string type and no value",
				input: struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *string `index:"2"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011000000000000000000000000000000000",
			},
			{
				name: "struct with *string type, no value and keepzero tag",
				input: struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *string `index:"2,keepzero"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "01104000000000000000000000000000000000",
			},

			// Tests for int type
			{
				name: "struct with int type and value set",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber int    `index:"2"`
				}{
					MTI:                  "0110",
					PrimaryAccountNumber: panInt,
				},
				expectedPackedString: "011040000000000000000000000000000000164242424242424242",
			},
			{
				name: "struct with int type and no value",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber int    `index:"2"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011000000000000000000000000000000000",
			},
			{
				name: "struct with int type, no value and keepzero tag",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber int    `index:"2,keepzero"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011040000000000000000000000000000000010",
			},

			// Tests for *int type
			{
				name: "struct with *int type and value set",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber *int   `index:"2"`
				}{
					MTI:                  "0110",
					PrimaryAccountNumber: &panInt,
				},
				expectedPackedString: "011040000000000000000000000000000000164242424242424242",
			},
			{
				name: "struct with *int type and no value",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber *int   `index:"2"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011000000000000000000000000000000000",
			},
			{
				name: "struct with *int type, no value and keepzero tag",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber *int   `index:"2,keepzero"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011040000000000000000000000000000000010",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				message := NewMessage(spec)
				err := message.Marshal(tt.input)
				require.NoError(t, err)

				packed, err := message.Pack()
				require.NoError(t, err)

				require.Equal(t, tt.expectedPackedString, string(packed))
			})
		}
	})

	t.Run("unpack", func(t *testing.T) {
		tests := []struct {
			name  string
			input any
		}{
			// Tests for string type
			{
				name: "struct with string type and value set",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber string `index:"2"`
				}{},
			},
			{
				name: "struct with string type and no value",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber string `index:"2"`
				}{},
			},
			{
				name: "struct with string type, no value and keepzero tag - length prefix is set to 0 and no value is following",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber string `index:"2,keepzero"`
				}{},
			},

			// Tests for *string type
			{
				name: "struct with *string type and value set",
				input: &struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *string `index:"2"`
				}{},
			},
			{
				name: "struct with *string type and no value",
				input: &struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *string `index:"2"`
				}{},
			},
			{
				name: "struct with *string type, no value and keepzero tag",
				input: &struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *string `index:"2,keepzero"`
				}{},
			},

			// Tests for int type
			{
				name: "struct with int type and value set",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber int    `index:"2"`
				}{},
			},
			{
				name: "struct with int type and no value",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber int    `index:"2"`
				}{},
			},
			{
				name: "struct with int type, no value and keepzero tag",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber int    `index:"2,keepzero"`
				}{},
			},

			// Tests for *int type
			{
				name: "struct with *int type and value set",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber *int   `index:"2"`
				}{},
			},
			{
				name: "struct with *int type and no value",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber *int   `index:"2"`
				}{},
			},
			{
				name: "struct with *int type, no value and keepzero tag",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber *int   `index:"2,keepzero"`
				}{},
			},
		}
		packed := []byte("011040000000000000000000000000000000164242424242424242")

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				message := NewMessage(spec)
				err := message.Unpack(packed)

				err = message.Unmarshal(tt.input)
				require.NoError(t, err)

				val := reflect.Indirect(reflect.ValueOf(tt.input))
				require.Equal(t, "0110", val.Field(0).String())

				tStr := val.Field(1).Type().String()
				switch tStr {
				case "int":
					require.Equal(t, int64(4242424242424242), val.Field(1).Int())
				case "*int":
					require.Equal(t, int64(4242424242424242), val.Field(1).Elem().Int())
				case "string":
					require.Equal(t, "4242424242424242", val.Field(1).String())
				case "*string":
					require.Equal(t, "4242424242424242", val.Field(1).Elem().String())
				}
			})
		}
	})
}
