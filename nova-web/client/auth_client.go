package client

import (
	"context"

	apiclient "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-client"
)

type AuthClient struct {
	*HTTPClient
}

type TokenResponse = apiclient.TokenResponse

func (c *AuthClient) Register(ctx context.Context, email, username, password, passwordConfirmation string) (*apiclient.RegisterResponse, error) {
	req := apiclient.NewRegisterRequest(email, username, password, passwordConfirmation)
	resp, httpResp, err := c.api.AuthAPI.Register(openAPIContext(ctx)).RegisterRequest(*req).Execute()
	if err != nil {
		return nil, convertError(httpResp, err)
	}
	return resp, nil
}

func (c *AuthClient) Login(ctx context.Context, username, password string) (*TokenResponse, error) {
	req := apiclient.NewLoginRequest(username, password)
	resp, httpResp, err := c.api.AuthAPI.Login(openAPIContext(ctx)).LoginRequest(*req).Execute()
	if err != nil {
		return nil, convertError(httpResp, err)
	}
	return resp, nil
}

// RefreshToken exchanges a refresh token for a new access + refresh token pair.
// It does not send an Authorization header (the gateway route is public).
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	req := apiclient.NewRefreshTokenRequest(refreshToken)
	resp, httpResp, err := c.api.AuthAPI.RefreshToken(ctx).RefreshTokenRequest(*req).Execute()
	if err != nil {
		return nil, convertError(httpResp, err)
	}
	return resp, nil
}

func (c *AuthClient) Logout(ctx context.Context, accessToken, refreshToken string) error {
	req := apiclient.NewLogoutRequest()
	if accessToken != "" {
		req.SetAccessToken(accessToken)
	}
	if refreshToken != "" {
		req.SetRefreshToken(refreshToken)
	}
	_, httpResp, err := c.api.AuthAPI.Logout(openAPIContext(ctx)).LogoutRequest(*req).Execute()
	if err != nil {
		return convertError(httpResp, err)
	}
	return nil
}

func (c *AuthClient) ChangePassword(ctx context.Context, oldPassword, newPassword, newPasswordConfirmation string) error {
	req := apiclient.NewChangePasswordRequest(oldPassword, newPassword, newPasswordConfirmation)
	_, httpResp, err := c.api.AuthAPI.ChangePassword(openAPIContext(ctx)).ChangePasswordRequest(*req).Execute()
	if err != nil {
		return convertError(httpResp, err)
	}
	return nil
}
