package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

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
				CoverURL string `json:"cover_url"`
				Content  string `json:"content"`
			} `json:"post"`
			UpdateMask string `json:"update_mask"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Post.Title != "new title" || body.Post.Content != "new content" || body.Post.CoverURL != "post-covers/2026/06/cover.png" {
			t.Fatalf("unexpected post body: %+v", body.Post)
		}
		if body.UpdateMask != "title,content,cover_url" {
			t.Fatalf("update_mask = %q, want %q", body.UpdateMask, "title,content,cover_url")
		}

		resp := testResponse(http.StatusOK, `{"post":{"id":"42","author_id":"7","author_name":"alice","title":"new title","cover_url":"post-covers/2026/06/cover.png","content":"new content"}}`)
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
	if resp.CoverURL != "/uploads/post-covers/2026/06/cover.png" {
		t.Fatalf("cover url = %q", resp.CoverURL)
	}
}
