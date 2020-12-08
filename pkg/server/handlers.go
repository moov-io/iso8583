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
	"strings"

	"github.com/gorilla/mux"
	"github.com/moov-io/iso8583/pkg/lib"
	"github.com/moov-io/iso8583/pkg/utils"
)

func parseInputFromRequest(r *http.Request) (lib.Iso8583Message, error) {
	src, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var input bytes.Buffer
	if _, err = io.Copy(&input, src); err != nil {
		return nil, err
	}

	message, _ := lib.NewISO8583Message(&utils.ISO8583DataElementsVer1987)
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
	mf, err := parseInputFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = mf.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	outputBufferToWriter(w, []byte("valid file"), utils.MessageFormatIso8583)
}

// validator - print file with ascii or json format
func print(w http.ResponseWriter, r *http.Request) {
	message, err := parseInputFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	format := r.FormValue("format")
	output, err := messageToBuf(format, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.EqualFold(format, utils.MessageFormatIso8583) ||
		strings.EqualFold(format, utils.MessageFormatJson) ||
		strings.EqualFold(format, utils.MessageFormatXml) {
		outputBufferToWriter(w, output, format)
	} else {
		http.Error(w, "invalid print format", http.StatusBadRequest)
	}
}

// convert - convert file with ascii or json format
func convert(w http.ResponseWriter, r *http.Request) {
	message, err := parseInputFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	format := r.FormValue("format")
	filename := "converted_file"
	output, err := messageToBuf(format, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	data := map[string]bool{"health": true}
	buf, err := json.Marshal(data)
	if err == nil {
		outputBufferToWriter(w, buf, utils.MessageFormatJson)
	}
}

// configure handlers
func ConfigureHandlers(r *mux.Router) error {
	r.HandleFunc("/health", health).Methods("GET")
	r.HandleFunc("/print", print).Methods("POST")
	r.HandleFunc("/validator", validator).Methods("POST")
	r.HandleFunc("/convert", convert).Methods("POST")
	return nil
}
