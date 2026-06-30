package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPostClientCreatePostDoesNotSendAuthorID(t *testing.T) {
	client := &PostClient{HTTPClient: newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/api/v1/posts" {
			t.Fatalf("path = %s, want /api/v1/posts", r.URL.Path)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(body), "author_id") {
			t.Fatalf("request body contains author_id: %s", body)
		}

		var decoded struct {
			Post struct {
				Title    string `json:"title"`
				CoverUrl string `json:"cover_url"`
				Content  string `json:"content"`
			} `json:"post"`
		}
		if err := json.Unmarshal(body, &decoded); err != nil {
			t.Fatal(err)
		}
		if decoded.Post.Title != "new title" || decoded.Post.Content != "new content" || decoded.Post.CoverUrl != "post-covers/2026/06/cover.png" {
			t.Fatalf("unexpected post body: %+v", decoded.Post)
		}

		resp := testResponse(http.StatusOK, `{"post":{"id":42,"author_id":7,"author_name":"alice","title":"new title","cover_url":"post-covers/2026/06/cover.png","content":"new content","status":"published","tags":[],"category_id":0,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}}`)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})}

	resp, err := client.CreatePost(context.Background(), "new title", "new content", "post-covers/2026/06/cover.png")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.AuthorId != 7 {
		t.Fatalf("unexpected response: %+v", resp)
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
				Title    string `json:"title"`
				CoverUrl string `json:"cover_url"`
				Content  string `json:"content"`
			} `json:"post"`
			UpdateFields []string `json:"update_fields"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Post.Title != "new title" || body.Post.Content != "new content" || body.Post.CoverUrl != "post-covers/2026/06/cover.png" {
			t.Fatalf("unexpected post body: %+v", body.Post)
		}
		if want := []string{"title", "content", "cover_url"}; !stringSlicesEqual(body.UpdateFields, want) {
			t.Fatalf("update_fields = %v, want %v", body.UpdateFields, want)
		}

		resp := testResponse(http.StatusOK, `{"post":{"id":42,"author_id":7,"author_name":"alice","title":"new title","cover_url":"post-covers/2026/06/cover.png","content":"new content","status":"published","tags":[],"category_id":0,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}}`)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})}

	coverURL := "post-covers/2026/06/cover.png"
	resp, err := client.UpdatePost(context.Background(), 42, "new title", "new content", &coverURL)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Id != 42 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.AuthorName != "alice" {
		t.Fatalf("author name = %q, want alice", resp.AuthorName)
	}
	if resp.CoverUrl != "post-covers/2026/06/cover.png" {
		t.Fatalf("cover url = %q", resp.CoverUrl)
	}
}
