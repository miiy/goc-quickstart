package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFileClientUploadAvatar(t *testing.T) {
	client := &FileClient{HTTPClient: newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/api/v1/files/upload/avatar" {
			t.Fatalf("path = %s, want /api/v1/files/upload/avatar", r.URL.Path)
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

		resp := testResponse(http.StatusOK, `{"user":{"id":"7","username":"alice","avatar":"avatars/2026/06/a.png"}}`)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})}

	resp, err := client.UploadAvatar(context.Background(), "avatar.png", strings.NewReader("png-bytes"))
	if err != nil {
		t.Fatal(err)
	}
	if int64(resp.Id) != 7 || resp.Avatar != "/uploads/avatars/2026/06/a.png" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestFileClientUploadPostCover(t *testing.T) {
	client := &FileClient{HTTPClient: newTestHTTPClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/api/v1/files/upload" {
			t.Fatalf("path = %s, want /api/v1/files/upload", r.URL.Path)
		}

		reader, err := r.MultipartReader()
		if err != nil {
			t.Fatal(err)
		}
		form, err := reader.ReadForm(1024)
		if err != nil {
			t.Fatal(err)
		}
		if got := form.Value["scene"]; len(got) != 1 || got[0] != "post_cover" {
			t.Fatalf("scene = %+v, want post_cover", got)
		}
		files := form.File["file"]
		if len(files) != 1 {
			t.Fatalf("file count = %d, want 1", len(files))
		}
		if files[0].Filename != "cover.png" {
			t.Fatalf("filename = %q, want cover.png", files[0].Filename)
		}

		resp := testResponse(http.StatusOK, `{"file":{"id":"9","scene":"FILE_SCENE_POST_COVER","object_key":"post-covers/2026/06/a.png","url":"/uploads/post-covers/2026/06/a.png","mime_type":"image/png","size":"9"}}`)
		resp.Header.Set("Content-Type", "application/json")
		return resp, nil
	})}

	resp, err := client.UploadPostCover(context.Background(), "cover.png", strings.NewReader("png-bytes"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.URL != "/uploads/post-covers/2026/06/a.png" {
		t.Fatalf("url = %q", resp.URL)
	}
}
