package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type AuthClient struct {
	*HTTPClient
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
	User        struct {
		Id       Int64String `json:"id"`
		Username string      `json:"username"`
	} `json:"user"`
}

type Int64String int64

func (v *Int64String) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*v = Int64String(n)
		return nil
	}

	var n int64
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*v = Int64String(n)
	return nil
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

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp.StatusCode, respBody)
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

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp.StatusCode, respBody)
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

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return parseError(resp.StatusCode, respBody)
	}

	return nil
}

func (c *AuthClient) ChangePassword(ctx context.Context, oldPassword, newPassword, newPasswordConfirmation string) error {
	url := fmt.Sprintf("%s/api/v1/auth/password", c.baseURL)
	reqBody := map[string]string{
		"old_password":              oldPassword,
		"new_password":              newPassword,
		"new_password_confirmation": newPasswordConfirmation,
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

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return parseError(resp.StatusCode, respBody)
	}

	return nil
}
