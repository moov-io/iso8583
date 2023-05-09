package field

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

var (
	compositeTestSpec = &Spec{
		Length:      6,
		Description: "Test Spec",
		Pref:        prefix.ASCII.Fixed,
		Pad:         padding.None,
		Tag: &TagSpec{
			Sort: sort.StringsByInt,
		},
		Subfields: map[string]Field{
			"1": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			"2": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			"3": NewNumeric(&Spec{
				Length:      2,
				Description: "Numeric Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
		},
	}

	compositeTestSpecWithTagPadding = &Spec{
		Length:      30,
		Description: "Test Spec",
		Pref:        prefix.ASCII.LL,
		Tag: &TagSpec{
			Length: 2,
			Enc:    encoding.ASCII,
			Pad:    padding.Left('0'),
			Sort:   sort.StringsByInt,
		},
		Subfields: map[string]Field{
			"1": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"2": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"3": NewNumeric(&Spec{
				Length:      2,
				Description: "Numeric Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"11": NewComposite(&Spec{
				Length:      6,
				Description: "Sub-Composite Field",
				Pref:        prefix.ASCII.LL,
				Tag: &TagSpec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pad:    padding.Left('0'),
					Sort:   sort.StringsByInt,
				},
				Subfields: map[string]Field{
					"1": NewString(&Spec{
						Length:      2,
						Description: "String Field",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.LL,
					}),
				},
			}),
		},
	}

	compositeTestSpecWithDefaultBitmap = &Spec{
		Length:      36,
		Description: "Test Spec",
		Pref:        prefix.ASCII.LL,
		Bitmap: NewBitmap(&Spec{
			Length:            8,
			Description:       "Bitmap",
			Enc:               encoding.BytesToASCIIHex,
			Pref:              prefix.Hex.Fixed,
			DisableAutoExpand: true,
		}),
		Subfields: map[string]Field{
			"1": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"2": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"3": NewNumeric(&Spec{
				Length:      2,
				Description: "Numeric Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"11": NewComposite(&Spec{
				Length:      6,
				Description: "Sub-Composite Field",
				Pref:        prefix.ASCII.LL,
				Tag: &TagSpec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pad:    padding.Left('0'),
					Sort:   sort.StringsByInt,
				},
				Subfields: map[string]Field{
					"1": NewString(&Spec{
						Length:      2,
						Description: "String Field",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.LL,
					}),
				},
			}),
		},
	}

	compositeTestSpecWithSizedBitmap = &Spec{
		Length:      30,
		Description: "Test Spec",
		Pref:        prefix.ASCII.LL,
		Bitmap: NewBitmap(&Spec{
			Length:            3,
			Description:       "Bitmap",
			Enc:               encoding.BytesToASCIIHex,
			Pref:              prefix.Hex.Fixed,
			DisableAutoExpand: true,
		}),
		Subfields: map[string]Field{
			"1": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"2": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"3": NewNumeric(&Spec{
				Length:      2,
				Description: "Numeric Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"11": NewComposite(&Spec{
				Length:      6,
				Description: "Sub-Composite Field",
				Pref:        prefix.ASCII.LL,
				Tag: &TagSpec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pad:    padding.Left('0'),
					Sort:   sort.StringsByInt,
				},
				Subfields: map[string]Field{
					"1": NewString(&Spec{
						Length:      2,
						Description: "String Field",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.LL,
					}),
				},
			}),
		},
	}

	compositeTestSpecWithoutTagPadding = &Spec{
		Length:      30,
		Description: "Test Spec",
		Pref:        prefix.ASCII.LL,
		Tag: &TagSpec{
			Length: 2,
			Enc:    encoding.ASCII,
			Sort:   sort.StringsByInt,
		},
		Subfields: map[string]Field{
			"01": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			"02": NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
		},
	}

	tlvTestSpec = &Spec{
		Length:      999,
		Description: "ICC Data – EMV Having Multiple Tags",
		Pref:        prefix.ASCII.LLL,
		Tag: &TagSpec{
			Enc:  encoding.BerTLVTag,
			Sort: sort.StringsByHex,
		},
		Subfields: map[string]Field{
			"9A": NewHex(&Spec{
				Description: "Transaction Date",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F02": NewHex(&Spec{
				Description: "Amount, Authorized (Numeric)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
		},
	}

	constructedBERTLVTestSpec = &Spec{
		Length:      999,
		Description: "ICC Data – EMV Having Multiple Tags",
		Pref:        prefix.ASCII.LLL,
		Tag: &TagSpec{
			Enc:  encoding.BerTLVTag,
			Sort: sort.StringsByHex,
		},
		Subfields: map[string]Field{
			"82": NewHex(&Spec{
				Description: "Application Interchange Profile",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F36": NewHex(&Spec{
				Description: "Currency Code, Application Reference",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F3B": NewComposite(&Spec{
				Description: "Currency Code, Application Reference",
				Pref:        prefix.BerTLV,
				Tag: &TagSpec{
					Enc:  encoding.BerTLVTag,
					Sort: sort.StringsByHex,
				},
				Subfields: map[string]Field{
					"9F45": NewHex(&Spec{
						Description: "Data Authentication Code",
						Enc:         encoding.Binary,
						Pref:        prefix.BerTLV,
					}),
				},
			}),
		},
	}
)

type CompositeTestData struct {
	F1  *String
	F2  *String
	F3  *Numeric
	F11 *SubCompositeData
}

type SubCompositeData struct {
	F1 *String
}

type CompositeTestDataWithoutTagPadding struct {
	F01 *String
	F02 *String
}

type CompositeTestDataWithoutTagPaddingWithIndexTag struct {
	F01 *String `index:"01"`
	F02 *String `index:"02"`
}

type TLVTestData struct {
	F9A   *Hex
	F9F02 *Hex
}

type ConstructedTLVTestData struct {
	F82   *Hex
	F9F36 *Hex
	F9F3B *SubConstructedTLVTestData
}

type SubConstructedTLVTestData struct {
	F9F45 *Hex
}

func TestComposite_SetData(t *testing.T) {
	t.Run("SetData returns an error on provision of primitive type data", func(t *testing.T) {
		composite := NewComposite(compositeTestSpec)
		err := composite.SetData("primitive str")
		require.EqualError(t, err, "data is not a pointer or nil")
	})
}

func TestCompositeFieldUnmarshal(t *testing.T) {
	t.Run("Unmarshal gets data for composite field", func(t *testing.T) {
		// first, we need to populate fields of composite field
		// we will do it by packing the field
		composite := NewComposite(tlvTestSpec)
		err := composite.SetData(&TLVTestData{
			F9A:   NewHexValue("210720"),
			F9F02: NewHexValue("000000000501"),
		})
		require.NoError(t, err)

		_, err = composite.Pack()
		require.NoError(t, err)

		data := &TLVTestData{}
		err = composite.Unmarshal(data)
		require.NoError(t, err)

		require.Equal(t, "210720", data.F9A.Value())
		require.Equal(t, "000000000501", data.F9F02.Value())
	})

	t.Run("Unmarshal gets data for composite field (constructed)", func(t *testing.T) {
		composite := NewComposite(constructedBERTLVTestSpec)
		err := composite.SetData(&ConstructedTLVTestData{
			F82:   NewHexValue("017F"),
			F9F36: NewHexValue("027F"),
			F9F3B: &SubConstructedTLVTestData{
				F9F45: NewHexValue("047F"),
			},
		})
		require.NoError(t, err)

		_, err = composite.Pack()
		require.NoError(t, err)

		data := &ConstructedTLVTestData{}
		err = composite.Unmarshal(data)
		require.NoError(t, err)

		require.Equal(t, "017F", data.F82.Value())
		require.Equal(t, "027F", data.F9F36.Value())
		require.Equal(t, "047F", data.F9F3B.F9F45.Value())
	})

	t.Run("Unmarshal gets data for composite field using field tag `index`", func(t *testing.T) {
		type tlvTestData struct {
			Date          *Hex `index:"9A"`
			TransactionID *Hex `index:"9F02"`
		}
		// first, we need to populate fields of composite field
		// we will do it by packing the field
		composite := NewComposite(tlvTestSpec)
		err := composite.SetData(&TLVTestData{
			F9A:   NewHexValue("210720"),
			F9F02: NewHexValue("000000000501"),
		})
		require.NoError(t, err)

		_, err = composite.Pack()
		require.NoError(t, err)

		data := &tlvTestData{}
		err = composite.Unmarshal(data)
		require.NoError(t, err)

		require.Equal(t, "210720", data.Date.Value())
		require.Equal(t, "000000000501", data.TransactionID.Value())
	})
}

func TestTLVPacking(t *testing.T) {
	t.Run("Pack correctly serializes data to bytes (general tlv)", func(t *testing.T) {
		data := &TLVTestData{
			F9A:   NewHexValue("210720"),
			F9F02: NewHexValue("000000000501"),
		}

		composite := NewComposite(tlvTestSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		// TLV Length: 0x30, 0x31, 0x34 (014)
		//
		// Tag: 0x9a (9A)
		// Length: 0x03 (3 bytes)
		// Value: 0x21, 0x07, 0x20 (210720)
		//
		// Tag: 0x9f, 0x02 (9F02)
		// Length: 0x06 (6 bytes)
		// Value: 0x00, 0x00, 0x00, 0x00, 0x05, 0x01 (000000000501)
		require.Equal(t, []byte{0x30, 0x31, 0x34, 0x9a, 0x3, 0x21, 0x7, 0x20, 0x9f, 0x2, 0x6, 0x0, 0x0, 0x0, 0x0, 0x5, 0x1}, packed)
	})

	t.Run("Unpack correctly deserialises bytes to the data struct", func(t *testing.T) {
		composite := NewComposite(tlvTestSpec)

		read, err := composite.Unpack([]byte{0x30, 0x31, 0x34, 0x9a, 0x3, 0x21, 0x7, 0x20, 0x9f, 0x2, 0x6, 0x0, 0x0, 0x0, 0x0, 0x5, 0x1})
		require.NoError(t, err)
		require.Equal(t, 17, read)

		data := &TLVTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "210720", data.F9A.Value())
		require.Equal(t, "000000000501", data.F9F02.Value())
	})

	t.Run("Unpack correctly deserialises bytes to the data struct skipping unexpected tags", func(t *testing.T) {
		// Turn on the skipping unexpected tags capability and turn it off at the end of test
		tlvTestSpec.Tag.SkipUnknownTLVTags = true
		defer func() {
			tlvTestSpec.Tag.SkipUnknownTLVTags = false
		}()

		// The field's specification contains the tags 9A and 9F02
		composite := NewComposite(tlvTestSpec)

		// This data contains tags 9A and 9F02 that are mapped in the specification, but also
		// contains tags 9F36 and 9F37 which aren't in the specification.
		read, err := composite.Unpack([]byte{0x30, 0x32, 0x36, 0x9f, 0x36, 0x2, 0x1, 0x57, 0x9a, 0x3, 0x21, 0x7, 0x20,
			0x9f, 0x2, 0x6, 0x0, 0x0, 0x0, 0x0, 0x5, 0x1, 0x9f, 0x37, 0x4, 0x9b, 0xad, 0xbc, 0xab})
		require.NoError(t, err)
		require.Equal(t, 29, read)

		data := &TLVTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "210720", data.F9A.Value())
		require.Equal(t, "000000000501", data.F9F02.Value())
	})

	t.Run("Pack correctly serializes data to bytes (constructed ber-tlv)", func(t *testing.T) {
		data := &ConstructedTLVTestData{
			F82:   NewHexValue("017f"),
			F9F36: NewHexValue("027f"),
			F9F3B: &SubConstructedTLVTestData{
				F9F45: NewHexValue("047f"),
			},
		}

		composite := NewComposite(constructedBERTLVTestSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		// TLV Length: 0x30, 0x31, 0x37 (017)
		//
		// Tag: 0x82 (82)
		// Length: 0x02 (2 bytes)
		// Value: 0x01, 0x7f (017F)
		//
		// Tag: 0x9f, 0x36 (9F36)
		// Length: 0x02 (2 bytes)
		// Value: 0x02, 0x7f (027f)
		//
		// Tag: 0x9f, 0x3b (9F3B)
		// Length: 0x05 (5 bytes)
		// Value:
		//  Tag: 0x9f, 0x45 (9F45)
		// 	Length: 0x02 (2 bytes)
		// 	Value: 0x04, 0x7f (047F)
		require.Equal(t, []byte{0x30, 0x31, 0x37, 0x82, 0x2, 0x1, 0x7f, 0x9f, 0x36, 0x2, 0x2, 0x7f, 0x9f, 0x3b, 0x5, 0x9f, 0x45, 0x2, 0x4, 0x7f}, packed)
	})

	t.Run("Unpack correctly deserialises bytes to the data struct (constructed ber-tlv)", func(t *testing.T) {
		composite := NewComposite(constructedBERTLVTestSpec)

		read, err := composite.Unpack([]byte{0x30, 0x31, 0x37, 0x82, 0x2, 0x1, 0x7f, 0x9f, 0x36, 0x2, 0x2, 0x7f, 0x9f, 0x3b, 0x5, 0x9f, 0x45, 0x2, 0x4, 0x7f})
		require.NoError(t, err)
		require.Equal(t, 20, read)

		data := &ConstructedTLVTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "017F", data.F82.Value())
		require.Equal(t, "027F", data.F9F36.Value())
		require.Equal(t, "047F", data.F9F3B.F9F45.Value())
	})

	t.Run("123Unpack correctly deserialises bytes to the data struct (constructed ber-tlv, unordered value)", func(t *testing.T) {
		composite := NewComposite(constructedBERTLVTestSpec)

		read, err := composite.Unpack([]byte{0x30, 0x31, 0x37, 0x9f, 0x36, 0x2, 0x2, 0x7f, 0x9f, 0x3b, 0x5, 0x9f, 0x45, 0x2, 0x4, 0x7f, 0x82, 0x2, 0x1, 0x7f})
		require.NoError(t, err)
		require.Equal(t, 20, read)

		data := &ConstructedTLVTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "017F", data.F82.Value())
		require.Equal(t, "027F", data.F9F36.Value())
		require.Equal(t, "047F", data.F9F3B.F9F45.Value())
	})

	t.Run("Unpack throws an error due unexpected tags", func(t *testing.T) {
		composite := NewComposite(tlvTestSpec)

		// This data contains tags 9A and 9F02 that are mapped in the specification, but also
		// contains tags 9F36 and 9F37 which aren't in the specification.
		_, err := composite.Unpack([]byte{0x30, 0x32, 0x36, 0x9f, 0x36, 0x2, 0x1, 0x57, 0x9a, 0x3, 0x21, 0x7, 0x20,
			0x9f, 0x2, 0x6, 0x0, 0x0, 0x0, 0x0, 0x5, 0x1, 0x9f, 0x37, 0x4, 0x9b, 0xad, 0xbc, 0xab})
		require.EqualError(t, err, "failed to unpack subfield 9F36: field not defined in Spec")
	})
}

