package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	programName = filepath.Base(os.Args[0])
	describeCmd = "describe"
)

func main() {
	describeCommand := flag.NewFlagSet("describe", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Work seamlessly with ISO 8583 from the command line.\n\nUsage:\n  %s <command> [flags]\n\n", programName)
		fmt.Fprintf(os.Stdout, "Available commands:\n")

		// TODO: we will print all commands when we have more than one
		fmt.Fprintf(os.Stdout, "  %s: display ISO 8583 file in a human-readable format\n", describeCmd)
		fmt.Fprintf(os.Stdout, "\n")
	}

	describeCommand.Usage = func() {
		fmt.Fprintf(os.Stdout, "Display ISO 8583 file in a human-readable format.\n\nUsage:\n  %s %s <files> [flags]\n\n", programName, describeCmd)
		fmt.Fprintf(os.Stdout, "Flags: \n")
		describeCommand.PrintDefaults()
		fmt.Fprintf(os.Stdout, "\n")
	}

	// TODO: we have to provide information about available specs
	specName := describeCommand.String("spec", "87", "name of built-in spec")

	flag.Parse()

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
		err := Describe(describeCommand.Args(), *specName)
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
