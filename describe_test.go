package iso8583

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

func TestDescribe(t *testing.T) {
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
			2: field.NewString(&field.Spec{
				Length:      19,
				Description: "Primary Account Number",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			3: field.NewComposite(&field.Spec{
				Length:      6,
				Description: "Processing Code",
				Pref:        prefix.ASCII.Fixed,
				Tag: &field.TagSpec{
					Sort: sort.StringsByInt,
				},
				Subfields: map[string]field.Field{
					"01": field.NewString(&field.Spec{
						Length:      2,
						Description: "Transaction Type",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
					"02": field.NewString(&field.Spec{
						Length:      2,
						Description: "From Account",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
					"03": field.NewString(&field.Spec{
						Length:      2,
						Description: "To Account",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
				},
			}),
		},
	}

	message := NewMessage(spec)
	message.MTI("0100")
	message.Field(2, "4242424242424242")
	message.Field(3, "123456")
	message.Pack() // to generate bitmap

	out := bytes.NewBuffer([]byte{})
	require.NotPanics(t, func() {
		Describe(message, out, DoNotFilterFields()...)
	})

	expectedOutput := `ISO 8583 Message:
MTI..........: 0100
Bitmap HEX...: 6000000000000000
Bitmap bits..:
    [1-8]01100000    [9-16]00000000   [17-24]00000000   [25-32]00000000
  [33-40]00000000   [41-48]00000000   [49-56]00000000   [57-64]00000000
F0   Message Type Indicator..: 0100
F2   Primary Account Number..: 4242424242424242
F3   Processing Code SUBFIELDS:
-------------------------------------------
F01  Transaction Type..: 12
F02  From Account......: 34
F03  To Account........: 56
------------------------------------------
`
	require.Equal(t, expectedOutput, out.String())
}

func Test_splitAndAnnotate(t *testing.T) {
	// test that splitAndAnnotate splits sequences of bits (0, 1) by spaces
	// then annotates each bit with its position in the bitmap
	// and adds a space or newline after every N bits (length of the sequence)
	tt := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
		{
			name:     "1 bit",
			input:    "1",
			expected: "[1-1]1",
		},
		{
			name:     "8 bits",
			input:    "11111111",
			expected: "[1-8]11111111",
		},
		{
			name:     "32 bits",
			input:    "11111111 11111111 11111111 11111111",
			expected: "[1-8]11111111 [9-16]11111111 [17-24]11111111 [25-32]11111111",
		},
		{
			name:     "64 bits",
			input:    "11111111 11111111 11111111 11111111 11111111 11111111 11111111 11111111",
			expected: "    [1-8]11111111    [9-16]11111111   [17-24]11111111   [25-32]11111111\n  [33-40]11111111   [41-48]11111111   [49-56]11111111   [57-64]11111111",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, splitAndAnnotate(tc.input))
		})
	}
}

func TestSortFieldIDs(t *testing.T) {
	var fields = map[string]field.Field{
		"z":   field.NewString(nil),
		"0":   field.NewString(nil),
		"1":   field.NewString(nil),
		"107": field.NewString(nil),
		"a":   field.NewString(nil),
		"17":  field.NewString(nil),
		"2":   field.NewString(nil),
		"4":   field.NewString(nil),
	}
	order := sortFieldIDs(fields)
	expected := []string{"0", "1", "2", "4", "17", "107", "a", "z"}
	if !reflect.DeepEqual(order, expected) {
		t.Errorf("Marshalled value should be \n\t%v\ninstead of \n\t%v", expected, order)
	}
}
