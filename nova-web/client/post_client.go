package client

import (
	"context"

	apiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

type PostClient struct {
	*HTTPClient
}

func (c *PostClient) ListPosts(ctx context.Context, page, pageSize int32) (*apiclient.ListPostsResponse, error) {
	resp, httpResp, err := c.api.PostsAPI.ListPosts(openAPIContext(ctx)).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		return nil, convertError(httpResp, err)
	}
	if resp == nil {
		return nil, nil
	}
	return resp, nil
}

func (c *PostClient) GetPost(ctx context.Context, id int64) (*apiclient.Post, error) {
	resp, httpResp, err := c.api.PostsAPI.GetPost(openAPIContext(ctx), id).Execute()
	if err != nil {
		return nil, convertError(httpResp, err)
	}
	if resp == nil {
		return nil, nil
	}
	return &resp.Post, nil
}

func (c *PostClient) CreatePost(ctx context.Context, title, content string, coverURL string) (*apiclient.Post, error) {
	post := apiclient.NewCreatePostInput(title, content)
	if coverURL != "" {
		post.SetCoverUrl(coverURL)
	}
	req := apiclient.NewCreatePostRequest(*post)
	resp, httpResp, err := c.api.PostsAPI.CreatePost(openAPIContext(ctx)).CreatePostRequest(*req).Execute()
	if err != nil {
		return nil, convertError(httpResp, err)
	}
	if resp == nil {
		return nil, nil
	}
	return &resp.Post, nil
}

func (c *PostClient) UpdatePost(ctx context.Context, id int64, title, content string, coverURL *string) (*apiclient.Post, error) {
	post := apiclient.NewUpdatePostInput()
	post.SetTitle(title)
	post.SetContent(content)
	updateFields := []string{"title", "content"}
	if coverURL != nil {
		post.SetCoverUrl(*coverURL)
		updateFields = append(updateFields, "cover_url")
	}
	req := apiclient.NewUpdatePostRequest(*post)
	req.SetUpdateFields(updateFields)
	resp, httpResp, err := c.api.PostsAPI.UpdatePost(openAPIContext(ctx), id).UpdatePostRequest(*req).Execute()
	if err != nil {
		return nil, convertError(httpResp, err)
	}
	if resp == nil {
		return nil, nil
	}
	return &resp.Post, nil
}

func (c *PostClient) DeletePost(ctx context.Context, id int64) error {
	_, httpResp, err := c.api.PostsAPI.DeletePost(openAPIContext(ctx), id).Execute()
	if err != nil {
		return convertError(httpResp, err)
	}
	return nil
}