func TestCompositePacking(t *testing.T) {
	t.Run("Pack returns an error on mismatch of subfield types", func(t *testing.T) {
		type TestDataIncorrectType struct {
			F1 *Numeric
		}
		composite := NewComposite(compositeTestSpec)
		err := composite.SetData(&TestDataIncorrectType{
			F1: NewNumericValue(1),
		})

		require.Error(t, err)
		require.EqualError(t, err, "failed to set data from field 1: data does not match required *String type")
	})

	t.Run("Pack returns error on failure of subfield packing", func(t *testing.T) {
		data := &CompositeTestData{
			// This subfield will return an error on F1.Pack() as its length
			// exceeds the max length defined in the spec.
			F1: NewStringValue("ABCD"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(compositeTestSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		_, err = composite.Pack()
		require.EqualError(t, err, "failed to pack subfield 1: failed to encode length: field length: 4 should be fixed: 2")
	})

	t.Run("Pack returns error when encoded data length is larger than specified fixed max length", func(t *testing.T) {
		invalidSpec := &Spec{
			// Base field length < summation of lengths of subfields
			// This will throw an error when encoding the field's length.
			Length: 4,
			Pref:   prefix.ASCII.Fixed,
			Tag: &TagSpec{
				Sort: sort.StringsByInt,
			},
			Subfields: map[string]Field{
				"1": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"2": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"3": NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(invalidSpec)
		err := composite.Marshal(data)
		require.NoError(t, err)

		_, err = composite.Pack()
		require.EqualError(t, err, "failed to encode length: field length: 6 should be fixed: 4")
	})

	t.Run("Pack correctly serializes data with padded tags to bytes", func(t *testing.T) {
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(compositeTestSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, "ABCD12", string(packed))
	})

	t.Run("Pack and unpack data with BCD encoding", func(t *testing.T) {
		var compositeSpecWithBCD = &Spec{
			Length:      2, // always in bytes
			Description: "Point of Service Entry Mode",
			Pref:        prefix.BCD.Fixed,
			Tag: &TagSpec{
				Sort: sort.StringsByInt,
			},
			Subfields: map[string]Field{
				"1": NewString(&Spec{
					Length:      2,
					Description: "PAN/Date Entry Mode",
					Enc:         encoding.BCD,
					Pref:        prefix.BCD.Fixed,
				}),
				"2": NewString(&Spec{
					Length:      2,
					Description: "PIN Entry Capability",
					Enc:         encoding.BCD,
					Pref:        prefix.BCD.Fixed,
				}),
			},
		}

		type data struct {
			PANEntryMode *String `index:"1"`
			PINEntryMode *String `index:"2"`
		}

		f := NewComposite(compositeSpecWithBCD)

		d := &data{
			PANEntryMode: NewStringValue("01"),
			PINEntryMode: NewStringValue("02"),
		}

		err := f.Marshal(d)
		require.NoError(t, err)

		packed, err := f.Pack()
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x02}, packed)

		// unpacking

		f = NewComposite(compositeSpecWithBCD)
		read, err := f.Unpack(packed)
		require.NoError(t, err)
		require.Equal(t, 2, read) // two bytes read

		d = &data{}
		err = f.Unmarshal(d)
		require.NoError(t, err)

		require.Equal(t, "01", d.PANEntryMode.Value())
		require.Equal(t, "02", d.PINEntryMode.Value())
	})

	t.Run("Unpack returns an error on mismatch of subfield types", func(t *testing.T) {
		type TestDataIncorrectType struct {
			F1 *Numeric
		}
		composite := NewComposite(compositeTestSpec)

		_, err := composite.Unpack([]byte("ABCD12"))
		require.NoError(t, err)

		data := &TestDataIncorrectType{}
		err = composite.Unmarshal(data)

		require.Error(t, err)
		require.EqualError(t, err, "failed to get data from field 1: data does not match required *String type")
	})

	t.Run("Unpack returns an error on failure of subfield to unpack bytes", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Last two characters must be an integer type. F3 fails to unpack.
		read, err := composite.Unpack([]byte("ABCDEF"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to unpack subfield 3: failed to set bytes: failed to convert into number")
		require.ErrorIs(t, err, strconv.ErrSyntax)
	})

	t.Run("Unpack returns an error on length of data exceeding max length", func(t *testing.T) {
		spec := &Spec{
			Length: 4,
			Pref:   prefix.ASCII.L,
			Tag: &TagSpec{
				Sort: sort.StringsByInt,
			},
			Subfields: map[string]Field{
				"1": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"2": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"3": NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompositeTestData{}

		composite := NewComposite(spec)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Length of denoted by prefix is too long, causing failure to decode length.
		read, err := composite.Unpack([]byte("7ABCD123"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to decode length: data length: 7 is larger than maximum 4")
	})

	t.Run("Unpack without error when not all subfields are set", func(t *testing.T) {
		spec := &Spec{
			Length: 4,
			Pref:   prefix.ASCII.L,
			Tag: &TagSpec{
				Sort: sort.StringsByInt,
			},
			Subfields: map[string]Field{
				"1": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"2": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"3": NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompositeTestData{}

		composite := NewComposite(spec)
		err := composite.SetData(data)
		require.NoError(t, err)

		// There is data only for first subfield
		read, err := composite.Unpack([]byte("2AB"))
		require.Equal(t, 3, read)
		require.NoError(t, err)
	})

	t.Run("Unpack returns an error on offset not matching data length", func(t *testing.T) {
		invalidSpec := &Spec{
			// Base field length < summation of lengths of subfields
			// This will throw an error when encoding the field's length.
			Length: 4,
			Pref:   prefix.ASCII.Fixed,
			Tag: &TagSpec{
				Sort: sort.StringsByInt,
			},
			Subfields: map[string]Field{
				"1": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"2": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"3": NewNumeric(&Spec{
					Length: 3,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}

		composite := NewComposite(invalidSpec)

		// Length of input too long, causing failure to decode length.
		read, err := composite.Unpack([]byte("ABCD123"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to unpack subfield 3: failed to decode content: not enough data to decode. expected len 3, got 0")
	})

	t.Run("Unpack correctly deserialises bytes to the data struct", func(t *testing.T) {
		composite := NewComposite(compositeTestSpec)

		read, err := composite.Unpack([]byte("ABCD12"))
		require.Equal(t, compositeTestSpec.Length, read)
		require.NoError(t, err)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, "CD", data.F2.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Nil(t, data.F11)
	})

	t.Run("SetBytes correctly deserialises bytes to the data struct", func(t *testing.T) {
		composite := NewComposite(compositeTestSpec)

		err := composite.SetBytes([]byte("ABCD12"))
		require.NoError(t, err)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, "CD", data.F2.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Nil(t, data.F11)
	})
}

func TestCompositePackingWithTags(t *testing.T) {
	t.Run("Pack returns error when encoded data length is larger than specified fixed max length", func(t *testing.T) {
		// Base field length < summation of (lengths of subfields + IDs).
		// This will throw an error when encoding the field's length.
		invalidSpec := &Spec{
			Length: 6,
			Pref:   prefix.ASCII.Fixed,
			Tag: &TagSpec{
				Length: 2,
				Enc:    encoding.ASCII,
				Pad:    padding.Left('0'),
				Sort:   sort.StringsByInt,
			},
			Subfields: map[string]Field{
				"1": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"2": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"3": NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(invalidSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		b, err := composite.Pack()
		require.Nil(t, b)
		require.Error(t, err)
		require.EqualError(t, err, "failed to encode length: field length: 12 should be fixed: 6")
	})

	t.Run("Pack returns error when encoded data length is larger than specified variable max length", func(t *testing.T) {
		invalidSpec := &Spec{
			Length: 8,
			Pref:   prefix.ASCII.LL,
			Tag: &TagSpec{
				Length: 2,
				Enc:    encoding.ASCII,
				Pad:    padding.Left('0'),
				Sort:   sort.StringsByInt,
			},
			Subfields: map[string]Field{
				"1": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"2": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"3": NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(invalidSpec)
		err := composite.Marshal(data)
		require.NoError(t, err)

		b, err := composite.Pack()
		require.Nil(t, b)
		require.EqualError(t, err, "failed to encode length: field length: 12 is larger than maximum: 8")
	})

	t.Run("Pack correctly serializes fully populated data to bytes", func(t *testing.T) {
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
			F11: &SubCompositeData{
				F1: NewStringValue("YZ"),
			},
		}

		composite := NewComposite(compositeTestSpecWithTagPadding)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		require.Equal(t, "280102AB0202CD03021211060102YZ", string(packed))
	})

	t.Run("Pack correctly serializes partially populated data to bytes", func(t *testing.T) {
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(compositeTestSpecWithTagPadding)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, "120102AB030212", string(packed))
	})

	t.Run("Pack correctly serializes fully populated unpadded tag data to bytes", func(t *testing.T) {
		data := &CompositeTestDataWithoutTagPadding{
			F01: NewStringValue("AB"),
			F02: NewStringValue("CD"),
		}

		composite := NewComposite(compositeTestSpecWithoutTagPadding)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		require.Equal(t, "120102AB0202CD", string(packed))
	})

	t.Run("Unpack returns an error on failure of subfield to unpack bytes", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpecWithTagPadding)
		err := composite.SetData(data)
		require.NoError(t, err)

		// F3 fails to unpack - it requires len to be defined instead of AB.
		read, err := composite.Unpack([]byte("180102AB0202CD03AB12"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to unpack subfield 3: failed to decode length: strconv.Atoi: parsing \"AB\": invalid syntax")
	})

	t.Run("Unpack returns an error on data having subfield ID not in spec", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpecWithTagPadding)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Index 2-3 should have '01' rather than '12'.
		read, err := composite.Unpack([]byte("181202AB0202CD030212"))
		require.Equal(t, 0, read)
		require.EqualError(t, err, "failed to unpack subfield 12: field not defined in Spec")
	})

	t.Run("Unpack returns an error on if subfield not defined in spec", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpecWithTagPadding)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Index 0, 1 should have '01' rather than 'ID'.
		read, err := composite.Unpack([]byte("18ID02AB0202CD030212"))
		require.Equal(t, 0, read)
		require.EqualError(t, err, "failed to unpack subfield ID: field not defined in Spec")
	})

	t.Run("Unpack correctly deserialises out of order composite subfields to the data struct", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithTagPadding)

		read, err := composite.Unpack([]byte("280202CD0302120102AB11060102YZ"))

		require.NoError(t, err)
		require.Equal(t, 30, read)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, "CD", data.F2.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Equal(t, "YZ", data.F11.F1.Value())
	})

	t.Run("Unpack correctly deserialises out of order composite subfields to the unpadded data struct", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithoutTagPadding)

		read, err := composite.Unpack([]byte("120202CD0102AB"))

		data := &CompositeTestDataWithoutTagPadding{}

		require.NoError(t, composite.Unmarshal(data))

		require.NoError(t, err)
		require.Equal(t, 14, read)

		require.Equal(t, "AB", data.F01.Value())
		require.Equal(t, "CD", data.F02.Value())
	})

	t.Run("Unpack correctly deserialises partial subfields to the data struct", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithTagPadding)

		read, err := composite.Unpack([]byte("120302120102AB"))

		require.NoError(t, err)
		require.Equal(t, 14, read)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Nil(t, data.F2)
		require.Equal(t, 12, data.F3.Value())
	})

	t.Run("Unpack correctly ignores excess bytes in excess of the length described by the prefix", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithTagPadding)

		// "11060102YZ" falls outside of the bounds of the 18 byte limit imposed
		// by the prefix. Therefore, F11 must be nil.
		read, err := composite.Unpack([]byte("180202CD0302120102AB11060102YZ"))

		require.NoError(t, err)
		require.Equal(t, 20, read)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, "CD", data.F2.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Nil(t, data.F11)
	})
}

func TestCompositePackingWithBitmap(t *testing.T) {
	t.Run("Pack returns error when encoded data length is different than specified fixed length", func(t *testing.T) {
		// Base field length < sum of subfields lengths.
		// This will throw an error when encoding the field's length.
		invalidSpec := &Spec{
			Length: 20,
			Pref:   prefix.ASCII.Fixed,
			Bitmap: NewBitmap(&Spec{
				Length:            8,
				Pref:              prefix.Binary.Fixed,
				Enc:               encoding.Binary,
				DisableAutoExpand: true,
			}),
			Subfields: map[string]Field{
				"1": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"2": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"3": NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(invalidSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		b, err := composite.Pack()
		require.Nil(t, b)
		require.Error(t, err)
		require.EqualError(t, err, "failed to encode length: field length: 14 should be fixed: 20")
	})

	t.Run("Pack returns error when encoded data length is larger than specified variable max length", func(t *testing.T) {
		invalidSpec := &Spec{
			Length: 5,
			Pref:   prefix.ASCII.LL,
			Bitmap: NewBitmap(&Spec{
				Length:            8,
				Pref:              prefix.Binary.Fixed,
				Enc:               encoding.Binary,
				DisableAutoExpand: true,
			}),
			Subfields: map[string]Field{
				"1": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"2": NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				"3": NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(invalidSpec)
		err := composite.Marshal(data)
		require.NoError(t, err)

		b, err := composite.Pack()
		require.Nil(t, b)
		require.EqualError(t, err, "failed to encode length: field length: 14 is larger than maximum: 5")
	})

	t.Run("Pack correctly serializes fully populated data to bytes with default bitmap", func(t *testing.T) {
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
			F11: &SubCompositeData{
				F1: NewStringValue("YZ"),
			},
		}

		composite := NewComposite(compositeTestSpecWithDefaultBitmap)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		require.Equal(t, "36E02000000000000002AB02CD0212060102YZ", string(packed))
	})

	t.Run("Pack correctly serializes partially populated data to bytes with default bitmap", func(t *testing.T) {
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(compositeTestSpecWithDefaultBitmap)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, "24A00000000000000002AB0212", string(packed))
	})

	t.Run("Pack correctly serializes fully populated data to bytes with sized bitmap on 3 bytes", func(t *testing.T) {
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
			F11: &SubCompositeData{
				F1: NewStringValue("YZ"),
			},
		}

		composite := NewComposite(compositeTestSpecWithSizedBitmap)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		require.Equal(t, "26E0200002AB02CD0212060102YZ", string(packed))
	})

	t.Run("Pack correctly serializes partially populated data to bytes with sized bitmap on 3 bytes", func(t *testing.T) {
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(compositeTestSpecWithSizedBitmap)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed, err := composite.Pack()
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, "14A0000002AB0212", string(packed))
	})

	t.Run("Unpack returns an error on failure of subfield to unpack bytes with default bitmap", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpecWithDefaultBitmap)
		err := composite.SetData(data)
		require.NoError(t, err)

		// F1 fails to unpack - it requires length to be defined instead of AB.
		read, err := composite.Unpack([]byte("30E020000000000000AB02AB060102YZ"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to unpack subfield 1 (String Field): failed to decode length: strconv.Atoi: parsing \"AB\": invalid syntax")
	})

	t.Run("Unpack returns an error on data having subfield ID not in spec with default bitmap", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpecWithDefaultBitmap)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Index 2-3 = 70 indicates the presence of field 4. This field is not defined on spec.
		read, err := composite.Unpack([]byte("32702000000000000002AB0212060102YZ"))
		require.Equal(t, 0, read)
		require.EqualError(t, err, "failed to unpack subfield 4: no specification found")
	})

	t.Run("Unpack correctly deserialises out of order composite subfields to the data struct with default bitmap", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithDefaultBitmap)

		read, err := composite.Unpack([]byte("36E02000000000000002AB02CD0212060102YZ"))

		require.NoError(t, err)
		require.Equal(t, 38, read)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, "CD", data.F2.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Equal(t, "YZ", data.F11.F1.Value())
	})

	t.Run("Unpack correctly deserialises partial subfields to the data struct with sized bitmap on 3 bytes", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithDefaultBitmap)

		read, err := composite.Unpack([]byte("24A00000000000000002AB0212"))

		require.NoError(t, err)
		require.Equal(t, 26, read)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Nil(t, data.F11)
	})

	t.Run("Unpack correctly ignores excess bytes in excess of the length described by the prefix with default bitmap", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithDefaultBitmap)

		// "060102YZ" falls outside of the bounds of the 24 byte limit imposed
		// by the prefix. Therefore, F11 must be nil.
		read, err := composite.Unpack([]byte("28E00000000000000002AB02CD0212060102YZ"))

		require.NoError(t, err)
		require.Equal(t, 30, read)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, "CD", data.F2.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Nil(t, data.F11)
	})

	t.Run("Unpack returns an error on failure of subfield to unpack bytes with sized bitmap on 3 bytes", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpecWithSizedBitmap)
		err := composite.SetData(data)
		require.NoError(t, err)

		// F1 fails to unpack - it requires length to be defined instead of AB.
		read, err := composite.Unpack([]byte("20E02000AB02CD060102YZ"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to unpack subfield 1 (String Field): failed to decode length: strconv.Atoi: parsing \"AB\": invalid syntax")
	})

	t.Run("Unpack returns an error on data having subfield ID not in spec with sized bitmap on 3 bytes", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpecWithSizedBitmap)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Index 2-3 = 70 indicates the presence of field 4. This field is not defined on spec.
		read, err := composite.Unpack([]byte("2270200002CD0212060102YZ"))
		require.Equal(t, 0, read)
		require.EqualError(t, err, "failed to unpack subfield 4: no specification found")
	})

	t.Run("Unpack correctly deserialises out of order composite subfields to the data struct with sized bitmap on 3 bytes", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithSizedBitmap)

		read, err := composite.Unpack([]byte("26E0200002AB02CD0212060102YZ"))

		require.NoError(t, err)
		require.Equal(t, 28, read)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, "CD", data.F2.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Equal(t, "YZ", data.F11.F1.Value())
	})

	t.Run("Unpack correctly deserialises partial subfields to the data struct with sized bitmap on 3 bytes", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithSizedBitmap)

		read, err := composite.Unpack([]byte("14A0000002AB0212"))

		require.NoError(t, err)
		require.Equal(t, 16, read)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Nil(t, data.F11)
	})

	t.Run("Unpack correctly ignores excess bytes in excess of the length described by the prefix with sized bitmap on 3 bytes", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithSizedBitmap)

		// "60102YZ" falls outside of the bounds of the 24 byte limit imposed
		// by the prefix. Therefore, F11 must be nil.
		read, err := composite.Unpack([]byte("18E0000002AB02CD0212060102YZ"))

		require.NoError(t, err)
		require.Equal(t, 20, read)

		data := &CompositeTestData{}
		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Equal(t, "CD", data.F2.Value())
		require.Equal(t, 12, data.F3.Value())
		require.Nil(t, data.F11)
	})
}

