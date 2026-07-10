package auth

import (
	"time"

	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type tokenResponse interface {
	GetTokenType() string
	GetAccessToken() string
	GetExpiresAt() *timestamppb.Timestamp
	GetUser() *authv1.User
	GetRefreshToken() string
	GetRefreshExpiresAt() *timestamppb.Timestamp
}

func protoToTokenResponse(resp tokenResponse) openapi.TokenResponse {
	return openapi.TokenResponse{
		TokenType:        resp.GetTokenType(),
		AccessToken:      resp.GetAccessToken(),
		ExpiresAt:        timestampTime(resp.GetExpiresAt()),
		User:             protoToAuthUser(resp.GetUser()),
		RefreshToken:     resp.GetRefreshToken(),
		RefreshExpiresAt: timestampTime(resp.GetRefreshExpiresAt()),
	}
}

func protoToAuthUser(user *authv1.User) openapi.AuthUser {
	if user == nil {
		return openapi.AuthUser{}
	}
	return openapi.AuthUser{
		Username: user.GetUsername(),
		Id:       user.GetId(),
	}
}

func timestampTime(value *timestamppb.Timestamp) time.Time {
	if value == nil {
		return time.Time{}
	}
	return value.AsTime()
}
