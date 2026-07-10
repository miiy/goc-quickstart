package auth

import (
	"bytes"
	"io"
	"net/http"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"

	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/jwtauth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthAPI struct {
	authClient authv1.AuthServiceClient
}

func NewAuthAPI(authClient authv1.AuthServiceClient) openapi.AuthAPI {
	return &AuthAPI{authClient: authClient}
}

func (api *AuthAPI) Register(c *gin.Context) {
	var req openapi.RegisterRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	resp, err := api.authClient.Register(c.Request.Context(), &authv1.RegisterRequest{
		Email:                req.Email,
		Username:             req.Username,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.RegisterResponse{User: protoToAuthUser(resp.GetUser())})
}

func (api *AuthAPI) UsernameCheck(c *gin.Context) {
	var req openapi.UsernameCheckRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	resp, err := api.authClient.UsernameCheck(c.Request.Context(), &authv1.UsernameCheckRequest{Value: req.Value})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.UsernameCheckResponse{Exist: resp.GetExist()})
}

func (api *AuthAPI) EmailCheck(c *gin.Context) {
	var req openapi.EmailCheckRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	resp, err := api.authClient.EmailCheck(c.Request.Context(), &authv1.EmailCheckRequest{Value: req.Value})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.EmailCheckResponse{Exist: resp.GetExist()})
}

func (api *AuthAPI) PhoneCheck(c *gin.Context) {
	var req openapi.PhoneCheckRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	resp, err := api.authClient.PhoneCheck(c.Request.Context(), &authv1.PhoneCheckRequest{Value: req.Value})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, openapi.PhoneCheckResponse{Exist: resp.GetExist()})
}

func (api *AuthAPI) Login(c *gin.Context) {
	var req openapi.LoginRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	resp, err := api.authClient.Login(c.Request.Context(), &authv1.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, protoToTokenResponse(resp))
}

func (api *AuthAPI) SendSmsCode(c *gin.Context) {
	var req openapi.SendSmsCodeRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	_, err := api.authClient.SendSmsCode(c.Request.Context(), &authv1.SendSmsCodeRequest{Phone: req.Phone})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]any{})
}

func (api *AuthAPI) PhoneAuth(c *gin.Context) {
	var req openapi.PhoneAuthRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	resp, err := api.authClient.PhoneAuth(c.Request.Context(), &authv1.PhoneAuthRequest{
		Phone: req.Phone,
		Code:  req.Code,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, protoToTokenResponse(resp))
}

func (api *AuthAPI) MpLogin(c *gin.Context) {
	var req openapi.MpLoginRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	resp, err := api.authClient.MpLogin(c.Request.Context(), &authv1.MpLoginRequest{Code: req.Code})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, protoToTokenResponse(resp))
}

func (api *AuthAPI) RefreshToken(c *gin.Context) {
	var req openapi.RefreshTokenRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	resp, err := api.authClient.RefreshToken(c.Request.Context(), &authv1.RefreshTokenRequest{RefreshToken: req.RefreshToken})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, protoToTokenResponse(resp))
}

func (api *AuthAPI) ChangePassword(c *gin.Context) {
	var req openapi.ChangePasswordRequest
	if !transport.BindJSON(c, &req) {
		return
	}

	_, err := api.authClient.ChangePassword(c.Request.Context(), &authv1.ChangePasswordRequest{
		OldPassword:             req.OldPassword,
		NewPassword:             req.NewPassword,
		NewPasswordConfirmation: req.NewPasswordConfirmation,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]any{})
}

func (api *AuthAPI) Logout(c *gin.Context) {
	req, ok := bindOptionalLogoutRequest(c)
	if !ok {
		return
	}
	if req.AccessToken == "" {
		req.AccessToken, _ = jwtauth.Token(c)
	}

	_, err := api.authClient.Logout(c.Request.Context(), &authv1.LogoutRequest{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]any{})
}

func bindOptionalLogoutRequest(c *gin.Context) (openapi.LogoutRequest, bool) {
	var req openapi.LogoutRequest
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		transport.WriteError(c, status.Error(codes.InvalidArgument, err.Error()))
		return req, false
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return req, true
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return req, transport.BindJSON(c, &req)
}
