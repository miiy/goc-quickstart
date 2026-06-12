package auth

import (
	"context"
	"errors"
	"strconv"

	authv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/auth/v1"
	"github.com/miiy/goc-quickstart/api-gateway/internal/transport"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	ginauth "github.com/miiy/goc/gin/middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type authenticatedUserClient interface {
	GetAuthenticatedUser(ctx context.Context, in *authv1.GetAuthenticatedUserRequest, opts ...grpc.CallOption) (*authv1.GetAuthenticatedUserResponse, error)
}

// JWTAuthenticationMiddleware assembles the JWT authentication middleware with api-gateway-specific behavior.
func JWTAuthenticationMiddleware(jwtAuth *gocauth.JWTAuth, authClient authenticatedUserClient) gin.HandlerFunc {
	return ginauth.JWTAuthenticationMiddleware(jwtAuth,
		ginauth.WithUnauthorized(transport.WriteUnauthorized),
		ginauth.WithUserResolver(authUserResolver(authClient)),
		ginauth.WithAfterAuth(func(c *gin.Context, claims *gocauth.UserClaims, _ string) {
			authUser := &gocauth.AuthenticatedUser{Username: claims.Username}
			if user, ok := ginauth.GetAuthUser(c); ok {
				if user.Username != "" {
					authUser.Username = user.Username
				}
				authUser.ID = user.ID
			}

			ctx := metadata.AppendToOutgoingContext(
				c.Request.Context(),
				gocauth.AuthenticatedUserIDMetadataKey, strconv.FormatInt(authUser.ID, 10),
				gocauth.AuthenticatedUsernameMetadataKey, authUser.Username,
			)
			c.Request = c.Request.WithContext(ctx)
		}),
	)
}

func authUserResolver(authClient authenticatedUserClient) ginauth.UserResolver {
	return func(ctx context.Context, claims *gocauth.UserClaims, _ string) (*gocauth.AuthenticatedUser, error) {
		if authClient == nil {
			return nil, errors.New("auth client not configured")
		}

		resp, err := authClient.GetAuthenticatedUser(ctx, &authv1.GetAuthenticatedUserRequest{
			Username: claims.Username,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				return nil, errors.New(st.Message())
			}
			return nil, err
		}
		if resp == nil || resp.GetUser() == nil {
			return nil, errors.New("authenticated user not found")
		}
		user := resp.GetUser()

		return &gocauth.AuthenticatedUser{
			ID:       user.Id,
			Username: user.Username,
		}, nil
	}
}
