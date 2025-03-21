package i2c

import (
	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
)

var field54 = field.NewComposite(&field.Spec{
	Length:      120,
	Description: "Additional Amounts",
	Tag: &field.TagSpec{
		Sort: sort.StringsByInt,
	},
	Pref: prefix.ASCII.LLL,
	Subfields: map[string]field.Field{
		"01": additionalAmountField,
		"02": additionalAmountField,
		"03": additionalAmountField,
		"04": additionalAmountField,
		"05": additionalAmountField,
		"06": additionalAmountField,
	},
})

var additionalAmountField = field.NewComposite(&field.Spec{
	Length:      20,
	Description: "Additional Amount",
	Tag: &field.TagSpec{
		Sort: sort.StringsByInt,
	},
	Pref: prefix.ASCII.Fixed,
	Subfields: map[string]field.Field{
		"01": field.NewString(&field.Spec{
			Length:      2,
			Description: "Account Type",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		"02": field.NewString(&field.Spec{
			Length:      2,
			Description: "Amount Type",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		"03": field.NewString(&field.Spec{
			Length:      3,
			Description: "Currency Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		"04": field.NewString(&field.Spec{
			Length:      1,
			Description: "Amount Sign",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		"05": field.NewString(&field.Spec{
			Length:      12,
			Description: "Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
	},
})

var ITooSeaSpec *iso8583.MessageSpec = &iso8583.MessageSpec{
	Name: "You get it?",
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Length:      8,
			Description: "Bitmap",
			Enc:         encoding.Binary,
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
			Length:      70,
			Description: "Card Acceptor Name/Location",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		48: field.NewString(&field.Spec{
			Length:      255,
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
		54: field54,
		57: field.NewString(&field.Spec{
			Length:      3,
			Description: "Reserved (National)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		59: field.NewString(&field.Spec{
			Length:      18,
			Description: "Geographic Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		61: field.NewString(&field.Spec{
			Length:      19,
			Description: "Point of Service Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		63: field.NewString(&field.Spec{
			Length:      70,
			Description: "Network Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
			Tag: &field.TagSpec{
				Sort: sort.StringsByInt,
			},
			Subfields: map[string]field.Field{
				"01": field.NewString(&field.Spec{
					Length:      4,
					Description: "Acquirer Network ID",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"02": field.NewString(&field.Spec{
					Length:      4,
					Description: "Issuer Network ID",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"03": field.NewString(&field.Spec{
					Length:      16,
					Description: "Transaction Identifier/Access Transaction Sequence Number",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"04": field.NewString(&field.Spec{
					Length:      9,
					Description: "Bank Net Reference Number",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"05": field.NewString(&field.Spec{
					Length:      1,
					Description: "Interchange Rate Indicator",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"06": field.NewString(&field.Spec{
					Length:      24,
					Description: "Acquirer Reference Number",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
				"07": field.NewString(&field.Spec{
					Length:      12,
					Description: "Network Type",
					Enc:         encoding.ASCII,
					Pref:        prefix.ASCII.Fixed,
				}),
			},
		}),
		65: field.NewBitmap(&field.Spec{
			Length:      8,
			Description: "SECONDARY BITMAP DATA",
			Enc:         encoding.Binary,
			Pref:        prefix.Hex.Fixed,
		}),
		70: field.NewString(&field.Spec{
			Length:      3,
			Description: "NETWORK MANAGEMENT INFORMATION CODE",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		72: field.NewString(&field.Spec{
			Length:      999,
			Description: "TRANSACTION ADDITIONAL DATA",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		80: field.NewString(&field.Spec{
			Length:      999,
			Description: "DISPUTE ACTION INFORMATION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		90: field.NewString(&field.Spec{
			Length:      44,
			Description: "Original Data Elements",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		102: field.NewString(&field.Spec{
			Length:      28,
			Description: "ACCOUNT IDENTIFICATION 1",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		110: field.NewString(&field.Spec{
			Length:      360,
			Description: "MINI STATEMENT DATA",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		111: field.NewString(&field.Spec{
			Length:      999,
			Description: "ADDITIONAL DATA",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		123: field.NewString(&field.Spec{
			Length:      255,
			Description: "VERIFICATION DATA",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		125: field.NewString(&field.Spec{
			Length:      999,
			Description: "SUPPORTING INFORMATION",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
	},
}
