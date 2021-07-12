package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/describe"
)

var specs = map[string]*iso8583.MessageSpec{
	"87": iso8583.Spec87,
}

func Describe(paths []string, specName string) error {
	spec := specs[specName]
	if spec == nil {
		return fmt.Errorf("unknown built-in spec %s", spec)
	}

	for _, path := range paths {
		message, err := createMessageFromFile(path, spec)
		if err != nil {
			return fmt.Errorf("creating message from file: %w", err)
		}

		err = describe.Message(os.Stdout, message)
		if err != nil {
			return fmt.Errorf("describing message: %w", err)
		}
	}

	return nil
}

func createMessageFromFile(path string, spec *iso8583.MessageSpec) (*iso8583.Message, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s: %v", path, err)
	}
	defer fd.Close()

	raw, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %v", path, err)
	}

	message := iso8583.NewMessage(spec)
	err = message.Unpack(raw)
	if err != nil {
		return nil, fmt.Errorf("unpacking ISO 8583 message: %v", err)
	}

	return message, nil
}
