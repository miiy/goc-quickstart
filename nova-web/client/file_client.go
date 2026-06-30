package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	apiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

type FileClient struct {
	*HTTPClient
}

func (c *FileClient) UploadPostCover(ctx context.Context, filename string, file io.Reader) (*apiclient.File, error) {
	return c.uploadFile(ctx, "post_cover", filename, file)
}

func (c *FileClient) UploadAvatar(ctx context.Context, filename string, file io.Reader) (*apiclient.User, error) {
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

	var wrapper apiclient.UpdateUserResponse
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.User, nil
}

func (c *FileClient) uploadFile(ctx context.Context, scene, filename string, file io.Reader) (*apiclient.File, error) {
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

	var wrapper apiclient.UploadFileResponse
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.File, nil
}
