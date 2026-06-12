package service

import (
	"context"
	"errors"

	pb "github.com/miiy/goc-quickstart/nova-user/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-user/internal/entity"
	"github.com/miiy/goc-quickstart/nova-user/internal/repository"
	gocauth "github.com/miiy/goc/auth"
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
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrUserNotFound     = errors.New("user not found")
	ErrPermissionDenied = errors.New("permission denied")
)

func NewUserServiceServer(logger *zap.Logger, repo repository.UserRepository) pb.UserServiceServer {
	return &UserService{
		logger: logger,
		repo:   repo,
	}
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}
	if err := requireSelf(ctx, req.Id); err != nil {
		return nil, err
	}

	user, err := s.repo.First(ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
		}
		s.logger.Error("repo.First", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetUserResponse{User: entityToProto(user)}, nil
}

func (s *UserService) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
	ids := normalizedUserIDs(req.GetIds())
	if len(ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}

	users, err := s.repo.FindByIDs(ctx, ids, "id", "username", "nickname", "avatar")
	if err != nil {
		s.logger.Error("repo.FindByIDs", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.BatchGetUsersResponse{
		Users: make([]*pb.User, 0, len(users)),
	}
	for _, user := range users {
		resp.Users = append(resp.Users, publicUserToProto(user))
	}
	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Id <= 0 || req.User == nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}
	if err := requireSelf(ctx, req.Id); err != nil {
		return nil, err
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

	return &pb.UpdateUserResponse{User: entityToProto(updated)}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	if _, err := authenticatedUser(ctx); err != nil {
		return nil, err
	}
	return nil, status.Error(codes.PermissionDenied, ErrPermissionDenied.Error())
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

func requireSelf(ctx context.Context, id int64) error {
	user, err := authenticatedUser(ctx)
	if err != nil {
		return err
	}
	if user.ID != id {
		return status.Error(codes.PermissionDenied, ErrPermissionDenied.Error())
	}
	return nil
}

func normalizedUserIDs(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func publicUserToProto(u *entity.User) *pb.User {
	return &pb.User{
		Id:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
	}
}

func entityToProto(u *entity.User) *pb.User {
	user := &pb.User{
		Id:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Email:     u.Email,
		Phone:     u.Phone,
		Status:    pb.UserStatus(u.Status),
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
