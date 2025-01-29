package notionapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/qonto/notionapi"
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
	return newTestClient(func(*http.Request) *http.Response {
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

func TestRateLimit(t *testing.T) {
	t.Run("should return error when rate limit is exceeded", func(t *testing.T) {
		c := newTestClient(func(*http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{"Retry-After": []string{"0"}},
			}
		})
		client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c), notionapi.WithRetry(2))
		_, err := client.Block.Get(context.Background(), "some_block_id")
		if err == nil {
			t.Errorf("Get() error = %v", err)
		}
		wantErr := "Retry request with 429 response failed after 2 retries"
		if err.Error() != wantErr {
			t.Errorf("Get() error = %v, wantErr %s", err, wantErr)
		}
	})

	t.Run("should make maxRetries attempts", func(t *testing.T) {
		attempts := 0
		maxRetries := 2
		c := newTestClient(func(*http.Request) *http.Response {
			attempts++
			return &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{"Retry-After": []string{"0"}},
			}
		})
		client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c), notionapi.WithRetry(maxRetries))
		_, err := client.Block.Get(context.Background(), "some_block_id")
		if err == nil {
			t.Errorf("Get() error = %v", err)
		}
		if attempts != maxRetries {
			t.Errorf("Get() attempts = %v, want %v", attempts, maxRetries)
		}
	})
}

func TestBasicAuthHeader(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		auth := request.Header.Get("Authorization")
		if auth != "Basic bXkgaWQgaGVyZTpzZWNyZXQgc2hoaA==" {
			t.Errorf("expected basic auth, got %q", auth)
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	srvURL, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("could not parse test server url")
	}

	c := newTestClient(func(req *http.Request) *http.Response {
		req.URL = srvURL
		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			t.Errorf("failed to make http request: %s", err.Error())
		}
		return resp
	})

	opts := []notionapi.ClientOption{
		notionapi.WithHTTPClient(c),
		notionapi.WithOAuthAppCredentials("my id here", "secret shhh"),
	}
	client := notionapi.NewClient("some_token", opts...)
	_, _ = client.Authentication.CreateToken(context.Background(), &notionapi.TokenCreateRequest{})
}
