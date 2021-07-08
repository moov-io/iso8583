package field

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

var (
	compositeTestSpec = &Spec{
		Length:      6,
		Description: "Test Spec",
		Pref:        prefix.ASCII.Fixed,
		Pad:         padding.None,
		Fields: map[int]Field{
			1: NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			2: NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			3: NewNumeric(&Spec{
				Length:      2,
				Description: "Numeric Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
		},
	}
	compositeTestSpecWithIDLength = &Spec{
		Length:      30,
		Description: "Test Spec",
		IDLength:    2,
		Pref:        prefix.ASCII.LL,
		Enc:         encoding.ASCII,
		Pad:         padding.None,
		Fields: map[int]Field{
			1: NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			2: NewString(&Spec{
				Length:      2,
				Description: "String Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			3: NewNumeric(&Spec{
				Length:      2,
				Description: "Numeric Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			4: NewComposite(&Spec{
				Length:      6,
				IDLength:    2,
				Description: "Sub-Composite Field",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
				Fields: map[int]Field{
					1: NewString(&Spec{
						Length:      2,
						Description: "String Field",
						Enc:         encoding.ASCII,
						Pref:        prefix.ASCII.LL,
					}),
				},
			}),
		},
	}
)

type CompsiteTestData struct {
	F1 *String
	F2 *String
	F3 *Numeric
	F4 *SubCompositeData
}

type SubCompositeData struct {
	F1 *String
}

func TestComposite_SetData(t *testing.T) {
	t.Run("SetData returns an error on provision of primitive type data", func(t *testing.T) {
		composite := NewComposite(compositeTestSpec)
		err := composite.SetData("primitive str")
		require.EqualError(t, err, "failed to set data as struct is expected, got: string")
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
		require.NoError(t, err)

		buf := bytes.NewBuffer([]byte{})
		_, err = composite.WriteTo(buf)
		require.EqualError(t, err, "failed to set data for field 1: data does not match required *String type")
	})

	t.Run("Pack returns error on failure of subfield packing", func(t *testing.T) {
		data := &CompsiteTestData{
			// This subfield will return an error on F1.Pack() as its length
			// exceeds the max length defined in the spec.
			F1: NewStringValue("ABCD"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(compositeTestSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		buf := bytes.NewBuffer([]byte{})
		_, err = composite.WriteTo(buf)
		require.EqualError(t, err, "failed to pack subfield 1: failed to encode length: field length: 4 should be fixed: 2")
	})

	t.Run("Pack returns error on failure to encode length", func(t *testing.T) {
		invalidSpec := &Spec{
			// Base field length < summation of lengths of subfields
			// This will throw an error when encoding the field's length.
			Length: 4,
			Pref:   prefix.ASCII.Fixed,
			Fields: map[int]Field{
				1: NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				2: NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				3: NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompsiteTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(invalidSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		buf := bytes.NewBuffer([]byte{})
		_, err = composite.WriteTo(buf)
		require.EqualError(t, err, "failed to encode length: field length: 6 should be fixed: 4")
	})

	t.Run("Pack correctly serializes data to bytes", func(t *testing.T) {
		data := &CompsiteTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(compositeTestSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed := bytes.NewBuffer([]byte{})
		_, err = composite.WriteTo(packed)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, "ABCD12", packed.String())
	})

	t.Run("ReadFrom returns an error on mismatch of subfield types", func(t *testing.T) {
		type TestDataIncorrectType struct {
			F1 *Numeric
		}
		composite := NewComposite(compositeTestSpec)
		err := composite.SetData(&TestDataIncorrectType{})
		require.NoError(t, err)

		read, err := composite.ReadFrom(strings.NewReader("ABCD12"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to set data for field 1: data does not match required *String type")
	})

	t.Run("ReadFrom returns an error on failure of subfield to unpack bytes", func(t *testing.T) {
		data := &CompsiteTestData{}

		composite := NewComposite(compositeTestSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Last two characters must be an integer type. F3 fails to unpack.
		read, err := composite.ReadFrom(strings.NewReader("ABCDEF"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to unpack subfield 3: failed to convert into number: strconv.Atoi: parsing \"EF\": invalid syntax")
	})

	t.Run("ReadFrom returns an error on length of data exceeding max length", func(t *testing.T) {
		spec := &Spec{
			Length: 4,
			Pref:   prefix.ASCII.L,
			Fields: map[int]Field{
				1: NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				2: NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				3: NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompsiteTestData{}

		composite := NewComposite(spec)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Length of denoted by prefix is too long, causing failure to decode length.
		read, err := composite.ReadFrom(strings.NewReader("7ABCD123"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to decode length: data length: 7 is larger than maximum 4")
	})

	t.Run("ReadFrom returns an error on offset not matching data length", func(t *testing.T) {
		invalidSpec := &Spec{
			// Base field length < summation of lengths of subfields
			// This will throw an error when encoding the field's length.
			Length: 4,
			Pref:   prefix.ASCII.Fixed,
			Fields: map[int]Field{
				1: NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				2: NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				3: NewNumeric(&Spec{
					Length: 3,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompsiteTestData{}

		composite := NewComposite(invalidSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Length of input too long, causing failure to decode length.
		read, err := composite.ReadFrom(strings.NewReader("ABCD123"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "data length: 4 does not match aggregate data read from decoded subfields: 7")
	})

	t.Run("ReadFrom correctly deserialises bytes to the data struct", func(t *testing.T) {
		data := &CompsiteTestData{}

		composite := NewComposite(compositeTestSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		read, err := composite.ReadFrom(strings.NewReader("ABCD12"))
		require.Equal(t, compositeTestSpec.Length, read)
		require.NoError(t, err)

		require.Equal(t, "AB", data.F1.Value)
		require.Equal(t, "CD", data.F2.Value)
		require.Equal(t, 12, data.F3.Value)
		require.Nil(t, data.F4)
	})
}

func TestCompositePackingWithID(t *testing.T) {
	t.Run("Pack returns error on failure to encode length", func(t *testing.T) {
		// Base field length < summation of (lengths of subfields + IDs).
		// This will throw an error when encoding the field's length.
		invalidSpec := &Spec{
			Length:   6,
			IDLength: 2,
			Pref:     prefix.ASCII.Fixed,
			Enc:      encoding.ASCII,
			Fields: map[int]Field{
				1: NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				2: NewString(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
				3: NewNumeric(&Spec{
					Length: 2,
					Enc:    encoding.ASCII,
					Pref:   prefix.ASCII.Fixed,
				}),
			},
		}
		data := &CompsiteTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(invalidSpec)
		err := composite.SetData(data)
		require.NoError(t, err)

		buf := bytes.NewBuffer([]byte{})
		_, err = composite.WriteTo(buf)
		require.Error(t, err)
		require.EqualError(t, err, "failed to encode length: field length: 12 should be fixed: 6")
	})

	t.Run("Pack correctly serializes fully populated data to bytes", func(t *testing.T) {
		data := &CompsiteTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
			F4: &SubCompositeData{
				F1: NewStringValue("YZ"),
			},
		}

		composite := NewComposite(compositeTestSpecWithIDLength)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed := bytes.NewBuffer([]byte{})
		_, err = composite.WriteTo(packed)

		require.NoError(t, err)
		require.Equal(t, "280102AB0202CD03021204060102YZ", packed.String())
	})

	t.Run("Pack correctly serializes partially populated data to bytes", func(t *testing.T) {
		data := &CompsiteTestData{
			F1: NewStringValue("AB"),
			F3: NewNumericValue(12),
		}

		composite := NewComposite(compositeTestSpecWithIDLength)
		err := composite.SetData(data)
		require.NoError(t, err)

		packed := bytes.NewBuffer([]byte{})
		_, err = composite.WriteTo(packed)

		require.NoError(t, err)
		require.Equal(t, "120102AB030212", packed.String())
	})

	t.Run("ReadFrom returns an error on failure of subfield to unpack bytes", func(t *testing.T) {
		data := &CompsiteTestData{}

		composite := NewComposite(compositeTestSpecWithIDLength)
		err := composite.SetData(data)
		require.NoError(t, err)

		// F3 fails to unpack - it requires len to be defined instead of AB.
		read, err := composite.ReadFrom(strings.NewReader("180102AB0202CD03AB12"))
		require.Equal(t, 0, read)
		require.Error(t, err)
		require.EqualError(t, err, "failed to unpack subfield 3: reading length: strconv.Atoi: parsing \"AB\": invalid syntax")
	})

	t.Run("ReadFrom returns an error on data having subfield ID not in spec", func(t *testing.T) {
		data := &CompsiteTestData{}

		composite := NewComposite(compositeTestSpecWithIDLength)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Index 2-3 should have '01' rather than '11'.
		read, err := composite.ReadFrom(strings.NewReader("181102AB0202CD030212"))
		require.Equal(t, 0, read)
		require.EqualError(t, err, "failed to unpack subfield 11: field not defined in Spec")
	})

	t.Run("ReadFrom returns an error on failure to unpack subfield ID", func(t *testing.T) {
		data := &CompsiteTestData{}

		composite := NewComposite(compositeTestSpecWithIDLength)
		err := composite.SetData(data)
		require.NoError(t, err)

		// Index 0, 1 should have '01' rather than 'ID'.
		read, err := composite.ReadFrom(strings.NewReader("18ID02AB0202CD030212"))
		require.Equal(t, 0, read)
		require.EqualError(t, err, "failed to convert subfield ID \"ID\" to int")
	})

	t.Run("ReadFrom correctly deserialises out of order composite subfields to the data struct", func(t *testing.T) {
		data := &CompsiteTestData{}

		composite := NewComposite(compositeTestSpecWithIDLength)
		err := composite.SetData(data)
		require.NoError(t, err)

		read, err := composite.ReadFrom(strings.NewReader("280202CD0302120102AB04060102YZ"))

		require.NoError(t, err)
		require.Equal(t, 30, read)

		require.Equal(t, "AB", data.F1.Value)
		require.Equal(t, "CD", data.F2.Value)
		require.Equal(t, 12, data.F3.Value)
		require.Equal(t, "YZ", data.F4.F1.Value)
	})

	t.Run("ReadFrom correctly deserialises partial subfields to the data struct", func(t *testing.T) {
		data := &CompsiteTestData{}

		composite := NewComposite(compositeTestSpecWithIDLength)
		err := composite.SetData(data)
		require.NoError(t, err)

		read, err := composite.ReadFrom(strings.NewReader("120302120102AB"))

		require.NoError(t, err)
		require.Equal(t, 14, read)

		require.Equal(t, "AB", data.F1.Value)
		require.Nil(t, data.F2)
		require.Equal(t, 12, data.F3.Value)
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
				Fields: map[int]Field{},
			},
		},
		{
			desc: "accepts nil Pad value",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Fields: map[int]Field{},
			},
		},
		{
			desc: "accepts None Pad value",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Pad:    padding.None,
				Fields: map[int]Field{},
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
			desc: "panics on non-None / non-nil Pad value being defined in spec",
			err:  "Composite spec only supports nil or None padding values",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Pad:    padding.Left('0'),
				Fields: map[int]Field{},
			},
		},
		{
			desc: "panics on nil Enc value being defined in spec if IDLength > 0",
			err:  "Composite spec requires an Enc to be defined if IDLength > 0",
			spec: &Spec{
				Length:   6,
				IDLength: 2,
				Pref:     prefix.ASCII.Fixed,
				Fields:   map[int]Field{},
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
