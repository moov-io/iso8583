
package padding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNonePadder(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
		{"with zero byte", []byte{'A', 0}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.data, None.Pad(tc.data, 4))
			require.Equal(t, tc.data, None.Unpad(tc.data))
		})
	}
	require.Nil(t, None.Inspect())
}
