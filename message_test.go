package iso8583

import (
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
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

			// this field will be ignored when packing and
			// unpacking, as bit 65 is a bitmap presence indicator
			65: field.NewString(&field.Spec{
				Length:      1,
				Description: "Settlement Code",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			// this is a field of the third bitmap
			130: field.NewString(&field.Spec{
				Length:      1,
				Description: "Additional Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
		},
	}

	// this test most probably will fail in regular mode,
	// and should fail when is run with -race flag
	t.Run("No data race when accessing fields concurrently", func(t *testing.T) {
		message := NewMessage(spec)

		var wg sync.WaitGroup

		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// calling GetString writes into the map of the
				// set fields
				message.GetString(0)
			}()
		}

		wg.Wait()
	})

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
		err = message.Unpack([]byte(want))
		require.NoError(t, err)

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

	t.Run("Do not pack fields that match the bitmap presence indicator", func(t *testing.T) {
		message := NewMessage(spec)
		message.MTI("0100")
		require.NoError(t, message.Field(65, "1"))
		require.NoError(t, message.Field(130, "1")) // field of third bitmap

		got, err := message.Pack()

		want := "01008000000000000000800000000000000040000000000000001"
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, want, string(got))

		message = NewMessage(spec)

		err = message.Unpack([]byte(want))
		require.NoError(t, err)

		s, err := message.GetMTI()
		require.NoError(t, err)
		require.Equal(t, "0100", s)

		s, err = message.GetString(65)
		require.NoError(t, err)
		require.Equal(t, "", s)

		s, err = message.GetString(130)
		require.NoError(t, err)
		require.Equal(t, "1", s)
	})

	t.Run("Does not fail when packing and unpacking message with three bitmaps", func(t *testing.T) {
		message := NewMessage(spec)
		message.MTI("0100")
		require.NoError(t, message.Field(130, "1")) // field of third bitmap

		got, err := message.Pack()

		require.NoError(t, err)
		require.NotNil(t, got)

		want := "01008000000000000000800000000000000040000000000000001"
		require.Equal(t, want, string(got))

		message = NewMessage(spec)
		err = message.Unpack([]byte(want))

		require.NoError(t, err)

		s, err := message.GetMTI()
		require.NoError(t, err)
		require.Equal(t, "0100", s)

		s, err = message.GetString(130)
		require.NoError(t, err)
		require.Equal(t, "1", s)
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

		data := &ISO87Data{}

		require.NoError(t, message.Unmarshal(data))

		require.Equal(t, "0100", data.F0.Value())
		require.Equal(t, "4242424242424242", data.F2.Value())
		require.Equal(t, "12", data.F3.F1.Value())
		require.Equal(t, "34", data.F3.F2.Value())
		require.Equal(t, "56", data.F3.F3.Value())
		require.Equal(t, "100", data.F4.Value())
	})

	t.Run("Test unpacking message with fields that have native types", func(t *testing.T) {
		type TestISOF3Data struct {
			F1 *string
			F2 string
			F3 string
		}

		type ISO87Data struct {
			F0 *string
			F2 string
			F3 *TestISOF3Data
			F4 string
		}

		message := NewMessage(spec)

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

		data := &ISO87Data{}

		require.NoError(t, message.Unmarshal(data))

		require.NotNil(t, data.F0)
		require.Equal(t, "0100", *data.F0)
		require.Equal(t, "4242424242424242", data.F2)
		require.NotNil(t, data.F3.F1)
		require.Equal(t, "12", *data.F3.F1)
		require.Equal(t, "34", data.F3.F2)
		require.Equal(t, "56", data.F3.F3)
		require.Equal(t, "100", data.F4)
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
		err := message.Marshal(&ISO87Data{
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

	t.Run("Test packing message with fields that have native types", func(t *testing.T) {
		type TestISOF3Data struct {
			F1 string
			F2 string
			F3 string
		}

		type ISO87Data struct {
			F0 *string
			F2 string
			F3 *TestISOF3Data
			F4 string
		}

		messageCode := "0100"
		message := NewMessage(spec)
		err := message.Marshal(&ISO87Data{
			F0: &messageCode,
			F2: "4242424242424242",
			F3: &TestISOF3Data{
				F1: "12",
				F2: "34",
				F3: "56",
			},
			F4: "100",
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
			// TLV
			55: field.NewComposite(&field.Spec{
				Length:      999,
				Description: "ICC Data – EMV Having Multiple Tags",
				Pref:        prefix.ASCII.LLL,
				Tag: &field.TagSpec{
					Enc:  encoding.BerTLVTag,
					Sort: sort.StringsByHex,
				},
				Subfields: map[string]field.Field{
					"9A": field.NewString(&field.Spec{
						Description: "Transaction Date",
						Enc:         encoding.Binary,
						Pref:        prefix.BerTLV,
					}),
					"9F02": field.NewString(&field.Spec{
						Description: "Amount, Authorized (Numeric)",
						Enc:         encoding.Binary,
						Pref:        prefix.BerTLV,
					}),
				},
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

	type TestISOF55Data struct {
		F9A   *field.String
		F9F02 *field.String
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
		F55  *TestISOF55Data
		F120 *field.String
	}

	t.Run("Pack data", func(t *testing.T) {
		message := NewMessage(spec)
		err := message.Marshal(&TestISOData{
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
			F52: field.NewStringValue(string([]byte{1, 2, 3, 4, 5, 6, 7, 8})),
			F53: field.NewNumericValue(1234000000000000),
			F55: &TestISOF55Data{
				F9A:   field.NewStringValue("210720"),
				F9F02: field.NewStringValue("000000000501"),
			},
			F120: field.NewStringValue("Another test text"),
		})
		require.NoError(t, err)

		message.MTI("0100")

		got, err := message.Pack()

		want := []byte{0x30, 0x31, 0x30, 0x30, 0xf2, 0x3c, 0x24, 0x81, 0x28, 0xe0, 0x9a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x31, 0x36, 0x34, 0x32, 0x37, 0x36, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x37, 0x37, 0x37, 0x30, 0x30, 0x30, 0x37, 0x30, 0x31, 0x31, 0x31, 0x31, 0x38, 0x34, 0x34, 0x30, 0x30, 0x30, 0x31, 0x32, 0x33, 0x31, 0x33, 0x31, 0x38, 0x34, 0x34, 0x30, 0x37, 0x30, 0x31, 0x31, 0x39, 0x30, 0x32, 0x6, 0x43, 0x39, 0x30, 0x31, 0x30, 0x32, 0x30, 0x36, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x33, 0x37, 0x34, 0x32, 0x37, 0x36, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x3d, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x39, 0x38, 0x37, 0x36, 0x35, 0x34, 0x33, 0x32, 0x31, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x32, 0x31, 0x31, 0x32, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x54, 0x65, 0x73, 0x74, 0x20, 0x74, 0x65, 0x78, 0x74, 0x64, 0x30, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x31, 0x32, 0x33, 0x34, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x33, 0x9a, 0x6, 0x32, 0x31, 0x30, 0x37, 0x32, 0x30, 0x9f, 0x2, 0xc, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x35, 0x30, 0x31, 0x30, 0x31, 0x37, 0x41, 0x6e, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x20, 0x74, 0x65, 0x73, 0x74, 0x20, 0x74, 0x65, 0x78, 0x74}

		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, want, got)
	})

	t.Run("Unpack data", func(t *testing.T) {
		message := NewMessage(spec)

		rawMsg := []byte{0x30, 0x31, 0x30, 0x30, 0xf2, 0x3c, 0x24, 0x81, 0x28, 0xe0, 0x9a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x31, 0x36, 0x34, 0x32, 0x37, 0x36, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x37, 0x37, 0x37, 0x30, 0x30, 0x30, 0x37, 0x30, 0x31, 0x31, 0x31, 0x31, 0x38, 0x34, 0x34, 0x30, 0x30, 0x30, 0x31, 0x32, 0x33, 0x31, 0x33, 0x31, 0x38, 0x34, 0x34, 0x30, 0x37, 0x30, 0x31, 0x31, 0x39, 0x30, 0x32, 0x6, 0x43, 0x39, 0x30, 0x31, 0x30, 0x32, 0x30, 0x36, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x33, 0x37, 0x34, 0x32, 0x37, 0x36, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x3d, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x39, 0x38, 0x37, 0x36, 0x35, 0x34, 0x33, 0x32, 0x31, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x32, 0x31, 0x31, 0x32, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x54, 0x65, 0x73, 0x74, 0x20, 0x74, 0x65, 0x78, 0x74, 0x64, 0x30, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x31, 0x32, 0x33, 0x34, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x33, 0x9a, 0x6, 0x32, 0x31, 0x30, 0x37, 0x32, 0x30, 0x9f, 0x2, 0xc, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x35, 0x30, 0x31, 0x30, 0x31, 0x37, 0x41, 0x6e, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x20, 0x74, 0x65, 0x73, 0x74, 0x20, 0x74, 0x65, 0x78, 0x74}

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

		data := &TestISOData{}
		require.NoError(t, message.Unmarshal(data))

		assert.Equal(t, "4276555555555555", data.F2.Value())
		assert.Equal(t, "00", data.F3.F1.Value())
		assert.Equal(t, "00", data.F3.F2.Value())
		assert.Equal(t, "00", data.F3.F3.Value())
		assert.Equal(t, int64(77700), data.F4.Value())
		assert.Equal(t, int64(701111844), data.F7.Value())
		assert.Equal(t, int64(123), data.F11.Value())
		assert.Equal(t, int64(131844), data.F12.Value())
		assert.Equal(t, int64(701), data.F13.Value())
		assert.Equal(t, int64(1902), data.F14.Value())
		assert.Equal(t, int64(643), data.F19.Value())
		assert.Equal(t, int64(901), data.F22.Value())
		assert.Equal(t, int64(2), data.F25.Value())
		assert.Equal(t, int64(123456), data.F32.Value())
		assert.Equal(t, "4276555555555555=12345678901234567890", data.F35.Value())
		assert.Equal(t, "987654321001", data.F37.Value())
		assert.Equal(t, "00000321", data.F41.Value())
		assert.Equal(t, "120000000000034", data.F42.Value())
		assert.Equal(t, "Test text", data.F43.Value())
		assert.Equal(t, int64(643), data.F49.Value())
		assert.Nil(t, data.F50)
		assert.Equal(t, string([]byte{1, 2, 3, 4, 5, 6, 7, 8}), data.F52.Value())
		assert.Equal(t, int64(1234000000000000), data.F53.Value())
		assert.Equal(t, "210720", data.F55.F9A.Value())
		assert.Equal(t, "000000000501", data.F55.F9F02.Value())
		assert.Equal(t, "Another test text", data.F120.Value())
	})

	t.Run("Pack invalid message returns error of *PackError type", func(t *testing.T) {
		message := NewMessage(spec)
		message.MTI("1")

		_, err := message.Pack()
		require.Error(t, err)

		var packErr *PackError
		require.ErrorAs(t, err, &packErr)
	})

	t.Run("Unpack nil", func(t *testing.T) {
		message := NewMessage(spec)

		err := message.Unpack(nil)

		require.Error(t, err)

		var unpackError *UnpackError
		require.ErrorAs(t, err, &unpackError)
	})

	t.Run("Unpack short mti", func(t *testing.T) {
		message := NewMessage(spec)

		rawMsg := []byte{0x30, 0x31}

		err := message.Unpack([]byte(rawMsg))

		require.Error(t, err)

		var unpackError *UnpackError
		require.ErrorAs(t, err, &unpackError)
		require.Equal(t, rawMsg, unpackError.RawMessage)
	})

	t.Run("Unpack data field error on final field returns partial message", func(t *testing.T) {
		message := NewMessage(spec)

		rawMsg := []byte{0x30, 0x31, 0x30, 0x30, 0xf2, 0x3c, 0x24, 0x81, 0x28, 0xe0, 0x9a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x31, 0x36, 0x34, 0x32, 0x37, 0x36, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x37, 0x37, 0x37, 0x30, 0x30, 0x30, 0x37, 0x30, 0x31, 0x31, 0x31, 0x31, 0x38, 0x34, 0x34, 0x30, 0x30, 0x30, 0x31, 0x32, 0x33, 0x31, 0x33, 0x31, 0x38, 0x34, 0x34, 0x30, 0x37, 0x30, 0x31, 0x31, 0x39, 0x30, 0x32, 0x6, 0x43, 0x39, 0x30, 0x31, 0x30, 0x32, 0x30, 0x36, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x33, 0x37, 0x34, 0x32, 0x37, 0x36, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x3d, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x39, 0x38, 0x37, 0x36, 0x35, 0x34, 0x33, 0x32, 0x31, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x32, 0x31, 0x31, 0x32, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x54, 0x65, 0x73, 0x74, 0x20, 0x74, 0x65, 0x78, 0x74, 0x64, 0x30, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x31, 0x32, 0x33, 0x34, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x33, 0x9a, 0x6, 0x32, 0x31, 0x30, 0x37, 0x32, 0x30, 0x9f, 0x2, 0xc, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x35, 0x30, 0x31, 0x30, 0x31, 0x37, 0x41, 0x6e, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x20, 0x74, 0x65, 0x73, 0x74, 0x20, 0x74, 0x65, 0x78}

		err := message.Unpack([]byte(rawMsg))

		require.Error(t, err)
		var unpackError *UnpackError
		require.ErrorAs(t, err, &unpackError)
		assert.ElementsMatch(t, unpackError.Fields, []string{"Field 120"})

		s, err := message.GetString(2)
		require.NoError(t, err)
		require.Equal(t, "4276555555555555", s)

		s, err = message.GetString(3)
		require.NoError(t, err)
		require.Equal(t, "000000", s)

		s, err = message.GetString(4)
		require.NoError(t, err)
		require.Equal(t, "77700", s)

		data := &TestISOData{}
		require.NoError(t, message.Unmarshal(data))

		assert.Equal(t, "4276555555555555", data.F2.Value())
		assert.Equal(t, "00", data.F3.F1.Value())
		assert.Equal(t, "00", data.F3.F2.Value())
		assert.Equal(t, "00", data.F3.F3.Value())
		assert.Equal(t, int64(77700), data.F4.Value())
		assert.Equal(t, int64(701111844), data.F7.Value())
		assert.Equal(t, int64(123), data.F11.Value())
		assert.Equal(t, int64(131844), data.F12.Value())
		assert.Equal(t, int64(701), data.F13.Value())
		assert.Equal(t, int64(1902), data.F14.Value())
		assert.Equal(t, int64(643), data.F19.Value())
		assert.Equal(t, int64(901), data.F22.Value())
		assert.Equal(t, int64(2), data.F25.Value())
		assert.Equal(t, int64(123456), data.F32.Value())
		assert.Equal(t, "4276555555555555=12345678901234567890", data.F35.Value())
		assert.Equal(t, "987654321001", data.F37.Value())
		assert.Equal(t, "00000321", data.F41.Value())
		assert.Equal(t, "120000000000034", data.F42.Value())
		assert.Equal(t, "Test text", data.F43.Value())
		assert.Equal(t, int64(643), data.F49.Value())
		assert.Nil(t, data.F50)
		assert.Equal(t, string([]byte{1, 2, 3, 4, 5, 6, 7, 8}), data.F52.Value())
		assert.Equal(t, int64(1234000000000000), data.F53.Value())
		assert.Equal(t, "210720", data.F55.F9A.Value())
		assert.Equal(t, "000000000501", data.F55.F9F02.Value())
		assert.Empty(t, data.F120)
	})

	t.Run("Unpack data field error on middle field with fixed prefix returns partial message", func(t *testing.T) {
		corruptSpecOut := &MessageSpec{
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
				3: field.NewNumeric(&field.Spec{
					Length:      3,
					Description: "Dodgy field",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				4: field.NewString(&field.Spec{
					Length:      6,
					Description: "Anything",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.LL,
				}),
			},
		}

		corruptSpecIn := &MessageSpec{
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
					Length:      3,
					Description: "Dodgy field",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				4: field.NewString(&field.Spec{
					Length:      6,
					Description: "Anything",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.LL,
				}),
			},
		}

		type TestCorruptISODataIn struct {
			F2 *field.String
			F3 *field.String
			F4 *field.String
		}

		type TestCorruptISODataOut struct {
			F2 *field.String
			F3 *field.Numeric
			F4 *field.String
		}

		message := NewMessage(corruptSpecIn)
		err := message.Marshal(&TestCorruptISODataIn{
			F2: field.NewStringValue("4276555555555555"),
			F3: field.NewStringValue("ABC"),
			F4: field.NewStringValue("123"),
		})
		require.NoError(t, err)

		message.MTI("0100")

		rawMsg, err := message.Pack()
		require.NoError(t, err)

		receivedMessage := NewMessage(corruptSpecOut)

		err = receivedMessage.Unpack([]byte(rawMsg))

		require.Error(t, err)
		var unpackErr *UnpackError
		require.ErrorAs(t, err, &unpackErr)
		require.Equal(t, []string{"Dodgy field"}, unpackErr.Fields)

		s, err := receivedMessage.GetString(2)
		require.NoError(t, err)
		require.Equal(t, "4276555555555555", s)

		s, err = receivedMessage.GetString(3)
		require.NoError(t, err)
		require.Equal(t, "0", s)

		s, err = receivedMessage.GetString(4)
		require.NoError(t, err)
		require.Equal(t, "123", s)

		data := &TestCorruptISODataOut{}
		err = message.Unmarshal(data)
		require.Error(t, err)

		assert.Equal(t, "4276555555555555", data.F2.Value())
		assert.Equal(t, int64(0), data.F3.Value())
		assert.Equal(t, "123", data.F4.Value())
	})

	// this test should check that BCD fields are packed and
	// unpacked correctly it's a confirmation that issue
	// https://github.com/moov-io/iso8583/issues/220 is fixed
	t.Run("Pack and Unpack BCD fields", func(t *testing.T) {
		var spec = &MessageSpec{
			Fields: map[int]field.Field{
				0: field.NewNumeric(&field.Spec{
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
				2: field.NewNumeric(&field.Spec{
					Length:      4,
					Description: "SomeFixedField",
					Enc:         encoding.BCD,
					Pref:        prefix.BCD.Fixed,
				}),
				3: field.NewNumeric(&field.Spec{
					Length:      3,
					Description: "SomeVarField",
					Enc:         encoding.BCD,
					Pref:        prefix.BCD.LLLL,
				}),
			},
		}

		msg := NewMessage(spec)

		msg.MTI("1234")
		msg.Field(2, "4567")
		msg.Field(3, "890")

		out, err := msg.Pack()
		require.NoError(t, err)

		got := hex.EncodeToString(out)

		expected := "1234" + // MTI
			"6000000000000000" + // Bitmap
			"4567" + // SomeFixedField
			"0003" + // LLLL in BCD
			"0890" // SomeVarField in BCD 0x08 0x90

		require.Equal(t, expected, got)

		in := NewMessage(spec)

		err = in.Unpack(out)
		require.NoError(t, err)

		result, _ := in.GetField(2).String()
		require.Equal(t, "4567", result)

		result, _ = in.GetField(3).String()
		require.Equal(t, "890", result)
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
			45: field.NewTrack1(&field.Spec{
				Length:      76,
				Description: "Track 1 Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
		},
	}

	type TestISOF3Data struct {
		F1 *field.String
		F2 *field.String
		F3 *field.String
	}

	type TestISOData struct {
		F0  *field.String
		F2  *field.String
		F3  *TestISOF3Data
		F4  *field.String
		F45 *field.Track1
	}

	t.Run("Test JSON encoding typed", func(t *testing.T) {
		expDate, err := time.Parse("0601", "9901")
		require.NoError(t, err)

		message := NewMessage(spec)
		err = message.Marshal(&TestISOData{
			F0: field.NewStringValue("0100"),
			F2: field.NewStringValue("4242424242424242"),
			F3: &TestISOF3Data{
				F1: field.NewStringValue("12"),
				F2: field.NewStringValue("34"),
				F3: field.NewStringValue("56"),
			},
			F4: field.NewStringValue("100"),
			F45: &field.Track1{
				FixedLength:          true,
				FormatCode:           "B",
				PrimaryAccountNumber: "1234567890123445",
				ServiceCode:          "120",
				DiscretionaryData:    "0000000000000**XXX******",
				ExpirationDate:       &expDate,
				Name:                 "PADILLA/L.",
			},
		})
		require.NoError(t, err)

		want := `{"0":"0100","1":"7000000000080000","2":"4242424242424242","3":{"1":"12","2":"34","3":"56"},"4":"100","45":{"fixed_length":true,"format_code":"B","primary_account_number":"1234567890123445","name":"PADILLA/L.","expiration_date":"1999-01-01T00:00:00Z","service_code":"120","discretionary_data":"0000000000000**XXX******"}}`

		got, err := json.Marshal(message)
		require.NoError(t, err)
		require.Equal(t, want, string(got))
	})

	t.Run("Test JSON encoding untyped", func(t *testing.T) {
		message := NewMessage(spec)
		message.MTI("0100")
		message.Field(2, "4242424242424242")
		message.Field(4, "100")

		want := `{"0":"0100","1":"5000000000000000","2":"4242424242424242","4":"100"}`

		got, err := json.Marshal(message)
		require.NoError(t, err)
		require.Equal(t, want, string(got))
	})

	t.Run("Test JSON encoding of unpacked fields typed", func(t *testing.T) {
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

		want := `{"0":"0100","1":"7000000000000000","2":"4242424242424242","3":{"1":"12","2":"34","3":"56"},"4":"100"}`

		message := NewMessage(spec)
		message.Marshal(&ISO87Data{})

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		require.NoError(t, message.Unpack([]byte(rawMsg)))
		got, err := json.Marshal(message)
		require.NoError(t, err)

		require.Equal(t, want, string(got))
	})

	t.Run("Test JSON encoding of unpacked fields untyped", func(t *testing.T) {
		want := `{"0":"0100","1":"7000000000000000","2":"4242424242424242","3":{"1":"12","2":"34","3":"56"},"4":"100"}`

		message := NewMessage(spec)

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		require.NoError(t, message.Unpack([]byte(rawMsg)))
		got, err := json.Marshal(message)
		require.NoError(t, err)

		require.Equal(t, want, string(got))
	})

	t.Run("Test JSON decoding typed", func(t *testing.T) {
		message := NewMessage(spec)

		input := []byte(`{"0":"0100","1":"7000000000000000","2":"4242424242424242","3":{"1":"12","2":"34","3":"56"},"4":"100"}`)

		want := &TestISOData{
			F0: field.NewStringValue("0100"),
			F2: field.NewStringValue("4242424242424242"),
			F3: &TestISOF3Data{
				F1: field.NewStringValue("12"),
				F2: field.NewStringValue("34"),
				F3: field.NewStringValue("56"),
			},
			F4: field.NewStringValue("100"),
		}

		require.NoError(t, json.Unmarshal(input, message))

		data := &TestISOData{}
		require.NoError(t, message.Unmarshal(data))

		require.Equal(t, want.F0.Value(), data.F0.Value())
		require.Equal(t, want.F2.Value(), data.F2.Value())
		require.Equal(t, want.F3.F1.Value(), data.F3.F1.Value())
		require.Equal(t, want.F3.F2.Value(), data.F3.F2.Value())
		require.Equal(t, want.F3.F3.Value(), data.F3.F3.Value())
		require.Equal(t, want.F4.Value(), data.F4.Value())
	})

	t.Run("Test JSON decoding untyped", func(t *testing.T) {
		message := NewMessage(spec)

		input := `{"0":"0100","1":"5000000000000000","2":"4242424242424242","4":"100"}`

		err := json.Unmarshal([]byte(input), message)
		require.NoError(t, err)

		mti, err := message.GetMTI()
		require.NoError(t, err)
		require.Equal(t, "0100", mti)

		f2, err := message.GetString(2)
		require.NoError(t, err)
		require.Equal(t, "4242424242424242", f2)

		f4, err := message.GetString(4)
		require.NoError(t, err)
		require.Equal(t, "100", f4)
	})

}

func TestMessageClone(t *testing.T) {
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
			52: field.NewBinary(&field.Spec{
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
			// TLV
			55: field.NewComposite(&field.Spec{
				Length:      999,
				Description: "ICC Data – EMV Having Multiple Tags",
				Pref:        prefix.ASCII.LLL,
				Tag: &field.TagSpec{
					Enc:  encoding.BerTLVTag,
					Sort: sort.StringsByHex,
				},
				Subfields: map[string]field.Field{
					"9A": field.NewString(&field.Spec{
						Description: "Transaction Date",
						Enc:         encoding.Binary,
						Pref:        prefix.BerTLV,
					}),
					"9F02": field.NewString(&field.Spec{
						Description: "Amount, Authorized (Numeric)",
						Enc:         encoding.Binary,
						Pref:        prefix.BerTLV,
					}),
				},
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

	type TestISOF55Data struct {
		F9A   *field.String
		F9F02 *field.String
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
		F52  *field.Binary
		F53  *field.Numeric
		F55  *TestISOF55Data
		F120 *field.String
	}

	message := NewMessage(spec)
	data2 := &TestISOData{
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
		F52: field.NewBinaryValue([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
		F53: field.NewNumericValue(1234000000000000),
		F55: &TestISOF55Data{
			F9A:   field.NewStringValue("210720"),
			F9F02: field.NewStringValue("000000000501"),
		},
		F120: field.NewStringValue("Another test text"),
	}
	require.NoError(t, message.Marshal(data2))

	message.MTI("0100")

	got, err := message.Pack()

	want := []byte{0x30, 0x31, 0x30, 0x30, 0xf2, 0x3c, 0x24, 0x81, 0x28, 0xe0, 0x9a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x31, 0x36, 0x34, 0x32, 0x37, 0x36, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x37, 0x37, 0x37, 0x30, 0x30, 0x30, 0x37, 0x30, 0x31, 0x31, 0x31, 0x31, 0x38, 0x34, 0x34, 0x30, 0x30, 0x30, 0x31, 0x32, 0x33, 0x31, 0x33, 0x31, 0x38, 0x34, 0x34, 0x30, 0x37, 0x30, 0x31, 0x31, 0x39, 0x30, 0x32, 0x6, 0x43, 0x39, 0x30, 0x31, 0x30, 0x32, 0x30, 0x36, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x33, 0x37, 0x34, 0x32, 0x37, 0x36, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x35, 0x3d, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x39, 0x38, 0x37, 0x36, 0x35, 0x34, 0x33, 0x32, 0x31, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x32, 0x31, 0x31, 0x32, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x34, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x54, 0x65, 0x73, 0x74, 0x20, 0x74, 0x65, 0x78, 0x74, 0x64, 0x30, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x31, 0x32, 0x33, 0x34, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x33, 0x9a, 0x6, 0x32, 0x31, 0x30, 0x37, 0x32, 0x30, 0x9f, 0x2, 0xc, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x35, 0x30, 0x31, 0x30, 0x31, 0x37, 0x41, 0x6e, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x20, 0x74, 0x65, 0x73, 0x74, 0x20, 0x74, 0x65, 0x78, 0x74}

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, want, got)

	message2, err := message.Clone()
	require.NoError(t, err)

	require.Equal(t, message.spec, message2.spec)
	require.Equal(t, message.Bitmap(), message2.Bitmap())

	mti, err := message.GetMTI()
	require.NoError(t, err)

	mti2, err := message2.GetMTI()
	require.NoError(t, err)

	require.Equal(t, mti, mti2)

	messageData := &TestISOData{}
	message2Data := &TestISOData{}

	require.NoError(t, message.Unmarshal(messageData))
	require.NoError(t, message2.Unmarshal(message2Data))

	require.Equal(t, messageData.F2.Value(), message2Data.F2.Value())
	require.Equal(t, messageData.F3.F1.Value(), message2Data.F3.F1.Value())
	require.Equal(t, messageData.F3.F2.Value(), message2Data.F3.F2.Value())
	require.Equal(t, messageData.F3.F3.Value(), message2Data.F3.F3.Value())
	require.Equal(t, messageData.F4.Value(), message2Data.F4.Value())
	require.Equal(t, messageData.F7.Value(), message2Data.F7.Value())
	require.Equal(t, messageData.F11.Value(), message2Data.F11.Value())
	require.Equal(t, messageData.F12.Value(), message2Data.F12.Value())
	require.Equal(t, messageData.F13.Value(), message2Data.F13.Value())
	require.Equal(t, messageData.F14.Value(), message2Data.F14.Value())
	require.Equal(t, messageData.F19.Value(), message2Data.F19.Value())
	require.Equal(t, messageData.F22.Value(), message2Data.F22.Value())
	require.Equal(t, messageData.F25.Value(), message2Data.F25.Value())
	require.Equal(t, messageData.F32.Value(), message2Data.F32.Value())
	require.Equal(t, messageData.F35.Value(), message2Data.F35.Value())
	require.Equal(t, messageData.F37.Value(), message2Data.F37.Value())
	require.Equal(t, messageData.F41.Value(), message2Data.F41.Value())
	require.Equal(t, messageData.F42.Value(), message2Data.F42.Value())
	require.Equal(t, messageData.F43.Value(), message2Data.F43.Value())
	require.Equal(t, messageData.F49.Value(), message2Data.F49.Value())
	require.Equal(t, messageData.F52.Value(), message2Data.F52.Value())
	require.Equal(t, messageData.F53.Value(), message2Data.F53.Value())
	require.Equal(t, messageData.F55.F9A.Value(), message2Data.F55.F9A.Value())
	require.Equal(t, messageData.F55.F9F02.Value(), message2Data.F55.F9F02.Value())
	require.Equal(t, messageData.F120.Value(), message2Data.F120.Value())

	message3, err := message2.Clone()
	require.NoError(t, err)

	require.Equal(t, message2.spec, message3.spec)
	require.Equal(t, message2.Bitmap(), message3.Bitmap())

	mti3, err := message.GetMTI()
	require.NoError(t, err)

	require.Equal(t, mti2, mti3)

	message3Data := &TestISOData{}
	require.NoError(t, message3.Unmarshal(message3Data))

	require.Equal(t, message2Data.F2.Value(), message3Data.F2.Value())
	require.Equal(t, message2Data.F3.F1.Value(), message3Data.F3.F1.Value())
	require.Equal(t, message2Data.F3.F2.Value(), message3Data.F3.F2.Value())
	require.Equal(t, message2Data.F3.F3.Value(), message3Data.F3.F3.Value())
	require.Equal(t, message2Data.F4.Value(), message3Data.F4.Value())
	require.Equal(t, message2Data.F7.Value(), message3Data.F7.Value())
	require.Equal(t, message2Data.F11.Value(), message3Data.F11.Value())
	require.Equal(t, message2Data.F12.Value(), message3Data.F12.Value())
	require.Equal(t, message2Data.F13.Value(), message3Data.F13.Value())
	require.Equal(t, message2Data.F14.Value(), message3Data.F14.Value())
	require.Equal(t, message2Data.F19.Value(), message3Data.F19.Value())
	require.Equal(t, message2Data.F22.Value(), message3Data.F22.Value())
	require.Equal(t, message2Data.F25.Value(), message3Data.F25.Value())
	require.Equal(t, message2Data.F32.Value(), message3Data.F32.Value())
	require.Equal(t, message2Data.F35.Value(), message3Data.F35.Value())
	require.Equal(t, message2Data.F37.Value(), message3Data.F37.Value())
	require.Equal(t, message2Data.F41.Value(), message3Data.F41.Value())
	require.Equal(t, message2Data.F42.Value(), message3Data.F42.Value())
	require.Equal(t, message2Data.F43.Value(), message3Data.F43.Value())
	require.Equal(t, message2Data.F49.Value(), message3Data.F49.Value())
	require.Equal(t, message2Data.F52.Value(), message3Data.F52.Value())
	require.Equal(t, message2Data.F53.Value(), message3Data.F53.Value())
	require.Equal(t, message2Data.F55.F9A.Value(), message3Data.F55.F9A.Value())
	require.Equal(t, message2Data.F55.F9F02.Value(), message3Data.F55.F9F02.Value())
	require.Equal(t, message2Data.F120.Value(), message3Data.F120.Value())
}

func TestMessageMarshaling(t *testing.T) {
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
		err := message.Marshal(&ISO87Data{
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
		require.NoError(t, message.Unmarshal(data))

		require.Equal(t, "0100", data.F0.Value())
		require.Equal(t, "4242424242424242", data.F2.Value())
		require.Equal(t, "12", data.F3.F1.Value())
		require.Equal(t, "34", data.F3.F2.Value())
		require.Equal(t, "56", data.F3.F3.Value())
		require.Equal(t, "100", data.F4.Value())
	})

	t.Run("Marshal nil returns nil", func(t *testing.T) {
		message := NewMessage(spec)

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		err = message.Marshal(nil)
		require.NoError(t, err)
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
		require.NoError(t, message.Marshal(data))

		rawMsg, err := message.Pack()
		require.NoError(t, err)

		expected := []byte("01007000000000000000164242424242424242123456000000000100")
		require.Equal(t, expected, rawMsg)
	})

	t.Run("Marshal when no idex is set for the fields", func(t *testing.T) {
		type ISO87Data struct {
			MTI                  *field.String
			PrimaryAccountNumber *field.String
			Amount               *field.String
		}

		data := &ISO87Data{
			MTI:                  field.NewStringValue("0100"),
			PrimaryAccountNumber: field.NewStringValue("4242424242424242"),
			Amount:               field.NewStringValue("100"),
		}

		message := NewMessage(spec)
		require.NoError(t, message.Marshal(data))

		rawMsg, err := message.Pack()
		require.NoError(t, err)
		// only bitmap is packed => 8 zero bytes in hex
		require.Equal(t, strings.Repeat("0", 16), string(rawMsg))
	})

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
		require.NoError(t, message.Unmarshal(data))

		require.Equal(t, "0100", data.F0.Value())
		require.Equal(t, "4242424242424242", data.F2.Value())
		require.Equal(t, "12", data.F3.F1.Value())
		require.Equal(t, "34", data.F3.F2.Value())
		require.Equal(t, "56", data.F3.F3.Value())
		require.Equal(t, "100", data.F4.Value())
	})

	t.Run("Unmarshal into nil", func(t *testing.T) {
		message := NewMessage(spec)

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))
		require.NoError(t, err)

		require.Error(t, message.Unmarshal(nil))
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
		require.NoError(t, message.Unmarshal(data))

		require.Equal(t, "0100", data.MTI.Value())
		require.Equal(t, "4242424242424242", data.PrimaryAccountNumber.Value())
		require.Equal(t, "12", data.AdditionalData.One.Value())
		require.Equal(t, "34", data.AdditionalData.Two.Value())
		require.Equal(t, "56", data.AdditionalData.Three.Value())
		require.Equal(t, "100", data.Amount.Value())
	})

	t.Run("Unmarshal skips fields with no index", func(t *testing.T) {
		type ISO87Data struct {
			MTI                  *field.String
			PrimaryAccountNumber *field.String
			Amount               *field.String
		}

		message := NewMessage(spec)

		rawMsg := []byte("01007000000000000000164242424242424242123456000000000100")
		err := message.Unpack([]byte(rawMsg))

		require.NoError(t, err)

		data := &ISO87Data{}
		require.NoError(t, message.Unmarshal(data))
	})
}

func FuzzUnpack(f *testing.F) {
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

	f.Add([]byte("01007000000000000000164242424242424242123456000000000100")) // Use f.Add to provide a seed corpus

	f.Fuzz(func(t *testing.T, orig []byte) {
		message := NewMessage(spec)
		// we only care when it panics
		message.Unpack(orig)
	})
}

func TestStructWithTypes(t *testing.T) {
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
		panByte := []byte("4242424242424242")

		tests := []struct {
			name                 string
			input                interface{}
			expectedPackedString string
			isError              bool
			errorString          string
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

			// Tests for []byte type
			{
				name: "struct with []byte type and value set",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber []byte `index:"2"`
				}{
					MTI:                  "0110",
					PrimaryAccountNumber: panByte,
				},
				expectedPackedString: "011040000000000000000000000000000000164242424242424242",
				isError:              true,
				errorString:          "failed to set value to field 2: data does not match required *String or (string, *string, int, *int) type",
			},
			{
				name: "struct with []byte type and no value",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber []byte `index:"2"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011000000000000000000000000000000000",
			},
			{
				name: "struct with []byte type, no value and keepzero tag",
				input: struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber []byte `index:"2,keepzero"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011040000000000000000000000000000000010",
				isError:              true,
				errorString:          "failed to set value to field 2: data does not match required *String or (string, *string, int, *int) type",
			},

			// Tests for *[]byte type
			{
				name: "struct with *[]byte type and value set",
				input: struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *[]byte `index:"2"`
				}{
					MTI:                  "0110",
					PrimaryAccountNumber: &panByte,
				},
				expectedPackedString: "011040000000000000000000000000000000164242424242424242",
				isError:              true,
				errorString:          "failed to set value to field 2: data does not match required *String or (string, *string, int, *int) type",
			},
			{
				name: "struct with *[]byte type and no value",
				input: struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *[]byte `index:"2"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011000000000000000000000000000000000",
				isError:              false, // there is not any modification
			},
			{
				name: "struct with *[]byte type, no value and keepzero tag",
				input: struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *[]byte `index:"2,keepzero"`
				}{
					MTI: "0110",
				},
				expectedPackedString: "011040000000000000000000000000000000010",
				isError:              true,
				errorString:          "failed to set value to field 2: data does not match required *String or (string, *string, int, *int) type",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				message := NewMessage(spec)
				err := message.Marshal(tt.input)
				if tt.isError {
					require.Error(t, err)
					require.Equal(t, tt.errorString, err.Error())
					return
				}

				require.NoError(t, err)

				packed, err := message.Pack()
				require.NoError(t, err)

				require.Equal(t, tt.expectedPackedString, string(packed))
			})
		}
	})

	t.Run("unpack", func(t *testing.T) {
		tests := []struct {
			name        string
			input       any
			isError     bool
			errorString string
		}{
			// Tests for string type
			{
				name: "struct with string type",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber string `index:"2"`
				}{},
			},
			{
				name: "struct with string type with keepzero tag",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber string `index:"2,keepzero"`
				}{},
			},

			// Tests for *string type
			{
				name: "struct with *string type",
				input: &struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *string `index:"2"`
				}{},
			},
			{
				name: "struct with *string type with keepzero tag",
				input: &struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *string `index:"2,keepzero"`
				}{},
			},

			// Tests for int type
			{
				name: "struct with int type",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber int    `index:"2"`
				}{},
			},
			{
				name: "struct with int type with keepzero tag",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber int    `index:"2,keepzero"`
				}{},
			},

			// Tests for *int type
			{
				name: "struct with *int type",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber *int   `index:"2"`
				}{},
			},
			{
				name: "struct with *int type with keepzero tag",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber *int   `index:"2,keepzero"`
				}{},
			},

			// Tests for []byte type
			{
				name: "struct with []byte type",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber []byte `index:"2"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got []uint8",
			},
			{
				name: "struct with []byte type with keepzero tag",
				input: &struct {
					MTI                  string `index:"0"`
					PrimaryAccountNumber []byte `index:"2,keepzero"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got []uint8",
			},

			// Tests for *[]byte type
			{
				name: "struct with []byte type",
				input: &struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *[]byte `index:"2"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got *[]uint8",
			},
			{
				name: "struct with []byte type with keepzero tag",
				input: &struct {
					MTI                  string  `index:"0"`
					PrimaryAccountNumber *[]byte `index:"2,keepzero"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got *[]uint8",
			},

			// Tests for []string type
			{
				name: "struct with []string type",
				input: &struct {
					MTI                  string   `index:"0"`
					PrimaryAccountNumber []string `index:"2"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got []string",
			},
			{
				name: "struct with []string type with keepzero tag",
				input: &struct {
					MTI                  string   `index:"0"`
					PrimaryAccountNumber []string `index:"2,keepzero"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got []string",
			},

			// Tests for *[]string type
			{
				name: "struct with *[]string type",
				input: &struct {
					MTI                  string    `index:"0"`
					PrimaryAccountNumber *[]string `index:"2"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got *[]string",
			},
			{
				name: "struct with *[]string type with keepzero tag",
				input: &struct {
					MTI                  string    `index:"0"`
					PrimaryAccountNumber *[]string `index:"2,keepzero"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got *[]string",
			},

			// Tests for map[string]string type
			{
				name: "struct with map[string]string type",
				input: &struct {
					MTI                  string            `index:"0"`
					PrimaryAccountNumber map[string]string `index:"2"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported reflect.Value type: map",
			},
			{
				name: "struct with map[string]string type with keepzero tag",
				input: &struct {
					MTI                  string            `index:"0"`
					PrimaryAccountNumber map[string]string `index:"2,keepzero"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported reflect.Value type: map",
			},

			// Tests for *map[string]string type
			{
				name: "struct with *map[string]string type",
				input: &struct {
					MTI                  string             `index:"0"`
					PrimaryAccountNumber *map[string]string `index:"2"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got *map[string]string",
			},
			{
				name: "struct with *map[string]string type with keepzero tag",
				input: &struct {
					MTI                  string             `index:"0"`
					PrimaryAccountNumber *map[string]string `index:"2,keepzero"`
				}{},
				isError:     true,
				errorString: "failed to get value from field 2: unsupported type: expected *String, *string, or reflect.Value, got *map[string]string",
			},
		}
		packed := []byte("011040000000000000000000000000000000164242424242424242")

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				message := NewMessage(spec)

				err := message.Unpack(packed)
				require.NoError(t, err)

				err = message.Unmarshal(tt.input)
				if tt.isError {
					require.Error(t, err)
					require.Equal(t, tt.errorString, err.Error())
					return
				}

				require.NoError(t, err)

				val := reflect.Indirect(reflect.ValueOf(tt.input))

				require.Equal(t, "0110", val.Field(0).String())
				switch val.Field(1).Type().String() {
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
