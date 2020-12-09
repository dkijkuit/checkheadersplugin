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
var urlDecode = true

var testcert = `
Subject%3D%22C%3DNL%2CST%3DST-TEST%2CL%3DCity%2CO%3DOrganization%2CCN%3Dcommon-name%22%3BIssuer%3D%22DC%3Dnl%2CDC%3Ddomainpart1%2CDC%3Ddomainpart2%2CCN%3DSomeKindOfCa%22%3BNB%3D%221589744159%22%3BNA%3D%221765837153%22%3BSAN%3D%22somkindofdomain.domain.thing.test%22
`

func TestHeadersMatch(t *testing.T) {
	requestHeaders := map[string]string{
		"test1":                            "testvalue1",
		"test2":                            "testvalue2",
		"test3":                            "testvalue3",
		"test4":                            "value4",
		"X-Forwarded-Tls-Client-Cert-Info": testcert,
		"testMultipleContainsValues":       "value5_or_value1_or_value_2_or_value_3",
	}

	executeTest(t, requestHeaders, http.StatusOK)
}

func TestHeadersOneMatch(t *testing.T) {
	requestHeaders := map[string]string{
		"test1":                            "testvalue1",
		"test2":                            "testvalue2",
		"test3":                            "testvalue3",
		"test4":                            "value4",
		"X-Forwarded-Tls-Client-Cert-Info": testcert,
		"testMultipleContainsValues":       "test_or_value2",
	}

	executeTest(t, requestHeaders, http.StatusOK)
}

func TestHeadersNotMatch(t *testing.T) {
	requestHeaders := map[string]string{
		"test1":                            "wrongvalue1",
		"test2":                            "wrongvalue2",
		"test3":                            "wrongvalue3",
		"test4":                            "correctvalue4",
		"X-Forwarded-Tls-Client-Cert-Info": "wrongvalue",
		"testMultipleContainsValues":       "wrongvalues",
	}

	executeTest(t, requestHeaders, http.StatusForbidden)
}

func TestHeadersNotRequired(t *testing.T) {
	requestHeaders := map[string]string{
		"test1":                            "testvalue1",
		"test2":                            "testvalue2",
		"test4":                            "ue4",
		"X-Forwarded-Tls-Client-Cert-Info": testcert,
		"testMultipleContainsValues":       "value5_or_value1_or_value_2_or_value_3",
	}

	executeTest(t, requestHeaders, http.StatusOK)
}

func executeTest(t *testing.T, requestHeaders map[string]string, expectedResultCode int) {
	cfg := checkheaders.CreateConfig()
	cfg.Headers = []checkheaders.SingleHeader{
		{
			Name:      "test1",
			MatchType: string(checkheaders.MatchOne),
			Values:    []string{"testvalue1"},
		},
		{
			Name:      "test2",
			MatchType: string(checkheaders.MatchOne),
			Values:    []string{"testvalue2"},
			Required:  &required,
		},
		{
			Name:      "test3",
			MatchType: string(checkheaders.MatchOne),
			Values:    []string{"testvalue3"},
			Required:  &not_required,
		},
		{
			Name:      "test4",
			MatchType: string(checkheaders.MatchOne),
			Values:    []string{"ue4"},
			Required:  &required,
			Contains:  &contains,
		},
		{
			Name: "X-Forwarded-Tls-Client-Cert-Info",
			Values: []string{
				"CN=common-name",
				"SAN=\"somkindofdomain.domain.thing.test\"",
			},
			MatchType: string(checkheaders.MatchAll),
			Required:  &required,
			Contains:  &contains,
			URLDecode: &urlDecode,
		},
		{
			Name: "testMultipleContainsValues",
			Values: []string{
				"value1",
				"or_value2",
			},
			MatchType: string(checkheaders.MatchOne),
			Required:  &required,
			Contains:  &contains,
			URLDecode: &urlDecode,
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
