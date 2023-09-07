package iso8583_test

import (
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/specs"
	"github.com/stretchr/testify/require"
)

func TestStructWithTypes(t *testing.T) {
	type authRequestData struct {
		MTI                  string `index:"0"`
		PrimaryAccountNumber string `index:"2"`
		ProcessingCode       int    `index:"3"`
		TransactionAmount    int    `index:"4,keepzero"` // we will set message field value to 0
	}

	t.Run("pack", func(t *testing.T) {
		authRequest := &authRequestData{
			MTI:                  "0110",
			PrimaryAccountNumber: "4242424242424242",
			ProcessingCode:       200000,
		}

		message := iso8583.NewMessage(specs.Spec87ASCII)
		err := message.Marshal(authRequest)
		require.NoError(t, err)

		packed, err := message.Pack()
		require.NoError(t, err)

		require.Equal(t, "011070000000000000000000000000000000164242424242424242200000000000000000", string(packed))
	})
}
