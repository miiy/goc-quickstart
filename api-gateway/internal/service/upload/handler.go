package upload

import (
	"io"
	"net/http"

	uploadv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/upload/v1"
	userv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/user/v1"
	"github.com/miiy/goc-quickstart/api-gateway/internal/transport"
	"github.com/miiy/goc/gin"
	ginauth "github.com/miiy/goc/gin/middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const maxAvatarUploadSize = 2 << 20

func (m *Module) avatar(c *gin.Context) {
	if m.client == nil || m.userClient == nil {
		transport.WriteError(c, status.Error(codes.Unavailable, "upload service not configured"))
		return
	}

	userID, ok := ginauth.GetAuthUserID(c)
	if !ok || userID <= 0 {
		transport.WriteError(c, status.Error(codes.Unauthenticated, "invalid authenticated user"))
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		transport.WriteError(c, status.Error(codes.InvalidArgument, "avatar file is required"))
		return
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, maxAvatarUploadSize+1))
	if err != nil {
		transport.WriteError(c, status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	if len(content) == 0 {
		transport.WriteError(c, status.Error(codes.InvalidArgument, "avatar file is required"))
		return
	}
	if len(content) > maxAvatarUploadSize {
		transport.WriteError(c, status.Error(codes.ResourceExhausted, "avatar file is too large"))
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(content)
	}

	uploadResp, err := m.client.CreateUpload(c.Request.Context(), &uploadv1.CreateUploadRequest{
		Scene:    uploadv1.UploadScene_UPLOAD_SCENE_AVATAR,
		Filename: header.Filename,
		MimeType: mimeType,
		Content:  content,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	if uploadResp.GetUpload() == nil || uploadResp.GetUpload().GetUrl() == "" {
		transport.WriteError(c, status.Error(codes.Internal, "empty upload url"))
		return
	}

	resp, err := m.userClient.UpdateUser(c.Request.Context(), &userv1.UpdateUserRequest{
		Id: userID,
		User: &userv1.User{
			Avatar: uploadResp.GetUpload().GetUrl(),
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"avatar"}},
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	transport.WriteProto(c, resp)
}
