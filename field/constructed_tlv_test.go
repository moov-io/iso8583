package field

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

var (
	constructedTestSpec = &Spec{
		Length:      3,
		Description: "Terminal Capabilities",
		Enc:         encoding.Binary,
		Pref:        prefix.BerTLV,
		Tag: &TagSpec{
			Tag:  "9F33",
			Enc:  encoding.BerTLVTag,
			Sort: sort.StringsByHex,
		},
		Subfields: map[string]Field{
			"01": NewPrimitiveTLV(&Spec{
				Length:      2,
				Description: "Application Interchange Profile",
				Enc:         encoding.ASCII,
				Pref:        prefix.BerTLV,
				Tag: &TagSpec{
					Tag:  "82",
					Enc:  encoding.BerTLVTag,
					Sort: sort.StringsByHex,
				},
			}),
			"02": NewPrimitiveTLV(&Spec{
				Length:      2,
				Description: "Application Transaction Counter",
				Enc:         encoding.ASCII,
				Pref:        prefix.BerTLV,
				Tag: &TagSpec{
					Tag:  "9F36",
					Enc:  encoding.BerTLVTag,
					Sort: sort.StringsByHex,
				},
			}),
			"03": NewConstructedTLV(&Spec{
				Length:      8,
				Description: "Currency Code, Application Reference",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
				Tag: &TagSpec{
					Tag:  "9F3B",
					Enc:  encoding.BerTLVTag,
					Sort: sort.StringsByHex,
				},
				Subfields: map[string]Field{
					"04": NewPrimitiveTLV(&Spec{
						Length:      2,
						Description: "Data Authentication Code",
						Enc:         encoding.ASCII,
						Pref:        prefix.BerTLV,
						Tag: &TagSpec{
							Tag:  "9F45",
							Enc:  encoding.BerTLVTag,
							Sort: sort.StringsByHex,
						},
					}),
				},
			}),
		},
	}
)

type ConstructedCurrencyCode struct {
	F04 *PrimitiveTLV
}

type ConstructedTerminal struct {
	F01 *PrimitiveTLV
	F02 *PrimitiveTLV
	F03 *ConstructedCurrencyCode
}

func TestConstructedTLVPacking(t *testing.T) {
	t.Run("Pack returns an null tlv when setting mismatched struct", func(t *testing.T) {
		type TestDataIncorrectType struct {
			F1 *PrimitiveTLV
		}
		tlv := NewConstructedTLV(constructedTestSpec)
		err := tlv.SetData(&TestDataIncorrectType{
			F1: NewPrimitiveTLVValue([]byte{0x0, 0x93}),
		})

		require.NoError(t, err)
		require.Equal(t, 0, len(tlv.getSubfields()))

		packed, err := tlv.Pack()
		require.NoError(t, err)
		require.Equal(t, []byte{0x9F, 0x33, 0x00}, packed)
	})

	t.Run("Pack returns nested tlv struct", func(t *testing.T) {

		tlv := NewConstructedTLV(constructedTestSpec)

		err := tlv.SetData(&ConstructedTerminal{
			F01: NewPrimitiveTLVValue([]byte{0x01, 0x7f}),
			F02: NewPrimitiveTLVValue([]byte{0x02, 0x7f}),
			F03: &ConstructedCurrencyCode{
				F04: NewPrimitiveTLVValue([]byte{0x04, 0x7f}),
			},
		})
		require.NoError(t, err)

		packed, err := tlv.Pack()

		require.NoError(t, err)
		require.Equal(t, "9F33118202017F9F3602027F9F3B059F4502047F", fmt.Sprintf("%X", packed))
	})

	t.Run("Pack returns an error on failure of invalid value", func(t *testing.T) {
		tlv := NewConstructedTLV(constructedTestSpec)

		err := tlv.SetData(&ConstructedTerminal{
			F01: NewPrimitiveTLVValue([]byte{0x01, 0xff}),
			F02: NewPrimitiveTLVValue([]byte{0x02, 0xff}),
			F03: &ConstructedCurrencyCode{
				F04: NewPrimitiveTLVValue([]byte{0x04, 0xff}),
			},
		})
		require.NoError(t, err)

		_, err = tlv.Pack()
		require.Error(t, err)
		require.Equal(
			t,
			"failed to pack subfield 01: failed to encode content: failed to perform ASCII encoding",
			err.Error())
	})
}

