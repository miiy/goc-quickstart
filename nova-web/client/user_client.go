package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserClient struct {
	*HTTPClient
}

type UserResponse struct {
	Id       Int64String `json:"id"`
	Username string      `json:"username"`
	Nickname string      `json:"nickname"`
	Avatar   string      `json:"avatar"`
	Email    string      `json:"email"`
}

type userEnvelope struct {
	User *UserResponse `json:"user"`
}

func (c *UserClient) GetUser(ctx context.Context, id int64) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp.StatusCode, body)
	}

	return decodeUserResponse(body)
}

func (c *UserClient) UpdateProfile(ctx context.Context, id int64, nickname, email string) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d", c.baseURL, id)
	reqBody := struct {
		User struct {
			Nickname string `json:"nickname"`
			Email    string `json:"email"`
		} `json:"user"`
		UpdateMask string `json:"update_mask"`
	}{
		UpdateMask: "nickname,email",
	}
	reqBody.User.Nickname = nickname
	reqBody.User.Email = email

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(bodyBytes))
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

	return decodeUserResponse(respBody)
}

func decodeUserResponse(body []byte) (*UserResponse, error) {
	var envelope userEnvelope
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.User != nil {
		return envelope.User, nil
	}

	var result UserResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