func TestCompositeHandlesValidSpecs(t *testing.T) {
	tests := []struct {
		desc string
		spec *Spec
	}{
		{
			desc: "accepts nil Enc value",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Tag: &TagSpec{
					Sort: sort.StringsByInt,
				},
				Subfields: map[string]Field{},
			},
		},
		{
			desc: "accepts nil Pad value",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Tag: &TagSpec{
					Sort: sort.StringsByInt,
				},
				Subfields: map[string]Field{},
			},
		},
		{
			desc: "accepts None Pad value",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Pad:    padding.None,
				Tag: &TagSpec{
					Sort: sort.StringsByInt,
				},
				Subfields: map[string]Field{},
			},
		},
		{
			desc: "accepts Bitmap on spec and no tag",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Bitmap: NewBitmap(&Spec{
					Length:            8,
					Pref:              prefix.Binary.Fixed,
					Enc:               encoding.Binary,
					DisableAutoExpand: true,
				}),
				Subfields: map[string]Field{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("NewComposite() %v", tc.desc), func(t *testing.T) {
			f := NewComposite(tc.spec)
			require.Equal(t, tc.spec, f.Spec())
		})
		t.Run(fmt.Sprintf("Composite.SetSpec() %v", tc.desc), func(t *testing.T) {
			f := &Composite{}
			f.SetSpec(tc.spec)
			require.Equal(t, tc.spec, f.Spec())
		})
	}
}

