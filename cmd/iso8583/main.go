package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/moov-io/iso8583"
)

var (
	programName = filepath.Base(os.Args[0])
	describeCmd = "describe"
)

func main() {
	versionFlag := flag.Bool("version", false, "show version")
	describeCommand := flag.NewFlagSet("describe", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Work seamlessly with ISO 8583 from the command line.\n\nUsage:\n  %s <command> [flags]\n\n", programName)
		fmt.Fprintf(os.Stdout, "Available commands:\n")

		// TODO: we will print all commands when we have more than one
		fmt.Fprintf(os.Stdout, "  %s: display ISO 8583 file in a human-readable format\n", describeCmd)
		fmt.Fprintf(os.Stdout, "\n")
	}

	describeCommand.Usage = func() {
		fmt.Fprintf(os.Stdout, "Display ISO 8583 file in a human-readable format.\n\nUsage:\n  %s %s [flags] <files> \n\n", programName, describeCmd)
		fmt.Fprintf(os.Stdout, "Flags: \n")
		describeCommand.PrintDefaults()
		fmt.Fprintf(os.Stdout, "\n")
	}

	var specNames []string
	for name := range availableSpecs {
		specNames = append(specNames, name)
	}
	availableSpecNames := strings.Join(specNames, ", ")

	specName := describeCommand.String("spec", "87ascii", fmt.Sprintf("name of built-in spec: %s", availableSpecNames))
	specFileName := describeCommand.String("spec-file", "", "path to customized specification file in JSON format")

	flag.Parse()

	if *versionFlag {
		fmt.Fprintf(os.Stdout, "Version: %s\n\n", iso8583.Version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := flag.Arg(0)

	switch command {
	case describeCmd:
		describeArgs := os.Args[2:]
		if len(describeArgs) == 0 {
			describeCommand.Usage()
			os.Exit(1)
		}

		describeCommand.Parse(os.Args[2:])

		var err error
		if specFileName != nil && *specFileName != "" {
			err = DescribeWithSpecFile(describeCommand.Args(), *specFileName)
		} else if availableSpecs[*specName] != nil {
			err = Describe(describeCommand.Args(), *specName)
		} else {
			fmt.Fprintf(os.Stdout, "Unknown spec: %s\n\n", *specName)
			fmt.Fprintf(os.Stdout, "Supported specs: %s\n\n", availableSpecNames)
			os.Exit(1)
		}

		if err != nil {
			fmt.Fprintf(os.Stdout, "Error describing files: %s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stdout, "Uknown command: %s\n\n", command)
		flag.Usage()
		os.Exit(1)
	}

	if describeCommand.Parsed() {
		files := describeCommand.Args()
		if len(files) == 0 {
			describeCommand.Usage()
			os.Exit(1)
		}
	}
}
