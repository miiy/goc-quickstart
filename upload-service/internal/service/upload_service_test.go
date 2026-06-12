package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pb "github.com/miiy/goc-quickstart/upload-service/gen/go/blog/upload/v1"
	"github.com/miiy/goc-quickstart/upload-service/internal/config"
	"github.com/miiy/goc-quickstart/upload-service/internal/entity"
	gocauth "github.com/miiy/goc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockUploadRepository struct {
	upload *entity.Upload
	err    error
}

func (m *mockUploadRepository) Create(ctx context.Context, upload *entity.Upload) error {
	if m.err != nil {
		return m.err
	}
	upload.ID = 1
	m.upload = upload
	return nil
}

func uploadTestContext(userID int64) context.Context {
	return gocauth.InjectAuthenticatedUser(context.Background(), &gocauth.AuthenticatedUser{
		ID:       userID,
		Username: "alice",
	})
}

func TestCreateUploadAvatarWritesFileAndRecord(t *testing.T) {
	dir := t.TempDir()
	repo := &mockUploadRepository{}
	service := NewUploadServiceServer(zap.NewNop(), repo, config.Storage{
		Root:          dir,
		PublicURL:     "http://cdn.test/uploads",
		MaxAvatarSize: 1024,
	}).(*UploadService)

	content := validPNGBytes()
	resp, err := service.CreateUpload(uploadTestContext(7), &pb.CreateUploadRequest{
		Scene:    pb.UploadScene_UPLOAD_SCENE_AVATAR,
		Filename: "avatar.png",
		MimeType: "image/png",
		Content:  content,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetUpload().GetOwnerId() != 7 {
		t.Fatalf("owner id = %d, want 7", resp.GetUpload().GetOwnerId())
	}
	if strings.Contains(resp.GetUpload().GetUrl(), "/7/") {
		t.Fatalf("avatar url exposes user id: %s", resp.GetUpload().GetUrl())
	}
	if !strings.HasPrefix(resp.GetUpload().GetObjectKey(), "avatars/") {
		t.Fatalf("object key = %q, want avatars prefix", resp.GetUpload().GetObjectKey())
	}
	if filepath.Ext(resp.GetUpload().GetObjectKey()) != ".png" {
		t.Fatalf("object ext = %q, want .png", filepath.Ext(resp.GetUpload().GetObjectKey()))
	}
	if repo.upload == nil || repo.upload.Checksum == "" {
		t.Fatalf("upload record not saved: %+v", repo.upload)
	}

	savedPath := filepath.Join(dir, filepath.FromSlash(resp.GetUpload().GetObjectKey()))
	saved, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(saved) != string(content) {
		t.Fatalf("saved content = %q, want %q", string(saved), string(content))
	}
}

func TestCreateUploadRejectsUnsupportedMime(t *testing.T) {
	service := NewUploadServiceServer(zap.NewNop(), &mockUploadRepository{}, config.Storage{
		Root:          t.TempDir(),
		MaxAvatarSize: 1024,
	}).(*UploadService)

	_, err := service.CreateUpload(uploadTestContext(7), &pb.CreateUploadRequest{
		Scene:    pb.UploadScene_UPLOAD_SCENE_AVATAR,
		MimeType: "text/plain",
		Content:  []byte("plain text"),
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("code = %v, want InvalidArgument, err=%v", status.Code(err), err)
	}
}

func TestCreateUploadRejectsOversizedAvatar(t *testing.T) {
	service := NewUploadServiceServer(zap.NewNop(), &mockUploadRepository{}, config.Storage{
		Root:          t.TempDir(),
		MaxAvatarSize: 1,
	}).(*UploadService)

	_, err := service.CreateUpload(uploadTestContext(7), &pb.CreateUploadRequest{
		Scene:    pb.UploadScene_UPLOAD_SCENE_AVATAR,
		MimeType: "image/png",
		Content:  validPNGBytes(),
	})
	if status.Code(err) != codes.ResourceExhausted {
		t.Fatalf("code = %v, want ResourceExhausted, err=%v", status.Code(err), err)
	}
}

func TestCreateUploadRequiresAuthenticatedUser(t *testing.T) {
	service := NewUploadServiceServer(zap.NewNop(), &mockUploadRepository{}, config.Storage{
		Root:          t.TempDir(),
		MaxAvatarSize: 1024,
	}).(*UploadService)

	_, err := service.CreateUpload(context.Background(), &pb.CreateUploadRequest{
		Scene:    pb.UploadScene_UPLOAD_SCENE_AVATAR,
		MimeType: "image/png",
		Content:  validPNGBytes(),
	})
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("code = %v, want Unauthenticated, err=%v", status.Code(err), err)
	}
}

func validPNGBytes() []byte {
	return []byte{
		0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n',
		0x00, 0x00, 0x00, 0x0d, 'I', 'H', 'D', 'R',
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00,
	}
}
