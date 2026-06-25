package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type fakeVerifyTokenClient struct {
	user     *authv1.User
	err      error
	gotToken string
}

func (f *fakeVerifyTokenClient) VerifyToken(ctx context.Context, req *authv1.VerifyTokenRequest, opts ...grpc.CallOption) (*authv1.VerifyTokenResponse, error) {
	f.gotToken = req.AccessToken
	if f.err != nil || f.user == nil {
		return nil, f.err
	}
	return &authv1.VerifyTokenResponse{User: f.user}, nil
}

func newProtectedEngine(client verifyTokenClient) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthenticationMiddleware(client))
	r.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}

func TestAuthenticationMiddlewareSuccess(t *testing.T) {
	client := &fakeVerifyTokenClient{user: &authv1.User{Id: 42, Username: "alice"}}
	engine := newProtectedEngine(client)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if client.gotToken != "some-token" {
		t.Fatalf("expected token forwarded as some-token, got %q", client.gotToken)
	}
}

func TestAuthenticationMiddlewareInjectsUserAndMetadata(t *testing.T) {
	client := &fakeVerifyTokenClient{user: &authv1.User{Id: 42, Username: "alice"}}

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(AuthenticationMiddleware(client))

	var gotUser *gocauth.AuthenticatedUser
	var gotMD metadata.MD
	engine.GET("/protected", func(c *gin.Context) {
		gotUser, _ = gocauth.ExtractAuthenticatedUser(c.Request.Context())
		gotMD, _ = metadata.FromOutgoingContext(c.Request.Context())
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer tok")
	engine.ServeHTTP(w, req)

	if gotUser == nil || gotUser.ID != "42" || gotUser.Username != "alice" {
		t.Fatalf("expected injected user {42 alice}, got %+v", gotUser)
	}
	if v := gotMD.Get(gocauth.AuthenticatedUserIDMetadataKey); len(v) == 0 || v[0] != "42" {
		t.Fatalf("expected metadata x-auth-user-id=42, got %v", v)
	}
	if v := gotMD.Get(gocauth.AuthenticatedUsernameMetadataKey); len(v) == 0 || v[0] != "alice" {
		t.Fatalf("expected metadata x-auth-username=alice, got %v", v)
	}
}

func TestAuthenticationMiddlewareRejectsMissingToken(t *testing.T) {
	client := &fakeVerifyTokenClient{user: &authv1.User{Id: 1, Username: "x"}}
	engine := newProtectedEngine(client)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthenticationMiddlewareRejectsVerifyError(t *testing.T) {
	client := &fakeVerifyTokenClient{err: status.Error(codes.Unauthenticated, "token revoked")}
	engine := newProtectedEngine(client)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer tok")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
