package iso8583

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
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
			3: field.NewNumeric(&field.Spec{
				Length:      6,
				Description: "Processing Code",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
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
		message.Field(2, "4242424242424242")
		message.Field(3, "123456")
		message.Field(4, "100")

		got, err := message.Pack()

		want := "01007000000000000000164242424242424242123456000000000100"
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, want, string(got))

		message = NewMessage(spec)
		message.Unpack([]byte(want))

		require.Equal(t, "0100", message.GetMTI())
		require.Equal(t, "4242424242424242", message.GetString(2))
		require.Equal(t, "123456", message.GetString(3))
		require.Equal(t, "100", message.GetString(4))
	})

	t.Run("Test unpacking with typed fields", func(t *testing.T) {
		type ISO87Data struct {
			F2 *field.String
			F3 *field.Numeric
			F4 *field.String
		}

		message := NewMessage(spec)
		message.SetData(&ISO87Data{})

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		require.Equal(t, "4242424242424242", message.GetString(2))
		require.Equal(t, "123456", message.GetString(3))
		require.Equal(t, "100", message.GetString(4))

		data := message.Data().(*ISO87Data)

		require.Equal(t, "4242424242424242", data.F2.Value)
		require.Equal(t, 123456, data.F3.Value)
		require.Equal(t, "100", data.F4.Value)
	})

	t.Run("Test packing with typed fields", func(t *testing.T) {
		type ISO87Data struct {
			F2 *field.String
			F3 *field.Numeric
			F4 *field.String
		}

		message := NewMessage(spec)
		message.MTI("0100")
		err := message.SetData(&ISO87Data{
			F2: field.NewStringValue("4242424242424242"),
			F3: field.NewNumericValue(123456),
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
			3: field.NewString(&field.Spec{
				Length:      6,
				Description: "Processing Code",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
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

	type TestISOData struct {
		F2   *field.String
		F3   *field.String
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
		F52  *field.String
		F53  *field.Numeric
		F120 *field.String
	}

	t.Run("Pack data", func(t *testing.T) {
		message := NewMessage(spec)
		message.SetData(&TestISOData{
			F2:   field.NewStringValue("4276555555555555"),
			F3:   field.NewStringValue("000000"),
			F4:   field.NewNumericValue(77700),
			F7:   field.NewNumericValue(701111844),
			F11:  field.NewNumericValue(123),
			F12:  field.NewNumericValue(131844),
			F13:  field.NewNumericValue(701),
			F14:  field.NewNumericValue(1902),
			F19:  field.NewNumericValue(643),
			F22:  field.NewNumericValue(901),
			F25:  field.NewNumericValue(2),
			F32:  field.NewNumericValue(123456),
			F35:  field.NewStringValue("4276555555555555=12345678901234567890"),
			F37:  field.NewStringValue("987654321001"),
			F41:  field.NewStringValue("00000321"),
			F42:  field.NewStringValue("120000000000034"),
			F43:  field.NewStringValue("Test text"),
			F49:  field.NewNumericValue(643),
			F52:  field.NewStringValue(string([]byte{1, 2, 3, 4, 5, 6, 7, 8})),
			F53:  field.NewNumericValue(1234000000000000),
			F120: field.NewStringValue("Another test text"),
		})

		message.MTI("0100")

		got, err := message.Pack()

		want := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, want, got)
	})

	t.Run("JSON marshal of data", func(t *testing.T) {
		message := NewMessage(spec)
		message.SetData(&TestISOData{
			F2:   field.NewStringValue("4276555555555555"),
			F3:   field.NewStringValue("000000"),
			F4:   field.NewNumericValue(77700),
			F7:   field.NewNumericValue(701111844),
			F11:  field.NewNumericValue(123),
			F12:  field.NewNumericValue(131844),
			F13:  field.NewNumericValue(701),
			F14:  field.NewNumericValue(1902),
			F19:  field.NewNumericValue(643),
			F22:  field.NewNumericValue(901),
			F25:  field.NewNumericValue(2),
			F32:  field.NewNumericValue(123456),
			F35:  field.NewStringValue("4276555555555555=12345678901234567890"),
			F37:  field.NewStringValue("987654321001"),
			F41:  field.NewStringValue("00000321"),
			F42:  field.NewStringValue("120000000000034"),
			F43:  field.NewStringValue("Test text"),
			F49:  field.NewNumericValue(643),
			F52:  field.NewStringValue(string([]byte{1, 2, 3, 4, 5, 6, 7, 8})),
			F53:  field.NewNumericValue(1234000000000000),
			F120: field.NewStringValue("Another test text"),
		})

		message.MTI("0100")
		message.Bitmap().SetBytes([]byte("11110010 00111100 00100100 10000001 00101000 11100000 10011000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000"))

		jsonBuf, err := json.MarshalIndent(message, "", "\t")
		require.NoError(t, err)
		require.NotNil(t, jsonBuf)

		want :=
			`{
	"bitmap": "11110010 00111100 00100100 10000001 00101000 11100000 10011000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 ",
	"field_11": "123",
	"field_12": "131844",
	"field_120": "Another test text",
	"field_13": "701",
	"field_14": "1902",
	"field_19": "643",
	"field_22": "901",
	"field_25": "2",
	"field_32": "123456",
	"field_35": "4276555555555555=12345678901234567890",
	"field_37": "987654321001",
	"field_4": "77700",
	"field_41": "00000321",
	"field_42": "120000000000034",
	"field_43": "Test text",
	"field_49": "643",
	"field_52": "\u0001\u0002\u0003\u0004\u0005\u0006\u0007\u0008",
	"field_53": "1234000000000000",
	"field_7": "701111844",
	"message_type_indicator": "0100",
	"primary_account_number": "4276555555555555",
	"processing_code": "000000"
}`

		require.Equal(t, want, string(jsonBuf))
	})

	t.Run("XML marshal of data", func(t *testing.T) {
		message := NewMessage(spec)
		message.SetData(&TestISOData{
			F2:   field.NewStringValue("4276555555555555"),
			F3:   field.NewStringValue("000000"),
			F4:   field.NewNumericValue(77700),
			F7:   field.NewNumericValue(701111844),
			F11:  field.NewNumericValue(123),
			F12:  field.NewNumericValue(131844),
			F13:  field.NewNumericValue(701),
			F14:  field.NewNumericValue(1902),
			F19:  field.NewNumericValue(643),
			F22:  field.NewNumericValue(901),
			F25:  field.NewNumericValue(2),
			F32:  field.NewNumericValue(123456),
			F35:  field.NewStringValue("4276555555555555=12345678901234567890"),
			F37:  field.NewStringValue("987654321001"),
			F41:  field.NewStringValue("00000321"),
			F42:  field.NewStringValue("120000000000034"),
			F43:  field.NewStringValue("Test text"),
			F49:  field.NewNumericValue(643),
			F52:  field.NewStringValue(string([]byte{1, 2, 3, 4, 5, 6, 7, 8})),
			F53:  field.NewNumericValue(1234000000000000),
			F120: field.NewStringValue("Another test text"),
		})

		message.MTI("0100")

		got, err := message.Pack()
		require.NoError(t, err)
		require.NotNil(t, got)

		xmlBuf, err := xml.MarshalIndent(message, "", "\t")
		require.NoError(t, err)
		require.NotNil(t, xmlBuf)

		want :=
			`<ISO8583>
	<MessageTypeIndicator>0100</MessageTypeIndicator>
	<Bitmap>11110010 00111100 00100100 10000001 00101000 11100000 10011000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 </Bitmap>
	<PrimaryAccountNumber>4276555555555555</PrimaryAccountNumber>
	<ProcessingCode>000000</ProcessingCode>
	<Field4>77700</Field4>
	<Field7>701111844</Field7>
	<Field11>123</Field11>
	<Field12>131844</Field12>
	<Field13>701</Field13>
	<Field14>1902</Field14>
	<Field19>643</Field19>
	<Field22>901</Field22>
	<Field25>2</Field25>
	<Field32>123456</Field32>
	<Field35>4276555555555555=12345678901234567890</Field35>
	<Field37>987654321001</Field37>
	<Field41>00000321</Field41>
	<Field42>120000000000034</Field42>
	<Field43>Test text</Field43>
	<Field49>643</Field49>
	<Field52>\u0001\u0002\u0003\u0004\u0005\u0006\u0007\u0008</Field52>
	<Field53>1234000000000000</Field53>
	<Field120>Another test text</Field120>
</ISO8583>`
		require.Equal(t, want, string(xmlBuf))
	})

	t.Run("Unpack data", func(t *testing.T) {
		message := NewMessage(spec)
		message.SetData(&TestISOData{})

		rawMsg := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		require.Equal(t, "4276555555555555", message.GetString(2))
		require.Equal(t, "000000", message.GetString(3))
		require.Equal(t, "77700", message.GetString(4))

		data := message.Data().(*TestISOData)

		assert.Equal(t, "4276555555555555", data.F2.Value)
		assert.Equal(t, "000000", data.F3.Value)
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
		assert.Equal(t, string([]byte{1, 2, 3, 4, 5, 6, 7, 8}), data.F52.Value)
		assert.Equal(t, 1234000000000000, data.F53.Value)
		assert.Equal(t, "Another test text", data.F120.Value)
	})

	t.Run("JSON unmarshal of data", func(t *testing.T) {

		want :=
			`{
	"bitmap": "11110010 00111100 00100100 10000001 00101000 11100000 10011000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 ",
	"field_11": "123",
	"field_12": "131844",
	"field_120": "Another test text",
	"field_13": "701",
	"field_14": "1902",
	"field_19": "643",
	"field_22": "901",
	"field_25": "2",
	"field_32": "123456",
	"field_35": "4276555555555555=12345678901234567890",
	"field_37": "987654321001",
	"field_4": "77700",
	"field_41": "00000321",
	"field_42": "120000000000034",
	"field_43": "Test text",
	"field_49": "643",
	"field_52": "\u0001\u0002\u0003\u0004\u0005\u0006\u0007\u0008",
	"field_53": "1234000000000000",
	"field_7": "701111844",
	"message_type_indicator": "0100",
	"primary_account_number": "4276555555555555",
	"processing_code": "000000"
}`

		message := NewMessage(spec)
		err := json.Unmarshal([]byte(want), message)
		require.NoError(t, err)

		expectMessage := NewMessage(spec)
		expectMessage.SetData(&TestISOData{
			F2:   field.NewStringValue("4276555555555555"),
			F3:   field.NewStringValue("000000"),
			F4:   field.NewNumericValue(77700),
			F7:   field.NewNumericValue(701111844),
			F11:  field.NewNumericValue(123),
			F12:  field.NewNumericValue(131844),
			F13:  field.NewNumericValue(701),
			F14:  field.NewNumericValue(1902),
			F19:  field.NewNumericValue(643),
			F22:  field.NewNumericValue(901),
			F25:  field.NewNumericValue(2),
			F32:  field.NewNumericValue(123456),
			F35:  field.NewStringValue("4276555555555555=12345678901234567890"),
			F37:  field.NewStringValue("987654321001"),
			F41:  field.NewStringValue("00000321"),
			F42:  field.NewStringValue("120000000000034"),
			F43:  field.NewStringValue("Test text"),
			F49:  field.NewNumericValue(643),
			F52:  field.NewStringValue(string([]byte{1, 2, 3, 4, 5, 6, 7, 8})),
			F53:  field.NewNumericValue(1234000000000000),
			F120: field.NewStringValue("Another test text"),
		})

		indexs := []int{2, 3, 4, 7, 11, 12, 13, 14, 19, 22, 25, 32, 35, 37, 41, 42, 43, 49, 52, 53, 120}
		for _, index := range indexs {
			assert.Equal(t, message.GetString(index), expectMessage.GetString(index))
		}

		jsonBuf, err := json.MarshalIndent(message, "", "\t")
		require.NoError(t, err)
		require.NotNil(t, jsonBuf)

		require.Equal(t, want, string(jsonBuf))
	})

	t.Run("XML unmarshal of data", func(t *testing.T) {

		want :=
			`<ISO8583>
	<MessageTypeIndicator>0100</MessageTypeIndicator>
	<Bitmap>11110010 00111100 00100100 10000001 00101000 11100000 10011000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 </Bitmap>
	<PrimaryAccountNumber>4276555555555555</PrimaryAccountNumber>
	<ProcessingCode>000000</ProcessingCode>
	<Field4>77700</Field4>
	<Field7>701111844</Field7>
	<Field11>123</Field11>
	<Field12>131844</Field12>
	<Field13>701</Field13>
	<Field14>1902</Field14>
	<Field19>643</Field19>
	<Field22>901</Field22>
	<Field25>2</Field25>
	<Field32>123456</Field32>
	<Field35>4276555555555555=12345678901234567890</Field35>
	<Field37>987654321001</Field37>
	<Field41>00000321</Field41>
	<Field42>120000000000034</Field42>
	<Field43>Test text</Field43>
	<Field49>643</Field49>
	<Field52>\u0001\u0002\u0003\u0004\u0005\u0006\u0007\u0008</Field52>
	<Field53>1234000000000000</Field53>
	<Field120>Another test text</Field120>
</ISO8583>`

		message := NewMessage(spec)
		err := xml.Unmarshal([]byte(want), message)
		require.NoError(t, err)

		expectMessage := NewMessage(spec)
		expectMessage.SetData(&TestISOData{
			F2:   field.NewStringValue("4276555555555555"),
			F3:   field.NewStringValue("000000"),
			F4:   field.NewNumericValue(77700),
			F7:   field.NewNumericValue(701111844),
			F11:  field.NewNumericValue(123),
			F12:  field.NewNumericValue(131844),
			F13:  field.NewNumericValue(701),
			F14:  field.NewNumericValue(1902),
			F19:  field.NewNumericValue(643),
			F22:  field.NewNumericValue(901),
			F25:  field.NewNumericValue(2),
			F32:  field.NewNumericValue(123456),
			F35:  field.NewStringValue("4276555555555555=12345678901234567890"),
			F37:  field.NewStringValue("987654321001"),
			F41:  field.NewStringValue("00000321"),
			F42:  field.NewStringValue("120000000000034"),
			F43:  field.NewStringValue("Test text"),
			F49:  field.NewNumericValue(643),
			F52:  field.NewStringValue(string([]byte{1, 2, 3, 4, 5, 6, 7, 8})),
			F53:  field.NewNumericValue(1234000000000000),
			F120: field.NewStringValue("Another test text"),
		})

		indexs := []int{2, 3, 4, 7, 11, 12, 13, 14, 19, 22, 25, 32, 35, 37, 41, 42, 43, 49, 52, 53, 120}
		for _, index := range indexs {
			assert.Equal(t, message.GetString(index), expectMessage.GetString(index))
		}

		xmlBuf, err := xml.MarshalIndent(message, "", "\t")
		require.NoError(t, err)
		require.NotNil(t, xmlBuf)

		require.Equal(t, want, string(xmlBuf))
	})
}
