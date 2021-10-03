package track

import (
	"fmt"
	"testing"
	"time"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
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
	FixedLength          bool
}

var (
	track1Spec = &field.Spec{
		Length:      76,
		Description: "Track 1 Data",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.LL,
	}

	track2Spec = &field.Spec{
		Length:      37,
		Description: "Track 2 Data",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.LL,
	}
)

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

		testTrackFields := func(t *testing.T, track *Track1, sample TestSample) {
			require.Equal(t, sample.FormatCode, track.FormatCode)
			require.Equal(t, sample.PrimaryAccountNumber, track.PrimaryAccountNumber)
			require.Equal(t, sample.ServiceCode, track.ServiceCode)
			require.Equal(t, sample.Name, track.Name)
			require.Equal(t, sample.DiscretionaryData, track.DiscretionaryData)

			if len(sample.ExpirationDate) > 0 {
				require.NotNil(t, track.ExpirationDate)
				require.Equal(t, sample.ExpirationDate, track.ExpirationDate.Format(expiryDateFormat))
			}
		}

		for id, sample := range samples {
			t.Run(fmt.Sprintf("sample %d", id), func(t *testing.T) {
				spec := &field.Spec{
					Length:      76,
					Description: "Track 1 Data",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.LL,
				}
				track := NewTrack1(spec)
				require.NotNil(t, track.Spec())

				// test SetBytes / Bytes
				track.FixedLength = true
				err := track.SetBytes([]byte(sample.Raw))
				require.NoError(t, err)

				testTrackFields(t, track, sample)

				buf, err := track.Bytes()
				require.NoError(t, err)
				require.Equal(t, sample.Raw, string(buf))

				// Test Pack / Unpack
				packBuf, err := track.Pack()
				require.NoError(t, err)
				require.Len(t, packBuf, len(sample.Raw)+2, "packed length must be 2 bytes longer as it has ASCII.LL prefix")

				unpackedTrack := NewTrack1(spec)
				require.NoError(t, err)

				_, err = unpackedTrack.Unpack(packBuf)
				require.NoError(t, err)

				testTrackFields(t, unpackedTrack, sample)
			})

		}
	})

	t.Run("Track 1 data with unfixed name length", func(t *testing.T) {
		spec := &field.Spec{
			Length:      76,
			Description: "Track 1 Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}
		tracker := NewTrack1(spec)
		require.NotNil(t, tracker.Spec())

		sample := TestSample{
			Raw:                  `B4242424242424242^SMITH JOHN Q^11052011000000000000`,
			FormatCode:           `B`,
			PrimaryAccountNumber: `4242424242424242`,
			ServiceCode:          `201`,
			DiscretionaryData:    `1000000000000`,
			ExpirationDate:       `1105`,
			Name:                 `SMITH JOHN Q`,
		}

		err := tracker.SetBytes([]byte(sample.Raw))
		require.NoError(t, err)

		buf, err := tracker.Bytes()
		require.NoError(t, err)
		require.Equal(t, sample.Raw, string(buf))

		str, err := tracker.String()
		require.NoError(t, err)
		require.Equal(t, sample.Raw, str)

		require.Equal(t, sample.FormatCode, tracker.FormatCode)
		require.Equal(t, sample.PrimaryAccountNumber, tracker.PrimaryAccountNumber)
		require.Equal(t, sample.ServiceCode, tracker.ServiceCode)
		require.Equal(t, sample.Name, tracker.Name)
		require.Equal(t, sample.DiscretionaryData, tracker.DiscretionaryData)
		if len(sample.ExpirationDate) > 0 {
			require.NotNil(t, tracker.ExpirationDate)
			require.Equal(t, sample.ExpirationDate, tracker.ExpirationDate.Format(expiryDateFormat))
		}

		packBuf, err := tracker.Pack()
		require.NoError(t, err)
		require.NotEqual(t, sample.Raw, string(packBuf))

		_, err = tracker.Unpack(packBuf)
		require.NoError(t, err)
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
			spec := &field.Spec{
				Length:      37,
				Description: "Track 2 Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}
			tracker := NewTrack2(spec)
			require.NotNil(t, tracker.Spec())

			err := tracker.SetBytes([]byte(sample.Raw))
			require.NoError(t, err)

			buf, err := tracker.Bytes()
			require.NoError(t, err)
			require.Equal(t, sample.Raw, string(buf))

			str, err := tracker.String()
			require.NoError(t, err)
			require.Equal(t, sample.Raw, str)

			require.Equal(t, sample.PrimaryAccountNumber, tracker.PrimaryAccountNumber)
			require.Equal(t, sample.ServiceCode, tracker.ServiceCode)
			require.Equal(t, sample.DiscretionaryData, tracker.DiscretionaryData)
			if len(sample.ExpirationDate) > 0 {
				require.NotNil(t, tracker.ExpirationDate)
				require.Equal(t, sample.ExpirationDate, tracker.ExpirationDate.Format(expiryDateFormat))
			}

			packBuf, err := tracker.Pack()
			require.NoError(t, err)
			require.NotEqual(t, sample.Raw, string(packBuf))

			_, err = tracker.Unpack(packBuf)
			require.NoError(t, err)
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
			spec := &field.Spec{
				Length:      104,
				Description: "Track 3 Data",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LLL,
			}
			tracker, err := NewTrack3(spec)
			require.NoError(t, err)
			require.NotNil(t, tracker.Spec())

			err = tracker.SetBytes([]byte(sample.Raw))
			require.NoError(t, err)

			buf, err := tracker.Bytes()
			require.NoError(t, err)
			require.Equal(t, sample.Raw, string(buf))

			str, err := tracker.String()
			require.NoError(t, err)
			require.Equal(t, sample.Raw, str)

			require.Equal(t, sample.FormatCode, tracker.FormatCode)
			require.Equal(t, sample.PrimaryAccountNumber, tracker.PrimaryAccountNumber)
			require.Equal(t, sample.DiscretionaryData, tracker.DiscretionaryData)

			packBuf, err := tracker.Pack()
			require.NoError(t, err)
			require.NotEqual(t, sample.Raw, string(packBuf))

			_, err = tracker.Unpack(packBuf)
			require.NoError(t, err)
		}
	})

	t.Run("Track1 value", func(t *testing.T) {

		raw := `B4242424242424242^SMITH JOHN Q^11052011000000000000`
		tracker, err := NewTrack1Value([]byte(raw), false)
		require.NoError(t, err)

		buf, err := tracker.Bytes()
		require.NoError(t, err)

		require.Equal(t, raw, string(buf))

		_, err = NewTrack2Value([]byte(raw))
		require.Error(t, err)

		err = tracker.SetData(&Track1{})
		require.NoError(t, err)

		err = tracker.SetData(&Track3{})
		require.Error(t, err)
	})

	t.Run("Track2 value", func(t *testing.T) {

		raw := `1234567890123445=99011200XXXX00000000`
		tracker, err := NewTrack2Value([]byte(raw))
		require.NoError(t, err)

		buf, err := tracker.Bytes()
		require.NoError(t, err)

		require.Equal(t, raw, string(buf))

		_, err = NewTrack1Value([]byte(raw), false)
		require.Error(t, err)

		err = tracker.SetData(&Track2{})
		require.NoError(t, err)

		err = tracker.SetData(&Track3{})
		require.Error(t, err)
	})

	t.Run("Track3 value", func(t *testing.T) {

		raw := `011234567890123445=000978100000000****8330*0000920000099010=************************==1=0000000*00000000`
		tracker, err := NewTrack3Value([]byte(raw))
		require.NoError(t, err)

		buf, err := tracker.Bytes()
		require.NoError(t, err)

		require.Equal(t, raw, string(buf))

		_, err = NewTrack1Value([]byte(raw), false)
		require.Error(t, err)

		err = tracker.SetData(&Track3{})
		require.NoError(t, err)

		err = tracker.SetData(&Track2{})
		require.Error(t, err)
	})
}

