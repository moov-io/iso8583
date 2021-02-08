// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/moov-io/iso8583/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	testSpecFilePath    = filepath.Join("..", "..", "test", "testdata", "specification_ver_1987.json")
	testMessageFilePath = filepath.Join("..", "..", "test", "testdata", "iso_reversal_message.dat")
	testInvalidFilePath = filepath.Join("..", "..", "test", "testdata", "iso_reversal_message_error_date.dat")
	testErrorFilePath   = filepath.Join("..", "..", "test", "testdata", "error_message.dat")
	testJsonFilePath    = filepath.Join("..", "..", "test", "testdata", "iso_reversal_message.json")
	testXmlFilePath     = filepath.Join("..", "..", "test", "testdata", "iso_reversal_message.xml")
)

func TestMain(m *testing.M) {
	initRootCmd()
	os.Exit(m.Run())
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func TestConvertWithoutInput(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "output", "--format", utils.MessageFormatJson)
	if err == nil {
		t.Errorf("invalid input file")
	}
}

func TestConvertWithInvalidParam(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "--input", testMessageFilePath, "--format", utils.MessageFormatJson)
	if err == nil {
		t.Errorf("requires output argument")
	}
}

func TestConvertJson(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "output", "--input", testMessageFilePath, "--format", utils.MessageFormatJson)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "convert", "output", "--input", testMessageFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatJson)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestConvertXml(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "output", "--input", testMessageFilePath, "--format", utils.MessageFormatXml)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "convert", "output", "--input", testMessageFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatXml)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestConvertIso8583(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "output", "--input", testMessageFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "convert", "output", "--input", testMessageFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestConvertUnknown(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "output", "--input", testMessageFilePath, "--format", "unknown")
	if err == nil {
		t.Errorf("don't support the format")
	}
}

func TestPrintIso8583(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testMessageFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "print", "--input", testMessageFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPrintJson(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testMessageFilePath, "--format", utils.MessageFormatJson)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "print", "--input", testMessageFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatJson)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPrintXml(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testMessageFilePath, "--format", utils.MessageFormatXml)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "print", "--input", testMessageFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatXml)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPrintUnknown(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testMessageFilePath, "--format", "unknown")
	if err == nil {
		t.Errorf("don't support the format")
	}
}

func TestValidator(t *testing.T) {
	_, err := executeCommand(rootCmd, "validator", "--input", testMessageFilePath)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "validator", "--input", testMessageFilePath, "--spec", testSpecFilePath)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUnknown(t *testing.T) {
	_, err := executeCommand(rootCmd, "unknown")
	if err == nil {
		t.Errorf("don't support unknown")
	}
}

func TestPrintWithInvalidData(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testInvalidFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "print", "--input", testInvalidFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestConvertWithInvalidData(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "output", "--input", testInvalidFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "convert", "output", "--input", testInvalidFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestValidatorWithInvalidData(t *testing.T) {
	_, err := executeCommand(rootCmd, "validator", "--input", testInvalidFilePath)
	if err == nil {
		t.Errorf("error data")
	}

	_, err = executeCommand(rootCmd, "validator", "--input", testInvalidFilePath, "--spec", testSpecFilePath)
	if err == nil {
		t.Errorf("error data")
	}
}

func TestPrintWithErrorData(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testErrorFilePath, "--format", utils.MessageFormatIso8583)
	if err == nil {
		t.Errorf("error data")
	}

	_, err = executeCommand(rootCmd, "print", "--input", testErrorFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err == nil {
		t.Errorf("error data")
	}
}

func TestConvertWithErrorData(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "output", "--input", testErrorFilePath, "--format", utils.MessageFormatIso8583)
	if err == nil {
		t.Errorf("error data")
	}

	_, err = executeCommand(rootCmd, "convert", "output", "--input", testErrorFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err == nil {
		t.Errorf("error data")
	}
}

func TestValidatorWithErrorData(t *testing.T) {
	_, err := executeCommand(rootCmd, "validator", "--input", testErrorFilePath)
	if err == nil {
		t.Errorf("error data")
	}

	_, err = executeCommand(rootCmd, "validator", "--input", testErrorFilePath, "--spec", testSpecFilePath)
	if err == nil {
		t.Errorf("error data")
	}
}

func TestPrintWithJsonData(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testJsonFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "print", "--input", testJsonFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestConvertWithJsonData(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "output", "--input", testJsonFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "convert", "output", "--input", testJsonFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestValidatorWithJsonData(t *testing.T) {
	_, err := executeCommand(rootCmd, "validator", "--input", testJsonFilePath)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "validator", "--input", testJsonFilePath, "--spec", testSpecFilePath)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPrintWithXmlData(t *testing.T) {
	_, err := executeCommand(rootCmd, "print", "--input", testXmlFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "print", "--input", testXmlFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestConvertWithXmlData(t *testing.T) {
	_, err := executeCommand(rootCmd, "convert", "output", "--input", testXmlFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "convert", "output", "--input", testXmlFilePath, "--spec", testSpecFilePath, "--format", utils.MessageFormatIso8583)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestValidatorWithXmlData(t *testing.T) {
	_, err := executeCommand(rootCmd, "validator", "--input", testXmlFilePath)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = executeCommand(rootCmd, "validator", "--input", testXmlFilePath, "--spec", testSpecFilePath)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestWebTest(t *testing.T) {
	_, err := executeCommand(rootCmd, "web", "--test=true")
	if err != nil {
		t.Errorf(err.Error())
	}
}
