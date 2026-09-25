// Package encoding_test demonstrates PackedBCDHex on the field it was
// created for: VisaNet Track 2 Data (DE35). All values are synthetic
// test data (PAN from the 4000 0000 0000 0002 test range); no real card
// or track data appears anywhere in this file.
package encoding_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

// The wire form of VisaNet DE35, per the VisaNet Authorization Only
// Online Messages specification:
//
//	Field 35 - Variable length 1 byte, binary + 37 N,
//	             4-bit BCD (unsigned packed); maximum 20 bytes
//
// The 1-byte binary LL prefix counts track 2 DIGITS; the digits that
// follow are unsigned packed BCD, so the data occupies ceil(LL/2)
// bytes and carries nibbles outside 0-9: the 'D' field separator and
// the 'F' padding of the discretionary data.
func TestVisaNetTrack2DE35(t *testing.T) {
	// synthetic track 2: test PAN, 'D' separator, YYMM expiry, service
	// code, 'F'-padded synthetic discretionary data (37 digits)
	track2 := "4000000000000002D26122211234567890FFF"
	require.Len(t, track2, 37)

	// LL prefix holds the DIGIT count (37), then 19 packed bytes
	spec := &field.Spec{
		Length:      37,
		Description: "Track 2 Data",
		Enc:         encoding.PackedBCDHex,
		Pref:        prefix.Binary.L,
	}

	// wire: 0x25 (=37 digits) + 19 packed BCD bytes = 20 bytes
	packedData, err := encoding.PackedBCDHex.Encode([]byte(track2))
	require.NoError(t, err)
	require.Len(t, packedData, 19)
	wire := append([]byte{0x25}, packedData...)

	// The decimal-only BCD encoder cannot unpack this field: the 'D'
	// separator and 'F' padding nibbles are undecodable.
	_, _, err = encoding.BCD.Decode(packedData, 37)
	require.Error(t, err)
	require.EqualError(t, err, "failed to perform BCD decoding")

	// The hex encoding moov previously suggested for Track2 reads LL
	// BYTES (37) instead of ceil(LL/2) (19), the 18-byte over-read that
	// misaligns every later field of a VisaNet message. Show it against
	// a buffer that continues with the next fields (EBCDIC spaces):
	withFollowing := append(append([]byte{}, packedData...), bytes.Repeat([]byte{0x40}, 18)...)
	_, read, err := encoding.ASCIIHexToBytes.Decode(withFollowing, 37)
	require.NoError(t, err)
	require.Equal(t, 37, read)

	// PackedBCDHex consumes ceil(LL/2) bytes and keeps the nibbles:
	// the 1-byte LL plus 19 data bytes.
	f := field.NewTrack2(spec)
	consumed, err := f.Unpack(wire)
	require.NoError(t, err)
	require.Equal(t, 20, consumed)

	require.Equal(t, "4000000000000002", f.PrimaryAccountNumber)
	require.Equal(t, "D", f.Separator)
	require.NotNil(t, f.ExpirationDate)
	require.Equal(t, 2026, f.ExpirationDate.Year())
	require.Equal(t, time.December, f.ExpirationDate.Month())
	require.Equal(t, "221", f.ServiceCode)
	require.Equal(t, "1234567890FFF", f.DiscretionaryData)

	// repacking reproduces the wire bytes exactly
	packed, err := f.Pack()
	require.NoError(t, err)
	require.Equal(t, wire, packed)
}
