package iso8583

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/franizus/iso8583/encoding"
	"github.com/franizus/iso8583/field"
	"github.com/franizus/iso8583/padding"
	"github.com/franizus/iso8583/prefix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
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
				Fields: map[int]field.Field{
					1: field.NewString(&field.Spec{
						Length:      2,
						Description: "Transaction Type",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
					2: field.NewString(&field.Spec{
						Length:      2,
						Description: "From Account",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
					3: field.NewString(&field.Spec{
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

	t.Run("Test packing and unpacking untyped fields", func(t *testing.T) {
		message := NewMessage(spec)
		message.MTI("0100")
		require.NoError(t, message.Field(2, "4242424242424242"))
		require.NoError(t, message.Field(3, "123456"))
		require.NoError(t, message.Field(4, "100"))

		got, err := message.Pack()

		want := "01007000000000000000164242424242424242123456000000000100"
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, want, string(got))

		message = NewMessage(spec)
		message.Unpack([]byte(want))

		s, err := message.GetMTI()
		require.NoError(t, err)
		require.Equal(t, "0100", s)

		s, err = message.GetString(2)
		require.NoError(t, err)
		require.Equal(t, "4242424242424242", s)

		s, err = message.GetString(3)
		require.NoError(t, err)
		require.Equal(t, "123456", s)

		s, err = message.GetString(4)
		require.NoError(t, err)
		require.Equal(t, "100", s)
	})

	t.Run("Test unpacking with typed fields", func(t *testing.T) {
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
		message.SetData(&ISO87Data{})

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		s, err := message.GetString(2)
		require.NoError(t, err)
		require.Equal(t, "4242424242424242", s)

		s, err = message.GetString(3)
		require.NoError(t, err)
		require.Equal(t, "123456", s)

		s, err = message.GetString(4)
		require.NoError(t, err)
		require.Equal(t, "100", s)

		data := message.Data().(*ISO87Data)

		require.Equal(t, "0100", data.F0.Value)
		require.Equal(t, "4242424242424242", data.F2.Value)
		require.Equal(t, "12", data.F3.F1.Value)
		require.Equal(t, "34", data.F3.F2.Value)
		require.Equal(t, "56", data.F3.F3.Value)
		require.Equal(t, "100", data.F4.Value)
	})

	t.Run("Test packing with typed fields", func(t *testing.T) {
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
		err := message.SetData(&ISO87Data{
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

		rawMsg, err := message.Pack()
		require.NoError(t, err)

		wantMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		require.Equal(t, wantMsg, rawMsg)
	})
}

func TestPackUnpack(t *testing.T) {
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
				Enc:         encoding.Binary,
				Pref:        prefix.ASCII.Fixed,
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
				Fields: map[int]field.Field{
					1: field.NewString(&field.Spec{
						Length:      2,
						Description: "Transaction Type",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
					2: field.NewString(&field.Spec{
						Length:      2,
						Description: "From Account",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
					3: field.NewString(&field.Spec{
						Length:      2,
						Description: "To Account",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
				},
			}),
			4: field.NewNumeric(&field.Spec{
				Length:      12,
				Description: "Field 4",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			7: field.NewNumeric(&field.Spec{
				Length:      10,
				Description: "Field 7",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			11: field.NewNumeric(&field.Spec{
				Length:      6,
				Description: "Field 11",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			12: field.NewNumeric(&field.Spec{
				Length:      6,
				Description: "Field 12",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			13: field.NewNumeric(&field.Spec{
				Length:      4,
				Description: "Field 13",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			14: field.NewNumeric(&field.Spec{
				Length:      4,
				Description: "Field 14",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			19: field.NewNumeric(&field.Spec{
				Length:      3,
				Description: "Field 19",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			22: field.NewNumeric(&field.Spec{
				Length:      3,
				Description: "Field 22",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			25: field.NewNumeric(&field.Spec{
				Length:      2,
				Description: "Field 25",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			32: field.NewNumeric(&field.Spec{
				Length:      11,
				Description: "Field 32",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			35: field.NewString(&field.Spec{
				Length:      37,
				Description: "Field 35",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			37: field.NewString(&field.Spec{
				Length:      12,
				Description: "Field 37",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			39: field.NewString(&field.Spec{
				Length:      2,
				Description: "Field 39",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			41: field.NewString(&field.Spec{
				Length:      8,
				Description: "Field 41",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			42: field.NewString(&field.Spec{
				Length:      15,
				Description: "Field 42",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			43: field.NewString(&field.Spec{
				Length:      40,
				Description: "Field 43",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left(' '),
			}),
			49: field.NewNumeric(&field.Spec{
				Length:      3,
				Description: "Field 49",
				Enc:         encoding.LBCD,
				Pref:        prefix.BCD.Fixed,
			}),
			50: field.NewNumeric(&field.Spec{
				Length:      3,
				Description: "Field 50",
				Enc:         encoding.LBCD,
				Pad:         padding.Left('0'),
			}),
			// this one should be binary...
			52: field.NewString(&field.Spec{
				Length:      8,
				Description: "Field 52",
				Enc:         encoding.Binary,
				Pref:        prefix.ASCII.Fixed,
			}),
			53: field.NewNumeric(&field.Spec{
				Length:      16,
				Description: "Field 53",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left(' '),
			}),
			120: field.NewString(&field.Spec{
				Length:      999,
				Description: "Field 120",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LLL,
			}),
		},
	}

	type TestISOF3Data struct {
		F1 *field.String
		F2 *field.String
		F3 *field.String
	}

	type TestISOData struct {
		F2   *field.String
		F3   *TestISOF3Data
		F4   *field.Numeric
		F7   *field.Numeric
		F11  *field.Numeric
		F12  *field.Numeric
		F13  *field.Numeric
		F14  *field.Numeric
		F19  *field.Numeric
		F22  *field.Numeric
		F25  *field.Numeric
		F32  *field.Numeric
		F35  *field.String
		F37  *field.String
		F39  *field.String
		F41  *field.String
		F42  *field.String
		F43  *field.String
		F49  *field.Numeric
		F50  *field.Numeric
		F52  *field.String
		F53  *field.Numeric
		F120 *field.String
	}

	t.Run("Pack data", func(t *testing.T) {
		message := NewMessage(spec)
		err := message.SetData(&TestISOData{
			F2: field.NewStringValue("4276555555555555"),
			F3: &TestISOF3Data{
				F1: field.NewStringValue("00"),
				F2: field.NewStringValue("00"),
				F3: field.NewStringValue("00"),
			},
			F4:  field.NewNumericValue(77700),
			F7:  field.NewNumericValue(701111844),
			F11: field.NewNumericValue(123),
			F12: field.NewNumericValue(131844),
			F13: field.NewNumericValue(701),
			F14: field.NewNumericValue(1902),
			F19: field.NewNumericValue(643),
			F22: field.NewNumericValue(901),
			F25: field.NewNumericValue(2),
			F32: field.NewNumericValue(123456),
			F35: field.NewStringValue("4276555555555555=12345678901234567890"),
			F37: field.NewStringValue("987654321001"),
			F41: field.NewStringValue("00000321"),
			F42: field.NewStringValue("120000000000034"),
			F43: field.NewStringValue("Test text"),
			F49: field.NewNumericValue(643),
			// F50 left nil to ensure that it has not been populated in the bitmap
			F52:  field.NewStringValue(string([]byte{1, 2, 3, 4, 5, 6, 7, 8})),
			F53:  field.NewNumericValue(1234000000000000),
			F120: field.NewStringValue("Another test text"),
		})
		require.NoError(t, err)

		message.MTI("0100")

		got, err := message.Pack()

		want := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, want, got)
	})

	t.Run("Unpack data", func(t *testing.T) {
		message := NewMessage(spec)
		message.SetData(&TestISOData{})

		rawMsg := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		s, err := message.GetString(2)
		require.NoError(t, err)
		require.Equal(t, "4276555555555555", s)

		s, err = message.GetString(3)
		require.NoError(t, err)
		require.Equal(t, "000000", s)

		s, err = message.GetString(4)
		require.NoError(t, err)
		require.Equal(t, "77700", s)

		data := message.Data().(*TestISOData)

		assert.Equal(t, "4276555555555555", data.F2.Value)
		assert.Equal(t, "00", data.F3.F1.Value)
		assert.Equal(t, "00", data.F3.F2.Value)
		assert.Equal(t, "00", data.F3.F3.Value)
		assert.Equal(t, 77700, data.F4.Value)
		assert.Equal(t, 701111844, data.F7.Value)
		assert.Equal(t, 123, data.F11.Value)
		assert.Equal(t, 131844, data.F12.Value)
		assert.Equal(t, 701, data.F13.Value)
		assert.Equal(t, 1902, data.F14.Value)
		assert.Equal(t, 643, data.F19.Value)
		assert.Equal(t, 901, data.F22.Value)
		assert.Equal(t, 2, data.F25.Value)
		assert.Equal(t, 123456, data.F32.Value)
		assert.Equal(t, "4276555555555555=12345678901234567890", data.F35.Value)
		assert.Equal(t, "987654321001", data.F37.Value)
		assert.Equal(t, "00000321", data.F41.Value)
		assert.Equal(t, "120000000000034", data.F42.Value)
		assert.Equal(t, "Test text", data.F43.Value)
		assert.Equal(t, 643, data.F49.Value)
		assert.Nil(t, data.F50)
		assert.Equal(t, string([]byte{1, 2, 3, 4, 5, 6, 7, 8}), data.F52.Value)
		assert.Equal(t, 1234000000000000, data.F53.Value)
		assert.Equal(t, "Another test text", data.F120.Value)
	})
}

func TestPackUnpack_ThirdBitmap(t *testing.T) {
	MessageSpecification := getTestSpec(t)

	type F62 struct {
		F2 *field.String
	}

	type F63 struct {
		F1 *field.String
		F2 *field.String
	}

	type IsoMessageData struct {
		F0   *field.String
		F2   *field.String
		F3   *field.String
		F4   *field.String
		F7   *field.String
		F11  *field.String
		F12  *field.String
		F13  *field.String
		F14  *field.String
		F18  *field.String
		F19  *field.String
		F22  *field.String
		F23  *field.String
		F25  *field.String
		F32  *field.String
		F35  *field.String
		F37  *field.String
		F41  *field.String
		F42  *field.String
		F43  *field.String
		F48  *field.String
		F49  *field.String
		F55  *field.String
		F60  *field.String
		F62  *F62
		F63  *F63
		F104 *field.String
		F135 *field.String
	}

	type DataF62 struct {
		F2 string
	}

	type DataF63 struct {
		F1 string
		F2 string
	}

	type Data struct {
		F2   string
		F3   string
		F4   string
		F7   string
		F11  string
		F12  string
		F13  string
		F14  string
		F18  string
		F19  string
		F22  string
		F23  string
		F25  string
		F32  string
		F35  string
		F37  string
		F41  string
		F42  string
		F43  string
		F48  string
		F49  string
		F55  string
		F60  string
		F62  *DataF62
		F63  *DataF63
		F104 string
		F135 string
	}

	data := &Data{
		F2:  "4761340000000019",
		F3:  "000000",
		F4:  "1000",
		F7:  "0119163908",
		F11: "98",
		F12: "113908",
		F13: "0119",
		F14: "2212",
		F18: "5999",
		F19: "0840",
		F22: "0510",
		F23: "001",
		F25: "00",
		F32: "12345678901",
		F35: "4761340000000019=221212312345129",
		F37: "201916000098",
		F41: "3.3.1   ",
		F42: "CARD ACCEPTOR  ",
		F43: "ACQUIRER NAME            CITY NAME    US",
		F48: "FIELD 48",
		F49: "840",
		F55: "0100569F3303204000950580000100009F37049BADBCAB9F100706010A03A000009F26080123456789ABCDEF9F360200FF820200009C01009F1A0208409A030101019F02060000000123005F2A0208409F0306000000000000",
		F60: "05000010",
		F62: &DataF62{F2: "00000000"},
		F63: &DataF63{
			F1: "0000",
			F2: "0002",
		},
		F104: "6900030101F3",
		F135: "1234567890",
	}

	t.Run("should pack message successfully", func(t *testing.T) {
		msg := NewMessage(MessageSpecification)
		msg.MTI("0100")

		f55, _ := hex.DecodeString(data.F55)
		f104, _ := hex.DecodeString(data.F104)
		f135, _ := hex.DecodeString(data.F135)

		msgData := &IsoMessageData{
			F2:  field.NewStringValue(data.F2),
			F3:  field.NewStringValue(data.F3),
			F4:  field.NewStringValue(data.F4),
			F7:  field.NewStringValue(data.F7),
			F11: field.NewStringValue(data.F11),
			F12: field.NewStringValue(data.F12),
			F13: field.NewStringValue(data.F13),
			F14: field.NewStringValue(data.F14),
			F18: field.NewStringValue(data.F18),
			F19: field.NewStringValue(data.F19),
			F22: field.NewStringValue(data.F22),
			F23: field.NewStringValue(data.F23),
			F25: field.NewStringValue(data.F25),
			F32: field.NewStringValue(data.F32),
			F35: field.NewStringValue(data.F35),
			F37: field.NewStringValue(data.F37),
			F41: field.NewStringValue(data.F41),
			F42: field.NewStringValue(data.F42),
			F43: field.NewStringValue(data.F43),
			F48: field.NewStringValue(data.F48),
			F49: field.NewStringValue(data.F49),
			F55: field.NewStringValue(string(f55)),
			F60: field.NewStringValue(data.F60),
			F62: &F62{
				F2: field.NewStringValue(data.F62.F2),
			},
			F63: &F63{
				F1: field.NewStringValue(data.F63.F1),
				F2: field.NewStringValue(data.F63.F2),
			},
			F104: field.NewStringValue(string(f104)),
			F135: field.NewStringValue(string(f135)),
		}

		msg.SetData(msgData)

		got, err := msg.Pack()

		want := "0100f23c668128e18216800000000100000002000000000000001047613400000000190000000000000010000119163908000098113908011922125999084005100001000b012345678901204761340000000019d221212312345129f2f0f1f9f1f6f0f0f0f0f9f8f34bf34bf1404040c3c1d9c440c1c3c3c5d7e3d6d94040c1c3d8e4c9d9c5d940d5c1d4c5404040404040404040404040c3c9e3e840d5c1d4c540404040e4e208c6c9c5d3c440f4f80840590100569f3303204000950580000100009f37049badbcab9f100706010a03a000009f26080123456789abcdef9f360200ff820200009c01009f1a0208409a030101019f02060000000123005f2a0208409f030600000000000004050000100c40000000000000000000000007c0000000000002066900030101f3051234567890"

		require.NoError(t, err)
		require.Equal(t, want, hex.EncodeToString(got))
	})

	t.Run("should unpack message successfully", func(t *testing.T) {
		message := NewMessage(MessageSpecification)
		msg := &IsoMessageData{}

		isoFrame := "0100f23c668128e18216800000000100000002000000000000001047613400000000190000000000000010000119163908000098113908011922125999084005100001000b012345678901204761340000000019d221212312345129f2f0f1f9f1f6f0f0f0f0f9f8f34bf34bf1404040c3c1d9c440c1c3c3c5d7e3d6d94040c1c3d8e4c9d9c5d940d5c1d4c5404040404040404040404040c3c9e3e840d5c1d4c540404040e4e208c6c9c5d3c440f4f80840590100569f3303204000950580000100009f37049badbcab9f100706010a03a000009f26080123456789abcdef9f360200ff820200009c01009f1a0208409a030101019f02060000000123005f2a0208409f030600000000000004050000100c40000000000000000000000007c0000000000002066900030101f3051234567890"

		msgBytes, err := hex.DecodeString(isoFrame)
		require.NoError(t, err)

		err = message.SetData(msg)
		require.NoError(t, err)

		err = message.Unpack(msgBytes)
		require.NoError(t, err)

		f55 := hex.EncodeToString([]byte(msg.F55.Value))
		f104 := hex.EncodeToString([]byte(msg.F104.Value))
		f135 := hex.EncodeToString([]byte(msg.F135.Value))

		require.Equal(t, data.F2, msg.F2.Value)
		require.Equal(t, data.F3, msg.F3.Value)
		require.Equal(t, data.F4, msg.F4.Value)
		require.Equal(t, data.F7, msg.F7.Value)
		require.Equal(t, data.F11, msg.F11.Value)
		require.Equal(t, data.F12, msg.F12.Value)
		require.Equal(t, data.F13, msg.F13.Value)
		require.Equal(t, data.F14, msg.F14.Value)
		require.Equal(t, data.F18, msg.F18.Value)
		require.Equal(t, data.F19, msg.F19.Value)
		require.Equal(t, data.F22, msg.F22.Value)
		require.Equal(t, data.F23, msg.F23.Value)
		require.Equal(t, data.F25, msg.F25.Value)
		require.Equal(t, data.F32, msg.F32.Value)
		require.Equal(t, data.F35, msg.F35.Value)
		require.Equal(t, data.F37, msg.F37.Value)
		require.Equal(t, data.F41, msg.F41.Value)
		require.Equal(t, data.F42, strings.ToUpper(msg.F42.Value))
		require.Equal(t, data.F43, strings.ToUpper(msg.F43.Value))
		require.Equal(t, data.F48, strings.ToUpper(msg.F48.Value))
		require.Equal(t, data.F49, msg.F49.Value)
		require.Equal(t, data.F55, strings.ToUpper(f55))
		require.Equal(t, data.F60, msg.F60.Value)
		require.Equal(t, data.F62.F2, msg.F62.F2.Value)
		require.Equal(t, data.F63.F1, msg.F63.F1.Value)
		require.Equal(t, data.F63.F2, msg.F63.F2.Value)
		require.Equal(t, data.F104, strings.ToUpper(f104))
		require.Equal(t, data.F135, f135)
	})

}

func TestMessageJSON(t *testing.T) {
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
				Fields: map[int]field.Field{
					1: field.NewString(&field.Spec{
						Length:      2,
						Description: "Transaction Type",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
					2: field.NewString(&field.Spec{
						Length:      2,
						Description: "From Account",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.Fixed,
					}),
					3: field.NewString(&field.Spec{
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

	type TestISOF3Data struct {
		F1 *field.String
		F2 *field.String
		F3 *field.String
	}

	type TestISOData struct {
		F2 *field.String
		F3 *TestISOF3Data
		F4 *field.String
	}

	t.Run("Test JSON encoding", func(t *testing.T) {
		message := NewMessage(spec)
		message.MTI("0100")
		err := message.SetData(&TestISOData{
			F2: field.NewStringValue("4242424242424242"),
			F3: &TestISOF3Data{
				F1: field.NewStringValue("12"),
				F2: field.NewStringValue("34"),
				F3: field.NewStringValue("56"),
			},
			F4: field.NewStringValue("100"),
		})
		require.NoError(t, err)

		want := `{"0":"0100","1":"700000000000000000000000000000000000000000000000","2":"4242424242424242","3":{"1":"12","2":"34","3":"56"},"4":"100"}`

		got, err := json.Marshal(message)
		require.NoError(t, err)
		require.Equal(t, want, string(got))
	})

	t.Run("Test JSON encoding untyped", func(t *testing.T) {
		message := NewMessage(spec)
		message.MTI("0100")
		message.Field(2, "4242424242424242")
		message.Field(4, "100")

		want := `{"0":"0100","1":"500000000000000000000000000000000000000000000000","2":"4242424242424242","4":"100"}`

		got, err := json.Marshal(message)
		require.NoError(t, err)
		require.Equal(t, want, string(got))
	})
}

func TestIsBitmapFlag(t *testing.T) {
	t.Run("should return true because 1 is the first bit of the first bitmap", func(t *testing.T) {
		got := IsBitmapFlag(1)
		require.True(t, got)
	})

	t.Run("should return true because 65 is the first bit of the second bitmap", func(t *testing.T) {
		got := IsBitmapFlag(65)
		require.True(t, got)
	})

	t.Run("should return true because is 129 the first bit of the third bitmap", func(t *testing.T) {
		got := IsBitmapFlag(129)
		require.True(t, got)
	})

	t.Run("should return false because 7 is not the first bit of the bitmap", func(t *testing.T) {
		got := IsBitmapFlag(7)
		require.False(t, got)
	})

	t.Run("should return false because 69 is not the first bit of the bitmap", func(t *testing.T) {
		got := IsBitmapFlag(69)
		require.False(t, got)
	})

	t.Run("should return false because 135 is not the first bit of the bitmap", func(t *testing.T) {
		got := IsBitmapFlag(135)
		require.False(t, got)
	})
}

func getTestSpec(t *testing.T) *MessageSpec {
	t.Helper()
	return &MessageSpec{
		Fields: map[int]field.Field{
			0: field.NewString(&field.Spec{
				Length:      4,
				Description: "Message Type Indicator",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			1: field.NewBitmap(&field.Spec{
				Description: "Bitmap",
				Enc:         encoding.Binary,
				Pref:        prefix.Binary.Fixed,
			}),
			2: field.NewString(&field.Spec{
				Length:      19,
				Description: "Primary Account Number",
				Enc:         encoding.BCD,
				Pref:        prefix.Binary.LL,
				CountT:      "1",
			}),
			3: field.NewString(&field.Spec{
				Length:      6,
				Description: "Processing Code",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			4: field.NewString(&field.Spec{
				Length:      12,
				Description: "Transaction Amount",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
				Pad:         padding.Left('0'),
			}),
			7: field.NewString(&field.Spec{
				Length:      10,
				Description: "Transmission Date & Time",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			9: field.NewString(&field.Spec{
				Length:      8,
				Description: "Conversion Rate, Settlement",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			11: field.NewString(&field.Spec{
				Length:      6,
				Description: "Systems Trace Audit Number (STAN)",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
				Pad:         padding.Left('0'),
			}),
			12: field.NewString(&field.Spec{
				Length:      6,
				Description: "Local Transaction Time",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			13: field.NewString(&field.Spec{
				Length:      4,
				Description: "Local Transaction Date",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			14: field.NewString(&field.Spec{
				Length:      4,
				Description: "Expiration Date",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			18: field.NewString(&field.Spec{
				Length:      4,
				Description: "Merchant Type",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			19: field.NewString(&field.Spec{
				Length:      4,
				Description: "Acquiring Institution Country Code",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			22: field.NewString(&field.Spec{
				Length:      4,
				Description: "Point of Service Entry Mode Code",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			23: field.NewString(&field.Spec{
				Length:      3,
				Description: "Card Sequence Number",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			25: field.NewString(&field.Spec{
				Length:      2,
				Description: "Point of Service Condition Code",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			32: field.NewString(&field.Spec{
				Length:      11,
				Description: "Acquiring Institutions ID Code",
				Enc:         encoding.BCD,
				Pref:        prefix.Binary.LL,
				CountT:      "1",
			}),
			35: field.NewString(&field.Spec{
				Length:      37,
				Description: "Track 2 Data",
				Enc:         encoding.BCD,
				Pref:        prefix.Binary.LL,
				CountT:      "1",
			}),
			37: field.NewString(&field.Spec{
				Length:      12,
				Description: "Retrieval Reference Number",
				Enc:         encoding.EBCDIC,
				Pref:        prefix.EBCDIC.Fixed,
			}),
			41: field.NewString(&field.Spec{
				Length:      8,
				Description: "Card Acceptor Terminal ID",
				Enc:         encoding.EBCDIC,
				Pref:        prefix.EBCDIC.Fixed,
			}),
			42: field.NewString(&field.Spec{
				Length:      15,
				Description: "Card Acceptor ID Code",
				Enc:         encoding.EBCDIC,
				Pref:        prefix.EBCDIC.Fixed,
				Pad:         padding.Left(' '),
			}),
			43: field.NewString(&field.Spec{
				Length:      40,
				Description: "Card Acceptor Name/Location",
				Enc:         encoding.EBCDIC,
				Pref:        prefix.EBCDIC.Fixed,
			}),
			48: field.NewString(&field.Spec{
				Length:      255,
				Description: "Additional Data-Private",
				Enc:         encoding.EBCDIC,
				Pref:        prefix.Binary.LLL,
			}),
			49: field.NewString(&field.Spec{
				Length:      3,
				Description: "Transaction Currency Code",
				Enc:         encoding.BCD,
				Pref:        prefix.BCD.Fixed,
			}),
			55: field.NewString(&field.Spec{
				Length:      510, // TODO fix on release
				Description: "Integrated Circuit Card (ICC) Related Data",
				Enc:         encoding.Binary,
				Pref:        prefix.Binary.LLL,
			}),
			60: field.NewString(&field.Spec{
				Length:      12,
				Description: "Additional POS Information",
				Enc:         encoding.BCD,
				Pref:        prefix.Binary.LL,
				CountT:      "2",
			}),
			62: field.NewComposite(&field.Spec{
				Length:      256,
				Description: "Custom Payment Service Fields (Bitmap Format)",
				Pref:        prefix.Binary.LL,
				HasBitmap:   true,
				Fields: map[int]field.Field{
					0: field.NewBitmap(&field.Spec{
						Length:      8,
						Description: "Bitmap 62",
						Enc:         encoding.Binary,
						Pref:        prefix.Binary.Fixed,
					}),
					2: field.NewString(&field.Spec{
						Length:      8,
						Description: "Transaction Identifier",
						Enc:         encoding.BCD,
						Pref:        prefix.BCD.Fixed,
					}),
				},
			}),
			63: field.NewComposite(&field.Spec{
				Length:      256,
				Description: "Reserved Private Field Bitmap",
				Pref:        prefix.Binary.LL,
				HasBitmap:   true,
				Fields: map[int]field.Field{
					0: field.NewBitmap(&field.Spec{
						Length:      3,
						Description: "Bitmap 63",
						Enc:         encoding.Binary,
						Pref:        prefix.Binary.Fixed,
					}),
					1: field.NewString(&field.Spec{
						Length:      4,
						Description: "Network ID",
						Enc:         encoding.BCD,
						Pref:        prefix.BCD.Fixed,
					}),
					2: field.NewString(&field.Spec{
						Length:      4,
						Description: "Time (Preauth Time Limit)",
						Enc:         encoding.BCD,
						Pref:        prefix.BCD.Fixed,
					}),
				},
			}),
			104: field.NewString(&field.Spec{
				Length:      255,
				Description: "Transaction Description & Transaction-Specific Data",
				Enc:         encoding.Binary,
				Pref:        prefix.Binary.LLL,
			}),
			135: field.NewString(&field.Spec{
				Length:      30,
				Description: "Issuer Discretionary Data",
				Enc:         encoding.Binary,
				Pref:        prefix.Binary.LLL,
			}),
		},
	}
}
