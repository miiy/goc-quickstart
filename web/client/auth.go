package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthClient struct {
	*HTTPClient
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
	User        struct {
		Username string `json:"username"`
	} `json:"user"`
}

type RegisterResponse struct {
	User struct {
		Username string `json:"username"`
	} `json:"user"`
}

func (c *AuthClient) Register(ctx context.Context, email, username, password, passwordConfirmation string) (*RegisterResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/register", c.baseURL)
	reqBody := map[string]string{
		"email":                 email,
		"username":              username,
		"password":              password,
		"password_confirmation": passwordConfirmation,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result RegisterResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AuthClient) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	url := fmt.Sprintf("%s/api/v1/auth/login", c.baseURL)
	reqBody := map[string]string{
		"username": username,
		"password": password,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result LoginResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AuthClient) Logout(ctx context.Context, accessToken string) error {
	url := fmt.Sprintf("%s/api/v1/auth/logout", c.baseURL)
	reqBody := map[string]string{
		"access_token": accessToken,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	return nil
}