func TestCompositePanicsOnSpecValidationFailures(t *testing.T) {
	tests := []struct {
		desc string
		err  string
		spec *Spec
	}{
		{
			desc: "panics on non-nil Enc value being defined in spec",
			err:  "Composite spec only supports a nil Enc value",
			spec: &Spec{
				Length:    6,
				Enc:       encoding.ASCII,
				Pref:      prefix.ASCII.Fixed,
				Subfields: map[string]Field{},
				Tag: &TagSpec{
					Sort: sort.StringsByInt,
				},
			},
		},
		{
			desc: "panics on non-None / non-nil Pad value being defined in spec",
			err:  "Composite spec only supports nil or None spec padding values",
			spec: &Spec{
				Length:    6,
				Pref:      prefix.ASCII.Fixed,
				Pad:       padding.Left('0'),
				Subfields: map[string]Field{},
				Tag: &TagSpec{
					Sort: sort.StringsByInt,
				},
			},
		},
		{
			desc: "panics on no Tag and no Bitmap being defined in spec",
			err:  "Composite spec only supports a definition of Bitmap or Tag, can't stand both or neither",
			spec: &Spec{
				Length:    6,
				Pref:      prefix.ASCII.Fixed,
				Subfields: map[string]Field{},
			},
		},
		{
			desc: "panics on both Tag and Bitmap being defined in spec",
			err:  "Composite spec only supports a definition of Bitmap or Tag, can't stand both or neither",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Tag: &TagSpec{
					Sort: sort.StringsByInt,
				},
				Bitmap: NewBitmap(&Spec{
					Length:            8,
					Pref:              prefix.Binary.Fixed,
					Enc:               encoding.Binary,
					DisableAutoExpand: true,
				}),
				Subfields: map[string]Field{},
			},
		},
		{
			desc: "panics on invalid int defined as a subfield key on a bitmapped composite definition",
			err:  "error parsing key from bitmapped subfield definition: strconv.Atoi: parsing \"invalid\": invalid syntax",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Bitmap: NewBitmap(&Spec{
					Length:            8,
					Pref:              prefix.Binary.Fixed,
					Enc:               encoding.Binary,
					DisableAutoExpand: true,
				}),
				Subfields: map[string]Field{
					"invalid": NewString(&Spec{
						Length:            1,
						Pref:              prefix.ASCII.Fixed,
						Enc:               encoding.ASCII,
						DisableAutoExpand: true,
					}),
				},
			},
		},
		{
			desc: "panics on an int lower than 1 defined as a subfield key on a bitmapped composite definition",
			err:  "Composite spec only supports integers greater than 0 as keys for bitmapped subfields definition",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Bitmap: NewBitmap(&Spec{
					Length:            8,
					Pref:              prefix.Binary.Fixed,
					Enc:               encoding.Binary,
					DisableAutoExpand: true,
				}),
				Subfields: map[string]Field{
					"0": NewString(&Spec{
						Length: 1,
						Pref:   prefix.ASCII.Fixed,
						Enc:    encoding.ASCII,
					}),
				},
			},
		},
		{
			desc: "panics on a bitmap with DisableAutoExpand = false",
			err:  "Composite spec only supports a bitmap with 'DisableAutoExpand = true'",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Bitmap: NewBitmap(&Spec{
					Length:            8,
					Pref:              prefix.Binary.Fixed,
					Enc:               encoding.Binary,
					DisableAutoExpand: false,
				}),
				Subfields: map[string]Field{},
			},
		},
		{
			desc: "panics on nil Tag.Sort",
			err:  "Composite spec requires a Tag.Sort function to define a Tag",
			spec: &Spec{
				Length:    6,
				Pref:      prefix.ASCII.Fixed,
				Subfields: map[string]Field{},
				Tag:       &TagSpec{},
			},
		},
		{
			desc: "panics on nil Enc value being defined in spec if Tag.Length > 0",
			err:  "Composite spec requires a Tag.Enc to be defined if Tag.Length > 0",
			spec: &Spec{
				Length:    6,
				Pref:      prefix.ASCII.Fixed,
				Subfields: map[string]Field{},
				Tag: &TagSpec{
					Length: 2,
					Pad:    padding.Left('0'),
					Sort:   sort.StringsByInt,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("NewComposite() %v", tc.desc), func(t *testing.T) {
			require.PanicsWithError(t, tc.err, func() {
				NewComposite(tc.spec)
			})
		})
		t.Run(fmt.Sprintf("Composite.SetSpec() %v", tc.desc), func(t *testing.T) {
			require.PanicsWithError(t, tc.err, func() {
				(&Composite{}).SetSpec(tc.spec)
			})
		})
	}
}

