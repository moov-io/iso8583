package iso8583_test

import (
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/specs"
	"github.com/stretchr/testify/require"
)

func TestStructWithTypes(t *testing.T) {
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
				expectedPackedString: "011000000000000000000000000000000000",
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
				expectedPackedString: "011000000000000000000000000000000000",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				message := iso8583.NewMessage(specs.Spec87ASCII)
				err := message.Marshal(tt.input)
				require.NoError(t, err)

				packed, err := message.Pack()
				require.NoError(t, err)

				require.Equal(t, tt.expectedPackedString, string(packed))
			})
		}
	})

	t.Run("unpack", func(t *testing.T) {
		type authRequest struct {
			MTI                  string `index:"0"`
			PrimaryAccountNumber string `index:"2"`
		}

		packed := []byte("011040000000000000000000000000000000164242424242424242")

		message := iso8583.NewMessage(specs.Spec87ASCII)
		err := message.Unpack(packed)
		require.NoError(t, err)

		data := authRequest{}
		err = message.Unmarshal(&data)
		require.NoError(t, err)
		require.Equal(t, "0110", data.MTI)
		require.Equal(t, "4242424242424242", data.PrimaryAccountNumber)
	})

	t.Run("unpack2", func(t *testing.T) {
		type authRequest struct {
			MTI                  *string       `index:"0"`
			PrimaryAccountNumber *field.String `index:"2"`
		}

		packed := []byte("011040000000000000000000000000000000164242424242424242")

		message := iso8583.NewMessage(specs.Spec87ASCII)
		err := message.Unpack(packed)
		require.NoError(t, err)

		data := authRequest{}
		err = message.Unmarshal(&data)
		require.NoError(t, err)
		require.Equal(t, "0110", *data.MTI)
		require.Equal(t, "4242424242424242", data.PrimaryAccountNumber.Value())
	})
}
