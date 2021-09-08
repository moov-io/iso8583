package track

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrack(t *testing.T) {

	t.Run("Track 1 data", func(t *testing.T) {
		tracker := NewTrackFirst(true)

		samples := []string{
			`B4815881002867896^YATES/EUGENE JOHN         ^21129821000123456789`,
			`B4815881002861896^YATES/EUGENE L            ^^^356858      00998000000`,
			`B4000340000000506^John/Doe                  ^13011110000123000`,
			`B1234567890123445^PADILLA/L.                ^99011200000000000000**XXX******`,
		}
		for _, sample := range samples {
			card, err := tracker.Read([]byte(sample))
			require.NoError(t, err)
			buf, err := tracker.Write(card)
			require.NoError(t, err)
			require.Equal(t, sample, string(buf))
		}
	})
	t.Run("Track 1 data with invalid name length", func(t *testing.T) {
		tracker := NewTrackFirst(false)

		samples := []string{
			`B4242424242424242^SMITH JOHN Q^11052011000000000000`,
		}
		for _, sample := range samples {
			card, err := tracker.Read([]byte(sample))
			require.NoError(t, err)
			buf, err := tracker.Write(card)
			require.NoError(t, err)
			require.Equal(t, sample, string(buf))
		}
	})

	t.Run("Track 2 data", func(t *testing.T) {
		tracker := NewTrackSecond()

		samples := []string{
			`4000340000000506=2512111123400001230`,
			`4242424242424242=15052011000000000000`,
			`1234567890123445=99011200XXXX00000000`,
		}
		for _, sample := range samples {
			card, err := tracker.Read([]byte(sample))
			require.NoError(t, err)
			buf, err := tracker.Write(card)
			require.NoError(t, err)
			require.Equal(t, sample, string(buf))
		}
	})

	t.Run("Track 3 data", func(t *testing.T) {
		tracker := NewTrackThird()

		samples := []string{
			`011234567890123445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**`,
			`011234567890123445=724724100000000000030300XXXX040400099010=************************==1=0000000000000000`,
			`011234567890123445=000978100000000****8330*0000920000099010=************************==1=0000000*00000000`,
		}
		for _, sample := range samples {
			card, err := tracker.Read([]byte(sample))
			require.NoError(t, err)
			buf, err := tracker.Write(card)
			require.NoError(t, err)
			require.Equal(t, sample, string(buf))
		}
	})
}
