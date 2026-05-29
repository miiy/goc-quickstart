package service

import (
	"context"
	"errors"

	pb "github.com/miiy/goc-quickstart/user-service/gen/go/blog/user/v1"
	"github.com/miiy/goc-quickstart/user-service/internal/entity"
	"github.com/miiy/goc-quickstart/user-service/internal/repository"
	"github.com/miiy/goc/pagination"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	logger *zap.Logger
	repo   repository.UserRepository
}

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrUserNotFound    = errors.New("user not found")
)

func NewUserServiceServer(logger *zap.Logger, repo repository.UserRepository) pb.UserServiceServer {
	return &UserService{
		logger: logger,
		repo:   repo,
	}
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}

	user, err := s.repo.First(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
		}
		s.logger.Error("repo.First", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return entityToProto(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	if req.Id <= 0 || req.User == nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}

	_, err := s.repo.First(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
		}
		s.logger.Error("repo.First", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	user := &entity.User{
		Nickname: req.User.Nickname,
		Avatar:   req.User.Avatar,
		Email:    req.User.Email,
		Phone:    req.User.Phone,
		Status:   int64(req.User.Status),
	}

	var columns []string
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		columns = protoPathsToDBColumns(req.UpdateMask.Paths)
	}

	rowsAffected, err := s.repo.Update(ctx, req.Id, user, columns...)
	if err != nil {
		s.logger.Error("repo.Update", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
	}

	updated, err := s.repo.First(ctx, req.Id)
	if err != nil {
		s.logger.Error("repo.First after update", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return entityToProto(updated), nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	page := int64(req.Page)
	pageSize := int64(req.PageSize)

	// 列表只查询需要的字段，不查询 password
	columns := []string{"id", "username", "nickname", "avatar", "email", "phone", "status", "created_at", "updated_at"}

	users, total, err := s.repo.List(ctx, page, pageSize, columns...)
	if err != nil {
		s.logger.Error("repo.List", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbUsers := make([]*pb.User, 0, len(users))
	for _, u := range users {
		pbUsers = append(pbUsers, entityToProto(u))
	}

	pg := pagination.NewPagination(page, pageSize, total)

	return &pb.ListUsersResponse{
		Total:       pg.Total,
		TotalPages:  int32(pg.LastPage),
		PageSize:    int32(pg.PerPage),
		CurrentPage: int32(pg.CurrentPage),
		Users:       pbUsers,
	}, nil
}

func entityToProto(u *entity.User) *pb.User {
	user := &pb.User{
		Id:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Email:    u.Email,
		Phone:    u.Phone,
		Status:   pb.UserStatus(u.Status),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}

	if u.EmailVerifiedTime != nil {
		user.EmailVerifiedTime = timestamppb.New(u.EmailVerifiedTime.Time)
	}

	if u.DeletedAt.Valid {
		user.DeletedAt = timestamppb.New(u.DeletedAt.Time)
	}

	return user
}

func protoPathsToDBColumns(paths []string) []string {
	columns := make([]string, 0, len(paths))
	for _, p := range paths {
		switch p {
		case "nickname":
			columns = append(columns, "nickname")
		case "avatar":
			columns = append(columns, "avatar")
		case "email":
			columns = append(columns, "email")
		case "phone":
			columns = append(columns, "phone")
		case "status":
			columns = append(columns, "status")
		}
	}
	return columns
}
