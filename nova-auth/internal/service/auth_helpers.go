package service

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"strings"

	pb "github.com/miiy/goc-quickstart/nova-auth/gen/go/nova/auth/v1"
	"github.com/miiy/goc-quickstart/nova-auth/internal/entity"
	"github.com/miiy/goc/logger/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type smsLoggerFunc func(string)

func (f smsLoggerFunc) Info(msg string) {
	f(msg)
}

func (s *AuthService) userExist(ctx context.Context, field, value string) (bool, error) {
	exist, err := s.authRepo.UserExist(ctx, field, value)
	if err != nil {
		s.logger.Error("authRepo.UserExist", zap.Error(err))
		return false, status.Error(codes.Internal, err.Error())
	}
	return exist, nil
}

func randUserName() (string, error) {
	var buf [8]byte
	if _, err := cryptorand.Read(buf[:]); err != nil {
		return "", err
	}
	return "user_" + hex.EncodeToString(buf[:]), nil
}

func defaultNickname(username string) string {
	return strings.TrimSpace(username)
}

func toPBUser(user *entity.User) *pb.User {
	if user == nil {
		return nil
	}
	return &pb.User{
		Id:       user.ID,
		Username: user.Username,
	}
}
