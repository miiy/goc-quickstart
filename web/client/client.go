package client

import (
	"net/http"
	"time"
)

// HTTPClient wraps the HTTP calls to gateway
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
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
