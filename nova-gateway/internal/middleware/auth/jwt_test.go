package auth

import (
	"context"
	"testing"

	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	gocauth "github.com/miiy/goc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeAuthenticatedUserClient struct {
	user        *authv1.User
	err         error
	gotUsername string
}

func (f *fakeAuthenticatedUserClient) GetAuthenticatedUser(ctx context.Context, req *authv1.GetAuthenticatedUserRequest, opts ...grpc.CallOption) (*authv1.GetAuthenticatedUserResponse, error) {
	f.gotUsername = req.Username
	if f.err != nil || f.user == nil {
		return nil, f.err
	}
	return &authv1.GetAuthenticatedUserResponse{User: f.user}, nil
}

func TestAuthUserResolverReturnsAuthenticatedUser(t *testing.T) {
	client := &fakeAuthenticatedUserClient{
		user: &authv1.User{
			Id:       42,
			Username: "alice",
		},
	}

	user, err := authUserResolver(client)(context.Background(), &gocauth.UserClaims{Username: "alice"}, "token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client.gotUsername != "alice" {
		t.Fatalf("expected username alice, got %q", client.gotUsername)
	}
	if user.ID != 42 || user.Username != "alice" {
		t.Fatalf("unexpected authenticated user: %+v", user)
	}
}

func TestAuthUserResolverReturnsStatusMessage(t *testing.T) {
	client := &fakeAuthenticatedUserClient{
		err: status.Error(codes.Unauthenticated, "user disabled"),
	}

	_, err := authUserResolver(client)(context.Background(), &gocauth.UserClaims{Username: "alice"}, "token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "user disabled" {
		t.Fatalf("expected user disabled, got %q", err.Error())
	}
}
