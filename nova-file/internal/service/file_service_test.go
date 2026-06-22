package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pb "github.com/miiy/goc-quickstart/nova-file/gen/go/nova/file/v1"
	"github.com/miiy/goc-quickstart/nova-file/internal/config"
	"github.com/miiy/goc-quickstart/nova-file/internal/entity"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/logger/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockFileRepository struct {
	file *entity.File
	err  error
}

func (m *mockFileRepository) Create(ctx context.Context, file *entity.File) error {
	if m.err != nil {
		return m.err
	}
	file.ID = 1
	m.file = file
	return nil
}

func fileTestContext(userID int64) context.Context {
	return gocauth.InjectAuthenticatedUser(context.Background(), &gocauth.AuthenticatedUser{
		ID:       userID,
		Username: "alice",
	})
}

func TestUploadFileAvatarWritesFileAndRecord(t *testing.T) {
	dir := t.TempDir()
	repo := &mockFileRepository{}
	service := NewFileServiceServer(zap.NewNop(), repo, config.Storage{
		Root:          dir,
		PublicURL:     "http://cdn.test/files",
		MaxAvatarSize: 1024,
	}).(*FileService)

	content := validPNGBytes()
	resp, err := service.UploadFile(fileTestContext(7), &pb.UploadFileRequest{
		Scene:    pb.FileScene_FILE_SCENE_AVATAR,
		Filename: "avatar.png",
		MimeType: "image/png",
		Content:  content,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetFile().GetOwnerId() != 7 {
		t.Fatalf("owner id = %d, want 7", resp.GetFile().GetOwnerId())
	}
	if strings.Contains(resp.GetFile().GetUrl(), "/7/") {
		t.Fatalf("avatar url exposes user id: %s", resp.GetFile().GetUrl())
	}
	if !strings.HasPrefix(resp.GetFile().GetObjectKey(), "avatars/") {
		t.Fatalf("object key = %q, want avatars prefix", resp.GetFile().GetObjectKey())
	}
	if wantURL := "http://cdn.test/files/" + resp.GetFile().GetObjectKey(); resp.GetFile().GetUrl() != wantURL {
		t.Fatalf("url = %q, want %q", resp.GetFile().GetUrl(), wantURL)
	}
	if filepath.Ext(resp.GetFile().GetObjectKey()) != ".png" {
		t.Fatalf("object ext = %q, want .png", filepath.Ext(resp.GetFile().GetObjectKey()))
	}
	if repo.file == nil || repo.file.Checksum == "" {
		t.Fatalf("file record not saved: %+v", repo.file)
	}

	savedPath := filepath.Join(dir, filepath.FromSlash(resp.GetFile().GetObjectKey()))
	saved, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(saved) != string(content) {
		t.Fatalf("saved content = %q, want %q", string(saved), string(content))
	}
}

func TestUploadFilePostCoverWritesFileAndRecord(t *testing.T) {
	dir := t.TempDir()
	repo := &mockFileRepository{}
	service := NewFileServiceServer(zap.NewNop(), repo, config.Storage{
		Root:             dir,
		MaxAvatarSize:    1024,
		MaxPostCoverSize: 2048,
	}).(*FileService)

	content := validPNGBytes()
	resp, err := service.UploadFile(fileTestContext(7), &pb.UploadFileRequest{
		Scene:    pb.FileScene_FILE_SCENE_POST_COVER,
		Filename: "cover.png",
		MimeType: "image/png",
		Content:  content,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetFile().GetScene() != pb.FileScene_FILE_SCENE_POST_COVER {
		t.Fatalf("scene = %v, want post cover", resp.GetFile().GetScene())
	}
	if !strings.HasPrefix(resp.GetFile().GetObjectKey(), "post-covers/") {
		t.Fatalf("object key = %q, want post-covers prefix", resp.GetFile().GetObjectKey())
	}
	if repo.file == nil || repo.file.Scene != entity.FileScenePostCover {
		t.Fatalf("file record not saved with post cover scene: %+v", repo.file)
	}

	savedPath := filepath.Join(dir, filepath.FromSlash(resp.GetFile().GetObjectKey()))
	if _, err := os.Stat(savedPath); err != nil {
		t.Fatal(err)
	}
}

func TestUploadFileRejectsUnsupportedMime(t *testing.T) {
	service := NewFileServiceServer(zap.NewNop(), &mockFileRepository{}, config.Storage{
		Root:          t.TempDir(),
		MaxAvatarSize: 1024,
	}).(*FileService)

	_, err := service.UploadFile(fileTestContext(7), &pb.UploadFileRequest{
		Scene:    pb.FileScene_FILE_SCENE_AVATAR,
		MimeType: "text/plain",
		Content:  []byte("plain text"),
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("code = %v, want InvalidArgument, err=%v", status.Code(err), err)
	}
}

func TestUploadFileRejectsOversizedPostCover(t *testing.T) {
	service := NewFileServiceServer(zap.NewNop(), &mockFileRepository{}, config.Storage{
		Root:             t.TempDir(),
		MaxAvatarSize:    1024,
		MaxPostCoverSize: 1,
	}).(*FileService)

	_, err := service.UploadFile(fileTestContext(7), &pb.UploadFileRequest{
		Scene:    pb.FileScene_FILE_SCENE_POST_COVER,
		MimeType: "image/png",
		Content:  validPNGBytes(),
	})
	if status.Code(err) != codes.ResourceExhausted {
		t.Fatalf("code = %v, want ResourceExhausted, err=%v", status.Code(err), err)
	}
}

func TestUploadFileRejectsOversizedAvatar(t *testing.T) {
	service := NewFileServiceServer(zap.NewNop(), &mockFileRepository{}, config.Storage{
		Root:          t.TempDir(),
		MaxAvatarSize: 1,
	}).(*FileService)

	_, err := service.UploadFile(fileTestContext(7), &pb.UploadFileRequest{
		Scene:    pb.FileScene_FILE_SCENE_AVATAR,
		MimeType: "image/png",
		Content:  validPNGBytes(),
	})
	if status.Code(err) != codes.ResourceExhausted {
		t.Fatalf("code = %v, want ResourceExhausted, err=%v", status.Code(err), err)
	}
}

func TestUploadFileRequiresAuthenticatedUser(t *testing.T) {
	service := NewFileServiceServer(zap.NewNop(), &mockFileRepository{}, config.Storage{
		Root:          t.TempDir(),
		MaxAvatarSize: 1024,
	}).(*FileService)

	_, err := service.UploadFile(context.Background(), &pb.UploadFileRequest{
		Scene:    pb.FileScene_FILE_SCENE_AVATAR,
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