func TestCompositeJSONConversion(t *testing.T) {
	json := `{"1":"AB","3":12,"11":{"1":"YZ"}}`

	t.Run("MarshalJSON typed", func(t *testing.T) {
		data := &CompositeTestData{
			F1: NewStringValue("AB"),
			F3: NewNumericValue(12),
			F11: &SubCompositeData{
				F1: NewStringValue("YZ"),
			},
		}

		composite := NewComposite(compositeTestSpecWithTagPadding)
		require.NoError(t, composite.SetData(data))

		actual, err := composite.MarshalJSON()
		require.NoError(t, err)

		require.JSONEq(t, json, string(actual))
	})

	t.Run("UnmarshalJSON typed", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpecWithTagPadding)
		err := composite.SetData(data)
		require.NoError(t, err)

		require.NoError(t, composite.UnmarshalJSON([]byte(json)))

		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "AB", data.F1.Value())
		require.Nil(t, data.F2)
		require.Equal(t, 12, data.F3.Value())
		require.Equal(t, "YZ", data.F11.F1.Value())
	})

	t.Run("MarshalJSON untyped", func(t *testing.T) {
		composite := NewComposite(compositeTestSpecWithTagPadding)
		require.NoError(t, composite.SetBytes([]byte("0102AB03021211060102YZ")))

		actual, err := composite.MarshalJSON()
		require.NoError(t, err)

		require.JSONEq(t, json, string(actual))
	})

	t.Run("UnmarshalJSON untyped", func(t *testing.T) {
		data := &CompositeTestData{}

		composite := NewComposite(compositeTestSpecWithTagPadding)
		require.NoError(t, composite.SetData(data))

		require.NoError(t, composite.UnmarshalJSON([]byte(json)))

		s, err := composite.String()
		require.NoError(t, err)
		require.Equal(t, "0102AB03021211060102YZ", s)
	})
}