func TestTrack1TypedAPI(t *testing.T) {
	t.Run("Returns an error on mismatch of track type", func(t *testing.T) {
		track := NewTrack1(track1Spec)
		err := track.SetData(field.NewStringValue("hello"))
		require.EqualError(t, err, "data does not match required *Track type")
	})

	t.Run("Pack correctly serializes data to bytes", func(t *testing.T) {
		expDate, err := time.Parse("0601", "9901")
		require.NoError(t, err)

		data := &Track1{
			FixedLength:          true,
			FormatCode:           "B",
			PrimaryAccountNumber: "1234567890123445",
			ServiceCode:          "120",
			DiscretionaryData:    "0000000000000**XXX******",
			ExpirationDate:       &expDate,
			Name:                 "PADILLA/L.",
		}

		track := NewTrack1(track1Spec)
		err = track.SetData(data)
		require.NoError(t, err)

		// test assigned fields
		require.Equal(t, "B", track.FormatCode)
		require.Equal(t, "1234567890123445", track.PrimaryAccountNumber)
		require.Equal(t, "120", track.ServiceCode)
		require.Equal(t, "0000000000000**XXX******", track.DiscretionaryData)
		require.Equal(t, expDate, *track.ExpirationDate)
		require.Equal(t, "PADILLA/L.", track.Name)

		packed, err := track.Pack()
		require.NoError(t, err)
		require.Equal(t, []byte("76B1234567890123445^PADILLA/L.                ^99011200000000000000**XXX******"), packed)
	})

	t.Run("Unpack correctly deserializes bytes with length prefix to the data struct", func(t *testing.T) {
		expDate, err := time.Parse("0601", "9901")
		require.NoError(t, err)

		data := &Track1{}

		track := NewTrack1(track1Spec)
		err = track.SetData(data)
		require.NoError(t, err)

		// bytes with LL prefix 76
		_, err = track.Unpack([]byte("76B1234567890123445^PADILLA/L.                ^99011200000000000000**XXX******"))
		require.NoError(t, err)

		// test assigned fields
		require.Equal(t, "B", data.FormatCode)
		require.Equal(t, "1234567890123445", data.PrimaryAccountNumber)
		require.Equal(t, "120", data.ServiceCode)
		require.Equal(t, "0000000000000**XXX******", data.DiscretionaryData)
		require.Equal(t, expDate, *data.ExpirationDate)
		require.Equal(t, "PADILLA/L.", data.Name)
	})

	t.Run("SetBytes correctly deserializes and assigns data", func(t *testing.T) {
		expDate, err := time.Parse("0601", "9901")
		require.NoError(t, err)

		data := &Track1{}

		track := NewTrack1(track1Spec)
		err = track.SetData(data)
		require.NoError(t, err)

		// bytes with LL prefix
		err = track.SetBytes([]byte("B1234567890123445^PADILLA/L.                ^99011200000000000000**XXX******"))
		require.NoError(t, err)

		// test assigned fields
		require.Equal(t, "B", data.FormatCode)
		require.Equal(t, "1234567890123445", data.PrimaryAccountNumber)
		require.Equal(t, "120", data.ServiceCode)
		require.Equal(t, "0000000000000**XXX******", data.DiscretionaryData)
		require.Equal(t, expDate, *data.ExpirationDate)
		require.Equal(t, "PADILLA/L.", data.Name)
	})
}

