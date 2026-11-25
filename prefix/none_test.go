// SPDX-License-Identifier: Apache-2.0

package prefix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoneFixed(t *testing.T) {
	data := []byte{0x01, 0x00, 0xFF}
	encoded, err := None.Fixed.EncodeLength(2, len(data))
	require.NoError(t, err)
	require.Empty(t, encoded)

	length, read, err := None.Fixed.DecodeLength(2, data)
	require.NoError(t, err)
	require.Equal(t, len(data), length)
	require.Zero(t, read)
}
