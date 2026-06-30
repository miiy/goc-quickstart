package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	"github.com/miiy/goc/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakeAuthClient struct {
	loginReq  *authv1.LoginRequest
	logoutReq *authv1.LogoutRequest
}

func (f *fakeAuthClient) VerifyToken(ctx context.Context, in *authv1.VerifyTokenRequest, opts ...grpc.CallOption) (*authv1.VerifyTokenResponse, error) {
	return nil, nil
}

func (f *fakeAuthClient) Register(ctx context.Context, in *authv1.RegisterRequest, opts ...grpc.CallOption) (*authv1.RegisterResponse, error) {
	return &authv1.RegisterResponse{User: &authv1.User{Id: 7, Username: in.GetUsername()}}, nil
}

func (f *fakeAuthClient) UsernameCheck(ctx context.Context, in *authv1.UsernameCheckRequest, opts ...grpc.CallOption) (*authv1.UsernameCheckResponse, error) {
	return &authv1.UsernameCheckResponse{}, nil
}

func (f *fakeAuthClient) EmailCheck(ctx context.Context, in *authv1.EmailCheckRequest, opts ...grpc.CallOption) (*authv1.EmailCheckResponse, error) {
	return &authv1.EmailCheckResponse{}, nil
}

func (f *fakeAuthClient) PhoneCheck(ctx context.Context, in *authv1.PhoneCheckRequest, opts ...grpc.CallOption) (*authv1.PhoneCheckResponse, error) {
	return &authv1.PhoneCheckResponse{}, nil
}

func (f *fakeAuthClient) Login(ctx context.Context, in *authv1.LoginRequest, opts ...grpc.CallOption) (*authv1.LoginResponse, error) {
	f.loginReq = in
	now := timestamppb.New(time.Date(2026, 6, 29, 10, 0, 0, 0, time.UTC))
	return &authv1.LoginResponse{
		TokenType:        "Bearer",
		AccessToken:      "access-token",
		ExpiresAt:        now,
		User:             &authv1.User{Id: 7, Username: "alice"},
		RefreshToken:     "refresh-token",
		RefreshExpiresAt: now,
	}, nil
}

func (f *fakeAuthClient) SendSmsCode(ctx context.Context, in *authv1.SendSmsCodeRequest, opts ...grpc.CallOption) (*authv1.SendSmsCodeResponse, error) {
	return &authv1.SendSmsCodeResponse{}, nil
}

func (f *fakeAuthClient) PhoneAuth(ctx context.Context, in *authv1.PhoneAuthRequest, opts ...grpc.CallOption) (*authv1.PhoneAuthResponse, error) {
	return &authv1.PhoneAuthResponse{}, nil
}

func (f *fakeAuthClient) MpLogin(ctx context.Context, in *authv1.MpLoginRequest, opts ...grpc.CallOption) (*authv1.MpLoginResponse, error) {
	return &authv1.MpLoginResponse{}, nil
}

func (f *fakeAuthClient) RefreshToken(ctx context.Context, in *authv1.RefreshTokenRequest, opts ...grpc.CallOption) (*authv1.RefreshTokenResponse, error) {
	return &authv1.RefreshTokenResponse{}, nil
}

func (f *fakeAuthClient) ChangePassword(ctx context.Context, in *authv1.ChangePasswordRequest, opts ...grpc.CallOption) (*authv1.ChangePasswordResponse, error) {
	return &authv1.ChangePasswordResponse{}, nil
}

func (f *fakeAuthClient) Logout(ctx context.Context, in *authv1.LogoutRequest, opts ...grpc.CallOption) (*authv1.LogoutResponse, error) {
	f.logoutReq = in
	return &authv1.LogoutResponse{}, nil
}

func TestLoginUsesOpenAPIRequestAndResponse(t *testing.T) {
	authClient := &fakeAuthClient{}
	api := NewAuthAPI(authClient)

	r := gin.New()
	r.POST("/api/v1/auth/login", api.Login)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"username":"alice","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if authClient.loginReq == nil || authClient.loginReq.GetUsername() != "alice" || authClient.loginReq.GetPassword() != "secret" {
		t.Fatalf("unexpected login request: %+v", authClient.loginReq)
	}

	var body struct {
		User struct {
			Id int64 `json:"id"`
		} `json:"user"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.User.Id != 7 || body.AccessToken != "access-token" || body.RefreshToken != "refresh-token" {
		t.Fatalf("unexpected login response: %+v", body)
	}
}

func TestLogoutAcceptsEmptyBodyAndUsesBearerToken(t *testing.T) {
	authClient := &fakeAuthClient{}
	api := NewAuthAPI(authClient)

	r := gin.New()
	r.POST("/api/v1/auth/logout", api.Logout)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer access-token")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if authClient.logoutReq == nil || authClient.logoutReq.GetAccessToken() != "access-token" {
		t.Fatalf("unexpected logout request: %+v", authClient.logoutReq)
	}
	if strings.TrimSpace(rec.Body.String()) != "{}" {
		t.Fatalf("body = %s, want {}", rec.Body.String())
	}
}
