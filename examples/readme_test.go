package examples_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/examples"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/network"
	"github.com/stretchr/testify/require"
)

// Types used by the Setting / Getting Message Data README snippets.
type Authorization struct {
	MTI                  string    `iso8583:"0"`
	PrimaryAccountNumber string    `iso8583:"2"`
	ProcessingCode       string    `iso8583:"3"`
	Amount               int64     `iso8583:"4"`
	STAN                 string    `iso8583:"11"`
	LocalTime            string    `iso8583:"12"`
	LocalDate            string    `iso8583:"13"`
	MerchantType         string    `iso8583:"18"`
	AcceptorInfo         *Acceptor `iso8583:"43"`
}

type Acceptor struct {
	Name    string `iso8583:"1"`
	City    string `iso8583:"2"`
	Country string `iso8583:"3"`
}

func sampleAuth() *Authorization {
	return &Authorization{
		MTI:                  "0100",
		PrimaryAccountNumber: "4242424242424242",
		ProcessingCode:       "000000",
		Amount:               9999,
		STAN:                 "000123",
		LocalTime:            "152059",
		LocalDate:            "0205",
		MerchantType:         "5411",
		AcceptorInfo: &Acceptor{
			Name:    "ACME Store",
			City:    "New York",
			Country: "US",
		},
	}
}

func TestReadmeSettingAndGettingMessageData(t *testing.T) {
	// Setting Message Data
	msg := iso8583.NewMessage(examples.Spec)

	auth := sampleAuth()

	err := msg.Marshal(auth)
	require.NoError(t, err)

	data, err := msg.Pack()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// Getting Message Data
	receivedData := data
	msg = iso8583.NewMessage(examples.Spec)

	err = msg.Unpack(receivedData)
	require.NoError(t, err)

	var got Authorization
	err = msg.Unmarshal(&got)
	require.NoError(t, err)

	require.Equal(t, int64(9999), got.Amount)
	require.NotNil(t, got.AcceptorInfo)
	require.Equal(t, "ACME Store", got.AcceptorInfo.Name)
	require.Equal(t, "New York", got.AcceptorInfo.City)

	fmt.Printf("Transaction amount: %d\n", got.Amount)
	fmt.Printf("Merchant: %s, %s\n", got.AcceptorInfo.Name, got.AcceptorInfo.City)

	var buf bytes.Buffer
	require.NoError(t, iso8583.Describe(msg, &buf))
	require.Contains(t, buf.String(), "ACME Store")
}

func TestReadmeInspectingFieldsFilter(t *testing.T) {
	message := iso8583.NewMessage(examples.Spec)
	message.MTI("0100")
	require.NoError(t, message.Field(2, "4242424242424242"))
	require.NoError(t, message.Field(4, "100"))

	filterAll := func(in string, data field.Field) string {
		runesInString := utf8.RuneCountInString(in)
		return strings.Repeat("*", runesInString)
	}

	var filtered bytes.Buffer
	require.NoError(t, iso8583.Describe(message, &filtered, iso8583.FilterField("2", filterAll)))
	require.Contains(t, filtered.String(), "****************")

	var unfiltered bytes.Buffer
	require.NoError(t, iso8583.Describe(message, &unfiltered, iso8583.DoNotFilterFields()...))
	require.Contains(t, unfiltered.String(), "4242424242424242")
}

func TestReadmeJSONIncludesBitmap(t *testing.T) {
	message := iso8583.NewMessage(iso8583.Spec87)
	message.MTI("0100")
	require.NoError(t, message.Field(2, "4242424242424242"))
	require.NoError(t, message.Field(3, "123456"))
	require.NoError(t, message.Field(4, "100"))

	jsonMessage, err := json.Marshal(message)
	require.NoError(t, err)
	require.Contains(t, string(jsonMessage), `"1":`)
	require.Contains(t, string(jsonMessage), `"0":"0100"`)
}

func TestReadmeNetworkHeaders(t *testing.T) {
	// Write a header + packed body, then read them back as the README snippet does.
	msg := iso8583.NewMessage(examples.Spec)
	require.NoError(t, msg.Marshal(sampleAuth()))
	packed, err := msg.Pack()
	require.NoError(t, err)

	var conn bytes.Buffer
	header := network.NewBCD2BytesHeader()
	header.SetLength(len(packed))
	_, err = header.WriteTo(&conn)
	require.NoError(t, err)
	n, err := conn.Write(packed)
	require.NoError(t, err)
	require.Equal(t, len(packed), n)

	header = network.NewBCD2BytesHeader()
	_, err = header.ReadFrom(&conn)
	require.NoError(t, err)

	buf := make([]byte, header.Length())
	read, err := io.ReadFull(&conn, buf)
	require.NoError(t, err)
	if read != header.Length() {
		t.Fatalf("read %d bytes, want %d", read, header.Length())
	}

	message := iso8583.NewMessage(examples.Spec)
	require.NoError(t, message.Unpack(buf))
	mti, err := message.GetMTI()
	require.NoError(t, err)
	require.Equal(t, "0100", mti)
}
