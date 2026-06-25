package file

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	filev1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/file/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const maxAvatarUploadSize = 2 << 20
const maxPostCoverUploadSize = 5 << 20

func (m *Module) avatar(c *gin.Context) {
	if m.fileClient == nil || m.userClient == nil {
		transport.WriteError(c, status.Error(codes.Unavailable, "file service not configured"))
		return
	}

	userID, ok := authctx.CurrentUserInt64ID(c)
	if !ok || userID <= 0 {
		transport.WriteError(c, status.Error(codes.Unauthenticated, "invalid authenticated user"))
		return
	}

	content, header, err := readMultipartFile(c, "avatar", maxAvatarUploadSize)
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(content)
	}

	fileResp, err := m.fileClient.UploadFile(c.Request.Context(), &filev1.UploadFileRequest{
		Scene:    filev1.FileScene_FILE_SCENE_AVATAR,
		Filename: header.Filename,
		MimeType: mimeType,
		Content:  content,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	if fileResp.GetFile() == nil || fileResp.GetFile().GetObjectKey() == "" {
		transport.WriteError(c, status.Error(codes.Internal, "empty file object key"))
		return
	}

	resp, err := m.userClient.UpdateUser(c.Request.Context(), &userv1.UpdateUserRequest{
		Id: userID,
		User: &userv1.User{
			Avatar: fileResp.GetFile().GetObjectKey(),
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"avatar"}},
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	transport.WriteProto(c, resp)
}

func (m *Module) upload(c *gin.Context) {
	if m.fileClient == nil {
		transport.WriteError(c, status.Error(codes.Unavailable, "file service not configured"))
		return
	}

	scene, err := uploadSceneFromForm(c.PostForm("scene"))
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	content, header, err := readMultipartFile(c, "file", maxUploadSizeForScene(scene))
	if err != nil {
		transport.WriteError(c, err)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(content)
	}

	resp, err := m.fileClient.UploadFile(c.Request.Context(), &filev1.UploadFileRequest{
		Scene:    scene,
		Filename: header.Filename,
		MimeType: mimeType,
		Content:  content,
	})
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func uploadSceneFromForm(value string) (filev1.FileScene, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "post_cover", "post-cover", "cover", "file_scene_post_cover":
		return filev1.FileScene_FILE_SCENE_POST_COVER, nil
	default:
		return filev1.FileScene_FILE_SCENE_UNSPECIFIED, status.Error(codes.InvalidArgument, "unsupported file scene")
	}
}

func maxUploadSizeForScene(scene filev1.FileScene) int64 {
	switch scene {
	case filev1.FileScene_FILE_SCENE_POST_COVER:
		return maxPostCoverUploadSize
	default:
		return maxPostCoverUploadSize
	}
}

func readMultipartFile(c *gin.Context, field string, maxSize int64) ([]byte, *multipart.FileHeader, error) {
	file, header, err := c.Request.FormFile(field)
	if err != nil {
		return nil, nil, status.Errorf(codes.InvalidArgument, "%s file is required", field)
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, maxSize+1))
	if err != nil {
		return nil, nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if len(content) == 0 {
		return nil, nil, status.Errorf(codes.InvalidArgument, "%s file is required", field)
	}
	if int64(len(content)) > maxSize {
		return nil, nil, status.Errorf(codes.ResourceExhausted, "%s file is too large", field)
	}
	return content, header, nil
}
