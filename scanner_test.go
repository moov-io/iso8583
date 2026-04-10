package iso8583

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	iso8583errors "github.com/moov-io/iso8583/errors"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageScanner(t *testing.T) {
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
			11: field.NewString(&field.Spec{
				Length:      6,
				Description: "STAN",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
		},
	}

	// Pack a full message for use in tests
	packMessage := func(t *testing.T) []byte {
		t.Helper()
		msg := NewMessage(spec)
		msg.MTI("0200")
		require.NoError(t, msg.Field(2, "4242424242424242"))
		require.NoError(t, msg.Field(3, "123456"))
		require.NoError(t, msg.Field(4, "100"))
		require.NoError(t, msg.Field(11, "123"))

		packed, err := msg.Pack()
		require.NoError(t, err)
		return packed
	}

	t.Run("scan MTI only", func(t *testing.T) {
		packed := packMessage(t)

		s := NewMessageScanner(spec, packed)

		f, err := s.ScanField(0)
		require.NoError(t, err)

		mti, err := f.String()
		require.NoError(t, err)
		assert.Equal(t, "0200", mti)
	})

	t.Run("scan MTI then STAN", func(t *testing.T) {
		packed := packMessage(t)

		s := NewMessageScanner(spec, packed)

		f, err := s.ScanField(0)
		require.NoError(t, err)
		mti, err := f.String()
		require.NoError(t, err)
		assert.Equal(t, "0200", mti)

		f, err = s.ScanField(11)
		require.NoError(t, err)
		stan, err := f.String()
		require.NoError(t, err)
		assert.Equal(t, "123", stan)
	})

	t.Run("scan multiple fields in order", func(t *testing.T) {
		packed := packMessage(t)

		s := NewMessageScanner(spec, packed)

		f, err := s.ScanField(0)
		require.NoError(t, err)
		mti, err := f.String()
		require.NoError(t, err)
		assert.Equal(t, "0200", mti)

		f, err = s.ScanField(2)
		require.NoError(t, err)
		pan, err := f.String()
		require.NoError(t, err)
		assert.Equal(t, "4242424242424242", pan)

		f, err = s.ScanField(4)
		require.NoError(t, err)
		amount, err := f.String()
		require.NoError(t, err)
		assert.Equal(t, "100", amount)
	})

	t.Run("scan only STAN", func(t *testing.T) {
		packed := packMessage(t)

		s := NewMessageScanner(spec, packed)

		f, err := s.ScanField(11)
		require.NoError(t, err)

		stan, err := f.String()
		require.NoError(t, err)
		assert.Equal(t, "123", stan)
	})

	t.Run("error when scanning backwards", func(t *testing.T) {
		packed := packMessage(t)

		s := NewMessageScanner(spec, packed)

		_, err := s.ScanField(0)
		require.NoError(t, err)

		_, err = s.ScanField(4)
		require.NoError(t, err)

		_, err = s.ScanField(2)
		require.Error(t, err)

		var unpackErr *iso8583errors.UnpackError
		require.ErrorAs(t, err, &unpackErr)
		assert.Contains(t, unpackErr.Err.Error(), "forward-only")
	})

	t.Run("error when scanning same field twice", func(t *testing.T) {
		packed := packMessage(t)

		s := NewMessageScanner(spec, packed)

		_, err := s.ScanField(0)
		require.NoError(t, err)

		_, err = s.ScanField(0)
		require.Error(t, err)
	})

	t.Run("error when field not in bitmap", func(t *testing.T) {
		// Pack message with only fields 2 and 3
		msg := NewMessage(spec)
		msg.MTI("0200")
		require.NoError(t, msg.Field(2, "4242424242424242"))
		require.NoError(t, msg.Field(3, "123456"))
		packed, err := msg.Pack()
		require.NoError(t, err)

		s := NewMessageScanner(spec, packed)

		_, err = s.ScanField(0)
		require.NoError(t, err)

		// Field 11 is not set in the bitmap
		_, err = s.ScanField(11)
		require.Error(t, err)

		var unpackErr *iso8583errors.UnpackError
		require.ErrorAs(t, err, &unpackErr)
		assert.Contains(t, unpackErr.Err.Error(), "not set in the bitmap")
	})

	t.Run("proxy use case: branch on MTI", func(t *testing.T) {
		packed := packMessage(t)

		s := NewMessageScanner(spec, packed)

		f, err := s.ScanField(0)
		require.NoError(t, err)

		mti, err := f.String()
		require.NoError(t, err)

		// Not a network message, so continue to get STAN
		assert.Equal(t, "0200", mti)

		f, err = s.ScanField(11)
		require.NoError(t, err)

		stan, err := f.String()
		require.NoError(t, err)
		assert.Equal(t, "123", stan)
	})
}
