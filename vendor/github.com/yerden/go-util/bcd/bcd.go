/*
Package bcd provides functions to encode byte arrays
to BCD (Binary-Coded Decimal) encoding and back.
*/
package bcd

import (
	"fmt"
)

// BCD is the configuration for Binary-Coded Decimal encoding.
type BCD struct {
	// Map of symbols to encode and decode routines.
	// Example:
	//    key 'a' -> value 0x9
	Map map[byte]byte

	// If true nibbles (4-bit part of a byte) will
	// be swapped, meaning bits 0123 will encode
	// first digit and bits 4567 will encode the
	// second.
	SwapNibbles bool

	// Filler nibble is used if the input has odd
	// number of bytes. Then the output's final nibble
	// will contain the specified nibble.
	Filler byte
}

var (
	// Standard 8-4-2-1 decimal-only encoding.
	Standard = &BCD{
		Map: map[byte]byte{
			'0': 0x0, '1': 0x1, '2': 0x2, '3': 0x3,
			'4': 0x4, '5': 0x5, '6': 0x6, '7': 0x7,
			'8': 0x8, '9': 0x9,
		},
		SwapNibbles: false,
		Filler:      0xf}

	// Excess-3 or Stibitz encoding.
	Excess3 = &BCD{
		Map: map[byte]byte{
			'0': 0x3, '1': 0x4, '2': 0x5, '3': 0x6,
			'4': 0x7, '5': 0x8, '6': 0x9, '7': 0xa,
			'8': 0xb, '9': 0xc,
		},
		SwapNibbles: false,
		Filler:      0x0}

	// TBCD (Telephony BCD) as in 3GPP TS 29.002.
	Telephony = &BCD{
		Map: map[byte]byte{
			'0': 0x0, '1': 0x1, '2': 0x2, '3': 0x3,
			'4': 0x4, '5': 0x5, '6': 0x6, '7': 0x7,
			'8': 0x8, '9': 0x9, '*': 0xa, '#': 0xb,
			'a': 0xc, 'b': 0xd, 'c': 0xe,
		},
		SwapNibbles: true,
		Filler:      0xf}

	// Aiken or 2421 code
	Aiken = &BCD{
		Map: map[byte]byte{
			'0': 0x0, '1': 0x1, '2': 0x2, '3': 0x3,
			'4': 0x4, '5': 0xb, '6': 0xc, '7': 0xd,
			'8': 0xe, '9': 0xf,
		},
		SwapNibbles: false,
		Filler:      0x5}
)

// Error values returned by API.
var (
	// ErrBadInput returned if input data cannot be encoded.
	ErrBadInput = fmt.Errorf("non-encodable data")
	// ErrBadBCD returned if input data cannot be decoded.
	ErrBadBCD = fmt.Errorf("Bad BCD data")
)
