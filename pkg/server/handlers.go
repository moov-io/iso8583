// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package server

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moov-io/iso8583/pkg/lib"
	"github.com/moov-io/iso8583/pkg/utils"
)

func outputError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func outputSuccess(w http.ResponseWriter, output string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": output,
	})
}

func parseSpecFromRequest(r *http.Request) (*utils.Specification, error) {
	specFile, _, err := r.FormFile("spec")
	if err != nil {
		return nil, err
	}
	defer specFile.Close()

	var spec bytes.Buffer
	if _, err = io.Copy(&spec, specFile); err != nil {
		return nil, err
	}
	return lib.NewSpecificationWithJson(spec.Bytes())
}

func parseInputFromRequest(r *http.Request) (lib.Iso8583Message, error) {
	inputFile, _, err := r.FormFile("input")
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	var input bytes.Buffer
	if _, err = io.Copy(&input, inputFile); err != nil {
		return nil, err
	}

	var spec *utils.Specification
	spec, err = parseSpecFromRequest(r)
	if err != nil {
		spec = &utils.ISO8583DataElementsVer1987
	}
	message, err := lib.NewISO8583Message(spec)
	if err != nil {
		return nil, err
	}

	data := input.Bytes()
	messageFormat := utils.MessageFormat(data)
	switch messageFormat {
	case utils.MessageFormatJson:
		err = json.Unmarshal(data, message)
	case utils.MessageFormatXml:
		err = xml.Unmarshal(data, message)
	case utils.MessageFormatIso8583:
		_, err = message.Load(data)
	}

	return message, err
}

func messageToBuf(format string, message lib.Iso8583Message) ([]byte, error) {
	var output []byte
	var err error
	switch format {
	case utils.MessageFormatJson:
		output, err = json.MarshalIndent(message, "", "\t")
	case utils.MessageFormatXml:
		output, err = xml.MarshalIndent(message, "", "\t")
	case utils.MessageFormatIso8583:
		output, err = message.Bytes()
	default:
		return nil, errors.New("invalid format")
	}
	return output, err
}

func outputBufferToWriter(w http.ResponseWriter, buf []byte, format string) {
	w.WriteHeader(http.StatusOK)
	switch format {
	case utils.MessageFormatJson:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(buf)
	case utils.MessageFormatXml:
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		xml.NewEncoder(w).Encode(buf)
	case utils.MessageFormatIso8583:
		w.Header().Set("Content-Type", "application/octet-stream; charset=utf-8")
		w.Write(buf)
	}
}

// validator - validate the file based on publication 1220
func validator(w http.ResponseWriter, r *http.Request) {
	message, err := parseInputFromRequest(r)
	if err != nil {
		outputError(w, http.StatusBadRequest, err)
		return
	}

	err = message.Validate()
	if err != nil {
		outputError(w, http.StatusNotImplemented, err)
		return
	}

	outputSuccess(w, "valid file")
	return
}

// validator - print file with ascii or json format
func print(w http.ResponseWriter, r *http.Request) {
	message, err := parseInputFromRequest(r)
	if err != nil {
		outputError(w, http.StatusBadRequest, err)
		return
	}

	format := r.FormValue("format")
	output, err := messageToBuf(format, message)
	if err != nil {
		outputError(w, http.StatusNotImplemented, err)
		return
	}

	outputBufferToWriter(w, output, format)
}

// convert - convert file with ascii or json format
func convert(w http.ResponseWriter, r *http.Request) {
	message, err := parseInputFromRequest(r)
	if err != nil {
		outputError(w, http.StatusBadRequest, err)
		return
	}

	format := r.FormValue("format")
	filename := "converted_file"
	output, err := messageToBuf(format, message)
	if err != nil {
		outputError(w, http.StatusNotImplemented, err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// health - health check
func health(w http.ResponseWriter, r *http.Request) {
	outputSuccess(w, "alive")
	return
}

// configure handlers
func ConfigureHandlers(r *mux.Router) error {
	r.HandleFunc("/health", health).Methods("GET")
	r.HandleFunc("/print", print).Methods("POST")
	r.HandleFunc("/validator", validator).Methods("POST")
	r.HandleFunc("/convert", convert).Methods("POST")
	return nil
}
