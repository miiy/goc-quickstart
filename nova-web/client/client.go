package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	apiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

// HTTPError preserves the upstream HTTP status while keeping the error payload
// in the same shape as the OpenAPI error response: {"error": {...}}.
type HTTPError struct {
	StatusCode int
	Response   apiclient.ErrorResponse
	Body       string
}

func NewHTTPError(statusCode int, message string) *HTTPError {
	return newHTTPError(statusCode, int32(statusCode), message, nil)
}

func newHTTPError(statusCode int, code int32, message string, body []byte) *HTTPError {
	message = strings.TrimSpace(message)
	if message == "" {
		message = http.StatusText(statusCode)
	}
	if message == "" {
		message = fmt.Sprintf("status: %d", statusCode)
	}
	return &HTTPError{
		StatusCode: statusCode,
		Response: apiclient.ErrorResponse{
			Error: apiclient.ErrorStatus{
				Code:    code,
				Message: message,
			},
		},
		Body: string(body),
	}
}

func (e *HTTPError) Error() string {
	if e == nil {
		return ""
	}
	if e.Response.Error.Message != "" {
		return e.Response.Error.Message
	}
	if e.Body == "" {
		return fmt.Sprintf("status: %d", e.StatusCode)
	}
	return fmt.Sprintf("status: %d, body: %s", e.StatusCode, e.Body)
}

func (e *HTTPError) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	if body := []byte(e.Body); json.Valid(body) {
		return body, nil
	}
	return json.Marshal(e.Response)
}

func IsStatus(err error, statusCode int) bool {
	var httpErr *HTTPError
	return errors.As(err, &httpErr) && httpErr.StatusCode == statusCode
}

// parseError extracts the error message from gateway response.
// Falls back to raw status+body if parsing fails.
func parseError(statusCode int, body []byte) error {
	var errResp apiclient.ErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
		return &HTTPError{
			StatusCode: statusCode,
			Response:   errResp,
			Body:       string(body),
		}
	}
	return newHTTPError(statusCode, int32(statusCode), string(body), body)
}

func convertError(resp *http.Response, err error) error {
	if err == nil {
		return nil
	}

	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}

	var openAPIErr *apiclient.GenericOpenAPIError
	if errors.As(err, &openAPIErr) && openAPIErr != nil {
		return convertOpenAPIError(statusCode, openAPIErr)
	}

	if statusCode > 0 {
		return NewHTTPError(statusCode, err.Error())
	}
	return err
}

func convertOpenAPIError(statusCode int, openAPIErr *apiclient.GenericOpenAPIError) error {
	if model, ok := openAPIErr.Model().(apiclient.ErrorResponse); ok && model.Error.Message != "" {
		return &HTTPError{
			StatusCode: statusCode,
			Response:   model,
			Body:       string(openAPIErr.Body()),
		}
	}
	if model, ok := openAPIErr.Model().(*apiclient.ErrorResponse); ok && model != nil && model.Error.Message != "" {
		return &HTTPError{
			StatusCode: statusCode,
			Response:   *model,
			Body:       string(openAPIErr.Body()),
		}
	}
	if body := openAPIErr.Body(); statusCode >= 300 && len(body) > 0 {
		return parseError(statusCode, body)
	}
	if statusCode >= 300 {
		return NewHTTPError(statusCode, openAPIErr.Error())
	}
	return errors.New(openAPIErr.Error())
}

// HTTPClient wraps the HTTP calls to gateway
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	api        *apiclient.APIClient
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

func openAPIContext(ctx context.Context) context.Context {
	if token, ok := AccessTokenFromContext(ctx); ok {
		return context.WithValue(ctx, apiclient.ContextAccessToken, token)
	}
	return ctx
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
	User *UserClient
	File *FileClient
}

func NewClients(gatewayAddr string) (*Clients, func(), error) {
	baseURL := gatewayBaseURL(gatewayAddr)
	rawHTTPClient := &http.Client{Timeout: 30 * time.Second}
	cfg := apiclient.NewConfiguration()
	cfg.HTTPClient = rawHTTPClient
	cfg.Servers = apiclient.ServerConfigurations{
		{
			URL: baseURL + "/api/v1",
		},
	}

	httpClient := &HTTPClient{
		baseURL:    baseURL,
		httpClient: rawHTTPClient,
		api:        apiclient.NewAPIClient(cfg),
	}

	return &Clients{
		Post: &PostClient{HTTPClient: httpClient},
		Auth: &AuthClient{HTTPClient: httpClient},
		User: &UserClient{HTTPClient: httpClient},
		File: &FileClient{HTTPClient: httpClient},
	}, func() {}, nil
}

func gatewayBaseURL(gatewayAddr string) string {
	gatewayAddr = strings.TrimRight(strings.TrimSpace(gatewayAddr), "/")
	if gatewayAddr == "" {
		return "http://"
	}
	if u, err := url.Parse(gatewayAddr); err == nil && u.Scheme != "" {
		return gatewayAddr
	}
	return "http://" + gatewayAddr
}