func TestTLVJSONConversion(t *testing.T) {
	json := `{"9A":"210720","9F02":"000000000501"}`

	t.Run("MarshalJSON TLV Data Ok", func(t *testing.T) {
		data := &TLVTestData{
			F9A:   NewHexValue("210720"),
			F9F02: NewHexValue("000000000501"),
		}

		composite := NewComposite(tlvTestSpec)
		require.NoError(t, composite.Marshal(data))

		actual, err := composite.MarshalJSON()
		require.NoError(t, err)

		require.JSONEq(t, json, string(actual))
	})

	t.Run("UnmarshalJSON TLV data skipping unexpected tags", func(t *testing.T) {
		// Turn on the skipping unexpected tags capability and turn it off at the end of test
		tlvTestSpec.Tag.SkipUnknownTLVTags = true
		defer func() {
			tlvTestSpec.Tag.SkipUnknownTLVTags = false
		}()

		// This data contains tags 9A and 9F02 that are mapped in the specification, but also
		// contains tag 9F37 which isn't in the specification.
		json_tags := `{"9A":"210720","9F02":"000000000501", "9F37": "9badbcab"}`

		data := &TLVTestData{}

		composite := NewComposite(tlvTestSpec)
		err := composite.Marshal(data)
		require.NoError(t, err)

		require.NoError(t, composite.UnmarshalJSON([]byte(json_tags)))

		require.NoError(t, composite.Unmarshal(data))

		require.Equal(t, "210720", data.F9A.Value())
		require.Equal(t, "000000000501", data.F9F02.Value())
	})

	t.Run("UnmarshalJSON TLV data throws an error due unexpected tags", func(t *testing.T) {
		// This data contains tags 9A and 9F02 that are mapped in the specification, but also
		// contains tag 9F37 which isn't in the specification.
		json_tags := `{"9A":"210720","9F02":"000000000501", "9F37": "9badbcab"}`

		data := &TLVTestData{}

		composite := NewComposite(tlvTestSpec)
		err := composite.Marshal(data)
		require.NoError(t, err)

		err = composite.UnmarshalJSON([]byte(json_tags))
		require.Error(t, err)
		require.EqualError(t, err, "failed to unmarshal subfield 9F37: received subfield not defined in spec")
	})
}

