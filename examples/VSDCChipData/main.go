// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/sort"
)

//	The sample is for VSDC chip data usage

//	This field 55 VSDC chip data usage contains three subfields after the length subfield.
//	Positions:
//		1 2 3 4 ... 255
//  Fields
//	 - Subfield 1: length Byte, a one-byte binary subfield  that contains the number of bytes in this field after the length subfield
//	 - Subfield 2: dataset ID, a one-byte binary identifier
//	 - Subfield 3: dataset length, 2-byte binary subfield that contains the total length of all TLV elements that follow.
//	 - Subfield 4:
//		Chip Card TLV data elements
//		Tag Length Value Tag Length Value

type Dataset struct {
	F9A   *field.String
	F9F02 *field.String
}

type TestISOF55Data struct {
	F1 *field.String
	F2 *field.String
	F3 *field.String
	F4 *Dataset
}

func main() {

	datasetSpec := &field.Spec{
		Description: "Chip Card TLV data elements",
		Pref:        prefix.None.Fixed,
		Tag: &field.TagSpec{
			Sort: sort.StringsByHex,
			Enc:  encoding.BerTLVTag,
		},
		Subfields: map[string]field.Field{
			"9A": field.NewString(&field.Spec{
				Description: "Transaction Date",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
			"9F02": field.NewString(&field.Spec{
				Description: "Amount, Authorized (Numeric)",
				Enc:         encoding.Binary,
				Pref:        prefix.BerTLV,
			}),
		},
	}

	f55Spec := &field.Spec{
		Length:      999,
		Description: "ICC Data – EMV Having Multiple Tags",
		Pref:        prefix.ASCII.LLL,
		Pad:         padding.None,
		Tag: &field.TagSpec{
			Sort: sort.StringsByInt,
		},
		Subfields: map[string]field.Field{
			"1": field.NewString(&field.Spec{
				Description: "Length Subfield",
				Enc:         encoding.ASCIIHexToBytes,
				Length:      1,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			"2": field.NewString(&field.Spec{
				Description: "Dataset ID Subfield",
				Enc:         encoding.ASCIIHexToBytes,
				Length:      1,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			"3": field.NewString(&field.Spec{
				Description: "Dataset Length Subfield",
				Enc:         encoding.ASCIIHexToBytes,
				Length:      2,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			"4": field.NewComposite(datasetSpec),
		},
	}

	dataSet := Dataset{
		F9A:   field.NewStringValue("210720"),
		F9F02: field.NewStringValue("000000000501"),
	}

	// Getting size
	datasetField := field.NewComposite(datasetSpec)
	err := datasetField.Marshal(&dataSet)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf, err := datasetField.Pack()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Making F55 field
	datasetSize := len(buf)

	f55Data := TestISOF55Data{
		F1: field.NewStringValue(fmt.Sprintf("%02x", datasetSize+3)),
		F2: field.NewStringValue("01"),
		F3: field.NewStringValue(fmt.Sprintf("%04x", datasetSize)),
		F4: &dataSet,
	}

	// creating field
	f55 := field.NewComposite(f55Spec)

	// Setting value
	err = f55.Marshal(&f55Data)

	if err != nil {
		fmt.Println(err)
		return
	}

	// get binary value of the field
	rawValue, err := f55.Pack()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("0x%X", rawValue))

	fmt.Println("\n EMV Having Multiple Tags (Bit 55) \n")
	fmt.Println("ICC Data length: ", fmt.Sprintf("0x%X", rawValue[0:3]))
	fmt.Println(".........length: ", fmt.Sprintf("0x%X", rawValue[3:4]))
	fmt.Println(".............id: ", fmt.Sprintf("0x%X", rawValue[4:5]))
	fmt.Println(".dataset length: ", fmt.Sprintf("0x%X", rawValue[5:7]))
	fmt.Println("..dataset(tlvs): ", fmt.Sprintf("0x%X", rawValue[7:]))

}