func TestConstructedTLVUnpacking(t *testing.T) {

	hexString := "9F33118202017F9F3602027F9F3B059F4502047F"
	raw, err := encoding.BerTLVTag.Encode([]byte(hexString))
	require.NoError(t, err)

	t.Run("Unpack decode raw bytes without any struct mapping", func(t *testing.T) {
		tlv := NewConstructedTLV(constructedTestSpec)

		read, err := tlv.Unpack(raw)
		require.NoError(t, err)
		require.Equal(t, 20, read)

		packed, err := tlv.Pack()
		require.NoError(t, err)
		require.Equal(t, hexString, fmt.Sprintf("%X", packed))
	})

	t.Run("Unpack decode raw bytes with struct mapping", func(t *testing.T) {
		tlv := NewConstructedTLV(constructedTestSpec)

		dummy := &ConstructedTerminal{
			F01: NewPrimitiveTLVValue(nil),
			F02: NewPrimitiveTLVValue(nil),
			F03: &ConstructedCurrencyCode{
				F04: NewPrimitiveTLVValue(nil),
			},
		}
		err := tlv.SetData(dummy)

		read, err := tlv.Unpack(raw)
		require.NoError(t, err)
		require.Equal(t, 20, read)

		packed, err := tlv.Pack()
		require.NoError(t, err)
		require.Equal(t, hexString, fmt.Sprintf("%X", packed))

		jsonStr, err := json.Marshal(dummy)
		require.NoError(t, err)
		require.Equal(t, `{"F01":"017F","F02":"027F","F03":{"F04":"047F"}}`, string(jsonStr))
	})
}

func TestConstructedTLVGetSetBytes(t *testing.T) {

	hexString := "9F33118202017F9F3602027F9F3B059F4502047F"
	raw, err := encoding.BerTLVTag.Encode([]byte(hexString))
	require.NoError(t, err)

	valueString := "8202017F9F3602027F9F3B059F4502047F"
	value, err := encoding.BerTLVTag.Encode([]byte(valueString))
	require.NoError(t, err)

	tlv := NewConstructedTLV(constructedTestSpec)
	err = tlv.SetBytes(value)
	require.NoError(t, err)

	packed, err := tlv.Pack()
	require.NoError(t, err)
	require.Equal(t, hexString, fmt.Sprintf("%X", packed))

	gotValue, err := tlv.Bytes()
	require.NoError(t, err)
	require.Equal(t, valueString, fmt.Sprintf("%X", gotValue))

	err = tlv.SetBytes(raw)
	require.Error(t, err)
	require.Equal(t, "failed to unpack subfield 82: tag mismatch: want to read 82, got 9F33", err.Error())
}

func TestConstructedTLVGetValue(t *testing.T) {

	tlv := NewConstructedTLV(constructedTestSpec)

	err := tlv.SetData(&ConstructedTerminal{
		F01: NewPrimitiveTLVValue([]byte{0x01, 0x7f}),
		F02: NewPrimitiveTLVValue([]byte{0x02, 0x7f}),
		F03: &ConstructedCurrencyCode{
			F04: NewPrimitiveTLVValue([]byte{0x04, 0x7f}),
		},
	})
	require.NoError(t, err)

	value, err := tlv.GetValue("9F45")
	require.NoError(t, err)
	require.Equal(t, "047F", fmt.Sprintf("%X", value))

}

func TestConstructedTLVSetValue(t *testing.T) {

	tlv := NewConstructedTLV(constructedTestSpec)

	err := tlv.SetData(&ConstructedTerminal{
		F01: NewPrimitiveTLVValue([]byte{0x01, 0x7f}),
		F02: NewPrimitiveTLVValue([]byte{0x02, 0x7f}),
		F03: &ConstructedCurrencyCode{
			F04: NewPrimitiveTLVValue([]byte{0x04, 0x7f}),
		},
	})
	require.NoError(t, err)

	err = tlv.SetValue("9F45", []byte{0x04, 0x6f})
	require.NoError(t, err)

	value, err := tlv.GetValue("9F45")
	require.NoError(t, err)
	require.Equal(t, "046F", fmt.Sprintf("%X", value))

}
