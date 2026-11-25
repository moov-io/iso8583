// SPDX-License-Identifier: Apache-2.0

package errors_test

import (
	"errors"
	"testing"

	iso8583errors "github.com/moov-io/iso8583/errors"
	"github.com/stretchr/testify/require"
)

func TestPackError(t *testing.T) {
	cause := errors.New("test packing error")
	err := &iso8583errors.PackError{Err: cause}
	require.EqualError(t, err, cause.Error())
	require.Same(t, cause, errors.Unwrap(err))
	require.ErrorIs(t, err, cause)
}
