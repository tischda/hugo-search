package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"
)

var testCfg = &Config{
	hugoPath:  "test",
	indexPath: "test/indexes/search.bleve",
	verbose:   false,
}

const TEST_INDEX_NAME = "test.bleve"

type Response struct {
	Status  string   `json:"status"`
	Indexes []string `json:"indexes"`
}

func TestHttpServer(t *testing.T) {

	// prepare index
	buildIndexFromSite(testCfg)

	index := registerIndex(testCfg.indexPath, TEST_INDEX_NAME, testCfg.verbose)
	defer unregisterIndex(index, TEST_INDEX_NAME)

	// http recorder
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://localhost/api", nil)

	// http handler
	handler := getCorsHandler(TEST_INDEX_NAME)
	handler.ServeHTTP(recorder, request)

	expected := TEST_INDEX_NAME

	rawJSON := recorder.Body.String()
	var response *Response
	json.Unmarshal([]byte(rawJSON), &response)
	actual := response.Indexes[0]

	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}
