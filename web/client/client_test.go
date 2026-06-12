package client

import (
	"context"
	"encoding/json"
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

func TestPostClientUpdatePostSendsFieldMask(t *testing.T) {
	client := &PostClient{HTTPClient: newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPut)
		}
		if r.URL.Path != "/api/v1/posts/42" {
			t.Fatalf("path = %s, want /api/v1/posts/42", r.URL.Path)
		}

		var body struct {
			Post struct {
				Title   string `json:"title"`
				Content string `json:"content"`
			} `json:"post"`
			UpdateMask string `json:"update_mask"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Post.Title != "new title" || body.Post.Content != "new content" {
			t.Fatalf("unexpected post body: %+v", body.Post)
		}
		if body.UpdateMask != "title,content" {
			t.Fatalf("update_mask = %q, want %q", body.UpdateMask, "title,content")
		}

		resp := testResponse(http.StatusOK, `{"post":{"id":"42","author_id":"7","author_name":"alice","title":"new title","content":"new content"}}`)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})}

	resp, err := client.UpdatePost(context.Background(), 42, "new title", "new content")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Id != 42 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.AuthorName != "alice" {
		t.Fatalf("author name = %q, want alice", resp.AuthorName)
	}
}

func TestAuthClientLoginParsesStringUserID(t *testing.T) {
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

		resp := testResponse(http.StatusOK, `{"access_token":"access-token","token_type":"Bearer","expires_at":"2026-01-01T00:00:00Z","user":{"id":"7","username":"alice"}}`)
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
	if int64(resp.User.Id) != 7 || resp.User.Username != "alice" {
		t.Fatalf("unexpected user: %+v", resp.User)
	}
}

func TestAuthClientChangePassword(t *testing.T) {
	client := &AuthClient{HTTPClient: newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
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

		return testResponse(http.StatusOK, `{}`), nil
	})}

	if err := client.ChangePassword(context.Background(), "old", "new", "new"); err != nil {
		t.Fatal(err)
	}
}

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

		resp := testResponse(http.StatusOK, `{"user":{"id":"7","username":"alice","nickname":"Alice","email":"alice@example.com"}}`)
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
}

func TestUserClientUploadAvatar(t *testing.T) {
	client := &UserClient{HTTPClient: newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/api/v1/uploads/avatar" {
			t.Fatalf("path = %s, want /api/v1/uploads/avatar", r.URL.Path)
		}

		reader, err := r.MultipartReader()
		if err != nil {
			t.Fatal(err)
		}
		form, err := reader.ReadForm(1024)
		if err != nil {
			t.Fatal(err)
		}
		files := form.File["avatar"]
		if len(files) != 1 {
			t.Fatalf("avatar files = %d, want 1", len(files))
		}
		if files[0].Filename != "avatar.png" {
			t.Fatalf("filename = %q, want avatar.png", files[0].Filename)
		}
		f, err := files[0].Open()
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != "png-bytes" {
			t.Fatalf("content = %q, want png-bytes", string(content))
		}

		resp := testResponse(http.StatusOK, `{"user":{"id":"7","username":"alice","avatar":"http://cdn.test/uploads/a.png"}}`)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})}

	resp, err := client.UploadAvatar(context.Background(), "avatar.png", strings.NewReader("png-bytes"))
	if err != nil {
		t.Fatal(err)
	}
	if int64(resp.Id) != 7 || resp.Avatar != "http://cdn.test/uploads/a.png" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
