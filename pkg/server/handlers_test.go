// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package server_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/moov-io/iso8583/pkg/server"
	"github.com/moov-io/iso8583/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testFileName = "iso_reversal_message.dat"
var testInvalidFileName = "iso_reversal_message_error_date.dat"
var testErrorFileName = "error_message.dat"

type HandlersTest struct {
	suite.Suite
	testServer *mux.Router
}

func (suite *HandlersTest) makeRequest(method, url, body string) (*httptest.ResponseRecorder, *http.Request) {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	assert.Equal(suite.T(), nil, err)
	recorder := httptest.NewRecorder()
	return recorder, request
}

func (suite *HandlersTest) getWriter(name string) (*multipart.Writer, *bytes.Buffer) {
	path := filepath.Join("..", "..", "test", "testdata", name)
	file, err := os.Open(path)
	assert.Equal(suite.T(), nil, err)
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	assert.Equal(suite.T(), nil, err)
	_, err = io.Copy(part, file)
	assert.Equal(suite.T(), nil, err)
	return writer, body
}

func (suite *HandlersTest) getErrWriter(name string) (*multipart.Writer, *bytes.Buffer) {
	path := filepath.Join("..", "..", "test", "testdata", name)
	file, err := os.Open(path)
	assert.Equal(suite.T(), nil, err)
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("err", filepath.Base(path))
	assert.Equal(suite.T(), nil, err)
	_, err = io.Copy(part, file)
	assert.Equal(suite.T(), nil, err)
	return writer, body
}

func (suite *HandlersTest) SetupTest() {
	var err error
	suite.testServer = mux.NewRouter()
	err = server.ConfigureHandlers(suite.testServer)
	assert.Equal(suite.T(), nil, err)
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTest))
}

func (suite *HandlersTest) TestUnknownRequest() {
	recorder, request := suite.makeRequest(http.MethodGet, "/unknown", "")
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)
}

func (suite *HandlersTest) TestHealth() {
	recorder, request := suite.makeRequest(http.MethodGet, "/health", "")
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestJsonPrint() {
	writer, body := suite.getWriter(testFileName)
	err := writer.WriteField("format", utils.MessageFormatJson)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestXmlPrint() {
	writer, body := suite.getWriter(testFileName)
	err := writer.WriteField("format", utils.MessageFormatXml)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestMessagePrint() {
	writer, body := suite.getWriter(testFileName)
	err := writer.WriteField("format", utils.MessageFormatIso8583)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestJsonConvert() {
	writer, body := suite.getWriter(testFileName)
	err := writer.WriteField("format", utils.MessageFormatJson)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/convert", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestXmlConvert() {
	writer, body := suite.getWriter(testFileName)
	err := writer.WriteField("format", utils.MessageFormatXml)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/convert", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestMessageConvert() {
	writer, body := suite.getWriter(testFileName)
	err := writer.WriteField("format", utils.MessageFormatIso8583)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/convert", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestValidator() {
	writer, body := suite.getWriter(testFileName)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/validator", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestPrintWithInvalidForm() {
	writer, body := suite.getErrWriter(testFileName)
	err := writer.WriteField("format", utils.MessageFormatJson)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *HandlersTest) TestConvertWithInvalidForm() {
	writer, body := suite.getErrWriter(testFileName)
	err := writer.WriteField("format", utils.MessageFormatJson)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/convert", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *HandlersTest) TestConvertWithInvalidData() {
	writer, body := suite.getWriter(testInvalidFileName)
	err := writer.WriteField("format", utils.MessageFormatJson)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/convert", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestPrintWithInvalidData() {
	writer, body := suite.getWriter(testInvalidFileName)
	err := writer.WriteField("format", utils.MessageFormatJson)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *HandlersTest) TestValidatorWithInvalidData() {
	writer, body := suite.getWriter(testInvalidFileName)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/validator", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusNotImplemented, recorder.Code)
}

func (suite *HandlersTest) TestPrintWithErrorData() {
	writer, body := suite.getWriter(testErrorFileName)
	err := writer.WriteField("format", utils.MessageFormatJson)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/print", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *HandlersTest) TestConvertWithErrorData() {
	writer, body := suite.getWriter(testErrorFileName)
	err := writer.WriteField("format", utils.MessageFormatJson)
	assert.Equal(suite.T(), nil, err)
	err = writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/convert", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *HandlersTest) TestValidatorWithErrorData() {
	writer, body := suite.getWriter(testErrorFileName)
	err := writer.Close()
	assert.Equal(suite.T(), nil, err)
	recorder, request := suite.makeRequest(http.MethodPost, "/validator", body.String())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	suite.testServer.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}
