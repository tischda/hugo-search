package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testIndexName = "test.bleve"

type Response struct {
	Status  string   `json:"status"`
	Indexes []string `json:"indexes"`
}

func TestHttpServer(t *testing.T) {

	// prepare index
	buildIndexFromSite(testHugoPath, testIndexPath)
	index := registerIndex(testIndexPath, testIndexName)
	defer index.Close()

	// http recorder
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://localhost/api", nil)

	// http handler
	handler := getCorsHandler(testIndexName)
	handler.ServeHTTP(recorder, request)

	expected := testIndexName

	rawJSON := recorder.Body.String()
	var response *Response
	json.Unmarshal([]byte(rawJSON), &response)
	actual := response.Indexes[0]

	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}
