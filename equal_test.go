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

var equalTestSpec = &MessageSpec{
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
		4: field.NewNumeric(&field.Spec{
			Length:      12,
			Description: "Transaction Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
	},
}

func newEqualTestMessage(t *testing.T) *Message {
	t.Helper()

	m := NewMessage(equalTestSpec)
	m.MTI("0100")
	require.NoError(t, m.Field(2, "4242424242424242"))
	require.NoError(t, m.MarshalPath("3.1", "00"))
	require.NoError(t, m.MarshalPath("3.2", "20"))
	require.NoError(t, m.MarshalPath("3.3", "30"))
	require.NoError(t, m.Field(4, "100"))

	return m
}

func requireEqual(t *testing.T, want bool, a, b *Message) {
	t.Helper()

	got, err := Equal(a, b)
	require.NoError(t, err)
	require.Equal(t, want, got)

	// Equal must be symmetric.
	got, err = Equal(b, a)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestEqual(t *testing.T) {
	t.Run("messages with the same fields are equal", func(t *testing.T) {
		requireEqual(t, true, newEqualTestMessage(t), newEqualTestMessage(t))
	})

	t.Run("a message equals itself", func(t *testing.T) {
		m := newEqualTestMessage(t)
		requireEqual(t, true, m, m)
	})

	t.Run("a message equals its packed and unpacked copy", func(t *testing.T) {
		original := newEqualTestMessage(t)
		packed, err := original.Pack()
		require.NoError(t, err)

		unpacked := NewMessage(equalTestSpec)
		require.NoError(t, unpacked.Unpack(packed))

		requireEqual(t, true, original, unpacked)
	})

	t.Run("bitmap state does not affect equality", func(t *testing.T) {
		packed := newEqualTestMessage(t)
		_, err := packed.Pack() // computes and stores the bitmap
		require.NoError(t, err)

		requireEqual(t, true, packed, newEqualTestMessage(t))
	})

	t.Run("different MTI", func(t *testing.T) {
		other := newEqualTestMessage(t)
		other.MTI("0110")
		requireEqual(t, false, newEqualTestMessage(t), other)
	})

	t.Run("different field value", func(t *testing.T) {
		other := newEqualTestMessage(t)
		require.NoError(t, other.Field(4, "101"))
		requireEqual(t, false, newEqualTestMessage(t), other)
	})

	t.Run("field set in only one message", func(t *testing.T) {
		other := newEqualTestMessage(t)
		other.UnsetField(2)
		requireEqual(t, false, newEqualTestMessage(t), other)
	})

	t.Run("different composite subfield value", func(t *testing.T) {
		other := newEqualTestMessage(t)
		require.NoError(t, other.MarshalPath("3.2", "10"))
		requireEqual(t, false, newEqualTestMessage(t), other)
	})

	t.Run("composite subfield set in only one message", func(t *testing.T) {
		other := newEqualTestMessage(t)
		require.NoError(t, other.UnsetPath("3.3"))
		requireEqual(t, false, newEqualTestMessage(t), other)
	})

	t.Run("nil messages", func(t *testing.T) {
		requireEqual(t, true, nil, nil)
		requireEqual(t, false, nil, newEqualTestMessage(t))
	})

	t.Run("empty messages are equal", func(t *testing.T) {
		requireEqual(t, true, NewMessage(equalTestSpec), NewMessage(equalTestSpec))
	})
}
