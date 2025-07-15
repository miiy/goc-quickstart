package auth

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	authpb "github.com/miiy/goc-quickstart/auth-service/gen/go/shop/auth/v1"
	gauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/auth/jwt"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Setup custom auth.
func AuthFunc(jwtAuth *jwt.JWTAuth, userProvider gauth.UserProvider) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}
		claims, err := jwtAuth.ParseToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}
		user, err := userProvider.FirstByIdentifier(ctx, claims.Username)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token")
		}

		newCtx := gauth.InjectAuthenticatedUser(ctx, user)
		return newCtx, nil
	}
}

// auth match
func AuthMatchFunc(ctx context.Context, c interceptors.CallMeta) bool {
	// health check
	if healthpb.Health_ServiceDesc.ServiceName == c.Service {
		return false
	}

	// auth service
	var fullMethodNames []string
	for _, v := range []string{"Login", "MpLogin", "Register", "UsernameCheck", "EmailCheck", "PhoneCheck"} {
		fullMethodNames = append(fullMethodNames, fmt.Sprintf("/%s/%s", authpb.Auth_ServiceDesc.ServiceName, v))
	}

	for _, v := range fullMethodNames {
		if c.FullMethod() == v {
			return false
		}
	}

	return true
}
