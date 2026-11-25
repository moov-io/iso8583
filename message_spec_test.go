// SPDX-License-Identifier: Apache-2.0

package iso8583

import (
	"maps"
	"testing"

	"github.com/moov-io/iso8583/field"
	"github.com/stretchr/testify/require"
)

func TestMessageSpecValidate(t *testing.T) {
	require.NoError(t, Spec87.Validate())
	for _, tc := range []struct {
		name string
		id   int
		want string
	}{
		{"missing MTI", 0, "MTI field (0) is required"},
		{"missing bitmap", 1, "Bitmap field (1) is required"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			spec := &MessageSpec{Fields: maps.Clone(Spec87.Fields)}
			delete(spec.Fields, tc.id)
			require.ErrorContains(t, spec.Validate(), tc.want)
		})
	}
	t.Run("string instead of bitmap", func(t *testing.T) {
		spec := &MessageSpec{Fields: maps.Clone(Spec87.Fields)}
		spec.Fields[1] = field.NewStringValue("0000000000000000")
		require.ErrorContains(t, spec.Validate(), "Bitmap field (1) must be of type *field.Bitmap")
	})
}
