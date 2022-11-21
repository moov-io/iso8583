package iso8583

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldFilter(t *testing.T) {

	spec := Spec87

	message := NewMessage(spec)

	message.MTI("0100")
	err := message.Field(2, "4242424242424242")
	require.NoError(t, err)

	err = message.Field(3, "123456")
	require.NoError(t, err)

	err = message.Field(4, "100")
	require.NoError(t, err)

	err = message.Field(20, "4242424242424242")
	require.NoError(t, err)

	err = message.Field(35, "4000340000000506=2512111123400001230")
	require.NoError(t, err)

	err = message.Field(36, "011234567890123445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**")
	require.NoError(t, err)

	err = message.Field(45, "B4815881002861896^YATES/EUGENE L            ^^^356858      00998000000")
	require.NoError(t, err)

	err = message.Field(52, "12345678")
	require.NoError(t, err)

	err = message.Field(55, "ICC Data – EMV Having Multiple Tags")
	require.NoError(t, err)

	filters := DefaultFilter()

	out := bytes.NewBuffer([]byte{})
	require.NotPanics(t, func() {
		Describe(message, out, filters...)
	})

	expectedOutput := `ISO 8583 v1987 ASCII Message:
MTI........................................: 0100
Bitmap.....................................: 000000000000000000000000000000000000000000000000
Bitmap bits................................: 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
F000 Message Type Indicator................: 0100
F002 Primary Account Number................: 4242****4242
F003 Processing Code.......................: 123456
F004 Transaction Amount....................: 100
F020 PAN Extended Country Code.............: 4242****4242
F035 Track 2 Data..........................: 4000****0506=2512111123400001230
F036 Track 3 Data..........................: 011234****3445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**
F045 Track 1 Data..........................: B4815****1896^YATES/EUGENE L^^^356858      00998000000
F052 PIN Data..............................: 12****78
F055 ICC Data – EMV Having Multiple Tags...: ICC  ... Tags
`
	require.Equal(t, expectedOutput, out.String())

	out.Reset()
	require.NotPanics(t, func() {
		Describe(message, out)
	})

	expectedOutput = `ISO 8583 v1987 ASCII Message:
MTI........................................: 0100
Bitmap.....................................: 000000000000000000000000000000000000000000000000
Bitmap bits................................: 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
F000 Message Type Indicator................: 0100
F002 Primary Account Number................: 4242424242424242
F003 Processing Code.......................: 123456
F004 Transaction Amount....................: 100
F020 PAN Extended Country Code.............: 4242424242424242
F035 Track 2 Data..........................: 4000340000000506=2512111123400001230
F036 Track 3 Data..........................: 011234567890123445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**
F045 Track 1 Data..........................: B4815881002861896^YATES/EUGENE L            ^^^356858      00998000000
F052 PIN Data..............................: 12345678
F055 ICC Data – EMV Having Multiple Tags...: ICC Data – EMV Having Multiple Tags
`
	require.Equal(t, expectedOutput, out.String())

}
