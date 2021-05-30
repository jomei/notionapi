package notionapi_test

import (
	"net/http"
	"os"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// newTestClient returns *http.Client with Transport replaced to avoid making real calls
func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

// newMockedClient returns *http.Client which responds with content from given file
func newMockedClient(t *testing.T, requestMockFile string, statusCode int) *http.Client {
	return newTestClient(func(req *http.Request) *http.Response {
		b, err := os.Open(requestMockFile)
		if err != nil {
			t.Fatal(err)
		}

		resp := &http.Response{
			StatusCode: statusCode,
			Body:       b,
			Header:     make(http.Header),
		}
		return resp
	})
}
