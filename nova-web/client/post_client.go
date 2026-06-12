package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PostClient struct {
	*HTTPClient
}

type PostResponse struct {
	Id         int64    `json:"id,string"`
	AuthorId   int64    `json:"author_id,string"`
	Title      string   `json:"title"`
	CoverURL   string   `json:"cover_url"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	CategoryId int64    `json:"category_id,string"`
	CreatedAt  string   `json:"created_at"`
	AuthorName string   `json:"author_name"`
}

type PostListResponse struct {
	Total       int32           `json:"total,string"`
	TotalPages  int32           `json:"total_pages"`
	PageSize    int32           `json:"page_size"`
	CurrentPage int32           `json:"current_page"`
	Posts       []*PostResponse `json:"posts"`
}

type CreatePostRequest struct {
	Post PostInput `json:"post"`
}

type UpdatePostRequest struct {
	Post       PostInput `json:"post"`
	UpdateMask string    `json:"update_mask"`
}

type PostInput struct {
	Title    string `json:"title"`
	CoverURL string `json:"cover_url,omitempty"`
	Content  string `json:"content"`
	AuthorId int64  `json:"author_id"`
}

func (c *PostClient) ListPosts(ctx context.Context, page, pageSize int32) (*PostListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/posts?page=%d&page_size=%d", c.baseURL, page, pageSize)
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

	var result PostListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *PostClient) GetPost(ctx context.Context, id int64) (*PostResponse, error) {
	url := fmt.Sprintf("%s/api/v1/posts/%d", c.baseURL, id)
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

	var wrapper struct {
		Post *PostResponse `json:"post"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.Post, nil
}

func (c *PostClient) CreatePost(ctx context.Context, title, content string, authorId int64, coverURL string) (*PostResponse, error) {
	url := fmt.Sprintf("%s/api/v1/posts", c.baseURL)
	reqBody := CreatePostRequest{
		Post: PostInput{
			Title:    title,
			CoverURL: coverURL,
			Content:  content,
			AuthorId: authorId,
		},
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

	var wrapper struct {
		Post *PostResponse `json:"post"`
	}
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.Post, nil
}

func (c *PostClient) UpdatePost(ctx context.Context, id int64, title, content string, coverURL *string) (*PostResponse, error) {
	url := fmt.Sprintf("%s/api/v1/posts/%d", c.baseURL, id)
	reqBody := UpdatePostRequest{
		Post: PostInput{
			Title:   title,
			Content: content,
		},
		UpdateMask: "title,content",
	}
	if coverURL != nil {
		reqBody.Post.CoverURL = *coverURL
		reqBody.UpdateMask += ",cover_url"
	}

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

	var wrapper struct {
		Post *PostResponse `json:"post"`
	}
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.Post, nil
}

func (c *PostClient) DeletePost(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/api/v1/posts/%d", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return parseError(resp.StatusCode, body)
	}

	return nil
}
