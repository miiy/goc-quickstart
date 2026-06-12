package upload

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	uploadv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/upload/v1"
	userv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/user/v1"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type fakeUploadClient struct {
	req *uploadv1.CreateUploadRequest
}

func (f *fakeUploadClient) CreateUpload(ctx context.Context, in *uploadv1.CreateUploadRequest, opts ...grpc.CallOption) (*uploadv1.CreateUploadResponse, error) {
	f.req = in
	return &uploadv1.CreateUploadResponse{
		Upload: &uploadv1.Upload{Url: "http://cdn.test/uploads/avatars/2026/06/avatar.png"},
	}, nil
}

type fakeAvatarUserClient struct {
	req *userv1.UpdateUserRequest
}

func (f *fakeAvatarUserClient) GetUser(ctx context.Context, in *userv1.GetUserRequest, opts ...grpc.CallOption) (*userv1.GetUserResponse, error) {
	return nil, nil
}

func (f *fakeAvatarUserClient) BatchGetUsers(ctx context.Context, in *userv1.BatchGetUsersRequest, opts ...grpc.CallOption) (*userv1.BatchGetUsersResponse, error) {
	return nil, nil
}

func (f *fakeAvatarUserClient) UpdateUser(ctx context.Context, in *userv1.UpdateUserRequest, opts ...grpc.CallOption) (*userv1.UpdateUserResponse, error) {
	f.req = in
	return &userv1.UpdateUserResponse{
		User: &userv1.User{Id: in.GetId(), Avatar: in.GetUser().GetAvatar()},
	}, nil
}

func (f *fakeAvatarUserClient) ListUsers(ctx context.Context, in *userv1.ListUsersRequest, opts ...grpc.CallOption) (*userv1.ListUsersResponse, error) {
	return nil, nil
}

func TestAvatarUploadsAndUpdatesCurrentUser(t *testing.T) {
	uploadClient := &fakeUploadClient{}
	userClient := &fakeAvatarUserClient{}
	module := NewModule(uploadClient, userClient)

	r := gin.New()
	r.POST("/api/v1/uploads/avatar", func(c *gin.Context) {
		ctx := gocauth.InjectAuthenticatedUser(c.Request.Context(), &gocauth.AuthenticatedUser{ID: 7, Username: "alice"})
		c.Request = c.Request.WithContext(ctx)
		module.avatar(c)
	})

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("avatar", "avatar.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte{0x89, 'P', 'N', 'G'}); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/uploads/avatar", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if uploadClient.req == nil || uploadClient.req.GetScene() != uploadv1.UploadScene_UPLOAD_SCENE_AVATAR {
		t.Fatalf("unexpected upload request: %+v", uploadClient.req)
	}
	if userClient.req == nil || userClient.req.GetId() != 7 {
		t.Fatalf("unexpected user request: %+v", userClient.req)
	}
	if got := userClient.req.GetUser().GetAvatar(); got != "http://cdn.test/uploads/avatars/2026/06/avatar.png" {
		t.Fatalf("avatar = %q", got)
	}
	if !fieldMaskEqual(userClient.req.GetUpdateMask(), []string{"avatar"}) {
		t.Fatalf("update mask = %+v", userClient.req.GetUpdateMask())
	}
}

func fieldMaskEqual(mask *fieldmaskpb.FieldMask, paths []string) bool {
	if mask == nil || len(mask.Paths) != len(paths) {
		return false
	}
	for i, path := range paths {
		if mask.Paths[i] != path {
			return false
		}
	}
	return true
}
