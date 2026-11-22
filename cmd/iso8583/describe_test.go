package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/specs"
	"github.com/stretchr/testify/require"
)

func packedTestMessage(t *testing.T) string {
	t.Helper()

	msg := iso8583.NewMessage(specs.Spec87ASCII)
	msg.MTI("0100")
	require.NoError(t, msg.Field(2, "4242424242424242"))
	require.NoError(t, msg.Field(3, "123456"))
	require.NoError(t, msg.Field(4, "100"))
	packed, err := msg.Pack()
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "msg.bin")
	require.NoError(t, os.WriteFile(path, packed, 0o644))
	return path
}

func TestDescribeWithBuiltInSpec(t *testing.T) {
	path := packedTestMessage(t)
	require.NoError(t, Describe([]string{path}, "87ascii"))
}

func TestDescribeWithSpecFileJSON(t *testing.T) {
	path := packedTestMessage(t)
	spec := filepath.Join("..", "..", "examples", "specs", "spec87ascii.json")
	require.NoError(t, DescribeWithSpecFile([]string{path}, spec))
}

func TestDescribeWithSpecFileYAML(t *testing.T) {
	path := packedTestMessage(t)
	for _, name := range []string{"spec87ascii.yaml"} {
		spec := filepath.Join("..", "..", "examples", "specs", name)
		require.NoError(t, DescribeWithSpecFile([]string{path}, spec), "spec file %s", name)
	}

	// Also accept the short .yml extension
	yml := filepath.Join(t.TempDir(), "spec.yml")
	yamlBytes, err := os.ReadFile(filepath.Join("..", "..", "examples", "specs", "spec87ascii.yaml"))
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(yml, yamlBytes, 0o644))
	require.NoError(t, DescribeWithSpecFile([]string{path}, yml))
}

func TestCreateSpecFromFileRejectsInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	require.NoError(t, os.WriteFile(path, []byte("not-json"), 0o644))
	_, err := createSpecFromFile(path)
	require.Error(t, err)
}
