// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/iso8583/pkg/lib"
	"github.com/moov-io/iso8583/pkg/server"
	"github.com/moov-io/iso8583/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	messageFile       string
	specificationFile string

	iso8583message      []byte
	specificationBuffer []byte
)

var WebCmd = &cobra.Command{
	Use:   "web",
	Short: "Launches web server",
	Long:  "Launches web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		env := &server.Environment{
			Logger: logging.NewDefaultLogger().WithKeyValue("app", "iso8583"),
		}

		env, err := server.NewEnvironment(env)
		if err != nil {
			env.Logger.Fatal().LogError("Error loading up environment.", err)
			os.Exit(1)
		}
		defer env.Shutdown()

		env.Logger.Info().Log("Starting services")
		test, _ := cmd.Flags().GetBool("test")
		if !test {
			shutdown := env.RunServers(true)
			defer shutdown()
		}
		return nil
	},
}

var Validate = &cobra.Command{
	Use:   "validator",
	Short: "Validate iso8583 message",
	Long:  "Validate an incoming iso8583 message",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var message lib.Iso8583Message
		spec, err := lib.NewSpecificationWithJson(specificationBuffer)
		if err != nil {
			spec = &utils.ISO8583DataElementsVer1987
		}
		message, err = lib.NewISO8583Message(spec)
		if err != nil {
			return err
		}

		messageFormat := utils.MessageFormat(iso8583message)
		switch messageFormat {
		case utils.MessageFormatJson:
			err = json.Unmarshal(iso8583message, message)
		case utils.MessageFormatXml:
			err = xml.Unmarshal(iso8583message, message)
		case utils.MessageFormatIso8583:
			_, err = message.Load(iso8583message)
		}
		if err != nil {
			return err
		}

		return message.Validate()
	},
}

var Print = &cobra.Command{
	Use:   "print",
	Short: "Print iso8583 message",
	Long:  "Print an incoming iso8583 message with special format (options: iso8583, json, xml)",
	RunE: func(cmd *cobra.Command, args []string) error {
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}
		if format != utils.MessageFormatJson && format != utils.MessageFormatXml && format != utils.MessageFormatIso8583 {
			return errors.New("don't support the format")
		}

		var message lib.Iso8583Message
		spec, err := lib.NewSpecificationWithJson(specificationBuffer)
		if err != nil {
			spec = &utils.ISO8583DataElementsVer1987
		}
		message, err = lib.NewISO8583Message(spec)
		if err != nil {
			return err
		}

		messageFormat := utils.MessageFormat(iso8583message)
		switch messageFormat {
		case utils.MessageFormatJson:
			err = json.Unmarshal(iso8583message, message)
		case utils.MessageFormatXml:
			err = xml.Unmarshal(iso8583message, message)
		case utils.MessageFormatIso8583:
			_, err = message.Load(iso8583message)
		}
		if err != nil {
			return err
		}

		var output []byte
		switch format {
		case utils.MessageFormatJson:
			output, err = json.MarshalIndent(message, "", "\t")
		case utils.MessageFormatXml:
			output, err = xml.MarshalIndent(message, "", "\t")
		case utils.MessageFormatIso8583:
			output, err = message.Bytes()
		}
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

var Convert = &cobra.Command{
	Use:   "convert [output]",
	Short: "Convert iso8583 message format",
	Long:  "Convert an incoming iso8583 message into another format (options: iso8583, json, xml)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires output argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}
		if format != utils.MessageFormatJson && format != utils.MessageFormatXml && format != utils.MessageFormatIso8583 {
			return errors.New("don't support the format")
		}

		var message lib.Iso8583Message
		spec, err := lib.NewSpecificationWithJson(specificationBuffer)
		if err != nil {
			spec = &utils.ISO8583DataElementsVer1987
		}
		message, err = lib.NewISO8583Message(spec)
		if err != nil {
			return err
		}

		messageFormat := utils.MessageFormat(iso8583message)
		switch messageFormat {
		case utils.MessageFormatJson:
			err = json.Unmarshal(iso8583message, message)
		case utils.MessageFormatXml:
			err = xml.Unmarshal(iso8583message, message)
		case utils.MessageFormatIso8583:
			_, err = message.Load(iso8583message)
		}
		if err != nil {
			return err
		}

		var output []byte
		switch format {
		case utils.MessageFormatJson:
			output, err = json.MarshalIndent(message, "", "\t")
		case utils.MessageFormatXml:
			output, err = xml.MarshalIndent(message, "", "\t")
		case utils.MessageFormatIso8583:
			output, err = message.Bytes()
		}
		if err != nil {
			return err
		}

		wFile, err := os.Create(args[0])
		if err != nil {
			return err
		}
		_, err = wFile.Write(output)
		wFile.Close()

		return err
	},
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		isWeb := false
		cmdNames := make([]string, 0)
		getName := func(c *cobra.Command) {}
		getName = func(c *cobra.Command) {
			if c == nil {
				return
			}
			cmdNames = append([]string{c.Name()}, cmdNames...)
			if c.Name() == "web" {
				isWeb = true
			}
			getName(c.Parent())
		}
		getName(cmd)

		if !isWeb {
			if messageFile == "" {
				path, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				messageFile = filepath.Join(path, "iso8583_message.dat")
			}
			_, err := os.Stat(messageFile)
			if os.IsNotExist(err) {
				return errors.New("invalid input file")
			}
			iso8583message, err = ioutil.ReadFile(messageFile)

			if err != nil {
				return err
			}

			if specificationFile == "" {
				path, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				specificationFile = filepath.Join(path, "iso8583_specification.json")
			}
			_, err = os.Stat(specificationFile)
			if err == nil {
				specificationBuffer, err = ioutil.ReadFile(specificationFile)
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func initRootCmd() {
	WebCmd.Flags().BoolP("test", "t", false, "test server")
	Convert.Flags().String("format", "iso8583", "format of iso8583 message(required)")
	Convert.MarkFlagRequired("format")
	Print.Flags().String("format", "iso8583", "print format")

	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().StringVar(&messageFile, "input", "", "iso8583 message (the message types are iso8583 raw message, xml, json. default is $PWD/iso8583_message.dat)")
	rootCmd.PersistentFlags().StringVar(&specificationFile, "spec", "", "specification file (default is $PWD/iso8583_specification.json)")
	rootCmd.AddCommand(WebCmd)
	rootCmd.AddCommand(Convert)
	rootCmd.AddCommand(Print)
	rootCmd.AddCommand(Validate)
}

func main() {
	initRootCmd()

	rootCmd.Execute()
}
