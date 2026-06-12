package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/miiy/goc-quickstart/upload-service/gen/go/blog/upload/v1"
	"github.com/miiy/goc-quickstart/upload-service/internal/config"
	"github.com/miiy/goc-quickstart/upload-service/internal/entity"
	"github.com/miiy/goc-quickstart/upload-service/internal/repository"
	gocauth "github.com/miiy/goc/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultStorageRoot     = "./storage/uploads"
	defaultPublicURL       = "/uploads"
	defaultMaxAvatarSize   = 2 << 20
	avatarObjectKeyPattern = "avatars/%04d/%02d/%s.%s"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrUnsupportedMime = errors.New("unsupported upload mime type")
	ErrUploadTooLarge  = errors.New("upload is too large")
)

type UploadService struct {
	pb.UnimplementedUploadServiceServer
	logger  *zap.Logger
	repo    repository.UploadRepository
	storage config.Storage
}

func NewUploadServiceServer(logger *zap.Logger, repo repository.UploadRepository, storage config.Storage) pb.UploadServiceServer {
	return &UploadService{
		logger:  logger,
		repo:    repo,
		storage: normalizeStorage(storage),
	}
}

func (s *UploadService) CreateUpload(ctx context.Context, req *pb.CreateUploadRequest) (*pb.CreateUploadResponse, error) {
	user, err := authenticatedUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.GetScene() != pb.UploadScene_UPLOAD_SCENE_AVATAR {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}
	if len(req.GetContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}
	if int64(len(req.GetContent())) > s.storage.MaxAvatarSize {
		return nil, status.Error(codes.ResourceExhausted, ErrUploadTooLarge.Error())
	}

	mimeType, ext, err := normalizeAvatarMime(req.GetMimeType(), req.GetContent())
	if err != nil {
		return nil, err
	}

	now := time.Now()
	randomName, err := randomHex(16)
	if err != nil {
		s.logger.Error("randomHex", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	objectKey := fmt.Sprintf(avatarObjectKeyPattern, now.Year(), int(now.Month()), randomName, ext)
	filePath, err := safeStoragePath(s.storage.Root, objectKey)
	if err != nil {
		s.logger.Error("safeStoragePath", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		s.logger.Error("os.MkdirAll", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := os.WriteFile(filePath, req.GetContent(), 0o644); err != nil {
		s.logger.Error("os.WriteFile", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	checksum := sha256.Sum256(req.GetContent())
	upload := &entity.Upload{
		OwnerID:   user.ID,
		OwnerType: entity.OwnerTypeUser,
		Scene:     entity.UploadSceneAvatar,
		ObjectKey: objectKey,
		URL:       publicURL(s.storage.PublicURL, objectKey),
		MimeType:  mimeType,
		Size:      int64(len(req.GetContent())),
		Checksum:  hex.EncodeToString(checksum[:]),
		Status:    entity.UploadStatusActive,
		CreatedBy: user.ID,
	}
	if err := s.repo.Create(ctx, upload); err != nil {
		_ = os.Remove(filePath)
		s.logger.Error("repo.Create", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateUploadResponse{Upload: entityToProto(upload)}, nil
}

func authenticatedUser(ctx context.Context) (*gocauth.AuthenticatedUser, error) {
	user, err := gocauth.ExtractAuthenticatedUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if user.ID <= 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid authenticated user")
	}
	return user, nil
}

func normalizeStorage(storage config.Storage) config.Storage {
	storage.Root = strings.TrimSpace(storage.Root)
	if storage.Root == "" {
		storage.Root = defaultStorageRoot
	}
	storage.PublicURL = strings.TrimSpace(storage.PublicURL)
	if storage.PublicURL == "" {
		storage.PublicURL = defaultPublicURL
	}
	if storage.MaxAvatarSize <= 0 {
		storage.MaxAvatarSize = defaultMaxAvatarSize
	}
	return storage
}

func normalizeAvatarMime(raw string, content []byte) (string, string, error) {
	mimeType := detectedMime(content)
	declared := strings.ToLower(strings.TrimSpace(strings.Split(raw, ";")[0]))
	if declared != "" && declared != "application/octet-stream" && declared != mimeType {
		return "", "", status.Error(codes.InvalidArgument, ErrUnsupportedMime.Error())
	}

	switch mimeType {
	case "image/jpeg":
		return mimeType, "jpg", nil
	case "image/png":
		return mimeType, "png", nil
	case "image/webp":
		return mimeType, "webp", nil
	default:
		return "", "", status.Error(codes.InvalidArgument, ErrUnsupportedMime.Error())
	}
}

func detectedMime(content []byte) string {
	mimeType := http.DetectContentType(content)
	return strings.ToLower(strings.TrimSpace(strings.Split(mimeType, ";")[0]))
}

func randomHex(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func safeStoragePath(root, objectKey string) (string, error) {
	root = filepath.Clean(root)
	rel := filepath.Clean(filepath.FromSlash(objectKey))
	if rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." || filepath.IsAbs(rel) {
		return "", fmt.Errorf("invalid object key: %s", objectKey)
	}

	full := filepath.Join(root, rel)
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	if absFull != absRoot && !strings.HasPrefix(absFull, absRoot+string(filepath.Separator)) {
		return "", fmt.Errorf("object key escapes storage root: %s", objectKey)
	}
	return full, nil
}

func publicURL(baseURL, objectKey string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(objectKey, "/")
}

func entityToProto(upload *entity.Upload) *pb.Upload {
	resp := &pb.Upload{
		Id:        upload.ID,
		OwnerId:   upload.OwnerID,
		OwnerType: upload.OwnerType,
		Scene:     pb.UploadScene(upload.Scene),
		ObjectKey: upload.ObjectKey,
		Url:       upload.URL,
		MimeType:  upload.MimeType,
		Size:      upload.Size,
		Checksum:  upload.Checksum,
		Status:    pb.UploadStatus(upload.Status),
		CreatedBy: upload.CreatedBy,
		CreatedAt: timestamppb.New(upload.CreatedAt),
		UpdatedAt: timestamppb.New(upload.UpdatedAt),
	}
	if upload.DeletedAt.Valid {
		resp.DeletedAt = timestamppb.New(upload.DeletedAt.Time)
	}
	return resp
}
