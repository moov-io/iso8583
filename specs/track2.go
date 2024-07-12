package specs

import (
	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
)

// we keep it for compatibility reasons
var Spec87Track2 *iso8583.MessageSpec = &iso8583.MessageSpec{
	Name: "ISO 8583 v1987 Track2",
	Fields: map[int]field.Field{
		35: field.NewTrack2(&field.Spec{
			Length:      37,
			Description: "TRACK 2 DATA",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.Binary.L,
		}),
	},
}
