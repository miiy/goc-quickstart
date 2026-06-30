package client

import (
	"context"

	apiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

type UserClient struct {
	*HTTPClient
}

func (c *UserClient) GetUser(ctx context.Context, id int64) (*apiclient.User, error) {
	resp, httpResp, err := c.api.UsersAPI.GetUser(openAPIContext(ctx), id).Execute()
	if err != nil {
		return nil, convertError(httpResp, err)
	}
	if resp == nil {
		return nil, nil
	}
	return &resp.User, nil
}

func (c *UserClient) UpdateUser(ctx context.Context, id int64, nickname, email string) (*apiclient.User, error) {
	user := apiclient.NewUserInput()
	user.SetNickname(nickname)
	user.SetEmail(email)
	req := apiclient.NewUpdateUserRequest(*user)
	req.SetUpdateFields([]string{"nickname", "email"})

	resp, httpResp, err := c.api.UsersAPI.UpdateUser(openAPIContext(ctx), id).UpdateUserRequest(*req).Execute()
	if err != nil {
		return nil, convertError(httpResp, err)
	}
	if resp == nil {
		return nil, nil
	}
	return &resp.User, nil
}
