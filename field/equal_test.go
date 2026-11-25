package field

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b Field
		want bool
	}{
		{"both nil", nil, nil, true},
		{"typed nil and nil", (*String)(nil), nil, true},
		{"nil and value", nil, NewStringValue("a"), false},
		{"same string", NewStringValue("abc"), NewStringValue("abc"), true},
		{"different string", NewStringValue("abc"), NewStringValue("abd"), false},
		{"same number", NewNumericValue(42), NewNumericValue(42), true},
		{"different number", NewNumericValue(42), NewNumericValue(43), false},
		{"same binary", NewBinaryValue([]byte{0x01, 0x02}), NewBinaryValue([]byte{0x01, 0x02}), true},
		{"values are compared, not types", NewStringValue("123"), NewNumericValue(123), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Equal(tt.a, tt.b)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)

			got, err = Equal(tt.b, tt.a)
			require.NoError(t, err)
			require.Equal(t, tt.want, got, "Equal must be symmetric")
		})
	}
}

func TestEqualComposite(t *testing.T) {
	newComposite := func(t *testing.T, values map[string]string) *Composite {
		t.Helper()
		c := NewComposite(compositeTestSpec)
		for path, v := range values {
			require.NoError(t, c.MarshalPath(path, v))
		}

		return c
	}

	a := newComposite(t, map[string]string{"1": "AB", "2": "CD"})

	t.Run("same subfields", func(t *testing.T) {
		got, err := Equal(a, newComposite(t, map[string]string{"1": "AB", "2": "CD"}))
		require.NoError(t, err)
		require.True(t, got)
	})

	t.Run("different subfield value", func(t *testing.T) {
		got, err := Equal(a, newComposite(t, map[string]string{"1": "AB", "2": "CE"}))
		require.NoError(t, err)
		require.False(t, got)
	})

	t.Run("subfield set in only one", func(t *testing.T) {
		got, err := Equal(a, newComposite(t, map[string]string{"1": "AB"}))
		require.NoError(t, err)
		require.False(t, got)
	})

	t.Run("composite and non-composite", func(t *testing.T) {
		got, err := Equal(a, NewStringValue("ABCD"))
		require.NoError(t, err)
		require.False(t, got)
	})
}
