package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"

	pb "github.com/miiy/goc-quickstart/auth-service/gen/go/blog/auth/v1"
	"github.com/miiy/goc-quickstart/auth-service/internal/entity"
	"github.com/miiy/goc-quickstart/auth-service/internal/repository"
	"github.com/miiy/goc/auth/jwt"
	"github.com/miiy/goc/contrib/wechat/miniprogram"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	logger    *zap.Logger
	authRepo  repository.AuthRepository
	tokenRepo repository.TokenRepository
	jwtAuth   *jwt.JWTAuth
	mp        *miniprogram.MiniProgram
}

const (
	AuthTokenKey = "user_token:%s" // user_token:md5({user_id})
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrPasswordsDiffer = errors.New("passwords differ")

	ErrUsernameOrEmailExist = errors.New("username or email already exist")

	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

func NewAuthServiceServer(logger *zap.Logger, authRepo repository.AuthRepository, tokenRepo repository.TokenRepository, jwtAuth *jwt.JWTAuth, mp *miniprogram.MiniProgram) pb.AuthServiceServer {
	return &AuthService{
		logger:    logger,
		authRepo:  authRepo,
		tokenRepo: tokenRepo,
		jwtAuth:   jwtAuth,
		mp:        mp,
	}
}

func registerValidate(req *pb.RegisterRequest) error {
	// trim space
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.PasswordConfirmation = strings.TrimSpace(req.PasswordConfirmation)

	// validate
	if req.Username == "" || req.Email == "" || req.Password == "" || req.PasswordConfirmation == "" {
		return ErrInvalidArgument
	}
	if req.Password != req.PasswordConfirmation {
		return ErrPasswordsDiffer
	}
	return nil
}
func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := registerValidate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	exist, err := s.authRepo.UserExist(ctx, entity.UserColumnUsername, req.Username)
	if err != nil {
		s.logger.Error("authRepo.UserExist", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exist {
		return nil, status.Error(codes.AlreadyExists, ErrUsernameOrEmailExist.Error())
	}

	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("bcrypt.GenerateFromPassword", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	user := entity.User{
		Username:          req.Username,
		Password:          string(hashPasswd),
		Email:             req.Email,
		EmailVerifiedTime: nil,
		Phone:             "",
		Status:            entity.UserStatusActive,
	}

	// register
	err = s.authRepo.Create(ctx, &user)
	if err != nil {
		s.logger.Error("authRepo.Create", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Username: user.Username,
		},
	}, nil
}

func (s *AuthService) UsernameCheck(ctx context.Context, req *pb.FieldCheckRequest) (*pb.FieldCheckResponse, error) {
	return s.userExist(ctx, entity.UserColumnUsername, req.Value)
}

func (s *AuthService) EmailCheck(ctx context.Context, req *pb.FieldCheckRequest) (*pb.FieldCheckResponse, error) {
	return s.userExist(ctx, entity.UserColumnEmail, req.Value)
}

func (s *AuthService) PhoneCheck(ctx context.Context, req *pb.FieldCheckRequest) (*pb.FieldCheckResponse, error) {
	return s.userExist(ctx, entity.UserColumnPhone, req.Value)
}

func (s *AuthService) userExist(ctx context.Context, field, value string) (*pb.FieldCheckResponse, error) {
	exist, err := s.authRepo.UserExist(ctx, field, value)
	if err != nil {
		s.logger.Error("authRepo.UserExist", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.FieldCheckResponse{
		Exist: exist,
	}, nil
}

func loginValidate(req *pb.LoginRequest) error {
	if req.Username == "" || req.Password == "" {
		return ErrInvalidArgument
	}
	return nil
}

// Login
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := loginValidate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.authRepo.FirstByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
		}
		s.logger.Error("authRepo.FirstByUsername", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.NotFound, ErrWrongPassword.Error())
	}

	claims := s.jwtAuth.CreateClaims(user.Username)
	token, err := s.jwtAuth.CreateTokenByClaims(claims)
	if err != nil {
		s.logger.Error("jwtAuth.CreateTokenByClaims", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	// store token
	expiresTime := claims.ExpiresAt.Time
	err = s.tokenRepo.Set(ctx, formatTokenKey(token), token, time.Now().Sub(expiresTime))
	if err != nil {
		s.logger.Error("tokenRepo.Set", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		TokenType:   "Bearer",
		AccessToken: token,
		ExpiresAt:   timestamppb.New(expiresTime),
		User: &pb.User{
			Username: user.Username,
		},
	}, nil
}

func mpLoginValidate(req *pb.MpLoginRequest) error {
	if req.Code == "" {
		return errors.New("code can not empty")
	}
	return nil
}

func (s *AuthService) MpLogin(ctx context.Context, req *pb.MpLoginRequest) (*pb.LoginResponse, error) {
	if err := mpLoginValidate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// get openid
	res, err := s.mp.Code2Session(ctx, req.Code)
	if err != nil {
		s.logger.Error("mp.Code2Session", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	var user *entity.User
	// get user by openid
	user, err = s.authRepo.FirstByMpOpenid(ctx, res.OpenID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("authRepo.FirstByMpOpenid", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	// if not found, create user
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		user = &entity.User{
			Username:          randUserName(),
			Password:          "",
			Email:             "",
			EmailVerifiedTime: nil,
			Phone:             "",
			Unionid:           res.UnionID,
			MpOpenid:          res.OpenID,
			MpSessionKey:      res.SessionKey,
			Status:            entity.UserStatusActive,
		}

		err = s.authRepo.Create(ctx, user)
		if err != nil {
			s.logger.Error("authRepo.Create", zap.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}

	}

	// create jwt token
	claims := s.jwtAuth.CreateClaims(user.Username)
	token, err := s.jwtAuth.CreateTokenByClaims(claims)
	if err != nil {
		s.logger.Error("jwtAuth.CreateTokenByClaims", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	// store token
	expiresTime := claims.ExpiresAt.Time
	ttl := expiresTime.Sub(time.Now())
	err = s.tokenRepo.Set(ctx, formatTokenKey(token), token, ttl)
	if err != nil {
		s.logger.Error("tokenRepo.Set", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		TokenType:   "Bearer",
		AccessToken: token,
		ExpiresAt:   timestamppb.New(expiresTime),
		User: &pb.User{
			Username: user.Username,
		},
	}, nil

}

const (
	SmsCodeKey = "sms_code:%s" // sms_code:{phone}
)

// SendSmsCode sends a verification code to the given phone number
func (s *AuthService) SendSmsCode(ctx context.Context, req *pb.SendSmsCodeRequest) (*pb.SendSmsCodeResponse, error) {
	req.Phone = strings.TrimSpace(req.Phone)
	if req.Phone == "" {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	// TODO: integrate with SMS provider to send code
	// for now, store a mock code "123456" in redis
	code := "123456"
	err := s.tokenRepo.Set(ctx, fmt.Sprintf(SmsCodeKey, req.Phone), code, 5*time.Minute)
	if err != nil {
		s.logger.Error("tokenRepo.Set sms code", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SendSmsCodeResponse{}, nil
}

// PhoneAuth authenticates via phone + sms code, auto registers if not exists
func (s *AuthService) PhoneAuth(ctx context.Context, req *pb.PhoneAuthRequest) (*pb.PhoneAuthResponse, error) {
	req.Phone = strings.TrimSpace(req.Phone)
	req.Code = strings.TrimSpace(req.Code)
	if req.Phone == "" || req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "phone and code are required")
	}

	// verify code
	storedCode, err := s.tokenRepo.Get(ctx, fmt.Sprintf(SmsCodeKey, req.Phone))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "code expired or not sent")
	}
	if storedCode != req.Code {
		return nil, status.Error(codes.InvalidArgument, "wrong code")
	}
	// delete used code
	_ = s.tokenRepo.Del(ctx, fmt.Sprintf(SmsCodeKey, req.Phone))

	// find or create user
	user, err := s.authRepo.FirstByPhone(ctx, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("authRepo.FirstByPhone", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &entity.User{
			Username:          randUserName(),
			Password:          "",
			Email:             "",
			EmailVerifiedTime: nil,
			Phone:             req.Phone,
			Status:            entity.UserStatusActive,
		}
		if err := s.authRepo.Create(ctx, user); err != nil {
			s.logger.Error("authRepo.Create", zap.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// create jwt token
	claims := s.jwtAuth.CreateClaims(user.Username)
	token, err := s.jwtAuth.CreateTokenByClaims(claims)
	if err != nil {
		s.logger.Error("jwtAuth.CreateTokenByClaims", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	expiresTime := claims.ExpiresAt.Time
	err = s.tokenRepo.Set(ctx, formatTokenKey(token), token, expiresTime.Sub(time.Now()))
	if err != nil {
		s.logger.Error("tokenRepo.Set", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PhoneAuthResponse{
		TokenType:   "Bearer",
		AccessToken: token,
		ExpiresAt:   timestamppb.New(expiresTime),
		User: &pb.User{
			Username: user.Username,
		},
	}, nil
}

// RefreshToken
// 1. validate old token
// 2. delete old token
// 3. create new token
func (s *AuthService) RefreshToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// validate old token and create new token
	oldClaims, err := s.jwtAuth.ParseToken(request.AccessToken)
	if err != nil {
		s.logger.Error("jwtAuth.ParseToken", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	// delete old token
	err = s.tokenRepo.Del(ctx, formatTokenKey(request.AccessToken))
	if err != nil {
		s.logger.Error("tokenRepo.Del", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	// create new token
	claims := s.jwtAuth.CreateClaims(oldClaims.Username)
	token, err := s.jwtAuth.CreateTokenByClaims(claims)
	if err != nil {
		s.logger.Error("jwtAuth.CreateTokenByClaims", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	expiresTime := claims.ExpiresAt.Time
	return &pb.RefreshTokenResponse{
		TokenType:   "Bearer",
		AccessToken: token,
		ExpiresAt:   timestamppb.New(expiresTime),
		User:        &pb.User{Username: claims.Username},
	}, nil
}

// Logout
// 1. delete token
func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := s.tokenRepo.Del(ctx, formatTokenKey(req.AccessToken))
	if err != nil {
		s.logger.Error("tokenRepo.Del", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LogoutResponse{}, nil
}

func formatTokenKey(token string) string {
	return fmt.Sprintf(AuthTokenKey, fmt.Sprintf("%x", md5.Sum([]byte(token))))
}

func randUserName() string {
	return fmt.Sprintf("用户_%d", time.Now().Unix())
}
