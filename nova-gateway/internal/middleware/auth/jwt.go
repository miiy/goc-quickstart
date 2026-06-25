package auth

import (
	"context"
	"errors"
	"strconv"

	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/jwtauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"
)

// verifyTokenClient is the subset of nova-auth AuthServiceClient used by the middleware.
type verifyTokenClient interface {
	VerifyToken(ctx context.Context, in *authv1.VerifyTokenRequest, opts ...grpc.CallOption) (*authv1.VerifyTokenResponse, error)
}

// AuthenticationMiddleware authenticates each request via nova-auth's VerifyToken
// (signature + revocation + active user in one RPC). The generic flow (token
// extraction, user/metadata injection, 401 handling) is provided by goc's
// jwtauth.Authenticate; only the token-to-user resolution is nova-specific.
func AuthenticationMiddleware(authClient verifyTokenClient) gin.HandlerFunc {
	return jwtauth.Authenticate(tokenUserResolver(authClient),
		jwtauth.WithUnauthorized(transport.WriteUnauthorized),
		jwtauth.WithMetadataPropagation())
}

// tokenUserResolver resolves a bearer token to an authenticated user by calling
// nova-auth's VerifyToken.
func tokenUserResolver(authClient verifyTokenClient) jwtauth.UserResolver {
	return func(ctx context.Context, token string) (*gocauth.AuthenticatedUser, error) {
		if authClient == nil {
			return nil, errors.New("auth client not configured")
		}
		resp, err := authClient.VerifyToken(ctx, &authv1.VerifyTokenRequest{AccessToken: token})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				return nil, errors.New(st.Message())
			}
			return nil, err
		}
		if resp == nil || resp.GetUser() == nil {
			return nil, errors.New("authenticated user not found")
		}
		return &gocauth.AuthenticatedUser{
			ID:       strconv.FormatInt(resp.User.Id, 10),
			Username: resp.User.Username,
		}, nil
	}
}
