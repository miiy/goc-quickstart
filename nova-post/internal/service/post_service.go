package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	"buf.build/go/protovalidate"
	pb "github.com/miiy/goc-quickstart/nova-post/gen/go/nova/post/v1"
	"github.com/miiy/goc-quickstart/nova-post/internal/entity"
	"github.com/miiy/goc-quickstart/nova-post/internal/repository"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/logger/zap"
	"github.com/miiy/goc/pagination"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostService struct {
	pb.UnimplementedPostServiceServer
	logger       *zap.Logger
	repo         repository.PostRepository
	categoryRepo repository.CategoryRepository
}

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrPostNotFound     = errors.New("post not found")
	ErrPermissionDenied = errors.New("permission denied")
)

func NewPostServiceServer(logger *zap.Logger, repo repository.PostRepository, categoryRepo repository.CategoryRepository) pb.PostServiceServer {
	return &PostService{
		logger:       logger,
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *PostService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	post, err := s.repo.First(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ErrPostNotFound.Error())
		}
		s.logger.Error("repo.First", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetPostResponse{
		Post: entityToProto(post),
	}, nil
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if strings.TrimSpace(req.Post.Title) == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	userID, err := authenticatedUserID(ctx)
	if err != nil {
		return nil, err
	}

	tags, _ := json.Marshal(req.Post.Tags)

	post := &entity.Post{
		UserId:      userID,
		Title:       strings.TrimSpace(req.Post.Title),
		Summary:     strings.TrimSpace(req.Post.Summary),
		CoverUrl:    normalizeCoverObjectKey(req.Post.CoverUrl),
		Content:     req.Post.Content,
		Status:      int64(req.Post.Status),
		Tags:        string(tags),
		CategoryId:  req.Post.CategoryId,
		PublishedAt: protoTimestampTime(req.Post.PublishedAt),
	}
	defaultPublishedAt(post)

	if err := s.repo.Create(ctx, post); err != nil {
		s.logger.Error("repo.Create", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreatePostResponse{
		Post: entityToProto(post),
	}, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, err := authenticatedUserID(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.First(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ErrPostNotFound.Error())
		}
		s.logger.Error("repo.First", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if existing.UserId != userID {
		return nil, status.Error(codes.PermissionDenied, ErrPermissionDenied.Error())
	}

	tags, _ := json.Marshal(req.Post.Tags)

	post := &entity.Post{
		UserId:      existing.UserId,
		Title:       strings.TrimSpace(req.Post.Title),
		Summary:     strings.TrimSpace(req.Post.Summary),
		CoverUrl:    normalizeCoverObjectKey(req.Post.CoverUrl),
		Content:     req.Post.Content,
		Status:      int64(req.Post.Status),
		Tags:        string(tags),
		CategoryId:  req.Post.CategoryId,
		PublishedAt: protoTimestampTime(req.Post.PublishedAt),
	}
	markPendingReview(post)

	var columns []string
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		columns = protoPathsToDBColumns(req.UpdateMask.Paths)
		if len(columns) == 0 {
			return nil, status.Error(codes.InvalidArgument, "no updatable fields")
		}
	}
	columns = ensureColumn(columns, "status")

	rowsAffected, err := s.repo.Update(ctx, req.Id, post, columns...)
	if err != nil {
		s.logger.Error("repo.Update", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, ErrPostNotFound.Error())
	}

	updated, err := s.repo.First(ctx, req.Id)
	if err != nil {
		s.logger.Error("repo.First after update", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdatePostResponse{
		Post: entityToProto(updated),
	}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, err := authenticatedUserID(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.First(ctx, req.Id, "id", "user_id")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ErrPostNotFound.Error())
		}
		s.logger.Error("repo.First", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if existing.UserId != userID {
		return nil, status.Error(codes.PermissionDenied, ErrPermissionDenied.Error())
	}

	rowsAffected, err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		s.logger.Error("repo.Delete", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, ErrPostNotFound.Error())
	}

	return &pb.DeletePostResponse{}, nil
}

func authenticatedUserID(ctx context.Context) (int64, error) {
	user, err := gocauth.ExtractAuthenticatedUser(ctx)
	if err != nil {
		return 0, status.Error(codes.Unauthenticated, err.Error())
	}
	id, err := user.Int64ID()
	if err != nil {
		return 0, status.Error(codes.Unauthenticated, err.Error())
	}
	return id, nil
}

func (s *PostService) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	filter := &repository.ListFilter{
		UserId:     req.UserId,
		CategoryId: req.CategoryId,
		Tag:        req.Tag,
		Status:     int64(req.Status),
	}

	// List responses skip body content but keep display metadata needed by public cards.
	columns := []string{"id", "user_id", "title", "summary", "cover_url", "status", "tags", "category_id", "published_at", "created_at", "updated_at"}

	posts, total, err := s.repo.List(ctx, filter, page, pageSize, columns...)
	if err != nil {
		s.logger.Error("repo.List", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbPosts := make([]*pb.Post, 0, len(posts))
	for _, p := range posts {
		pbPosts = append(pbPosts, entityToProto(p))
	}

	pg := pagination.NewPagination(page, pageSize, total)

	return &pb.ListPostsResponse{
		Total:       pg.Total,
		TotalPages:  int32(pg.LastPage),
		PageSize:    int32(pg.PerPage),
		CurrentPage: int32(pg.CurrentPage),
		Posts:       pbPosts,
	}, nil
}

func entityToProto(p *entity.Post) *pb.Post {
	var tags []string
	_ = json.Unmarshal([]byte(p.Tags), &tags)

	protoPost := &pb.Post{
		Id:         p.ID,
		UserId:     p.UserId,
		Title:      p.Title,
		Summary:    p.Summary,
		CoverUrl:   p.CoverUrl,
		Content:    p.Content,
		Status:     pb.PostStatus(p.Status),
		Tags:       tags,
		CategoryId: p.CategoryId,
		CreatedAt:  timestamppb.New(p.CreatedAt),
		UpdatedAt:  timestamppb.New(p.UpdatedAt),
	}
	if p.PublishedAt.Valid {
		protoPost.PublishedAt = timestamppb.New(p.PublishedAt.Time)
	}

	if p.DeletedAt.Valid {
		protoPost.DeletedAt = timestamppb.New(p.DeletedAt.Time)
	}

	return protoPost
}

func protoPathsToDBColumns(paths []string) []string {
	columns := make([]string, 0, len(paths))
	for _, p := range paths {
		switch p {
		case "title":
			columns = append(columns, "title")
		case "summary":
			columns = append(columns, "summary")
		case "content":
			columns = append(columns, "content")
		case "cover_url":
			columns = append(columns, "cover_url")
		case "published_at":
			columns = append(columns, "published_at")
		case "status":
			columns = append(columns, "status")
		case "tags":
			columns = append(columns, "tags")
		case "category_id":
			columns = append(columns, "category_id")
		}
	}
	return columns
}

func protoTimestampTime(value *timestamppb.Timestamp) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: value.AsTime(), Valid: true}
}

func normalizeCoverObjectKey(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	if parsed, err := url.Parse(value); err == nil && parsed.IsAbs() {
		if !isUploadsPath(parsed.Path) {
			return value
		}
		value = parsed.Path
	}

	value = strings.TrimLeft(value, "/")
	return strings.TrimPrefix(value, "uploads/")
}

func isUploadsPath(value string) bool {
	return value == "/uploads" || strings.HasPrefix(value, "/uploads/")
}

// Published posts always get a stable display time when callers omit one.
func defaultPublishedAt(post *entity.Post) {
	if post.Status != entity.PostStatusPublished || post.PublishedAt.Valid {
		return
	}
	post.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true}
}

func markPendingReview(post *entity.Post) {
	post.Status = entity.PostStatusPendingReview
}

func ensureColumn(columns []string, column string) []string {
	if len(columns) > 0 && !hasColumn(columns, column) {
		return append(columns, column)
	}
	return columns
}

func hasColumn(columns []string, column string) bool {
	for _, v := range columns {
		if v == column {
			return true
		}
	}
	return false
}
