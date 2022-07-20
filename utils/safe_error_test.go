package utils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeErrorDoesNotExposeSensitiveDetails(t *testing.T) {
	sensitiveError := errors.New("1234 failed to parse")
	safeError := &SafeError{
		Msg: "bad PIN",
		Err: sensitiveError,
	}
	var err error = safeError

	errMsg := fmt.Sprintf("%s", err)
	require.Equal(t, "bad PIN", errMsg)
	require.Equal(t, "bad PIN", safeError.Error())
}

func TestSafeErrorAllowsMatching(t *testing.T) {
	sensitiveError := errors.New("1234 failed to parse")
	safeError := &SafeError{
		Msg: "bad PIN",
		Err: sensitiveError,
	}
	var err error = safeError

	require.ErrorIs(t, err, sensitiveError)
	require.Equal(t, sensitiveError, safeError.Unwrap())
}

func TestSafeErrorUnsafeError(t *testing.T) {
	sensitiveError := errors.New("1234 failed to parse")
	safeError := &SafeError{
		Msg: "bad PIN",
		Err: sensitiveError,
	}

	require.Equal(t, "bad PIN: 1234 failed to parse", safeError.UnsafeError())
}
