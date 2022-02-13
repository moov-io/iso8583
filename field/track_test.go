package field

import (
	"fmt"
	"testing"
	"time"

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
	FixedLength          bool
}

var (
	track1Spec = &Spec{
		Length:      76,
		Description: "Track 1 Data",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.LL,
	}

	track2Spec = &Spec{
		Length:      37,
		Description: "Track 2 Data",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.LL,
	}

	track3Spec = &Spec{
		Length:      104,
		Description: "Track 3 Data",
		Enc:         encoding.ASCII,
		Pref:        prefix.ASCII.LLL,
	}
)

func TestTrack1(t *testing.T) {
	t.Run("Track 1 untyped", func(t *testing.T) {
		samples := []TestSample{
			{
				Raw:                  `B4815881002861896^YATES/EUGENE L            ^^^356858      00998000000`,
				FormatCode:           `B`,
				PrimaryAccountNumber: `4815881002861896`,
				DiscretionaryData:    `356858      00998000000`,
				Name:                 `YATES/EUGENE L`,
				FixedLength:          true,
			},
			{
				Raw:                  `B1234567890123445^PADILLA/L.                ^99011200000000000000**XXX******`,
				FormatCode:           `B`,
				PrimaryAccountNumber: `1234567890123445`,
				ServiceCode:          `120`,
				DiscretionaryData:    `0000000000000**XXX******`,
				ExpirationDate:       `9901`,
				Name:                 `PADILLA/L.`,
				FixedLength:          true,
			},
			{
				Raw:                  `B4242424242424242^SMITH JOHN Q^11052011000000000000`,
				FormatCode:           `B`,
				PrimaryAccountNumber: `4242424242424242`,
				ServiceCode:          `201`,
				DiscretionaryData:    `1000000000000`,
				ExpirationDate:       `1105`,
				Name:                 `SMITH JOHN Q`,
				FixedLength:          false,
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
				spec := &Spec{
					Length:      76,
					Description: "Track 1 Data",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.LL,
				}
				track := NewTrack1(spec)
				require.NotNil(t, track.Spec())

				// test SetBytes / Bytes
				track.FixedLength = sample.FixedLength
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

	t.Run("Track 1 typed", func(t *testing.T) {
		var (
			raw           = []byte("B1234567890123445^PADILLA/L.                ^99011200000000000000**XXX******")
			rawWithPrefix = []byte("76B1234567890123445^PADILLA/L.                ^99011200000000000000**XXX******")
		)

		t.Run("Returns an error on mismatch of track type", func(t *testing.T) {
			track := NewTrack1(track1Spec)
			err := track.SetData(NewStringValue("hello"))
			require.EqualError(t, err, "data does not match required *Track type")
		})

		t.Run("Unmarshal gets track values into data parameter", func(t *testing.T) {
			expDate, err := time.Parse("0601", "9901")
			require.NoError(t, err)

			track := NewTrack1(track1Spec)
			err = track.SetData(&Track1{
				FixedLength:          true,
				FormatCode:           "B",
				PrimaryAccountNumber: "1234567890123445",
				ServiceCode:          "120",
				DiscretionaryData:    "0000000000000**XXX******",
				ExpirationDate:       &expDate,
				Name:                 "PADILLA/L.",
			})
			require.NoError(t, err)

			data := &Track1{}

			err = track.Unmarshal(data)

			require.NoError(t, err)
			require.Equal(t, "B", data.FormatCode)
			require.Equal(t, "1234567890123445", data.PrimaryAccountNumber)
			require.Equal(t, "120", data.ServiceCode)
			require.Equal(t, "0000000000000**XXX******", data.DiscretionaryData)
			require.Equal(t, expDate, *data.ExpirationDate)
			require.Equal(t, "PADILLA/L.", data.Name)
		})

		t.Run("Pack correctly serializes data struct to bytes", func(t *testing.T) {
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
			require.Equal(t, rawWithPrefix, packed)
		})

		t.Run("Unpack correctly deserializes bytes with length prefix to the data struct", func(t *testing.T) {
			expDate, err := time.Parse("0601", "9901")
			require.NoError(t, err)

			data := &Track1{}

			track := NewTrack1(track1Spec)
			err = track.SetData(data)
			require.NoError(t, err)

			_, err = track.Unpack(rawWithPrefix)
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

			err = track.SetBytes(raw)
			require.NoError(t, err)

			// test assigned fields
			require.Equal(t, "B", data.FormatCode)
			require.Equal(t, "1234567890123445", data.PrimaryAccountNumber)
			require.Equal(t, "120", data.ServiceCode)
			require.Equal(t, "0000000000000**XXX******", data.DiscretionaryData)
			require.Equal(t, expDate, *data.ExpirationDate)
			require.Equal(t, "PADILLA/L.", data.Name)
		})
	})
}

func TestTrack2TypedAPI(t *testing.T) {
	var (
		raw                        = []byte("4000340000000506=2512111123400001230")
		rawWithDSeparator          = []byte("4000340000000506D2512111123400001230")
		rawWithPrefix              = []byte("364000340000000506=2512111123400001230")
		rawWithPrefixAndDSeparator = []byte("364000340000000506D2512111123400001230")
	)
	t.Run("Track 2 untyped", func(t *testing.T) {
		testCases := []struct {
			Bytes     []byte
			Separator string
			Packed    []byte
		}{
			{
				Bytes:     raw,
				Separator: "=",
				Packed:    rawWithPrefix,
			},
			{
				Bytes:     rawWithDSeparator,
				Separator: "D",
				Packed:    rawWithPrefixAndDSeparator,
			},
		}
		for _, tc := range testCases {
			tracker := NewTrack2(track2Spec)
			require.NotNil(t, tracker.Spec())

			err := tracker.SetBytes(tc.Bytes)
			require.NoError(t, err)

			buf, err := tracker.Bytes()
			require.NoError(t, err)
			require.Equal(t, tc.Bytes, buf)

			str, err := tracker.String()
			require.NoError(t, err)
			require.Equal(t, string(tc.Bytes), str)

			require.Equal(t, "4000340000000506", tracker.PrimaryAccountNumber)
			require.Equal(t, tc.Separator, tracker.Separator)
			require.Equal(t, "111", tracker.ServiceCode)
			require.Equal(t, "123400001230", tracker.DiscretionaryData)
			require.Equal(t, "2512", tracker.ExpirationDate.Format(expiryDateFormat))

			packBuf, err := tracker.Pack()
			require.NoError(t, err)
			require.Equal(t, tc.Packed, packBuf)

			_, err = tracker.Unpack(packBuf)
			require.NoError(t, err)
		}
	})

	t.Run("Track 2 typed", func(t *testing.T) {
		t.Run("Returns an error on mismatch of track type", func(t *testing.T) {
			track := NewTrack2(track2Spec)
			err := track.SetData(NewStringValue("hello"))
			require.EqualError(t, err, "data does not match required *Track type")
		})

		t.Run("Unmarshal gets track values into data parameter", func(t *testing.T) {
			expDate, err := time.Parse("0601", "9901")
			require.NoError(t, err)

			track := NewTrack2(track2Spec)
			err = track.SetData(&Track2{
				PrimaryAccountNumber: "4000340000000506",
				Separator:            "D",
				ServiceCode:          "111",
				DiscretionaryData:    "123400001230",
				ExpirationDate:       &expDate,
			})
			require.NoError(t, err)

			data := &Track2{}

			err = track.Unmarshal(data)

			require.NoError(t, err)
			require.Equal(t, "4000340000000506", data.PrimaryAccountNumber)
			require.Equal(t, "D", data.Separator)
			require.Equal(t, "111", data.ServiceCode)
			require.Equal(t, "123400001230", data.DiscretionaryData)
			require.Equal(t, expDate, *data.ExpirationDate)
		})

		t.Run("Pack correctly serializes data to bytes", func(t *testing.T) {
			testCases := []struct {
				Separator    string
				ExpectedPack []byte
			}{
				{
					Separator:    "=",
					ExpectedPack: rawWithPrefix,
				},
				{
					Separator:    "D",
					ExpectedPack: rawWithPrefixAndDSeparator,
				},
				{
					Separator:    "",
					ExpectedPack: rawWithPrefix,
				},
			}

			for _, tc := range testCases {
				expDate, err := time.Parse("0601", "2512")
				require.NoError(t, err)

				data := &Track2{
					PrimaryAccountNumber: `4000340000000506`,
					Separator:            tc.Separator,
					ServiceCode:          `111`,
					DiscretionaryData:    `123400001230`,
					ExpirationDate:       &expDate,
				}

				track := NewTrack2(track2Spec)
				err = track.SetData(data)
				require.NoError(t, err)

				// test assigned fields
				require.Equal(t, "4000340000000506", track.PrimaryAccountNumber)
				require.Equal(t, tc.Separator, track.Separator)
				require.Equal(t, "111", track.ServiceCode)
				require.Equal(t, "123400001230", track.DiscretionaryData)
				require.Equal(t, expDate, *track.ExpirationDate)

				packed, err := track.Pack()
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedPack, packed)
			}
		})

		t.Run("Unpack correctly deserializes bytes with length prefix to the data struct", func(t *testing.T) {
			testCases := []struct {
				Separator string
				Bytes     []byte
			}{
				{
					Separator: "=",
					Bytes:     rawWithPrefix,
				},
				{
					Separator: "D",
					Bytes:     rawWithPrefixAndDSeparator,
				},
			}

			for _, tc := range testCases {
				expDate, err := time.Parse("0601", "2512")
				require.NoError(t, err)

				data := &Track2{}

				track := NewTrack2(track2Spec)
				err = track.SetData(data)
				require.NoError(t, err)

				_, err = track.Unpack(tc.Bytes)
				require.NoError(t, err)

				// test assigned fields
				require.Equal(t, "4000340000000506", data.PrimaryAccountNumber)
				require.Equal(t, tc.Separator, data.Separator)
				require.Equal(t, "111", data.ServiceCode)
				require.Equal(t, "123400001230", data.DiscretionaryData)
				require.Equal(t, expDate, *data.ExpirationDate)
			}
		})

		t.Run("SetBytes correctly deserializes and assigns data", func(t *testing.T) {
			testCases := []struct {
				Separator string
				TrackData []byte
			}{
				{
					Separator: "=",
					TrackData: raw,
				},
				{
					Separator: "D",
					TrackData: rawWithDSeparator,
				},
			}

			for _, tc := range testCases {
				expDate, err := time.Parse("0601", "2512")
				require.NoError(t, err)

				data := &Track2{}

				track := NewTrack2(track2Spec)
				err = track.SetData(data)
				require.NoError(t, err)

				err = track.SetBytes(tc.TrackData)
				require.NoError(t, err)

				// test assigned fields
				require.Equal(t, "4000340000000506", data.PrimaryAccountNumber)
				require.Equal(t, tc.Separator, data.Separator)
				require.Equal(t, "111", data.ServiceCode)
				require.Equal(t, "123400001230", data.DiscretionaryData)
				require.Equal(t, expDate, *data.ExpirationDate)
			}
		})
	})
}

func TestTrack3TypedAPI(t *testing.T) {
	t.Run("Track 3 untyped", func(t *testing.T) {
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
				Pref:        prefix.ASCII.LLL,
			}
			tracker := NewTrack3(spec)
			require.NotNil(t, tracker.Spec())

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
			require.Equal(t, sample.DiscretionaryData, tracker.DiscretionaryData)

			packBuf, err := tracker.Pack()
			require.NoError(t, err)
			require.NotEqual(t, sample.Raw, string(packBuf))

			_, err = tracker.Unpack(packBuf)
			require.NoError(t, err)
		}
	})

	t.Run("Track 3 typed", func(t *testing.T) {
		var (
			raw           = []byte("011234567890123445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**")
			rawWithPrefix = []byte("104011234567890123445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**")
		)
		t.Run("Returns an error on mismatch of track type", func(t *testing.T) {
			track := NewTrack3(track3Spec)
			err := track.SetData(NewStringValue("hello"))
			require.EqualError(t, err, "data does not match required *Track type")
		})

		t.Run("Unmarshal gets track values into data parameter", func(t *testing.T) {
			track := NewTrack3(track3Spec)
			err := track.SetData(&Track3{
				FormatCode:           `01`,
				PrimaryAccountNumber: `1234567890123445`,
				DiscretionaryData:    `724724000000000****00300XXXX020200099010=********************==1=100000000000000000**`,
			})
			require.NoError(t, err)

			data := &Track3{}

			err = track.Unmarshal(data)

			require.NoError(t, err)
			require.Equal(t, "01", data.FormatCode)
			require.Equal(t, "1234567890123445", data.PrimaryAccountNumber)
			require.Equal(t, "724724000000000****00300XXXX020200099010=********************==1=100000000000000000**", data.DiscretionaryData)
		})

		t.Run("Pack correctly serializes data to bytes", func(t *testing.T) {
			data := &Track3{
				FormatCode:           `01`,
				PrimaryAccountNumber: `1234567890123445`,
				DiscretionaryData:    `724724000000000****00300XXXX020200099010=********************==1=100000000000000000**`,
			}

			track := NewTrack3(track3Spec)
			err := track.SetData(data)
			require.NoError(t, err)

			// test assigned fields
			require.Equal(t, "01", track.FormatCode)
			require.Equal(t, "1234567890123445", track.PrimaryAccountNumber)
			require.Equal(t, "724724000000000****00300XXXX020200099010=********************==1=100000000000000000**", track.DiscretionaryData)

			packed, err := track.Pack()
			require.NoError(t, err)
			require.Equal(t, rawWithPrefix, packed)
		})

		t.Run("Unpack correctly deserializes bytes with length prefix to the data struct", func(t *testing.T) {
			data := &Track3{}

			track := NewTrack3(track3Spec)
			err := track.SetData(data)
			require.NoError(t, err)

			_, err = track.Unpack(rawWithPrefix)
			require.NoError(t, err)

			// test assigned fields
			require.Equal(t, "01", data.FormatCode)
			require.Equal(t, "1234567890123445", data.PrimaryAccountNumber)
			require.Equal(t, "724724000000000****00300XXXX020200099010=********************==1=100000000000000000**", track.DiscretionaryData)
		})

		t.Run("SetBytes correctly deserializes and assigns data", func(t *testing.T) {
			data := &Track3{}

			track := NewTrack3(track3Spec)
			err := track.SetData(data)
			require.NoError(t, err)

			err = track.SetBytes(raw)
			require.NoError(t, err)

			// test assigned fields
			require.Equal(t, "01", data.FormatCode)
			require.Equal(t, "1234567890123445", data.PrimaryAccountNumber)
			require.Equal(t, "724724000000000****00300XXXX020200099010=********************==1=100000000000000000**", data.DiscretionaryData)
		})
	})
}
