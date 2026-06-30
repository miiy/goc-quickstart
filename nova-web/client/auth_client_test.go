package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAuthClientLoginParsesUserID(t *testing.T) {
	client := &AuthClient{HTTPClient: newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/api/v1/auth/login" {
			t.Fatalf("path = %s, want /api/v1/auth/login", r.URL.Path)
		}

		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Username != "alice" || body.Password != "secret" {
			t.Fatalf("unexpected login body: %+v", body)
		}

		resp := testResponse(http.StatusOK, `{"access_token":"access-token","token_type":"Bearer","expires_at":"2026-01-01T00:00:00Z","refresh_token":"refresh-token","refresh_expires_at":"2026-01-08T00:00:00Z","user":{"id":7,"username":"alice"}}`)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})}

	resp, err := client.Login(context.Background(), "alice", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if resp.AccessToken != "access-token" {
		t.Fatalf("access token = %q, want access-token", resp.AccessToken)
	}
	if resp.User.Id != 7 || resp.User.Username != "alice" {
		t.Fatalf("unexpected user: %+v", resp.User)
	}
}

func TestAuthClientChangePassword(t *testing.T) {
	client := &AuthClient{HTTPClient: newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPut)
		}
		if r.URL.Path != "/api/v1/auth/password" {
			t.Fatalf("path = %s, want /api/v1/auth/password", r.URL.Path)
		}

		var body struct {
			OldPassword             string `json:"old_password"`
			NewPassword             string `json:"new_password"`
			NewPasswordConfirmation string `json:"new_password_confirmation"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.OldPassword != "old" || body.NewPassword != "new" || body.NewPasswordConfirmation != "new" {
			t.Fatalf("unexpected password body: %+v", body)
		}

		resp := testResponse(http.StatusOK, `{}`)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})}

	if err := client.ChangePassword(context.Background(), "old", "new", "new"); err != nil {
		t.Fatal(err)
	}
}
