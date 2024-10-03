package emv

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/stretchr/testify/require"
)

func TestEmv(t *testing.T) {
	iccData := `4f07a00000000410105f2a0208405f2d047275656e82021980950500000000009a032410029c01009f02060000000001009f03060000000000009f0607a00000000410109f090200029f1a0208409f1e0834543734353038359f21030909329f33030000e89f34030000009f3501229f360200539f370409bc21069f3901919f4104000000069f530100df81290830f0f00030f0ff00dfee2601d1dfef4c06002700000000dfef4d28fd4b4f392e1278361252d85649e405b430b9cb5e57d5211ba81d00242dbbd987ffe099fbf92b422eff810581ac500a4d6173746572436172648407a00000000410109f6d02000156a1292a353438392a2a2a2a2a2a2a2a313433375e202f5e323330333230312a2a2a2a2a2a2a2a2a2a2a2a2a56c1307004c1af307c59842b4edd3e745d27e783f4fefefb36698568ee921be4b2d6ae066b1a8cac8aba5c29453da922bc7ed89f6ba1135489cccccccc1437d2303201cccccccccccccc9f6bc118c5e921e707dabbfc762d47986509668fd22f9b3dbd10cbd6ff81063cdf812ac11870b80a15602160f54503ce3f5d9116997adb85fcfcd5388edf812bc110e912f0f4c0d82d1075ccc984f9864127df8115060000000000ffffee012cdf300100df31c110ab6c422bbfb95378504cbc6683306641df32c11081592240b4e352d17af39b923ad115caffee120acdcdcd0701453de0000d`
	rawData, err := hex.DecodeString(iccData)
	require.NoError(t, err)

	msg := iso8583.NewMessage(MessageSpec)
	msg.MTI("0100")
	msg.BinaryField(55, rawData)

	// this will print the all EMV tags in readable format
	// like this (note, that first F is not part of the tag, it's just a Filed prefix):
	// F55  ICC Data SUBFIELDS:
	// -------------------------------------------
	// F4F  Application Identifier (AID) – card............: A0000000041010
	// F5F2A Transaction Currency Code.....................: 0840
	// F5F2D Language Preference...........................: 7275656E
	// F82  Application Interchange Profile................: 1980
	// F95  Terminal Verification Results..................: 0000000000
	// F9A  Transaction Date...............................: 241002
	// F9C  Transaction Type...............................: 00
	// F9F02 Amount, Authorised (Numeric)..................: 100

	iso8583.Describe(msg, os.Stdout)

	// now we can extract values we can use
	iccField := msg.GetField(55)
	data := &NativeData{}
	iccField.Unmarshal(data)
}
