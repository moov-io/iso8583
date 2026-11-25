package main

import (
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/moov-io/iso8583"
)

const describeCmd = "describe"

var programName = filepath.Base(os.Args[0])

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet(programName, flag.ContinueOnError)
	flags.SetOutput(stderr)
	versionFlag := flags.Bool("version", false, "show version")
	flags.Usage = func() {
		fmt.Fprintf(stdout, "Work seamlessly with ISO 8583 from the command line.\n\nUsage:\n  %s <command> [flags]\n\n", programName)
		fmt.Fprintf(stdout, "Available commands:\n  %s: display ISO 8583 file in a human-readable format\n\n", describeCmd)
	}

	if err := flags.Parse(args); err != nil {
		return flagExitCode(err)
	}
	if *versionFlag {
		fmt.Fprintf(stdout, "Version: %s\n\n", iso8583.Version)
		return 0
	}
	if flags.NArg() == 0 {
		flags.Usage()
		return 1
	}
	command := flags.Arg(0)
	if command != describeCmd {
		fmt.Fprintf(stdout, "Unknown command: %s\n\n", command)
		flags.Usage()
		return 1
	}

	return runDescribe(flags.Args()[1:], stdout, stderr)
}

func runDescribe(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet(describeCmd, flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		fmt.Fprintf(stdout, "Display ISO 8583 file in a human-readable format.\n\nUsage:\n  %s %s [flags] <files> \n\n", programName, describeCmd)
		fmt.Fprintln(stdout, "Flags: ")
		flags.PrintDefaults()
		fmt.Fprintln(stdout)
	}
	availableSpecNames := strings.Join(slices.Sorted(maps.Keys(availableSpecs)), ", ")
	specName := flags.String("spec", "87ascii", fmt.Sprintf("name of built-in spec: %s", availableSpecNames))
	specFileName := flags.String("spec-file", "", "path to customized specification file in JSON or YAML format")

	if err := flags.Parse(args); err != nil {
		return flagExitCode(err)
	}

	files := flags.Args()
	var err error
	if *specFileName != "" {
		err = describeWithSpecFile(files, *specFileName, stdout)
	} else if spec := availableSpecs[*specName]; spec != nil {
		err = describeMessage(files, spec, stdout)
	} else {
		fmt.Fprintf(stdout, "Unknown spec: %s\n\nSupported specs: %s\n\n", *specName, availableSpecNames)
		return 1
	}
	if err != nil {
		fmt.Fprintf(stdout, "Error describing files: %s\n", err)
		return 1
	}
	if len(files) == 0 {
		flags.Usage()
		return 1
	}
	return 0
}

func flagExitCode(err error) int {
	if err == flag.ErrHelp {
		return 0
	}
	return 2
}
