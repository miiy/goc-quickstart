package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type FileClient struct {
	*HTTPClient
}

type FileResponse struct {
	Id        int64  `json:"id,string"`
	Scene     string `json:"scene"`
	ObjectKey string `json:"object_key"`
	URL       string `json:"url"`
	MimeType  string `json:"mime_type"`
	Size      int64  `json:"size,string"`
}

func (c *FileClient) UploadPostCover(ctx context.Context, filename string, file io.Reader) (*FileResponse, error) {
	return c.uploadFile(ctx, "post_cover", filename, file)
}

func (c *FileClient) UploadAvatar(ctx context.Context, filename string, file io.Reader) (*UserResponse, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("avatar", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/files/upload/avatar", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

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

func (c *FileClient) uploadFile(ctx context.Context, scene, filename string, file io.Reader) (*FileResponse, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("scene", scene); err != nil {
		return nil, err
	}
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/files/upload", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
		File *FileResponse `json:"file"`
	}
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, err
	}
	normalizeFileResponse(wrapper.File)
	return wrapper.File, nil
}

func normalizeFileResponse(file *FileResponse) {
	if file == nil {
		return
	}
	if file.URL == "" {
		file.URL = UploadsURL(file.ObjectKey)
	}
}
