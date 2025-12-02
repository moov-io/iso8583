package iso8583_test

import (
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
	"github.com/stretchr/testify/require"
)

func BenchmarkMarshaling(b *testing.B) {
	b.ReportAllocs()

	b.StopTimer()
	data := getTestMessageData()

	// thes that we can Marshal without errors before starting the benchmark
	msg := iso8583.NewMessage(benchmarkSpec)
	err := msg.Marshal(data)
	require.NoError(b, err)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		msg := iso8583.NewMessage(benchmarkSpec)
		msg.Marshal(data)
	}
}

func BenchmarkUnpacking(b *testing.B) {
	b.ReportAllocs()

	b.StopTimer()
	// prepare packed message
	msg0 := iso8583.NewMessage(benchmarkSpec)
	data := getTestMessageData()
	err := msg0.Marshal(data)
	require.NoError(b, err)
	packed, err := msg0.Pack()
	require.NoError(b, err)

	// test that we can Unpack without errors before starting the benchmark
	msg := iso8583.NewMessage(benchmarkSpec)
	err = msg.Unpack(packed)
	require.NoError(b, err)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		msg := iso8583.NewMessage(benchmarkSpec)
		msg.Unpack(packed)
	}
}

var benchmarkSpec *iso8583.MessageSpec = &iso8583.MessageSpec{
	Name: "benchmark spec",
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Length:      16,
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		}),
		2: field.NewString(&field.Spec{
			Length:      19,
			Description: "Primary Account Number",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		3: field.NewNumeric(&field.Spec{
			Length:      6,
			Description: "Processing Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		4: field.NewString(&field.Spec{
			Length:      12,
			Description: "Transaction Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		5: field.NewString(&field.Spec{
			Length:      12,
			Description: "Settlement Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		6: field.NewString(&field.Spec{
			Length:      12,
			Description: "Billing Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		7: field.NewString(&field.Spec{
			Length:      10,
			Description: "Transmission Date & Time",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		8: field.NewString(&field.Spec{
			Length:      8,
			Description: "Billing Fee Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		9: field.NewString(&field.Spec{
			Length:      8,
			Description: "Settlement Conversion Rate",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		10: field.NewString(&field.Spec{
			Length:      8,
			Description: "Cardholder Billing Conversion Rate",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		11: field.NewString(&field.Spec{
			Length:      6,
			Description: "Systems Trace Audit Number (STAN)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		12: field.NewString(&field.Spec{
			Length:      6,
			Description: "Local Transaction Time",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		13: field.NewString(&field.Spec{
			Length:      4,
			Description: "Local Transaction Date",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		14: field.NewString(&field.Spec{
			Length:      4,
			Description: "Expiration Date",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		15: field.NewString(&field.Spec{
			Length:      4,
			Description: "Settlement Date",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		16: field.NewString(&field.Spec{
			Length:      4,
			Description: "Currency Conversion Date",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		17: field.NewString(&field.Spec{
			Length:      4,
			Description: "Capture Date",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		18: field.NewString(&field.Spec{
			Length:      4,
			Description: "Merchant Type",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		19: field.NewString(&field.Spec{
			Length:      3,
			Description: "Acquiring Institution Country Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		20: field.NewString(&field.Spec{
			Length:      3,
			Description: "PAN Extended Country Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		21: field.NewString(&field.Spec{
			Length:      3,
			Description: "Forwarding Institution Country Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		22: field.NewString(&field.Spec{
			Length:      3,
			Description: "Point of Sale (POS) Entry Mode",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		23: field.NewString(&field.Spec{
			Length:      3,
			Description: "Card Sequence Number (CSN)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		24: field.NewString(&field.Spec{
			Length:      3,
			Description: "Function Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		25: field.NewString(&field.Spec{
			Length:      2,
			Description: "Point of Service Condition Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		26: field.NewString(&field.Spec{
			Length:      2,
			Description: "Point of Service PIN Capture Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		27: field.NewString(&field.Spec{
			Length:      1,
			Description: "Authorizing Identification Response Length",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		28: field.NewString(&field.Spec{
			Length:      9,
			Description: "Transaction Fee Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		29: field.NewString(&field.Spec{
			Length:      9,
			Description: "Settlement Fee Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		30: field.NewString(&field.Spec{
			Length:      9,
			Description: "Transaction Processing Fee Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		31: field.NewString(&field.Spec{
			Length:      9,
			Description: "Settlement Processing Fee Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		32: field.NewString(&field.Spec{
			Length:      11,
			Description: "Acquiring Institution Identification Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		33: field.NewString(&field.Spec{
			Length:      11,
			Description: "Forwarding Institution Identification Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		34: field.NewString(&field.Spec{
			Length:      28,
			Description: "Extended Primary Account Number",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		35: field.NewString(&field.Spec{
			Length:      37,
			Description: "Track 2 Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		36: field.NewString(&field.Spec{
			Length:      104,
			Description: "Track 3 Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		37: field.NewString(&field.Spec{
			Length:      12,
			Description: "Retrieval Reference Number",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		38: field.NewString(&field.Spec{
			Length:      6,
			Description: "Authorization Identification Response",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		39: field.NewString(&field.Spec{
			Length:      2,
			Description: "Response Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		40: field.NewString(&field.Spec{
			Length:      3,
			Description: "Service Restriction Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		41: field.NewString(&field.Spec{
			Length:      8,
			Description: "Card Acceptor Terminal Identification",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		42: field.NewString(&field.Spec{
			Length:      15,
			Description: "Card Acceptor Identification Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		43: field.NewString(&field.Spec{
			Length:      40,
			Description: "Card Acceptor Name/Location",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		44: field.NewString(&field.Spec{
			Length:      99,
			Description: "Additional Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		45: field.NewString(&field.Spec{
			Length:      76,
			Description: "Track 1 Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		46: field.NewString(&field.Spec{
			Length:      999,
			Description: "Additional data (ISO)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		47: field.NewString(&field.Spec{
			Length:      999,
			Description: "Additional data (National)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		48: field.NewString(&field.Spec{
			Length:      999,
			Description: "Additional data (Private)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		49: field.NewString(&field.Spec{
			Length:      3,
			Description: "Transaction Currency Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		50: field.NewString(&field.Spec{
			Length:      3,
			Description: "Settlement Currency Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		51: field.NewString(&field.Spec{
			Length:      3,
			Description: "Cardholder Billing Currency Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		52: field.NewString(&field.Spec{
			Length:      8,
			Description: "PIN Data",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.Fixed,
		}),
		53: field.NewString(&field.Spec{
			Length:      16,
			Description: "Security Related Control Information",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		54: field.NewString(&field.Spec{
			Length:      120,
			Description: "Additional Amounts",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		55: field.NewString(&field.Spec{
			Length:      999,
			Description: "ICC Data – EMV Having Multiple Tags",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		56: field.NewString(&field.Spec{
			Length:      999,
			Description: "Reserved (ISO)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		57: field.NewString(&field.Spec{
			Length:      999,
			Description: "Reserved (National)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		58: field.NewString(&field.Spec{
			Length:      999,
			Description: "Reserved (National)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		59: field.NewString(&field.Spec{
			Length:      999,
			Description: "Reserved (National)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		60: field.NewString(&field.Spec{
			Length:      999,
			Description: "Reserved (National)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		61: field.NewString(&field.Spec{
			Length:      999,
			Description: "Reserved (Private)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		62: field.NewString(&field.Spec{
			Length:      999,
			Description: "Reserved (Private)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		63: field.NewString(&field.Spec{
			Length:      999,
			Description: "Reserved (Private)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		64: field.NewString(&field.Spec{
			Length:      8,
			Description: "Message Authentication Code (MAC)",
			Enc:         encoding.Binary,
			Pref:        prefix.Binary.Fixed,
		}),
		66: field.NewComposite(&field.Spec{
			Length:      999,
			Description: "Private Data",
			Pref:        prefix.ASCII.LLL,
			Tag: &field.TagSpec{
				Sort: sort.StringsByInt,
			},
			Subfields: map[string]field.Field{
				"1": field.NewString(&field.Spec{
					Length:      2,
					Description: "Transaction Type",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"2": field.NewComposite(&field.Spec{
					Length:      999,
					Description: "Additional Data",
					Pref:        prefix.ASCII.LLL,
					Tag: &field.TagSpec{
						Sort: sort.StringsByInt,
					},
					Subfields: map[string]field.Field{
						"1": field.NewString(&field.Spec{
							Length:      4,
							Description: "Date",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
						"2": field.NewNumeric(&field.Spec{
							Length:      6,
							Description: "Time",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
						"3": field.NewString(&field.Spec{
							Length:      3,
							Description: "Reserved",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
						"4": field.NewString(&field.Spec{
							Length:      5,
							Description: "Reserved",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
						"5": field.NewString(&field.Spec{
							Length:      3,
							Description: "Reserved",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
						"6": field.NewString(&field.Spec{
							Length:      3,
							Description: "Reserved",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
						"7": field.NewString(&field.Spec{
							Length:      3,
							Description: "Reserved",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
						"8": field.NewString(&field.Spec{
							Length:      3,
							Description: "Reserved",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
						"9": field.NewString(&field.Spec{
							Length:      3,
							Description: "Reserved",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
						"10": field.NewString(&field.Spec{
							Length:      3,
							Description: "Reserved",
							Enc:         encoding.ASCII,
							Pref:        prefix.ASCII.Fixed,
						}),
					},
				}),
				"3": field.NewString(&field.Spec{
					Length:      3,
					Description: "Reserved",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"4": field.NewString(&field.Spec{
					Length:      5,
					Description: "Reserved",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"5": field.NewString(&field.Spec{
					Length:      3,
					Description: "Reserved",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"6": field.NewString(&field.Spec{
					Length:      3,
					Description: "Reserved",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"7": field.NewString(&field.Spec{
					Length:      3,
					Description: "Reserved",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"8": field.NewString(&field.Spec{
					Length:      3,
					Description: "Reserved",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"9": field.NewString(&field.Spec{
					Length:      3,
					Description: "Reserved",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"10": field.NewString(&field.Spec{
					Length:      3,
					Description: "Reserved",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
			},
		}),
		90: field.NewString(&field.Spec{
			Length:      42,
			Description: "Original Data Elements",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
	},
}

type additionalData struct {
	Date       string `iso8583:"1"`
	Time       string `iso8583:"2"`
	Reserved3  string `iso8583:"3"`
	Reserved4  string `iso8583:"4"`
	Reserved5  string `iso8583:"5"`
	Reserved6  string `iso8583:"6"`
	Reserved7  string `iso8583:"7"`
	Reserved8  string `iso8583:"8"`
	Reserved9  string `iso8583:"9"`
	Reserved10 string `iso8583:"10"`
}

type privateData struct {
	TransactionType string          `iso8583:"1"`
	AdditionalData  *additionalData `iso8583:"2"`
	Reserved3       string          `iso8583:"3"`
	Reserved4       string          `iso8583:"4"`
	Reserved5       string          `iso8583:"5"`
	Reserved6       string          `iso8583:"6"`
	Reserved7       string          `iso8583:"7"`
	Reserved8       string          `iso8583:"8"`
	Reserved9       string          `iso8583:"9"`
	Reserved10      string          `iso8583:"10"`
}

type messageData struct {
	MTI                        string       `iso8583:"0"`
	PrimaryAccountNumber       string       `iso8583:"2"`
	ProcessingCode             string       `iso8583:"3"`
	TransactionAmount          string       `iso8583:"4"`
	SettlementAmount           string       `iso8583:"5"`
	BillingAmount              string       `iso8583:"6"`
	TransmissionDateTime       string       `iso8583:"7"`
	BillingFeeAmount           string       `iso8583:"8"`
	SettlementConversionRate   string       `iso8583:"9"`
	CardholderBillingConvRate  string       `iso8583:"10"`
	SystemsTraceAuditNumber    string       `iso8583:"11"`
	LocalTransactionTime       string       `iso8583:"12"`
	LocalTransactionDate       string       `iso8583:"13"`
	ExpirationDate             string       `iso8583:"14"`
	SettlementDate             string       `iso8583:"15"`
	CurrencyConversionDate     string       `iso8583:"16"`
	CaptureDate                string       `iso8583:"17"`
	MerchantType               string       `iso8583:"18"`
	AcquiringInstCountryCode   string       `iso8583:"19"`
	PANExtendedCountryCode     string       `iso8583:"20"`
	ForwardingInstCountryCode  string       `iso8583:"21"`
	POSEntryMode               string       `iso8583:"22"`
	CardSequenceNumber         string       `iso8583:"23"`
	FunctionCode               string       `iso8583:"24"`
	POSConditionCode           string       `iso8583:"25"`
	POSPINCaptureCode          string       `iso8583:"26"`
	AuthorizingIDRespLength    string       `iso8583:"27"`
	TransactionFeeAmount       string       `iso8583:"28"`
	SettlementFeeAmount        string       `iso8583:"29"`
	TransactionProcessingFee   string       `iso8583:"30"`
	SettlementProcessingFee    string       `iso8583:"31"`
	AcquiringInstIDCode        string       `iso8583:"32"`
	ForwardingInstIDCode       string       `iso8583:"33"`
	ExtendedPrimaryAccountNum  string       `iso8583:"34"`
	Track2Data                 string       `iso8583:"35"`
	Track3Data                 string       `iso8583:"36"`
	RetrievalReferenceNumber   string       `iso8583:"37"`
	AuthorizationIDResponse    string       `iso8583:"38"`
	ResponseCode               string       `iso8583:"39"`
	ServiceRestrictionCode     string       `iso8583:"40"`
	CardAcceptorTerminalID     string       `iso8583:"41"`
	CardAcceptorIdentification string       `iso8583:"42"`
	CardAcceptorNameLocation   string       `iso8583:"43"`
	AdditionalData             string       `iso8583:"44"`
	Track1Data                 string       `iso8583:"45"`
	AdditionalDataISO          string       `iso8583:"46"`
	AdditionalDataNational     string       `iso8583:"47"`
	AdditionalDataPrivate      string       `iso8583:"48"`
	TransactionCurrencyCode    string       `iso8583:"49"`
	SettlementCurrencyCode     string       `iso8583:"50"`
	CardholderBillingCurrCode  string       `iso8583:"51"`
	PINData                    string       `iso8583:"52"`
	SecurityRelatedCtrlInfo    string       `iso8583:"53"`
	AdditionalAmounts          string       `iso8583:"54"`
	ICCDataEMV                 string       `iso8583:"55"`
	ReservedISO                string       `iso8583:"56"`
	ReservedNational           string       `iso8583:"57"`
	ReservedNational2          string       `iso8583:"58"`
	ReservedNational3          string       `iso8583:"59"`
	ReservedNational4          string       `iso8583:"60"`
	ReservedPrivate            string       `iso8583:"61"`
	ReservedPrivate2           string       `iso8583:"62"`
	ReservedPrivate3           string       `iso8583:"63"`
	MessageAuthenticationCode  string       `iso8583:"64"`
	PrivateData                *privateData `iso8583:"66"`
}

func getTestMessageData() *messageData {
	return &messageData{
		MTI:                        "0200",
		PrimaryAccountNumber:       "4111111111111111",                                                                                     // 16 digits (max 19)
		ProcessingCode:             "000000",                                                                                               // 6 digits
		TransactionAmount:          "000000010000",                                                                                         // 12 chars
		SettlementAmount:           "000000010000",                                                                                         // 12 chars
		BillingAmount:              "000000010000",                                                                                         // 12 chars
		TransmissionDateTime:       "0701123045",                                                                                           // 10 chars (MMDDhhmmss)
		BillingFeeAmount:           "00000100",                                                                                             // 8 chars
		SettlementConversionRate:   "12345678",                                                                                             // 8 chars
		CardholderBillingConvRate:  "12345678",                                                                                             // 8 chars
		SystemsTraceAuditNumber:    "123456",                                                                                               // 6 chars
		LocalTransactionTime:       "123045",                                                                                               // 6 chars (hhmmss)
		LocalTransactionDate:       "0701",                                                                                                 // 4 chars (MMDD)
		ExpirationDate:             "2512",                                                                                                 // 4 chars (YYMM)
		SettlementDate:             "0701",                                                                                                 // 4 chars
		CurrencyConversionDate:     "0701",                                                                                                 // 4 chars
		CaptureDate:                "0701",                                                                                                 // 4 chars
		MerchantType:               "5411",                                                                                                 // 4 chars
		AcquiringInstCountryCode:   "840",                                                                                                  // 3 chars (USA)
		PANExtendedCountryCode:     "840",                                                                                                  // 3 chars
		ForwardingInstCountryCode:  "840",                                                                                                  // 3 chars
		POSEntryMode:               "051",                                                                                                  // 3 chars
		CardSequenceNumber:         "001",                                                                                                  // 3 chars
		FunctionCode:               "200",                                                                                                  // 3 chars
		POSConditionCode:           "00",                                                                                                   // 2 chars
		POSPINCaptureCode:          "12",                                                                                                   // 2 chars
		AuthorizingIDRespLength:    "6",                                                                                                    // 1 char
		TransactionFeeAmount:       "000000100",                                                                                            // 9 chars
		SettlementFeeAmount:        "000000100",                                                                                            // 9 chars
		TransactionProcessingFee:   "000000100",                                                                                            // 9 chars
		SettlementProcessingFee:    "000000100",                                                                                            // 9 chars
		AcquiringInstIDCode:        "12345678901",                                                                                          // 11 chars max
		ForwardingInstIDCode:       "12345678901",                                                                                          // 11 chars max
		ExtendedPrimaryAccountNum:  "1234567890123456789012345678",                                                                         // 28 chars max
		Track2Data:                 "4111111111111111=25121011234567890",                                                                   // 37 chars max
		Track3Data:                 "0104111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111", // 104 chars max
		RetrievalReferenceNumber:   "123456789012",                                                                                         // 12 chars
		AuthorizationIDResponse:    "ABC123",                                                                                               // 6 chars
		ResponseCode:               "00",                                                                                                   // 2 chars
		ServiceRestrictionCode:     "000",                                                                                                  // 3 chars
		CardAcceptorTerminalID:     "TERM0001",                                                                                             // 8 chars
		CardAcceptorIdentification: "MERCHANT0000012",                                                                                      // 15 chars
		CardAcceptorNameLocation:   "Test Store    123 Main St    City  US840",                                                             // 40 chars
		AdditionalData:             "Additional response data field 44",                                                                    // up to 99 chars
		Track1Data:                 "B4111111111111111^CARDHOLDER/TEST^2512101123456789",                                                   // up to 76 chars
		AdditionalDataISO:          "ISO additional data field 46",                                                                         // up to 999 chars
		AdditionalDataNational:     "National additional data field 47",                                                                    // up to 999 chars
		AdditionalDataPrivate:      "Private additional data field 48",                                                                     // up to 999 chars
		TransactionCurrencyCode:    "840",                                                                                                  // 3 chars (USD)
		SettlementCurrencyCode:     "840",                                                                                                  // 3 chars
		CardholderBillingCurrCode:  "840",                                                                                                  // 3 chars
		PINData:                    "12345678",                                                                                             // 8 chars (hex)
		SecurityRelatedCtrlInfo:    "1234567890123456",                                                                                     // 16 chars
		AdditionalAmounts:          "Additional amounts data",                                                                              // up to 120 chars
		ICCDataEMV:                 "EMV chip data",                                                                                        // up to 999 chars
		ReservedISO:                "Reserved ISO",                                                                                         // up to 999 chars
		ReservedNational:           "Reserved National",                                                                                    // up to 999 chars
		ReservedNational2:          "Reserved National 2",                                                                                  // up to 999 chars
		ReservedNational3:          "Reserved National 3",                                                                                  // up to 999 chars
		ReservedNational4:          "Reserved National 4",                                                                                  // up to 999 chars
		ReservedPrivate:            "Reserved Private",                                                                                     // up to 999 chars
		ReservedPrivate2:           "Reserved Private 2",                                                                                   // up to 999 chars
		ReservedPrivate3:           "Reserved Private 3",                                                                                   // up to 999 chars
		MessageAuthenticationCode:  "12345678",                                                                                             // 8 chars (hex)
		PrivateData: &privateData{
			TransactionType: "01", // 2 chars
			AdditionalData: &additionalData{
				Date:       "0701",   // 4 chars
				Time:       "123045", // 6 chars
				Reserved3:  "RES",    // 3 chars
				Reserved4:  "RESV4",  // 5 chars
				Reserved5:  "RS5",    // 3 chars
				Reserved6:  "RS6",    // 3 chars
				Reserved7:  "RS7",    // 3 chars
				Reserved8:  "RS8",    // 3 chars
				Reserved9:  "RS9",    // 3 chars
				Reserved10: "R10",    // 3 chars
			},
			Reserved3:  "RV3",   // 3 chars
			Reserved4:  "RESV4", // 5 chars
			Reserved5:  "RV5",   // 3 chars
			Reserved6:  "RV6",   // 3 chars
			Reserved7:  "RV7",   // 3 chars
			Reserved8:  "RV8",   // 3 chars
			Reserved9:  "RV9",   // 3 chars
			Reserved10: "R10",   // 3 chars
		},
	}
}
