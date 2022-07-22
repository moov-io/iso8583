package field

import (
	"github.com/moov-io/iso8583/sort"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func TestPrimitiveTLVField(t *testing.T) {
	spec := &Spec{
		Length:      6,
		Description: "Amount, Authorised",
		Enc:         encoding.ASCII,
		Pref:        prefix.BerTLV,
		Pad:         padding.Left(' '),
		Tag: &TagSpec{
			Tag:  "9F02",
			Enc:  encoding.BerTLVTag,
			Sort: sort.StringsByHex,
		},
	}
	tlv := NewPrimitiveTLV(spec)

	sampleValue := []byte{0x0, 0x0, 0x0, 0x0, 0x63, 0x0}
	sample := []byte{0x9f, 0x2, 0x6, 0x0, 0x0, 0x0, 0x0, 0x63, 0x0}

	tlv.SetBytes(sampleValue)
	require.Equal(t, sampleValue, tlv.Value)

	packed, err := tlv.Pack()
	require.NoError(t, err)
	require.Equal(t, sample, packed)

	length, err := tlv.Unpack(sample)
	require.NoError(t, err)
	require.Equal(t, 9, length)

	b, err := tlv.Bytes()
	require.NoError(t, err)
	require.Equal(t, sampleValue, b)

	require.Equal(t, sampleValue, tlv.Value)

	tlv = NewPrimitiveTLV(spec)
	tlv.SetData(NewPrimitiveTLVValue(sampleValue))
	packed, err = tlv.Pack()
	require.NoError(t, err)
	require.Equal(t, sample, packed)

	tlv = NewPrimitiveTLV(spec)
	data := NewPrimitiveTLVValue(nil)
	tlv.SetData(data)
	length, err = tlv.Unpack(sample)
	require.NoError(t, err)
	require.Equal(t, 9, length)
	require.Equal(t, sampleValue, data.Value)

	tlv = NewPrimitiveTLV(spec)
	data = &PrimitiveTLV{}
	tlv.SetData(data)
	err = tlv.SetBytes(sampleValue)
	require.NoError(t, err)
	require.Equal(t, sampleValue, data.Value)

	sampleWithUnmatchedTag := []byte{0x9f, 0x3, 0x6, 0x0, 0x0, 0x0, 0x0, 0x63, 0x0}

	tlv = NewPrimitiveTLV(spec)
	length, err = tlv.Unpack(sampleWithUnmatchedTag)
	require.Error(t, err)
	require.Equal(t, 0, length)
	require.Equal(t, "tag mismatch: want to read 9F02, got 9F03", err.Error())

}

func TestPrimitiveTLVPack(t *testing.T) {
	t.Run("pack with null value", func(t *testing.T) {
		spec := &Spec{
			Length:      6,
			Description: "Amount, Authorised",
			Enc:         encoding.ASCII,
			Pref:        prefix.BerTLV,
			Tag: &TagSpec{
				Tag:  "9F02",
				Enc:  encoding.BerTLVTag,
				Sort: sort.StringsByHex,
			},
		}
		tlv := NewPrimitiveTLV(spec)
		pack, err := tlv.Pack()

		require.NoError(t, err)
		require.Equal(t, []byte{0x9f, 0x02, 0x00}, pack)
	})

	t.Run("pack with null value and padding", func(t *testing.T) {
		spec := &Spec{
			Length:      6,
			Description: "Amount, Authorised",
			Enc:         encoding.ASCII,
			Pref:        prefix.BerTLV,
			Pad:         padding.Left(0),
			Tag: &TagSpec{
				Tag:  "9F02",
				Enc:  encoding.BerTLVTag,
				Sort: sort.StringsByHex,
			},
		}
		tlv := NewPrimitiveTLV(spec)
		pack, err := tlv.Pack()

		require.NoError(t, err)
		require.Equal(t, []byte{0x9f, 0x2, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, pack)
	})

	t.Run("pack with hex string", func(t *testing.T) {
		spec := &Spec{
			Length:      6,
			Description: "Amount, Authorised",
			Enc:         encoding.Binary,
			Pref:        prefix.BerTLV,
			Tag: &TagSpec{
				Tag:  "9F02",
				Enc:  encoding.BerTLVTag,
				Sort: sort.StringsByHex,
			},
		}

		tlv := NewPrimitiveTLVHexString("9FA813")
		tlv.SetSpec(spec)

		pack, err := tlv.Pack()

		require.NoError(t, err)
		require.Equal(t, []byte{0x9f, 0x2, 0x3, 0x9f, 0xa8, 0x13}, pack)
	})
}

func TestPrimitiveTLVFieldUnmarshal(t *testing.T) {

	input := []byte{0x0, 0x0, 0x0, 0x0, 0x63, 0x0}

	tlv := NewPrimitiveTLVValue(input)

	val := &PrimitiveTLV{}

	err := tlv.Unmarshal(val)

	require.NoError(t, err)

	require.Equal(t, input, val.Value)
}

func TestPrimitiveTLVJSONUnmarshal(t *testing.T) {
	input := []byte(`"000000006300"`)

	str := NewPrimitiveTLV(&Spec{
		Length:      6,
		Description: "Amount, Authorised",
		Enc:         encoding.ASCII,
		Pref:        prefix.BerTLV,
		Tag: &TagSpec{
			Tag:  "9F02",
			Enc:  encoding.BerTLVTag,
			Sort: sort.StringsByHex,
		},
	})

	require.NoError(t, str.UnmarshalJSON(input))
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x63, 0x0}, str.Value)
}

func TestPrimitiveTLVJSONMarshal(t *testing.T) {

	input := []byte{0x0, 0x0, 0x0, 0x0, 0x63, 0x0}

	str := NewPrimitiveTLVValue(input)

	marshalledJSON, err := str.MarshalJSON()

	require.NoError(t, err)
	require.Equal(t, `"000000006300"`, string(marshalledJSON))
}
