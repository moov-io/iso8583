package emv

import (
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
)

var Spec = &field.Spec{
	Length:      999,
	Description: "ICC Data",
	Pref:        prefix.ASCII.LLL,
	Tag: &field.TagSpec{
		Sort: sort.StringsByHex,
		Enc:  encoding.BerTLVTag,
	},
	Subfields: map[string]field.Field{
		"9F01": field.NewString(&field.Spec{
			Description: "Acquirer Identifier",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F40": field.NewString(&field.Spec{
			Description: "Additional Terminal Capabilities",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"81": field.NewString(&field.Spec{
			Description: "Amount, Authorised (Binary)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F02": field.NewNumeric(&field.Spec{
			Description: "Amount, Authorised (Numeric)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F04": field.NewString(&field.Spec{
			Description: "Amount, Other (Binary)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F03": field.NewNumeric(&field.Spec{
			Description: "Amount, Other (Numeric)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F3A": field.NewString(&field.Spec{
			Description: "Amount, Reference Currency",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F26": field.NewString(&field.Spec{
			Description: "Application Cryptogram",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F42": field.NewString(&field.Spec{
			Description: "Application Currency Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F44": field.NewString(&field.Spec{
			Description: "Application Currency Exponent",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F05": field.NewString(&field.Spec{
			Description: "Application Discretionary Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F25": field.NewString(&field.Spec{
			Description: "Application Effective Date",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F24": field.NewString(&field.Spec{
			Description: "Application Expiration Date",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"94": field.NewString(&field.Spec{
			Description: "Application File Locator (AFL)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"4F": field.NewString(&field.Spec{
			Description: "Application Identifier (AID) – card",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F06": field.NewString(&field.Spec{
			Description: "Application Identifier (AID) – terminal",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"82": field.NewString(&field.Spec{
			Description: "Application Interchange Profile",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"50": field.NewString(&field.Spec{
			Description: "Application Label",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F12": field.NewString(&field.Spec{
			Description: "Application Preferred Name",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5A": field.NewString(&field.Spec{
			Description: "Application Primary Account Number (PAN)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F34": field.NewString(&field.Spec{
			Description: "Application Primary Account Number (PAN) Sequence Number",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"87": field.NewString(&field.Spec{
			Description: "Application Priority Indicator",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F3B": field.NewString(&field.Spec{
			Description: "Application Reference Currency",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F43": field.NewString(&field.Spec{
			Description: "Application Reference Currency Exponent",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F0A": field.NewString(&field.Spec{
			Description: "Application Selection Registered Proprietary Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"61": field.NewString(&field.Spec{
			Description: "Application Template",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F36": field.NewNumeric(&field.Spec{
			Description: "Application Transaction Counter",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F07": field.NewString(&field.Spec{
			Description: "Application Usage Control",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F08": field.NewString(&field.Spec{
			Description: "Application Version Number ICC",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F09": field.NewString(&field.Spec{
			Description: "Application Version Number Terminal",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"89": field.NewString(&field.Spec{
			Description: "Authorisation Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"8A": field.NewString(&field.Spec{
			Description: "Authorisation Response Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F54": field.NewString(&field.Spec{
			Description: "Bank Identifier Code (BIC)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F31": field.NewString(&field.Spec{
			Description: "Card BIT Group Template",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"8C": field.NewString(&field.Spec{
			Description: "Card Risk Management Data Object List 1 (CDOL1)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"8D": field.NewString(&field.Spec{
			Description: "Card Risk Management Data Object List 2 (CDOL2)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F20": field.NewString(&field.Spec{
			Description: "Cardholder Name",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F0B": field.NewString(&field.Spec{
			Description: "Cardholder Name Extended",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"8E": field.NewString(&field.Spec{
			Description: "Cardholder Verification Method (CVM) List",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F34": field.NewString(&field.Spec{
			Description: "Cardholder Verification Method (CVM) Results",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"8F": field.NewString(&field.Spec{
			Description: "Certification Authority Public Key Index ICC",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F22": field.NewString(&field.Spec{
			Description: "Certification Authority Public Key Index Terminal",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"83": field.NewString(&field.Spec{
			Description: "Command Template",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F27": field.NewString(&field.Spec{
			Description: "Cryptogram Information Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F45": field.NewString(&field.Spec{
			Description: "Data Authentication Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"84": field.NewString(&field.Spec{
			Description: "Dedicated File (DF) Name",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9D": field.NewString(&field.Spec{
			Description: "Directory Definition File (DDF) Name",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"73": field.NewString(&field.Spec{
			Description: "Directory Discretionary Template",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F49": field.NewString(&field.Spec{
			Description: "Dynamic Data Authentication Data Object List (DDOL)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"70": field.NewString(&field.Spec{
			Description: "EMV Proprietary Template",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"DF50": field.NewString(&field.Spec{
			Description: "Facial Try Counter",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"BF0C": field.NewString(&field.Spec{
			Description: "File Control Information (FCI) Issuer Discretionary Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"A5": field.NewString(&field.Spec{
			Description: "File Control Information (FCI) Proprietary Template",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"6F": field.NewString(&field.Spec{
			Description: "File Control Information (FCI) Template",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"DF51": field.NewString(&field.Spec{
			Description: "Finger Try Counter",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F4C": field.NewString(&field.Spec{
			Description: "ICC Dynamic Number",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F2D": field.NewString(&field.Spec{
			Description: "Integrated Circuit Card (ICC) PIN Encipherment Public Key Certificate",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F2E": field.NewString(&field.Spec{
			Description: "Integrated Circuit Card (ICC) PIN Encipherment Public Key Exponent",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F2F": field.NewString(&field.Spec{
			Description: "Integrated Circuit Card (ICC) PIN Encipherment Public Key Remainder",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F46": field.NewString(&field.Spec{
			Description: "Integrated Circuit Card (ICC) Public Key Certificate",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F47": field.NewString(&field.Spec{
			Description: "Integrated Circuit Card (ICC) Public Key Exponent",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F48": field.NewString(&field.Spec{
			Description: "Integrated Circuit Card (ICC) Public Key Remainder",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F1E": field.NewString(&field.Spec{
			Description: "Interface Device (IFD) Serial Number",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F53": field.NewString(&field.Spec{
			Description: "International Bank Account Number (IBAN)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F0D": field.NewString(&field.Spec{
			Description: "Issuer Action Code – Default",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F0E": field.NewString(&field.Spec{
			Description: "Issuer Action Code – Denial",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F0F": field.NewString(&field.Spec{
			Description: "Issuer Action Code – Online",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F10": field.NewString(&field.Spec{
			Description: "Issuer Application Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"91": field.NewString(&field.Spec{
			Description: "Issuer Authentication Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F11": field.NewString(&field.Spec{
			Description: "Issuer Code Table Index",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F28": field.NewString(&field.Spec{
			Description: "Issuer Country Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F55": field.NewString(&field.Spec{
			Description: "Issuer Country Code (alpha2 format)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F56": field.NewString(&field.Spec{
			Description: "Issuer Country Code (alpha3 format)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"42": field.NewString(&field.Spec{
			Description: "Issuer Identification Number (IIN)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F0C": field.NewString(&field.Spec{
			Description: "Issuer Identification Number Extended",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"90": field.NewString(&field.Spec{
			Description: "Issuer Public Key Certificate",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F32": field.NewString(&field.Spec{
			Description: "Issuer Public Key Exponent",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"92": field.NewString(&field.Spec{
			Description: "Issuer Public Key Remainder",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"86": field.NewString(&field.Spec{
			Description: "Issuer Script Command",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F18": field.NewString(&field.Spec{
			Description: "Issuer Script Identifier",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"71": field.NewString(&field.Spec{
			Description: "Issuer Script Template 1",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"72": field.NewString(&field.Spec{
			Description: "Issuer Script Template 2",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F50": field.NewString(&field.Spec{
			Description: "Issuer URL",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F2D": field.NewString(&field.Spec{
			Description: "Language Preference",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F13": field.NewString(&field.Spec{
			Description: "Last Online Application Transaction Counter (ATC) Register",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F4D": field.NewString(&field.Spec{
			Description: "Log Entry",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F4F": field.NewString(&field.Spec{
			Description: "Log Format",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F14": field.NewString(&field.Spec{
			Description: "Lower Consecutive Offline Limit",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F15": field.NewString(&field.Spec{
			Description: "Merchant Category Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F16": field.NewString(&field.Spec{
			Description: "Merchant Identifier",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F4E": field.NewString(&field.Spec{
			Description: "Merchant Name and Location",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F24": field.NewString(&field.Spec{
			Description: "Payment Account Reference (PAR)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F17": field.NewString(&field.Spec{
			Description: "Personal Identification Number (PIN) Try Counter",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F39": field.NewString(&field.Spec{
			Description: "Point-of-Service (POS) Entry Mode",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F38": field.NewString(&field.Spec{
			Description: "Processing Options Data Object List (PDOL)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"80": field.NewString(&field.Spec{
			Description: "Response Message Template Format 1",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"77": field.NewString(&field.Spec{
			Description: "Response Message Template Format 2",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F30": field.NewString(&field.Spec{
			Description: "Service Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"88": field.NewString(&field.Spec{
			Description: "Short File Identifier (SFI)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F4B": field.NewString(&field.Spec{
			Description: "Signed Dynamic Application Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"93": field.NewString(&field.Spec{
			Description: "Signed Static Application Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F4A": field.NewString(&field.Spec{
			Description: "Static Data Authentication Tag List",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F33": field.NewString(&field.Spec{
			Description: "Terminal Capabilities",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F1A": field.NewString(&field.Spec{
			Description: "Terminal Country Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F1B": field.NewString(&field.Spec{
			Description: "Terminal Floor Limit",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F1C": field.NewString(&field.Spec{
			Description: "Terminal Identification",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F1D": field.NewString(&field.Spec{
			Description: "Terminal Risk Management Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F35": field.NewString(&field.Spec{
			Description: "Terminal Type",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"95": field.NewString(&field.Spec{
			Description: "Terminal Verification Results",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F19": field.NewString(&field.Spec{
			Description: "Token Requestor ID",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F1F": field.NewString(&field.Spec{
			Description: "Track 1 Discretionary Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F20": field.NewString(&field.Spec{
			Description: "Track 2 Discretionary Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"57": field.NewString(&field.Spec{
			Description: "Track 2 Equivalent Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"98": field.NewString(&field.Spec{
			Description: "Transaction Certificate (TC) Hash Value",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"97": field.NewString(&field.Spec{
			Description: "Transaction Certificate Data Object List (TDOL)",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F2A": field.NewString(&field.Spec{
			Description: "Transaction Currency Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"5F36": field.NewString(&field.Spec{
			Description: "Transaction Currency Exponent",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9A": field.NewString(&field.Spec{
			Description: "Transaction Date",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"99": field.NewString(&field.Spec{
			Description: "Transaction Personal Identification Number (PIN) Data",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F3C": field.NewString(&field.Spec{
			Description: "Transaction Reference Currency Code",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F3D": field.NewString(&field.Spec{
			Description: "Transaction Reference Currency Exponent",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F41": field.NewString(&field.Spec{
			Description: "Transaction Sequence Counter",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9B": field.NewString(&field.Spec{
			Description: "Transaction Status Information",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F21": field.NewString(&field.Spec{
			Description: "Transaction Time",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9C": field.NewString(&field.Spec{
			Description: "Transaction Type",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F37": field.NewString(&field.Spec{
			Description: "Unpredictable Number",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
		"9F23": field.NewString(&field.Spec{
			Description: "Upper Consecutive Offline Limit",
			Enc:         encoding.ASCIIHexToBytes,
			Pref:        prefix.BerTLV,
		}),
	},
}
