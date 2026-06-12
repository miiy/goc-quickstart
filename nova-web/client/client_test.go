package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestHTTPClient(handler func(*http.Request) (*http.Response, error)) *HTTPClient {
	return &HTTPClient{
		baseURL: "http://gateway.test",
		httpClient: &http.Client{
			Transport: roundTripFunc(handler),
		},
	}
}

func testResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestHTTPClientDoAddsAccessTokenFromContext(t *testing.T) {
	client := newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("Authorization"); got != "Bearer access-token" {
			t.Fatalf("Authorization header = %q, want %q", got, "Bearer access-token")
		}
		return testResponse(http.StatusNoContent, ""), nil
	})

	ctx := WithAccessToken(context.Background(), "access-token")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://gateway.test/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestHTTPClientDoKeepsExistingAuthorizationHeader(t *testing.T) {
	client := newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("Authorization"); got != "Custom credential" {
			t.Fatalf("Authorization header = %q, want %q", got, "Custom credential")
		}
		return testResponse(http.StatusNoContent, ""), nil
	})

	ctx := WithAccessToken(context.Background(), "access-token")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://gateway.test/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Custom credential")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestParseErrorKeepsHTTPStatus(t *testing.T) {
	err := parseError(http.StatusUnauthorized, []byte(`{"error":{"code":16,"message":"invalid auth token"}}`))

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected HTTPError, got %T", err)
	}
	if httpErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", httpErr.StatusCode, http.StatusUnauthorized)
	}
	if err.Error() != "invalid auth token" {
		t.Fatalf("message = %q, want invalid auth token", err.Error())
	}
	if !IsStatus(err, http.StatusUnauthorized) {
		t.Fatal("expected IsStatus to match unauthorized")
	}
}