func TestTrack2TypedAPI(t *testing.T) {
	var (
		raw           = []byte("4000340000000506=2512111123400001230")
		rawWithPrefix = []byte("364000340000000506=2512111123400001230")
	)
	t.Run("Returns an error on mismatch of track type", func(t *testing.T) {
		track := NewTrack2(track2Spec)
		err := track.SetData(field.NewStringValue("hello"))
		require.EqualError(t, err, "data does not match required *Track type")
	})

	t.Run("Pack correctly serializes data to bytes", func(t *testing.T) {
		expDate, err := time.Parse("0601", "2512")
		require.NoError(t, err)

		data := &Track2{
			PrimaryAccountNumber: `4000340000000506`,
			ServiceCode:          `111`,
			DiscretionaryData:    `123400001230`,
			ExpirationDate:       &expDate,
		}

		track := NewTrack2(track2Spec)
		err = track.SetData(data)
		require.NoError(t, err)

		// test assigned fields
		require.Equal(t, "4000340000000506", track.PrimaryAccountNumber)
		require.Equal(t, "111", track.ServiceCode)
		require.Equal(t, "123400001230", track.DiscretionaryData)
		require.Equal(t, expDate, *track.ExpirationDate)

		packed, err := track.Pack()
		require.NoError(t, err)
		require.Equal(t, rawWithPrefix, packed)
	})

	t.Run("Unpack correctly deserializes bytes with length prefix to the data struct", func(t *testing.T) {
		expDate, err := time.Parse("0601", "2512")
		require.NoError(t, err)

		data := &Track2{}

		track := NewTrack2(track2Spec)
		err = track.SetData(data)
		require.NoError(t, err)

		// bytes with LL prefix 36
		_, err = track.Unpack(rawWithPrefix)
		require.NoError(t, err)

		// test assigned fields
		require.Equal(t, "4000340000000506", data.PrimaryAccountNumber)
		require.Equal(t, "111", data.ServiceCode)
		require.Equal(t, "123400001230", data.DiscretionaryData)
		require.Equal(t, expDate, *data.ExpirationDate)
	})

	t.Run("SetBytes correctly deserializes and assigns data", func(t *testing.T) {
		expDate, err := time.Parse("0601", "2512")
		require.NoError(t, err)

		data := &Track2{}

		track := NewTrack2(track2Spec)
		err = track.SetData(data)
		require.NoError(t, err)

		// bytes without LL prefix
		err = track.SetBytes(raw)
		require.NoError(t, err)

		// test assigned fields
		require.Equal(t, "4000340000000506", data.PrimaryAccountNumber)
		require.Equal(t, "111", data.ServiceCode)
		require.Equal(t, "123400001230", data.DiscretionaryData)
		require.Equal(t, expDate, *data.ExpirationDate)
	})
}