func TestComposite_getFieldIndexOrTag(t *testing.T) {
	t.Run("returns index from field name", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			F1 string
		}{}).Elem()

		index, err := getFieldIndexOrTag(st.Type().Field(0))

		require.NoError(t, err)
		require.Equal(t, "1", index)
	})

	t.Run("returns index from field tag instead of field name when both match", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			F1 string `index:"AB"`
		}{}).Elem()

		index, err := getFieldIndexOrTag(st.Type().Field(0))

		require.NoError(t, err)
		require.Equal(t, "AB", index)
	})

	t.Run("returns index from field tag", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			Name string `index:"abcd"`
			F    string `index:"02"`
		}{}).Elem()

		// get index from field Name
		index, err := getFieldIndexOrTag(st.Type().Field(0))

		require.NoError(t, err)
		require.Equal(t, "abcd", index)

		// get index from field F
		index, err = getFieldIndexOrTag(st.Type().Field(1))

		require.NoError(t, err)
		require.Equal(t, "02", index)
	})

	t.Run("returns empty string when no tag and field name does not match the pattern", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			Name string
		}{}).Elem()

		index, err := getFieldIndexOrTag(st.Type().Field(0))

		require.NoError(t, err)
		require.Empty(t, index)

		// single letter field without tag is ignored
		st = reflect.ValueOf(&struct {
			F string
		}{}).Elem()

		index, err = getFieldIndexOrTag(st.Type().Field(0))

		require.NoError(t, err)
		require.Empty(t, index)
	})
}
