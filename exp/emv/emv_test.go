package emv

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/field"
	"github.com/stretchr/testify/require"
)

func TestEmv(t *testing.T) {
	exampleICCData := `9F0206000000006300820258009F360200029F2608B9B2B58202D37033840FA000000152301010000100000000009F100801050000000000009F3303E0F0C09F1A020840950500000000009A031711209C01005F2A0208409F370459F58EB1`
	rawData, err := hex.DecodeString(exampleICCData)
	require.NoError(t, err)

	// we have to add LLL length before rawData to unpack it
	lenPrefix := fmt.Sprintf("%03d", len(rawData))
	rawData = append([]byte(lenPrefix), rawData...)

	emvField := field.NewComposite(Spec)
	_, err = emvField.Unpack(rawData)

	require.NoError(t, err)

	data := &Data{}

	err = emvField.Unmarshal(data)
	require.NoError(t, err)

	require.Equal(t, 6300, data.AmountAuthorisedNumeric.Value())
	require.Equal(t, "5800", data.ApplicationInterchangeProfile.Value())
	require.Equal(t, 2, data.ApplicationTransactionCounter.Value())
	require.Equal(t, "B9B2B58202D37033", data.ApplicationCryptogram.Value())

}
