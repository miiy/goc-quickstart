package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type errorResponse struct {
	Error struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// parseError extracts the error message from gateway response.
// Falls back to raw status+body if parsing fails.
func parseError(statusCode int, body []byte) error {
	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
		return fmt.Errorf("%s", errResp.Error.Message)
	}
	return fmt.Errorf("status: %d, body: %s", statusCode, string(body))
}

// HTTPClient wraps the HTTP calls to gateway
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

type accessTokenContextKey struct{}

func WithAccessToken(ctx context.Context, token string) context.Context {
	token = strings.TrimSpace(token)
	if token == "" {
		return ctx
	}
	return context.WithValue(ctx, accessTokenContextKey{}, token)
}

func AccessTokenFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	token, ok := ctx.Value(accessTokenContextKey{}).(string)
	token = strings.TrimSpace(token)
	return token, ok && token != ""
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	if token, ok := AccessTokenFromContext(req.Context()); ok && req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return c.httpClient.Do(req)
}

type Clients struct {
	Post *PostClient
	Auth *AuthClient
}

func NewClients(gatewayAddr string) (*Clients, func(), error) {
	httpClient := &HTTPClient{
		baseURL:    "http://" + gatewayAddr,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	return &Clients{
		Post: &PostClient{HTTPClient: httpClient},
		Auth: &AuthClient{HTTPClient: httpClient},
	}, func() {}, nil
}
