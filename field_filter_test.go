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

	out := bytes.NewBuffer([]byte{})
	require.NotPanics(t, func() {
		Describe(message, out)
	})

	expectedOutput := `ISO 8583 v1987 ASCII Message:
MTI..........: 0100
Bitmap HEX...: 0000000000000000
Bitmap bits..:
[1-8]00000000 [9-16]00000000 [17-24]00000000 [25-32]00000000
[33-40]00000000 [41-48]00000000 [49-56]00000000 [57-64]00000000
F0   Message Type Indicator...............: 0100
F2   Primary Account Number...............: 4242****4242
F3   Processing Code......................: 123456
F4   Transaction Amount...................: 100
F20  PAN Extended Country Code............: 4242****4242
F35  Track 2 Data.........................: 4000****0506=2512111123400001230
F36  Track 3 Data.........................: 011234****3445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**
F45  Track 1 Data.........................: B4815****1896^YATES/EUGENE L^^^356858      00998000000
F52  PIN Data.............................: 12****78
F55  ICC Data – EMV Having Multiple Tags..: ICC  ... Tags
`

	require.Equal(t, expectedOutput, out.String())

	out.Reset()
	require.NotPanics(t, func() {
		Describe(message, out, DoNotFilterFields()...)
	})

	expectedOutput = `ISO 8583 v1987 ASCII Message:
MTI..........: 0100
Bitmap HEX...: 0000000000000000
Bitmap bits..:
[1-8]00000000 [9-16]00000000 [17-24]00000000 [25-32]00000000
[33-40]00000000 [41-48]00000000 [49-56]00000000 [57-64]00000000
F0   Message Type Indicator...............: 0100
F2   Primary Account Number...............: 4242424242424242
F3   Processing Code......................: 123456
F4   Transaction Amount...................: 100
F20  PAN Extended Country Code............: 4242424242424242
F35  Track 2 Data.........................: 4000340000000506=2512111123400001230
F36  Track 3 Data.........................: 011234567890123445=724724000000000****00300XXXX020200099010=********************==1=100000000000000000**
F45  Track 1 Data.........................: B4815881002861896^YATES/EUGENE L            ^^^356858      00998000000
F52  PIN Data.............................: 12345678
F55  ICC Data – EMV Having Multiple Tags..: ICC Data – EMV Having Multiple Tags
`

	require.Equal(t, expectedOutput, out.String())
}
