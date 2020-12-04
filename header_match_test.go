package checkheadersplugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	checkheaders "github.com/dkijkuit/checkheadersplugin"
)

var required = true
var not_required = false
var contains = true

func TestHeadersMatch(t *testing.T) {
	requestHeaders := map[string]string{
		"test1": "testvalue1",
		"test2": "testvalue2",
		"test3": "testvalue3",
		"test4": "value4",
	}

	executeTest(t, requestHeaders, http.StatusOK)
}
func TestHeadersNotMatch(t *testing.T) {
	requestHeaders := map[string]string{
		"test1": "wrongvalue1",
		"test2": "wrongvalue2",
		"test3": "wrongvalue3",
		"test4": "correctvalue4",
	}

	executeTest(t, requestHeaders, http.StatusForbidden)
}

func TestHeadersNotRequired(t *testing.T) {
	requestHeaders := map[string]string{
		"test1": "testvalue1",
		"test2": "testvalue2",
		"test4": "ue4",
	}

	executeTest(t, requestHeaders, http.StatusOK)
}

func executeTest(t *testing.T, requestHeaders map[string]string, expectedResultCode int) {
	cfg := checkheaders.CreateConfig()
	cfg.Headers = []checkheaders.SingleHeader{
		{
			Name:  "test1",
			Value: "testvalue1",
		},
		{
			Name:     "test2",
			Value:    "testvalue2",
			Required: &required,
		},
		{
			Name:     "test3",
			Value:    "testvalue3",
			Required: &not_required,
		},
		{
			Name:     "test4",
			Value:    "ue4",
			Required: &required,
			Contains: &contains,
		},
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := checkheaders.New(ctx, next, cfg, "check-headers-plugin")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	for headerName, headerValue := range requestHeaders {
		req.Header.Add(headerName, headerValue)
	}

	handler.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != expectedResultCode {
		t.Errorf("Unexpected response status code: %d, expected: %d", recorder.Result().StatusCode, expectedResultCode)
	}
}
