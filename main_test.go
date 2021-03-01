package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var testJSON = "{\"store\":{\"book\":[{\"category\":\"reference\",\"author\":\"Nigel Rees\",\"title\":\"Sayings of the Century\",\"price\":8.95},{\"category\":\"fiction\",\"author\":\"Evelyn Waugh\",\"title\":\"Sword of Honour\",\"price\":12.99},{\"category\":\"fiction\",\"author\":\"Herman Melville\",\"title\":\"Moby Dick\",\"isbn\":\"0-553-21311-3\",\"price\":8.99},{\"category\":\"fiction\",\"author\":\"J. R. R. Tolkien\",\"title\":\"The Lord of the Rings\",\"isbn\":\"0-395-19395-8\",\"price\":22.99}],\"bicycle\":{\"color\":\"red\",\"price\":19.95}},\"expensive\":10}"

var probeTests = []struct {
	inData      string
	inField     string
	outHttpCode int
	outValue    string
}{
	{"{\"field\": 23}", "$.field", 200, "field 23"},
	{"{\"field\": 19}", "$.field", 200, "field 19"},
	{"{\"field\": true}", "$.field", 200, "field 1"},
	{"{\"field\": false}", "$.field", 200, "field 0"},
	{"{\"field\": 19}", "$.undefined", 404, "Jsonpath not found"},
	{testJSON, "$.store.book.price", 200, "store_book_price 4"},
}

func TestProbeHandler(t *testing.T) {

	for _, tt := range probeTests {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, tt.inData)
		}))
		defer ts.Close()

		u := fmt.Sprintf("http://example.com/probe?target=%s&jsonpath=%s", url.QueryEscape(ts.URL), url.QueryEscape(tt.inField))

		req := httptest.NewRequest("GET", u, nil)
		w := httptest.NewRecorder()

		probeHandler(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if tt.outHttpCode != resp.StatusCode {
			t.Error(fmt.Sprintf("HTTP Code mismatch - %d expected %d", resp.StatusCode, tt.outHttpCode))
		}

		if !strings.Contains(string(body), tt.outValue) {
			t.Error(fmt.Sprintf("Expected output: %s got\n%s", tt.outValue, string(body)))
		}
	}
}
