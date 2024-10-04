package emv

import (
	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
)

var (
	MessageSpec = &iso8583.MessageSpec{
		Fields: map[int]field.Field{
			0: field.NewString(&field.Spec{
				Length:      4,
				Description: "Message Type Indicator",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			1: field.NewBitmap(&field.Spec{
				Description: "Bitmap",
				Enc:         encoding.Binary,
				Pref:        prefix.Binary.Fixed,
			}),
			55: field.NewComposite(Spec),
		},
	}

	Spec = &field.Spec{
		Length:      999,
		Description: "ICC Data",
		Pref:        prefix.ASCII.LLL,
		Tag: &field.TagSpec{
			Sort:               sort.StringsByHex,
			Enc:                encoding.BerTLVTag,
			SkipUnknownTLVTags: true,
		},
		Subfields: map[string]field.Field{
			"9F01": field.NewHex(&field.Spec{
				Description: "Acquirer Identifier",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F40": field.NewHex(&field.Spec{
				Description: "Additional Terminal Capabilities",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"81": field.NewHex(&field.Spec{
				Description: "Amount, Authorised (Binary)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F02": field.NewHex(&field.Spec{
				Description: "Amount, Authorised (Numeric)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F04": field.NewHex(&field.Spec{
				Description: "Amount, Other (Binary)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F03": field.NewHex(&field.Spec{
				Description: "Amount, Other (Numeric)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F3A": field.NewHex(&field.Spec{
				Description: "Amount, Reference Currency",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F26": field.NewHex(&field.Spec{
				Description: "Application Cryptogram",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F42": field.NewHex(&field.Spec{
				Description: "Application Currency Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F44": field.NewHex(&field.Spec{
				Description: "Application Currency Exponent",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F05": field.NewHex(&field.Spec{
				Description: "Application Discretionary Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F25": field.NewHex(&field.Spec{
				Description: "Application Effective Date",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F24": field.NewHex(&field.Spec{
				Description: "Application Expiration Date",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"94": field.NewHex(&field.Spec{
				Description: "Application File Locator (AFL)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"4F": field.NewHex(&field.Spec{
				Description: "Application Identifier (AID) – card",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F06": field.NewHex(&field.Spec{
				Description: "Application Identifier (AID) – terminal",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"82": field.NewHex(&field.Spec{
				Description: "Application Interchange Profile",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"50": field.NewHex(&field.Spec{
				Description: "Application Label",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F12": field.NewHex(&field.Spec{
				Description: "Application Preferred Name",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5A": field.NewHex(&field.Spec{
				Description: "Application Primary Account Number (PAN)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F34": field.NewHex(&field.Spec{
				Description: "Application Primary Account Number (PAN) Sequence Number",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"87": field.NewHex(&field.Spec{
				Description: "Application Priority Indicator",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F3B": field.NewHex(&field.Spec{
				Description: "Application Reference Currency",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F43": field.NewHex(&field.Spec{
				Description: "Application Reference Currency Exponent",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F0A": field.NewHex(&field.Spec{
				Description: "Application Selection Registered Proprietary Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"61": field.NewHex(&field.Spec{
				Description: "Application Template",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F36": field.NewHex(&field.Spec{
				Description: "Application Transaction Counter",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F07": field.NewHex(&field.Spec{
				Description: "Application Usage Control",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F08": field.NewHex(&field.Spec{
				Description: "Application Version Number ICC",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F09": field.NewHex(&field.Spec{
				Description: "Application Version Number Terminal",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"89": field.NewHex(&field.Spec{
				Description: "Authorisation Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"8A": field.NewHex(&field.Spec{
				Description: "Authorisation Response Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F54": field.NewHex(&field.Spec{
				Description: "Bank Identifier Code (BIC)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F31": field.NewHex(&field.Spec{
				Description: "Card BIT Group Template",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"8C": field.NewHex(&field.Spec{
				Description: "Card Risk Management Data Object List 1 (CDOL1)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"8D": field.NewHex(&field.Spec{
				Description: "Card Risk Management Data Object List 2 (CDOL2)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F20": field.NewHex(&field.Spec{
				Description: "Cardholder Name",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F0B": field.NewHex(&field.Spec{
				Description: "Cardholder Name Extended",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"8E": field.NewHex(&field.Spec{
				Description: "Cardholder Verification Method (CVM) List",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F34": field.NewHex(&field.Spec{
				Description: "Cardholder Verification Method (CVM) Results",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"8F": field.NewHex(&field.Spec{
				Description: "Certification Authority Public Key Index ICC",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F22": field.NewHex(&field.Spec{
				Description: "Certification Authority Public Key Index Terminal",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"83": field.NewHex(&field.Spec{
				Description: "Command Template",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F27": field.NewHex(&field.Spec{
				Description: "Cryptogram Information Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F45": field.NewHex(&field.Spec{
				Description: "Data Authentication Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"84": field.NewHex(&field.Spec{
				Description: "Dedicated File (DF) Name",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9D": field.NewHex(&field.Spec{
				Description: "Directory Definition File (DDF) Name",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"73": field.NewHex(&field.Spec{
				Description: "Directory Discretionary Template",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F49": field.NewHex(&field.Spec{
				Description: "Dynamic Data Authentication Data Object List (DDOL)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"70": field.NewHex(&field.Spec{
				Description: "EMV Proprietary Template",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"DF50": field.NewHex(&field.Spec{
				Description: "Facial Try Counter",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"BF0C": field.NewHex(&field.Spec{
				Description: "File Control Information (FCI) Issuer Discretionary Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"A5": field.NewHex(&field.Spec{
				Description: "File Control Information (FCI) Proprietary Template",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"6F": field.NewHex(&field.Spec{
				Description: "File Control Information (FCI) Template",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"DF51": field.NewHex(&field.Spec{
				Description: "Finger Try Counter",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F4C": field.NewHex(&field.Spec{
				Description: "ICC Dynamic Number",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F2D": field.NewHex(&field.Spec{
				Description: "Integrated Circuit Card (ICC) PIN Encipherment Public Key Certificate",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F2E": field.NewHex(&field.Spec{
				Description: "Integrated Circuit Card (ICC) PIN Encipherment Public Key Exponent",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F2F": field.NewHex(&field.Spec{
				Description: "Integrated Circuit Card (ICC) PIN Encipherment Public Key Remainder",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F46": field.NewHex(&field.Spec{
				Description: "Integrated Circuit Card (ICC) Public Key Certificate",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F47": field.NewHex(&field.Spec{
				Description: "Integrated Circuit Card (ICC) Public Key Exponent",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F48": field.NewHex(&field.Spec{
				Description: "Integrated Circuit Card (ICC) Public Key Remainder",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F1E": field.NewHex(&field.Spec{
				Description: "Interface Device (IFD) Serial Number",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F53": field.NewHex(&field.Spec{
				Description: "International Bank Account Number (IBAN)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F0D": field.NewHex(&field.Spec{
				Description: "Issuer Action Code – Default",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F0E": field.NewHex(&field.Spec{
				Description: "Issuer Action Code – Denial",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F0F": field.NewHex(&field.Spec{
				Description: "Issuer Action Code – Online",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F10": field.NewHex(&field.Spec{
				Description: "Issuer Application Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"91": field.NewHex(&field.Spec{
				Description: "Issuer Authentication Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F11": field.NewHex(&field.Spec{
				Description: "Issuer Code Table Index",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F28": field.NewHex(&field.Spec{
				Description: "Issuer Country Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F55": field.NewHex(&field.Spec{
				Description: "Issuer Country Code (alpha2 format)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F56": field.NewHex(&field.Spec{
				Description: "Issuer Country Code (alpha3 format)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"42": field.NewHex(&field.Spec{
				Description: "Issuer Identification Number (IIN)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F0C": field.NewHex(&field.Spec{
				Description: "Issuer Identification Number Extended",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"90": field.NewHex(&field.Spec{
				Description: "Issuer Public Key Certificate",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F32": field.NewHex(&field.Spec{
				Description: "Issuer Public Key Exponent",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"92": field.NewHex(&field.Spec{
				Description: "Issuer Public Key Remainder",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"86": field.NewHex(&field.Spec{
				Description: "Issuer Script Command",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F18": field.NewHex(&field.Spec{
				Description: "Issuer Script Identifier",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"71": field.NewHex(&field.Spec{
				Description: "Issuer Script Template 1",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"72": field.NewHex(&field.Spec{
				Description: "Issuer Script Template 2",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F50": field.NewHex(&field.Spec{
				Description: "Issuer URL",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F2D": field.NewHex(&field.Spec{
				Description: "Language Preference",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F13": field.NewHex(&field.Spec{
				Description: "Last Online Application Transaction Counter (ATC) Register",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F4D": field.NewHex(&field.Spec{
				Description: "Log Entry",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F4F": field.NewHex(&field.Spec{
				Description: "Log Format",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F14": field.NewHex(&field.Spec{
				Description: "Lower Consecutive Offline Limit",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F15": field.NewHex(&field.Spec{
				Description: "Merchant Category Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F16": field.NewHex(&field.Spec{
				Description: "Merchant Identifier",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F4E": field.NewHex(&field.Spec{
				Description: "Merchant Name and Location",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F24": field.NewHex(&field.Spec{
				Description: "Payment Account Reference (PAR)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F17": field.NewHex(&field.Spec{
				Description: "Personal Identification Number (PIN) Try Counter",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F39": field.NewHex(&field.Spec{
				Description: "Point-of-Service (POS) Entry Mode",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F38": field.NewHex(&field.Spec{
				Description: "Processing Options Data Object List (PDOL)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"80": field.NewHex(&field.Spec{
				Description: "Response Message Template Format 1",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"77": field.NewHex(&field.Spec{
				Description: "Response Message Template Format 2",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F30": field.NewHex(&field.Spec{
				Description: "Service Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"88": field.NewHex(&field.Spec{
				Description: "Short File Identifier (SFI)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F4B": field.NewHex(&field.Spec{
				Description: "Signed Dynamic Application Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"93": field.NewHex(&field.Spec{
				Description: "Signed Static Application Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F4A": field.NewHex(&field.Spec{
				Description: "Static Data Authentication Tag List",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F33": field.NewHex(&field.Spec{
				Description: "Terminal Capabilities",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F1A": field.NewHex(&field.Spec{
				Description: "Terminal Country Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F1B": field.NewHex(&field.Spec{
				Description: "Terminal Floor Limit",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F1C": field.NewHex(&field.Spec{
				Description: "Terminal Identification",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F1D": field.NewHex(&field.Spec{
				Description: "Terminal Risk Management Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F35": field.NewHex(&field.Spec{
				Description: "Terminal Type",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"95": field.NewHex(&field.Spec{
				Description: "Terminal Verification Results",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F19": field.NewHex(&field.Spec{
				Description: "Token Requestor ID",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F1F": field.NewHex(&field.Spec{
				Description: "Track 1 Discretionary Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F20": field.NewHex(&field.Spec{
				Description: "Track 2 Discretionary Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"57": field.NewHex(&field.Spec{
				Description: "Track 2 Equivalent Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"98": field.NewHex(&field.Spec{
				Description: "Transaction Certificate (TC) Hash Value",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"97": field.NewHex(&field.Spec{
				Description: "Transaction Certificate Data Object List (TDOL)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F2A": field.NewHex(&field.Spec{
				Description: "Transaction Currency Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"5F36": field.NewHex(&field.Spec{
				Description: "Transaction Currency Exponent",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9A": field.NewHex(&field.Spec{
				Description: "Transaction Date",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"99": field.NewHex(&field.Spec{
				Description: "Transaction Personal Identification Number (PIN) Data",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F3C": field.NewHex(&field.Spec{
				Description: "Transaction Reference Currency Code",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F3D": field.NewHex(&field.Spec{
				Description: "Transaction Reference Currency Exponent",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F41": field.NewHex(&field.Spec{
				Description: "Transaction Sequence Counter",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9B": field.NewHex(&field.Spec{
				Description: "Transaction Status Information",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F21": field.NewHex(&field.Spec{
				Description: "Transaction Time",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9C": field.NewHex(&field.Spec{
				Description: "Transaction Type",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F37": field.NewHex(&field.Spec{
				Description: "Unpredictable Number",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F23": field.NewHex(&field.Spec{
				Description: "Upper Consecutive Offline Limit",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F5B": field.NewHex(&field.Spec{
				Description: "Issuer Script Results",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
		},
	}
)
