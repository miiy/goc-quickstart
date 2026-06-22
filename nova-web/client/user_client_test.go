package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestUserClientUpdateProfile(t *testing.T) {
	client := &UserClient{HTTPClient: newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPut)
		}
		if r.URL.Path != "/api/v1/users/7" {
			t.Fatalf("path = %s, want /api/v1/users/7", r.URL.Path)
		}

		var body struct {
			User struct {
				Nickname string `json:"nickname"`
				Email    string `json:"email"`
			} `json:"user"`
			UpdateMask string `json:"update_mask"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.User.Nickname != "Alice" || body.User.Email != "alice@example.com" {
			t.Fatalf("unexpected user body: %+v", body.User)
		}
		if body.UpdateMask != "nickname,email" {
			t.Fatalf("update_mask = %q, want nickname,email", body.UpdateMask)
		}

		resp := testResponse(http.StatusOK, `{"user":{"id":"7","username":"alice","nickname":"Alice","avatar":"avatars/2026/06/a.png","email":"alice@example.com"}}`)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})}

	resp, err := client.UpdateProfile(context.Background(), 7, "Alice", "alice@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if int64(resp.Id) != 7 || resp.Nickname != "Alice" || resp.Email != "alice@example.com" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.Avatar != "/uploads/avatars/2026/06/a.png" {
		t.Fatalf("avatar = %q", resp.Avatar)
	}
}
