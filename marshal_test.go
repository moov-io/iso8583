package iso8583

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
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
						Length:      2,
						Description: "To Account",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
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

	t.Run("Marshal message", func(t *testing.T) {
		type TestISOF3Data struct {
			F1 *field.String
			F2 *field.String
			F3 *field.String
		}

		type ISO87Data struct {
			F0 *field.String
			F2 *field.String
			F3 *TestISOF3Data
			F4 *field.String
		}

		message := NewMessage(spec)
		err := Marshal(message, &ISO87Data{
			F0: field.NewStringValue("0100"),
			F2: field.NewStringValue("4242424242424242"),
			F3: &TestISOF3Data{
				F1: field.NewStringValue("12"),
				F2: field.NewStringValue("34"),
				F3: field.NewStringValue("56"),
			},
			F4: field.NewStringValue("100"),
		})
		require.NoError(t, err)

		data := &ISO87Data{}
		require.NoError(t, Unmarshal(message, data))

		require.Equal(t, "0100", data.F0.Value)
		require.Equal(t, "4242424242424242", data.F2.Value)
		require.Equal(t, "12", data.F3.F1.Value)
		require.Equal(t, "34", data.F3.F2.Value)
		require.Equal(t, "56", data.F3.F3.Value)
		require.Equal(t, "100", data.F4.Value)
	})

	t.Run("Marshal nil", func(t *testing.T) {
		message := NewMessage(spec)

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		err = Marshal(message, nil)
		require.Error(t, err)
	})

	t.Run("Marshal using field tags", func(t *testing.T) {
		type TestISOF3Data struct {
			One   *field.String `index:"1"`
			Two   *field.String `index:"2"`
			Three *field.String `index:"3"`
		}

		type ISO87Data struct {
			MTI                  *field.String  `index:"0"`
			PrimaryAccountNumber *field.String  `index:"2"`
			AdditionalData       *TestISOF3Data `index:"3"`
			Amount               *field.String  `index:"4"`
		}

		data := &ISO87Data{
			MTI:                  field.NewStringValue("0100"),
			PrimaryAccountNumber: field.NewStringValue("4242424242424242"),
			AdditionalData: &TestISOF3Data{

				One:   field.NewStringValue("12"),
				Two:   field.NewStringValue("34"),
				Three: field.NewStringValue("56"),
			},
			Amount: field.NewStringValue("100"),
		}

		message := NewMessage(spec)
		require.NoError(t, Marshal(message, data))

		rawMsg, err := message.Pack()
		require.NoError(t, err)

		expected := []byte("01007000000000000000164242424242424242123456000000000100")
		require.Equal(t, expected, rawMsg)
	})
}

func TestUnmarshal(t *testing.T) {
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
						Length:      2,
						Description: "To Account",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
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

	t.Run("Unmarshal after unpacking", func(t *testing.T) {
		type TestISOF3Data struct {
			F1 *field.String
			F2 *field.String
			F3 *field.String
		}

		type ISO87Data struct {
			F0 *field.String
			F2 *field.String
			F3 *TestISOF3Data
			F4 *field.String
		}

		message := NewMessage(spec)

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		data := &ISO87Data{}
		err = Unmarshal(message, data)
		require.NoError(t, err)

		require.Equal(t, "0100", data.F0.Value)
		require.Equal(t, "4242424242424242", data.F2.Value)
		require.Equal(t, "12", data.F3.F1.Value)
		require.Equal(t, "34", data.F3.F2.Value)
		require.Equal(t, "56", data.F3.F3.Value)
		require.Equal(t, "100", data.F4.Value)
	})

	t.Run("Unmarshal into nil", func(t *testing.T) {
		message := NewMessage(spec)

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		err = Unmarshal(message, nil)
		require.Error(t, err)
	})

	t.Run("Unmarshal using field tags", func(t *testing.T) {
		type TestISOF3Data struct {
			One   *field.String `index:"1"`
			Two   *field.String `index:"2"`
			Three *field.String `index:"3"`
		}

		type ISO87Data struct {
			MTI                  *field.String  `index:"0"`
			PrimaryAccountNumber *field.String  `index:"2"`
			AdditionalData       *TestISOF3Data `index:"3"`
			Amount               *field.String  `index:"4"`
		}

		message := NewMessage(spec)

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		data := &ISO87Data{}
		err = Unmarshal(message, data)
		require.NoError(t, err)

		require.Equal(t, "0100", data.MTI.Value)
		require.Equal(t, "4242424242424242", data.PrimaryAccountNumber.Value)
		require.Equal(t, "12", data.AdditionalData.One.Value)
		require.Equal(t, "34", data.AdditionalData.Two.Value)
		require.Equal(t, "56", data.AdditionalData.Three.Value)
		require.Equal(t, "100", data.Amount.Value)
	})
}
