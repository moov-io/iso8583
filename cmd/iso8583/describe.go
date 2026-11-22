package main

import (
	"fmt"
	"io"
	"os"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/adapters"
	"github.com/moov-io/iso8583/specs"
)

var availableSpecs = map[string]*iso8583.MessageSpec{
	"87ascii": specs.Spec87ASCII,
	"87hex":   specs.Spec87Hex,
}

func describeMessage(paths []string, spec *iso8583.MessageSpec, adapter adapters.Adapter) error {
	for _, path := range paths {
		message, err := createMessageFromFile(path, spec, adapter)
		if err != nil {
			if message == nil {
				return fmt.Errorf("creating message from file: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Failed to create message from file: %v\n", err)
			fmt.Fprintf(os.Stdout, "Trying to describe file anyway...\n")
		}

		err = iso8583.Describe(message, os.Stdout)
		if err != nil {
			return fmt.Errorf("describing message: %w", err)
		}
	}
	return nil
}

func Describe(paths []string, specName string, useAdapter bool) error {
	spec := availableSpecs[specName]
	if spec == nil {
		return fmt.Errorf("unknown built-in spec %s", specName)
	}

	var adapter adapters.Adapter
	if useAdapter {
		a, err := adapters.AdapterForSpec(spec)
		if err != nil {
			return err
		}
		adapter = a
	}

	return describeMessage(paths, spec, adapter)
}

func DescribeWithSpecFile(paths []string, specFileName string) error {
	spec, err := createSpecFromFile(specFileName)
	if err != nil || spec == nil {
		return fmt.Errorf("creating spec from file: %w", err)
	}

	return describeMessage(paths, spec, nil)
}

func createMessageFromFile(path string, spec *iso8583.MessageSpec, adapter adapters.Adapter) (*iso8583.Message, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %s: %w", path, err)
	}
	defer fd.Close()

	raw, err := io.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	return adapters.UnpackWithAdapter(spec, raw, adapter)
}

func createSpecFromFile(path string) (*iso8583.MessageSpec, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %s: %w", path, err)
	}
	defer fd.Close()

	raw, err := io.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	return specs.ImportJSON(raw)
}
