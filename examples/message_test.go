package examples

import (
	"os"

	"github.com/moov-io/iso8583"
)

// Define types for the message fields
type Authorization struct {
	MTI                  string               `iso8583:"0"`  // Message Type Indicator
	PrimaryAccountNumber string               `iso8583:"2"`  // PAN
	ProcessingCode       string               `iso8583:"3"`  // Processing code
	Amount               int64                `iso8583:"4"`  // Transaction amount
	STAN                 string               `iso8583:"11"` // System Trace Audit Number
	ExpirationDate       string               `iso8583:"14"` // YYMM
	AcceptorInformation  *AcceptorInformation `iso8583:"43"` // Merchant details
}

type AcceptorInformation struct {
	Name    string `index:"1"`
	City    string `index:"2"`
	Country string `index:"3"`
}

func Example() {
	// Pack the message
	msg := iso8583.NewMessage(Spec)

	authData := &Authorization{
		MTI:                  "0100",
		PrimaryAccountNumber: "4242424242424242",
		ProcessingCode:       "000000",
		Amount:               2599,
		ExpirationDate:       "2201",
		AcceptorInformation: &AcceptorInformation{
			Name:    "Merchant Name",
			City:    "Denver",
			Country: "US",
		},
	}

	// Set the field values
	err := msg.Marshal(authData)
	if err != nil {
		panic(err)
	}

	// Pack the message
	packed, err := msg.Pack()
	if err != nil {
		panic(err)
	}

	// send packed message to the server
	// ...

	// Unpack the message
	msg = iso8583.NewMessage(Spec)
	err = msg.Unpack(packed)
	if err != nil {
		panic(err)
	}

	// Get the field values
	authData = &Authorization{}
	err = msg.Unmarshal(authData)
	if err != nil {
		panic(err)
	}

	// Print the entire message
	iso8583.Describe(msg, os.Stdout)
	//Output:
	// ISO 8583 v1987 ASCII Message:
	// MTI..........: 0100
	// Bitmap HEX...: 70040000002000000000000000000000
	// Bitmap bits..:
	//     [1-8]01110000    [9-16]00000100   [17-24]00000000   [25-32]00000000
	//   [33-40]00000000   [41-48]00100000   [49-56]00000000   [57-64]00000000
	//   [65-72]00000000   [73-80]00000000   [81-88]00000000   [89-96]00000000
	//  [97-104]00000000 [105-112]00000000 [113-120]00000000 [121-128]00000000
	// F0   Message Type Indicator..: 0100
	// F2   Primary Account Number..: 4242****4242
	// F3   Processing Code.........: 0
	// F4   Transaction Amount......: 2599
	// F14  Expiration Date.........: 2201
	// F43  Card Acceptor Name/Location SUBFIELDS:
	// -------------------------------------------
	// F1   Card Acceptor Name..........: Merchant Name
	// F2   Card Acceptor City..........: Denver
	// F3   Card Acceptor Country Code..: US
	// ------------------------------------------

}
