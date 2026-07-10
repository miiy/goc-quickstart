package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakeUsersUserClient struct {
	gotUsername string
	gotID       int64
	updateReq   *userv1.UpdateUserRequest
}

func (f *fakeUsersUserClient) GetUser(ctx context.Context, in *userv1.GetUserRequest, opts ...grpc.CallOption) (*userv1.GetUserResponse, error) {
	f.gotID = in.GetId()
	now := timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:        in.GetId(),
			Username:  "alice",
			Nickname:  "Alice",
			Avatar:    "avatars/alice.png",
			Email:     "alice@example.com",
			Phone:     "123",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

func (f *fakeUsersUserClient) GetUserByUsername(ctx context.Context, in *userv1.GetUserByUsernameRequest, opts ...grpc.CallOption) (*userv1.GetUserByUsernameResponse, error) {
	f.gotUsername = in.GetUsername()
	now := timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	return &userv1.GetUserByUsernameResponse{
		User: &userv1.User{
			Id:        7,
			Username:  in.GetUsername(),
			Nickname:  "Alice",
			Avatar:    "avatars/alice.png",
			Email:     "alice@example.com",
			Phone:     "123",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

func (f *fakeUsersUserClient) BatchGetUsers(ctx context.Context, in *userv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*userv1.BatchGetUsersResponse, error) {
	return nil, nil
}

func (f *fakeUsersUserClient) UpdateUser(ctx context.Context, in *userv1.UpdateUserRequest, opts ...grpc.CallOption) (*userv1.UpdateUserResponse, error) {
	f.updateReq = in
	now := timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	return &userv1.UpdateUserResponse{
		User: &userv1.User{
			Id:        in.GetId(),
			Username:  "alice",
			Nickname:  in.GetUser().GetNickname(),
			Email:     "alice@example.com",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

func (f *fakeUsersUserClient) ListUsers(ctx context.Context, in *userv1.ListUsersRequest, opts ...grpc.CallOption) (*userv1.ListUsersResponse, error) {
	return nil, nil
}

func TestGetUserUsesUsernameAndHidesPrivateFields(t *testing.T) {
	userClient := &fakeUsersUserClient{}
	api := &UsersAPI{userClient: userClient}

	r := gin.New()
	r.GET("/api/v1/users/:username", api.GetUser)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/alice", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if userClient.gotUsername != "alice" {
		t.Fatalf("username = %q, want alice", userClient.gotUsername)
	}

	var body map[string]map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	user := body["user"]
	if user["username"] != "alice" || user["nickname"] != "Alice" {
		t.Fatalf("unexpected public user: %+v", user)
	}
	if user["avatar"] != "/uploads/avatars/alice.png" {
		t.Fatalf("avatar = %q, want /uploads/avatars/alice.png", user["avatar"])
	}
	if _, ok := user["email"]; ok {
		t.Fatalf("public user leaked email: %+v", user)
	}
	if _, ok := user["phone"]; ok {
		t.Fatalf("public user leaked phone: %+v", user)
	}
}

func TestGetProfileUsesAuthenticatedUserID(t *testing.T) {
	userClient := &fakeUsersUserClient{}
	api := &UsersAPI{userClient: userClient}

	r := gin.New()
	r.GET("/api/v1/profile", func(c *gin.Context) {
		ctx := gocauth.InjectAuthenticatedUser(c.Request.Context(), &gocauth.AuthenticatedUser{ID: "7", Username: "alice"})
		c.Request = c.Request.WithContext(ctx)
		api.GetProfile(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if userClient.gotID != 7 {
		t.Fatalf("id = %d, want 7", userClient.gotID)
	}

	var body map[string]map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	user := body["user"]
	if user["email"] != "alice@example.com" || user["phone"] != "123" {
		t.Fatalf("unexpected profile user: %+v", user)
	}
	if user["avatar"] != "/uploads/avatars/alice.png" {
		t.Fatalf("avatar = %q, want /uploads/avatars/alice.png", user["avatar"])
	}
}

func TestUpdateProfileUsesAuthenticatedUserIDAndUpdateMask(t *testing.T) {
	userClient := &fakeUsersUserClient{}
	api := &UsersAPI{userClient: userClient}

	r := gin.New()
	r.PUT("/api/v1/profile", func(c *gin.Context) {
		ctx := gocauth.InjectAuthenticatedUser(c.Request.Context(), &gocauth.AuthenticatedUser{ID: "7", Username: "alice"})
		c.Request = c.Request.WithContext(ctx)
		api.UpdateProfile(c)
	})

	body := `{"user":{"nickname":"Alice New"},"update_fields":["nickname"]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if userClient.updateReq == nil {
		t.Fatal("UpdateUser was not called")
	}
	if userClient.updateReq.GetId() != 7 {
		t.Fatalf("id = %d, want 7", userClient.updateReq.GetId())
	}
	if got := userClient.updateReq.GetUser().GetNickname(); got != "Alice New" {
		t.Fatalf("nickname = %q, want Alice New", got)
	}
	if got := userClient.updateReq.GetUser().GetEmail(); got != "" {
		t.Fatalf("email = %q, want empty input guarded by update mask", got)
	}
	paths := userClient.updateReq.GetUpdateMask().GetPaths()
	if len(paths) != 1 || paths[0] != "nickname" {
		t.Fatalf("update mask paths = %v, want [nickname]", paths)
	}
}
