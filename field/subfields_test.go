package field

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

var (
	subfieldsTestSpec = &Spec{
		Length:      6,
		Description: "Test Spec",
		Pref:        prefix.ASCII.Fixed,
		Enc:         encoding.None,
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
)

type SubfieldTestData struct {
	F1 *String
	F2 *String
	F3 *Numeric
}

func TestSubfields_SetData(t *testing.T) {
	t.Run("SetData returns an error on mismatch of subfield types", func(t *testing.T) {
		type TestDataIncorrectType struct {
			F1 *Numeric
		}

		subfields := NewSubfields(subfieldsTestSpec)
		err := subfields.SetData(&TestDataIncorrectType{
			F1: NewNumericValue(1),
		})
		require.Error(t, err)
	})

	t.Run("SetData returns an error on provision of primitive type data", func(t *testing.T) {
		subfields := NewSubfields(subfieldsTestSpec)
		err := subfields.SetData("primitive str")
		require.Error(t, err)
	})
}

func TestSubfieldsPacking(t *testing.T) {
	t.Run("Pack returns error on failure of subfield packing", func(t *testing.T) {
		data := &SubfieldTestData{
			// This subfield will return an error on F1.Pack() as its length
			// exceeds the max length defined in the spec.
			F1: NewStringValue("ABCD"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		subfields := NewSubfields(subfieldsTestSpec)
		err := subfields.SetData(data)
		require.NoError(t, err)

		_, err = subfields.Pack()
		require.Error(t, err)
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
		data := &SubfieldTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		subfields := NewSubfields(invalidSpec)
		err := subfields.SetData(data)
		require.NoError(t, err)

		_, err = subfields.Pack()
		require.Error(t, err)
	})

	t.Run("Pack correctly serialises data to bytes", func(t *testing.T) {
		data := &SubfieldTestData{
			F1: NewStringValue("AB"),
			F2: NewStringValue("CD"),
			F3: NewNumericValue(12),
		}

		subfields := NewSubfields(subfieldsTestSpec)
		err := subfields.SetData(data)
		require.NoError(t, err)

		packed, err := subfields.Pack()
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, "ABCD12", string(packed))
	})

	t.Run("Unpack returns an error on failure of subfield to unpack bytes", func(t *testing.T) {
		data := &SubfieldTestData{}

		subfields := NewSubfields(subfieldsTestSpec)
		err := subfields.SetData(data)
		require.NoError(t, err)

		// Last two characters must be an integer type. F3 fails to unpack.
		read, err := subfields.Unpack([]byte("ABCDEF"))
		require.Equal(t, 0, read)
		require.Error(t, err)
	})

	t.Run("Unpack returns an error on length of data exceeding max length", func(t *testing.T) {
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
		data := &SubfieldTestData{}

		subfields := NewSubfields(spec)
		err := subfields.SetData(data)
		require.NoError(t, err)

		// Length of denoted by prefix is too long, causing failure to decode length.
		read, err := subfields.Unpack([]byte("7ABCD123"))
		require.Equal(t, 0, read)
		require.Error(t, err)
	})

	t.Run("Unpack returns an error on offset not matching data length", func(t *testing.T) {
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
		data := &SubfieldTestData{}

		subfields := NewSubfields(invalidSpec)
		err := subfields.SetData(data)
		require.NoError(t, err)

		// Length of input too long, causing failure to decode length.
		read, err := subfields.Unpack([]byte("ABCD123"))
		require.Equal(t, 0, read)
		require.Error(t, err)
	})

	t.Run("Unpack correctly deserialises bytes to the data struct", func(t *testing.T) {
		data := &SubfieldTestData{}

		subfields := NewSubfields(subfieldsTestSpec)
		err := subfields.SetData(data)
		require.NoError(t, err)

		read, err := subfields.Unpack([]byte("ABCD12"))
		require.Equal(t, subfieldsTestSpec.Length, read)
		require.NoError(t, err)

		require.Equal(t, "AB", data.F1.Value)
		require.Equal(t, "CD", data.F2.Value)
		require.Equal(t, 12, data.F3.Value)
	})
}

func TestSubfieldsHandlesValidSpecs(t *testing.T) {
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
			desc: "accepts None Enc value",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Enc:    encoding.None,
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
		t.Run(fmt.Sprintf("NewSubfields() %v", tc.desc), func(t *testing.T) {
			f := NewSubfields(tc.spec)
			require.Equal(t, tc.spec, f.Spec())
		})
		t.Run(fmt.Sprintf("Subfields.SetSpec() %v", tc.desc), func(t *testing.T) {
			f := &Subfields{}
			f.SetSpec(tc.spec)
			require.Equal(t, tc.spec, f.Spec())
		})
	}
}

func TestSubfieldsPanicsOnSpecValidationFailures(t *testing.T) {
	tests := []struct {
		desc string
		spec *Spec
	}{
		{
			desc: "panics on non-None / non-nil Enc value being defined in spec",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Enc:    encoding.ASCII,
				Pad:    padding.Left('0'),
				Fields: map[int]Field{},
			},
		},
		{
			desc: "panics on non-None / non-nil Pad value being defined in spec",
			spec: &Spec{
				Length: 6,
				Pref:   prefix.ASCII.Fixed,
				Enc:    encoding.ASCII,
				Pad:    padding.Left('0'),
				Fields: map[int]Field{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("NewSubfields() %v", tc.desc), func(t *testing.T) {
			require.Panics(t, func() {
				NewSubfields(tc.spec)
			})
		})
		t.Run(fmt.Sprintf("Subfields.SetSpec() %v", tc.desc), func(t *testing.T) {
			require.Panics(t, func() {
				(&Subfields{}).SetSpec(tc.spec)
			})
		})
	}
}
