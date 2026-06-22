package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"buf.build/go/protovalidate"
	pb "github.com/miiy/goc-quickstart/nova-auth/gen/go/nova/auth/v1"
	"github.com/miiy/goc-quickstart/nova-auth/internal/entity"
	"github.com/miiy/goc-quickstart/nova-auth/internal/repository"
	"github.com/miiy/goc/auth"
	"github.com/miiy/goc/contrib/sms"
	"github.com/miiy/goc/contrib/wechat/miniprogram"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/logger/zap"
	"github.com/miiy/goc/utils/password"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthService implements nova-auth's AuthServiceServer.
type RefreshTTL time.Duration

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	logger      *zap.Logger
	authRepo    repository.AuthRepository
	tokenRepo   repository.TokenRepository
	smsCodeRepo repository.SMSCodeRepository
	jwtAuth     *auth.JWTAuth
	mp          *miniprogram.MiniProgram
	smsSender   sms.Sender
	refreshTTL  time.Duration
}

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrPasswordsDiffer = errors.New("passwords differ")

	ErrUsernameOrEmailExist = errors.New("username or email already exist")

	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

func NewAuthServiceServer(logger *zap.Logger, authRepo repository.AuthRepository, tokenRepo repository.TokenRepository, smsCodeRepo repository.SMSCodeRepository, jwtAuth *auth.JWTAuth, mp *miniprogram.MiniProgram, refreshTTL RefreshTTL) pb.AuthServiceServer {
	return &AuthService{
		logger:      logger,
		authRepo:    authRepo,
		tokenRepo:   tokenRepo,
		smsCodeRepo: smsCodeRepo,
		jwtAuth:     jwtAuth,
		mp:          mp,
		smsSender: sms.NewLogSender(smsLoggerFunc(func(msg string) {
			if logger != nil {
				logger.Info(msg)
			}
		})),
		refreshTTL: time.Duration(refreshTTL),
	}
}

