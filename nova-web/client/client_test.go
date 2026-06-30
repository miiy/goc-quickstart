package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	apiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestHTTPClient(handler func(*http.Request) (*http.Response, error)) *HTTPClient {
	httpClient := &http.Client{
		Transport: roundTripFunc(handler),
	}
	cfg := apiclient.NewConfiguration()
	cfg.HTTPClient = httpClient
	cfg.Servers = apiclient.ServerConfigurations{
		{
			URL: "http://gateway.test/api/v1",
		},
	}
	return &HTTPClient{
		baseURL:    "http://gateway.test",
		httpClient: httpClient,
		api:        apiclient.NewAPIClient(cfg),
	}
}

func testResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
	if httpErr.Response.Error.Code != 16 {
		t.Fatalf("error code = %d, want 16", httpErr.Response.Error.Code)
	}
	if !IsStatus(err, http.StatusUnauthorized) {
		t.Fatal("expected IsStatus to match unauthorized")
	}

	body, err := json.Marshal(httpErr)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `{"error":{"code":16,"message":"invalid auth token"}}` {
		t.Fatalf("body = %s", body)
	}
}

func TestNewHTTPErrorUsesOpenAPIErrorShape(t *testing.T) {
	err := NewHTTPError(http.StatusUnauthorized, "unauthenticated")

	if err.Error() != "unauthenticated" {
		t.Fatalf("message = %q, want unauthenticated", err.Error())
	}
	if err.Response.Error.Code != http.StatusUnauthorized {
		t.Fatalf("error code = %d, want %d", err.Response.Error.Code, http.StatusUnauthorized)
	}

	body, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	if string(body) != `{"error":{"code":401,"message":"unauthenticated"}}` {
		t.Fatalf("body = %s", body)
	}
}
