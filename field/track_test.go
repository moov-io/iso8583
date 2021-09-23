package field

import (
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

type TestSample struct {
	Raw                  string
	Name                 string
	FormatCode           string
	PrimaryAccountNumber string
	ServiceCode          string
	DiscretionaryData    string
	ExpirationDate       string
}

func TestTrack(t *testing.T) {

	t.Run("Track 1 data with fixed name length", func(t *testing.T) {
		samples := []TestSample{
			{
				Raw:                  `B4815881002861896^YATES/EUGENE L            ^^^356858      00998000000`,
				FormatCode:           `B`,
				PrimaryAccountNumber: `4815881002861896`,
				DiscretionaryData:    `356858      00998000000`,
				Name:                 `YATES/EUGENE L`,
			},
			{
				Raw:                  `B1234567890123445^PADILLA/L.                ^99011200000000000000**XXX******`,
				FormatCode:           `B`,
				PrimaryAccountNumber: `1234567890123445`,
				ServiceCode:          `120`,
				DiscretionaryData:    `0000000000000**XXX******`,
				ExpirationDate:       `9901`,
				Name:                 `PADILLA/L.`,
			},
		}
		for _, sample := range samples {
			spec := &Spec{
				Length:      76,
				Description: "Track 1 Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}
			tracker, err := NewTrack(spec, Track1)
			require.NoError(t, err)

			tracker.FixedLength = true
			err = tracker.SetBytes([]byte(sample.Raw))
			require.NoError(t, err)

			buf, err := tracker.Bytes()
			require.NoError(t, err)
			require.Equal(t, sample.Raw, string(buf))

			require.Equal(t, sample.FormatCode, tracker.FormatCode)
			require.Equal(t, sample.PrimaryAccountNumber, tracker.PrimaryAccountNumber)
			require.Equal(t, sample.ServiceCode, tracker.ServiceCode)
			require.Equal(t, sample.Name, tracker.Name)
			require.Equal(t, sample.DiscretionaryData, tracker.DiscretionaryData)
			if len(sample.ExpirationDate) > 0 {
				require.NotNil(t, tracker.ExpirationDate)
				require.Equal(t, sample.ExpirationDate, tracker.ExpirationDate.Format(expiryDateFormat))
			}
		}
	})
	t.Run("Track 1 data with unfixed name length", func(t *testing.T) {
		spec := &Spec{
			Length:      76,
			Description: "Track 1 Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}
		tracker, err := NewTrack(spec, Track1)
		require.NoError(t, err)

		sample := TestSample{
			Raw:                  `B4242424242424242^SMITH JOHN Q^11052011000000000000`,
			FormatCode:           `B`,
			PrimaryAccountNumber: `4242424242424242`,
			ServiceCode:          `201`,
			DiscretionaryData:    `1000000000000`,
			ExpirationDate:       `1105`,
			Name:                 `SMITH JOHN Q`,
		}

		err = tracker.SetBytes([]byte(sample.Raw))
		require.NoError(t, err)

		buf, err := tracker.Bytes()
		require.NoError(t, err)
		require.Equal(t, sample.Raw, string(buf))

		require.Equal(t, sample.FormatCode, tracker.FormatCode)
		require.Equal(t, sample.PrimaryAccountNumber, tracker.PrimaryAccountNumber)
		require.Equal(t, sample.ServiceCode, tracker.ServiceCode)
		require.Equal(t, sample.Name, tracker.Name)
		require.Equal(t, sample.DiscretionaryData, tracker.DiscretionaryData)
		if len(sample.ExpirationDate) > 0 {
			require.NotNil(t, tracker.ExpirationDate)
			require.Equal(t, sample.ExpirationDate, tracker.ExpirationDate.Format(expiryDateFormat))
		}
	})

	t.Run("Track 2 data", func(t *testing.T) {
		samples := []TestSample{
			{
				Raw:                  `4000340000000506=2512111123400001230`,
				PrimaryAccountNumber: `4000340000000506`,
				ServiceCode:          `111`,
				DiscretionaryData:    `123400001230`,
				ExpirationDate:       `2512`,
			},
			{
				Raw:                  `1234567890123445=99011200XXXX00000000`,
				PrimaryAccountNumber: `1234567890123445`,
				ServiceCode:          `120`,
				DiscretionaryData:    `0XXXX00000000`,
				ExpirationDate:       `9901`,
			},
		}
		for _, sample := range samples {
			spec := &Spec{
				Length:      37,
				Description: "Track 2 Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}
			tracker, err := NewTrack(spec, Track2)
			require.NoError(t, err)

			err = tracker.SetBytes([]byte(sample.Raw))
			require.NoError(t, err)

			buf, err := tracker.Bytes()
			require.NoError(t, err)
			require.Equal(t, sample.Raw, string(buf))

			require.Equal(t, sample.FormatCode, tracker.FormatCode)
			require.Equal(t, sample.PrimaryAccountNumber, tracker.PrimaryAccountNumber)
			require.Equal(t, sample.ServiceCode, tracker.ServiceCode)
			require.Equal(t, sample.Name, tracker.Name)
			require.Equal(t, sample.DiscretionaryData, tracker.DiscretionaryData)
			if len(sample.ExpirationDate) > 0 {
				require.NotNil(t, tracker.ExpirationDate)
				require.Equal(t, sample.ExpirationDate, tracker.ExpirationDate.Format(expiryDateFormat))
			}
		}
	})

	t.Run("Track 3 data", func(t *testing.T) {

		samples := []TestSample{
			{
				Raw:                  `011234567890123445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**`,
				FormatCode:           `01`,
				PrimaryAccountNumber: `1234567890123445`,
				DiscretionaryData:    `724724000000000****00300XXXX020200099010=********************==1=100000000000000000**`,
			},
			{
				Raw:                  `011234567890123445=000978100000000****8330*0000920000099010=************************==1=0000000*00000000`,
				FormatCode:           `01`,
				PrimaryAccountNumber: `1234567890123445`,
				DiscretionaryData:    `000978100000000****8330*0000920000099010=************************==1=0000000*00000000`,
			},
		}
		for _, sample := range samples {
			spec := &Spec{
				Length:      104,
				Description: "Track 3 Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}
			tracker, err := NewTrack(spec, Track3)
			require.NoError(t, err)

			err = tracker.SetBytes([]byte(sample.Raw))
			require.NoError(t, err)

			buf, err := tracker.Bytes()
			require.NoError(t, err)
			require.Equal(t, sample.Raw, string(buf))

			require.Equal(t, sample.FormatCode, tracker.FormatCode)
			require.Equal(t, sample.PrimaryAccountNumber, tracker.PrimaryAccountNumber)
			require.Equal(t, sample.ServiceCode, tracker.ServiceCode)
			require.Equal(t, sample.Name, tracker.Name)
			require.Equal(t, sample.DiscretionaryData, tracker.DiscretionaryData)
			if len(sample.ExpirationDate) > 0 {
				require.NotNil(t, tracker.ExpirationDate)
				require.Equal(t, sample.ExpirationDate, tracker.ExpirationDate.Format(expiryDateFormat))
			}
		}
	})

	t.Run("Track value", func(t *testing.T) {

		raw := `B4242424242424242^SMITH JOHN Q^11052011000000000000`
		tracker, err := NewTrackValue([]byte(raw), Track1, false)
		require.NoError(t, err)

		buf, err := tracker.Bytes()
		require.NoError(t, err)

		require.Equal(t, raw, string(buf))
	})
}
