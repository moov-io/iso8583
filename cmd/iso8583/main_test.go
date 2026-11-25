// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/specs"
	"github.com/stretchr/testify/require"
)

func TestRunFlags(t *testing.T) {
	for _, tc := range []struct {
		name   string
		args   []string
		code   int
		stdout []string
		stderr []string
	}{
		{"version", []string{"-version"}, 0, []string{"Version: " + iso8583.Version + "\n"}, nil},
		{"no arguments", nil, 1, []string{"Usage:", "Available commands:", "describe:"}, nil},
		{"unknown command", []string{"unknown"}, 1, []string{"Unknown command: unknown", "Available commands:"}, nil},
		{"describe without files", []string{"describe"}, 1, []string{"describe [flags] <files>"}, []string{"-spec", "-spec-file"}},
		{"describe flags without files", []string{"describe", "-spec", "87hex"}, 1, []string{"describe [flags] <files>"}, []string{"-spec-file"}},
		{"unknown spec", []string{"describe", "-spec", "unknown", "msg.bin"}, 1, []string{"Unknown spec: unknown", "Supported specs:", "87ascii", "87hex"}, nil},
		{"help", []string{"-help"}, 0, []string{"Available commands:"}, nil},
		{"describe help", []string{"describe", "-help"}, 0, []string{"describe [flags] <files>"}, []string{"-spec-file"}},
		{"unknown flag", []string{"-invalid"}, 2, []string{"Available commands:"}, []string{"flag provided but not defined: -invalid"}},
		{"missing describe flag value", []string{"describe", "-spec"}, 2, []string{"describe [flags] <files>"}, []string{"flag needs an argument: -spec"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			require.Equal(t, tc.code, run(tc.args, &stdout, &stderr))
			for _, want := range tc.stdout {
				require.Contains(t, stdout.String(), want)
			}
			for _, want := range tc.stderr {
				require.Contains(t, stderr.String(), want)
			}
			if tc.stderr == nil {
				require.Empty(t, stderr.String())
			}
		})
	}
}

func TestRunDescribe(t *testing.T) {
	ascii := packedTestMessage(t)
	hex := filepath.Join(t.TempDir(), "hex.bin")
	const hexMessage = "30323030" + "2000000000000000" + "313233343536" // MTI 0200, bitmap with F3, processing code 123456.
	require.NoError(t, os.WriteFile(hex, []byte(hexMessage), 0o600))
	yaml := filepath.Join("..", "..", "examples", "specs", "spec87ascii.yaml")
	badSpec := filepath.Join(t.TempDir(), "invalid.yaml")
	require.NoError(t, os.WriteFile(badSpec, []byte("fields: [\n"), 0o600))
	missingSpec := filepath.Join(t.TempDir(), "missing.yaml")
	missingMessage := filepath.Join(t.TempDir(), "missing.bin")

	for _, tc := range []struct {
		name     string
		args     []string
		code     int
		specName string
		mti      string
		output   []string
	}{
		{name: "default ASCII spec", args: []string{ascii}, specName: specs.Spec87ASCII.Name, mti: "0100"},
		{name: "hex spec", args: []string{"-spec", "87hex", hex}, specName: specs.Spec87Hex.Name, mti: "0200"},
		{name: "external spec takes precedence", args: []string{"-spec", "unknown", "-spec-file", yaml, ascii}, specName: specs.Spec87ASCII.Name, mti: "0100"},
		{name: "malformed external spec", args: []string{"-spec-file", badSpec, ascii}, code: 1, output: []string{"Error describing files: creating spec from file:", "failed to YAML unmarshal bytes to MessageSpec"}},
		{name: "missing external spec", args: []string{"-spec-file", missingSpec, ascii}, code: 1, output: []string{"Error describing files: creating spec from file:", "opening file " + missingSpec}},
		{name: "missing message file", args: []string{missingMessage}, code: 1, output: []string{"Error describing files: creating message from file:", "opening file " + missingMessage}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			require.Equal(t, tc.code, run(append([]string{"describe"}, tc.args...), &stdout, &stderr), stdout.String())
			require.Empty(t, stderr.String())
			if tc.specName != "" {
				require.Contains(t, stdout.String(), tc.specName+" Message:")
				require.Regexp(t, `(?m)^MTI\.+: `+tc.mti+`$`, stdout.String())
				require.Regexp(t, `(?m)^F3\s+Processing Code\.+: 123456$`, stdout.String())
			}
			for _, want := range tc.output {
				require.Contains(t, stdout.String(), want)
			}
		})
	}
}

func TestRunDescribeMultipleFiles(t *testing.T) {
	first := packedTestMessage(t)
	second := writeTestMessage(t, "0200")
	var stdout, stderr bytes.Buffer
	require.Zero(t, run([]string{"describe", first, second}, &stdout, &stderr))
	require.Empty(t, stderr.String())

	var mtis []string
	for _, line := range strings.Split(stdout.String(), "\n") {
		if strings.HasPrefix(line, "MTI") {
			_, value, _ := strings.Cut(line, ": ")
			mtis = append(mtis, value)
		}
	}
	require.Equal(t, []string{"0100", "0200"}, mtis)
}

func TestRunDescribePartiallyUnpackedMessage(t *testing.T) {
	packed, err := os.ReadFile(packedTestMessage(t))
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), "partial.bin")
	require.NoError(t, os.WriteFile(path, packed[:len(packed)-1], 0o600))

	var stdout, stderr bytes.Buffer
	require.Zero(t, run([]string{"describe", path}, &stdout, &stderr), stdout.String())
	require.Empty(t, stderr.String())
	require.Contains(t, stdout.String(), "Failed to create message from file: unpacking message:")
	require.Contains(t, stdout.String(), "Trying to describe file anyway...")
	require.Regexp(t, `(?m)^MTI\.+: 0100$`, stdout.String())
	require.Regexp(t, `(?m)^F3\s+Processing Code\.+: 123456$`, stdout.String())
}

func TestMainVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("building the CLI binary")
	}
	const version = "test-version"
	bin := filepath.Join(t.TempDir(), "iso8583")
	goTool := filepath.Join(runtime.GOROOT(), "bin", "go")
	if runtime.GOOS == "windows" {
		bin += ".exe"
		goTool += ".exe"
	}
	build := exec.CommandContext(t.Context(), goTool, "build", "-buildvcs=false", "-ldflags", "-X github.com/moov-io/iso8583.Version="+version, "-o", bin, ".")
	output, err := build.CombinedOutput()
	require.NoError(t, err, "%s", output)

	output, err = exec.CommandContext(t.Context(), bin, "-version").CombinedOutput()
	require.NoError(t, err, "%s", output)
	require.Equal(t, "Version: "+version+"\n\n", string(output))
}