func (s *AuthService) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.verifyAccessToken(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}
	return &pb.VerifyTokenResponse{User: toPBUser(user)}, nil
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

	exist, err = s.authRepo.UserExist(ctx, entity.UserColumnEmail, req.Email)
	if err != nil {
		s.logger.Error("authRepo.UserExist", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exist {
		return nil, status.Error(codes.AlreadyExists, ErrUsernameOrEmailExist.Error())
	}

	hashPasswd, err := password.Hash(req.Password)
	if err != nil {
		s.logger.Error("password.Hash", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	user := entity.User{
		Username:          req.Username,
		Password:          hashPasswd,
		Nickname:          defaultNickname(req.Username),
		Email:             req.Email,
		EmailVerifiedTime: nil,
		Phone:             "",
		Status:            entity.UserStatusActive,
	}

	err = s.authRepo.Create(ctx, &user)
	if err != nil {
		s.logger.Error("authRepo.Create", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{User: toPBUser(&user)}, nil
}

func (s *AuthService) UsernameCheck(ctx context.Context, req *pb.UsernameCheckRequest) (*pb.UsernameCheckResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	exist, err := s.userExist(ctx, entity.UserColumnUsername, req.Value)
	if err != nil {
		return nil, err
	}
	return &pb.UsernameCheckResponse{Exist: exist}, nil
}

func (s *AuthService) EmailCheck(ctx context.Context, req *pb.EmailCheckRequest) (*pb.EmailCheckResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	exist, err := s.userExist(ctx, entity.UserColumnEmail, req.Value)
	if err != nil {
		return nil, err
	}
	return &pb.EmailCheckResponse{Exist: exist}, nil
}

func (s *AuthService) PhoneCheck(ctx context.Context, req *pb.PhoneCheckRequest) (*pb.PhoneCheckResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	exist, err := s.userExist(ctx, entity.UserColumnPhone, req.Value)
	if err != nil {
		return nil, err
	}
	return &pb.PhoneCheckResponse{Exist: exist}, nil
}

// Login authenticates a user with username and password.
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
	if user.Status != entity.UserStatusActive {
		return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.NotFound, ErrWrongPassword.Error())
	}

	pair, err := s.issueTokenPair(ctx, user.ID, user.Username)
	if err != nil {
		s.logger.Error("issueTokenPair", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		TokenType:        "Bearer",
		AccessToken:      pair.AccessToken,
		ExpiresAt:        timestamppb.New(pair.AccessExpiresAt),
		RefreshToken:     pair.RefreshToken,
		RefreshExpiresAt: timestamppb.New(pair.RefreshExpiresAt),
		User:             toPBUser(user),
	}, nil
}

// SendSmsCode sends a verification code to the given phone number.
func (s *AuthService) SendSmsCode(ctx context.Context, req *pb.SendSmsCodeRequest) (*pb.SendSmsCodeResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	req.Phone = strings.TrimSpace(req.Phone)
	if req.Phone == "" {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	code, err := sms.GenerateCode(6)
	if err != nil {
		s.logger.Error("generateSMSCode", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := s.smsSender.Send(ctx, req.Phone, code); err != nil {
		s.logger.Error("smsSender.Send", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := s.smsCodeRepo.Set(ctx, req.Phone, code, 5*time.Minute); err != nil {
		s.logger.Error("smsCodeRepo.Set", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SendSmsCodeResponse{}, nil
}

// PhoneAuth authenticates via phone + sms code, auto registers if not exists.
func (s *AuthService) PhoneAuth(ctx context.Context, req *pb.PhoneAuthRequest) (*pb.PhoneAuthResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	req.Phone = strings.TrimSpace(req.Phone)
	req.Code = strings.TrimSpace(req.Code)
	if req.Phone == "" || req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "phone and code are required")
	}

	storedCode, err := s.smsCodeRepo.Get(ctx, req.Phone)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "code expired or not sent")
	}
	if storedCode != req.Code {
		return nil, status.Error(codes.InvalidArgument, "wrong code")
	}
	_ = s.smsCodeRepo.Del(ctx, req.Phone)

	user, err := s.authRepo.FirstByPhone(ctx, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("authRepo.FirstByPhone", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		username, err := randUserName()
		if err != nil {
			s.logger.Error("randUserName", zap.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}
		user = &entity.User{
			Username:          username,
			Password:          "",
			Nickname:          defaultNickname(username),
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

	pair, err := s.issueTokenPair(ctx, user.ID, user.Username)
	if err != nil {
		s.logger.Error("issueTokenPair", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PhoneAuthResponse{
		TokenType:        "Bearer",
		AccessToken:      pair.AccessToken,
		ExpiresAt:        timestamppb.New(pair.AccessExpiresAt),
		RefreshToken:     pair.RefreshToken,
		RefreshExpiresAt: timestamppb.New(pair.RefreshExpiresAt),
		User:             toPBUser(user),
	}, nil
}

func (s *AuthService) MpLogin(ctx context.Context, req *pb.MpLoginRequest) (*pb.MpLoginResponse, error) {
	if err := mpLoginValidate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.mp.Code2Session(ctx, req.Code)
	if err != nil {
		s.logger.Error("mp.Code2Session", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	var user *entity.User
	user, err = s.authRepo.FirstByMpOpenid(ctx, res.OpenID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("authRepo.FirstByMpOpenid", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		username, err := randUserName()
		if err != nil {
			s.logger.Error("randUserName", zap.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}
		user = &entity.User{
			Username:          username,
			Password:          "",
			Nickname:          defaultNickname(username),
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

	pair, err := s.issueTokenPair(ctx, user.ID, user.Username)
	if err != nil {
		s.logger.Error("issueTokenPair", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.MpLoginResponse{
		TokenType:        "Bearer",
		AccessToken:      pair.AccessToken,
		ExpiresAt:        timestamppb.New(pair.AccessExpiresAt),
		RefreshToken:     pair.RefreshToken,
		RefreshExpiresAt: timestamppb.New(pair.RefreshExpiresAt),
		User:             toPBUser(user),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pair, user, err := s.rotateRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshTokenResponse{
		TokenType:        "Bearer",
		AccessToken:      pair.AccessToken,
		ExpiresAt:        timestamppb.New(pair.AccessExpiresAt),
		RefreshToken:     pair.RefreshToken,
		RefreshExpiresAt: timestamppb.New(pair.RefreshExpiresAt),
		User:             toPBUser(user),
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if err := changePasswordValidate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	authUser, err := auth.ExtractAuthenticatedUser(ctx)
	if err != nil || authUser.ID <= 0 {
		if err == nil {
			err = errors.New("invalid authenticated user")
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	user, err := s.authRepo.First(ctx, uint64(authUser.ID), "id", "password", "status")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
		}
		s.logger.Error("authRepo.First", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if user.Status != entity.UserStatusActive {
		return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrWrongPassword.Error())
	}

	hashPasswd, err := password.Hash(req.NewPassword)
	if err != nil {
		s.logger.Error("password.Hash", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	rowsAffected, err := s.authRepo.Update(ctx, uint64(authUser.ID), &entity.User{Password: hashPasswd}, "password")
	if err != nil {
		s.logger.Error("authRepo.Update", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
	}

	return &pb.ChangePasswordResponse{}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.revokeTokens(ctx, req.AccessToken, req.RefreshToken)
	return &pb.LogoutResponse{}, nil
}
