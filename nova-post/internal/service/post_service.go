package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

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
	logger *zap.Logger
	repo   repository.PostRepository
}

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrPostNotFound     = errors.New("post not found")
	ErrPermissionDenied = errors.New("permission denied")
)

func NewPostServiceServer(logger *zap.Logger, repo repository.PostRepository) pb.PostServiceServer {
	return &PostService{
		logger: logger,
		repo:   repo,
	}
}

func (s *PostService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
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
	if req.Post == nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}
	if strings.TrimSpace(req.Post.Title) == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	user, err := authenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	tags, _ := json.Marshal(req.Post.Tags)

	post := &entity.Post{
		AuthorId:   user.ID,
		Title:      strings.TrimSpace(req.Post.Title),
		CoverUrl:   strings.TrimSpace(req.Post.CoverUrl),
		Content:    req.Post.Content,
		Status:     int64(req.Post.Status),
		Tags:       string(tags),
		CategoryId: req.Post.CategoryId,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		s.logger.Error("repo.Create", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreatePostResponse{
		Post: entityToProto(post),
	}, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	if req.Id <= 0 || req.Post == nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}
	user, err := authenticatedUser(ctx)
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
	if existing.AuthorId != user.ID {
		return nil, status.Error(codes.PermissionDenied, ErrPermissionDenied.Error())
	}

	tags, _ := json.Marshal(req.Post.Tags)

	post := &entity.Post{
		AuthorId:   existing.AuthorId,
		Title:      strings.TrimSpace(req.Post.Title),
		CoverUrl:   strings.TrimSpace(req.Post.CoverUrl),
		Content:    req.Post.Content,
		Status:     int64(req.Post.Status),
		Tags:       string(tags),
		CategoryId: req.Post.CategoryId,
	}

	var columns []string
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		columns = protoPathsToDBColumns(req.UpdateMask.Paths)
		if len(columns) == 0 {
			return nil, status.Error(codes.InvalidArgument, "no updatable fields")
		}
	}

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
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}
	user, err := authenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.First(ctx, req.Id, "id", "author_id")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ErrPostNotFound.Error())
		}
		s.logger.Error("repo.First", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if existing.AuthorId != user.ID {
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

func (s *PostService) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	filter := &repository.ListFilter{
		AuthorId:   req.AuthorId,
		CategoryId: req.CategoryId,
		Tag:        req.Tag,
	}

	// 列表只查询需要的字段，不查询 content
	columns := []string{"id", "author_id", "title", "cover_url", "status", "tags", "category_id", "created_at", "updated_at"}

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
		AuthorId:   p.AuthorId,
		Title:      p.Title,
		CoverUrl:   p.CoverUrl,
		Content:    p.Content,
		Status:     pb.PostStatus(p.Status),
		Tags:       tags,
		CategoryId: p.CategoryId,
		CreatedAt:  timestamppb.New(p.CreatedAt),
		UpdatedAt:  timestamppb.New(p.UpdatedAt),
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
		case "content":
			columns = append(columns, "content")
		case "cover_url":
			columns = append(columns, "cover_url")
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
