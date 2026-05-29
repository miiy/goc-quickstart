package service

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/miiy/goc-quickstart/auth-service/gen/go/blog/auth/v1"
	"github.com/miiy/goc-quickstart/auth-service/internal/entity"
	"github.com/miiy/goc/auth"
	"github.com/miiy/goc/auth/jwt"
	"github.com/miiy/goc/contrib/wechat/miniprogram"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// MockAuthRepository implements repository.AuthRepository for testing
type MockAuthRepository struct {
	users  map[int64]*entity.User
	nextID int64
	err    error
}

func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{
		users:  make(map[int64]*entity.User),
		nextID: 1,
	}
}

func (m *MockAuthRepository) Create(ctx context.Context, user *entity.User) error {
	if m.err != nil {
		return m.err
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return nil
}

func (m *MockAuthRepository) Update(ctx context.Context, id uint64, user *entity.User, columns ...string) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	if _, ok := m.users[int64(id)]; !ok {
		return 0, nil
	}
	m.users[int64(id)] = user
	return 1, nil
}

func (m *MockAuthRepository) First(ctx context.Context, id uint64, columns ...string) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if user, ok := m.users[int64(id)]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockAuthRepository) FirstByUsername(ctx context.Context, username string, columns ...string) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockAuthRepository) FirstByMpOpenid(ctx context.Context, openid string, columns ...string) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, user := range m.users {
		if user.MpOpenid == openid {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockAuthRepository) FirstByPhone(ctx context.Context, phone string, columns ...string) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, user := range m.users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockAuthRepository) UserExist(ctx context.Context, column, value string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	for _, user := range m.users {
		switch column {
		case entity.UserColumnUsername:
			if user.Username == value {
				return true, nil
			}
		case entity.UserColumnEmail:
			if user.Email == value {
				return true, nil
			}
		case entity.UserColumnPhone:
			if user.Phone == value {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *MockAuthRepository) FirstByIdentifier(ctx context.Context, identifier string) (*auth.AuthenticatedUser, error) {
	user, err := m.FirstByUsername(ctx, identifier)
	if err != nil {
		return nil, err
	}
	return &auth.AuthenticatedUser{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

// MockTokenRepository implements repository.TokenRepository for testing
type MockTokenRepository struct {
	tokens map[string]string
	ttl    map[string]time.Duration
	err    error
}

func NewMockTokenRepository() *MockTokenRepository {
	return &MockTokenRepository{
		tokens: make(map[string]string),
		ttl:    make(map[string]time.Duration),
	}
}

func (m *MockTokenRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if m.err != nil {
		return m.err
	}
	m.tokens[key] = value.(string)
	m.ttl[key] = ttl
	return nil
}

func (m *MockTokenRepository) Get(ctx context.Context, key string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	if v, ok := m.tokens[key]; ok {
		return v, nil
	}
	return "", errors.New("not found")
}

func (m *MockTokenRepository) Del(ctx context.Context, key string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.tokens, key)
	return nil
}

func TestRegisterValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.RegisterRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: &pb.RegisterRequest{
				Username:             "testuser",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr: nil,
		},
		{
			name: "empty username",
			req: &pb.RegisterRequest{
				Username:             "",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr: ErrInvalidArgument,
		},
		{
			name: "passwords differ",
			req: &pb.RegisterRequest{
				Username:             "testuser",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password456",
			},
			wantErr: ErrPasswordsDiffer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registerValidate(tt.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("registerValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := jwt.NewJWTAuth(&jwt.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, jwtAuth, mp).(*AuthService)

	tests := []struct {
		name    string
		req     *pb.RegisterRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "successful registration",
			req: &pb.RegisterRequest{
				Username:             "testuser",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr: false,
		},
		{
			name: "duplicate username",
			req: &pb.RegisterRequest{
				Username:             "testuser",
				Email:                "another@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr: true,
			errCode: codes.AlreadyExists,
		},
		{
			name: "invalid request",
			req: &pb.RegisterRequest{
				Username:             "",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Register(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp == nil || resp.User == nil {
					t.Error("expected response with user")
				}
				if resp.User.Username != tt.req.Username {
					t.Errorf("expected username %s, got %s", tt.req.Username, resp.User.Username)
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := jwt.NewJWTAuth(&jwt.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Create a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	authRepo.users[1] = &entity.User{
		Username: "testuser",
		Password: string(hashedPassword),
		Email:    "test@example.com",
		Status:   entity.UserStatusActive,
	}

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, jwtAuth, mp).(*AuthService)

	tests := []struct {
		name    string
		req     *pb.LoginRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "successful login",
			req: &pb.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req: &pb.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "wrong password",
			req: &pb.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			wantErr: true,
			errCode: codes.NotFound,
		},
		{
			name: "empty credentials",
			req: &pb.LoginRequest{
				Username: "",
				Password: "",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Login(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp == nil {
					t.Error("expected response")
				}
				if resp.TokenType != "Bearer" {
					t.Errorf("expected token type Bearer, got %s", resp.TokenType)
				}
				if resp.AccessToken == "" {
					t.Error("expected access token")
				}
			}
		})
	}
}

func TestAuthService_UserExist(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := jwt.NewJWTAuth(&jwt.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Create a test user
	authRepo.users[1] = &entity.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "13800138000",
		Status:   entity.UserStatusActive,
	}

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, jwtAuth, mp).(*AuthService)

	tests := []struct {
		name      string
		method    string
		value     string
		wantExist bool
	}{
		{
			name:      "username exists",
			method:    "UsernameCheck",
			value:     "testuser",
			wantExist: true,
		},
		{
			name:      "username not exists",
			method:    "UsernameCheck",
			value:     "nonexistent",
			wantExist: false,
		},
		{
			name:      "email exists",
			method:    "EmailCheck",
			value:     "test@example.com",
			wantExist: true,
		},
		{
			name:      "phone exists",
			method:    "PhoneCheck",
			value:     "13800138000",
			wantExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *pb.FieldCheckResponse
			var err error

			switch tt.method {
			case "UsernameCheck":
				resp, err = service.UsernameCheck(context.Background(), &pb.FieldCheckRequest{Value: tt.value})
			case "EmailCheck":
				resp, err = service.EmailCheck(context.Background(), &pb.FieldCheckRequest{Value: tt.value})
			case "PhoneCheck":
				resp, err = service.PhoneCheck(context.Background(), &pb.FieldCheckRequest{Value: tt.value})
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if resp.Exist != tt.wantExist {
				t.Errorf("expected exist=%v, got %v", tt.wantExist, resp.Exist)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := jwt.NewJWTAuth(&jwt.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Store a test token
	testToken := "test-token-123"
	tokenRepo.Set(context.Background(), formatTokenKey(testToken), testToken, time.Hour)

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, jwtAuth, mp).(*AuthService)

	tests := []struct {
		name    string
		req     *pb.LogoutRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name:    "successful logout",
			req:     &pb.LogoutRequest{AccessToken: testToken},
			wantErr: false,
		},
		{
			name:    "logout with empty token",
			req:     &pb.LogoutRequest{AccessToken: ""},
			wantErr: false, // Logout doesn't validate empty token
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Logout(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := jwt.NewJWTAuth(&jwt.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, jwtAuth, mp).(*AuthService)

	// Create a valid token
	claims := jwtAuth.CreateClaims("testuser")
	validToken, _ := jwtAuth.CreateTokenByClaims(claims)

	tests := []struct {
		name     string
		req      *pb.RefreshTokenRequest
		wantErr  bool
		errCode  codes.Code
	}{
		{
			name:    "successful refresh",
			req:     &pb.RefreshTokenRequest{AccessToken: validToken},
			wantErr: false,
		},
		{
			name:    "refresh with invalid token",
			req:     &pb.RefreshTokenRequest{AccessToken: "invalid-token"},
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name:    "refresh with empty token",
			req:     &pb.RefreshTokenRequest{AccessToken: ""},
			wantErr: true,
			errCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.RefreshToken(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp == nil {
					t.Error("expected response")
				}
				if resp.TokenType != "Bearer" {
					t.Errorf("expected token type Bearer, got %s", resp.TokenType)
				}
				if resp.AccessToken == "" {
					t.Error("expected access token")
				}
				if resp.User == nil || resp.User.Username != "testuser" {
					t.Error("expected user with username testuser")
				}
			}
		})
	}
}

func TestAuthService_SendSmsCode(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := jwt.NewJWTAuth(&jwt.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, jwtAuth, mp).(*AuthService)

	tests := []struct {
		name    string
		req     *pb.SendSmsCodeRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name:    "successful send",
			req:     &pb.SendSmsCodeRequest{Phone: "13800138000"},
			wantErr: false,
		},
		{
			name:    "empty phone",
			req:     &pb.SendSmsCodeRequest{Phone: ""},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name:    "whitespace phone",
			req:     &pb.SendSmsCodeRequest{Phone: "   "},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.SendSmsCode(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAuthService_PhoneAuth(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := jwt.NewJWTAuth(&jwt.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Pre-store SMS code
	tokenRepo.Set(context.Background(), "sms_code:13800138000", "123456", 5*time.Minute)

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, jwtAuth, mp).(*AuthService)

	tests := []struct {
		name    string
		req     *pb.PhoneAuthRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "successful auth - new user",
			req: &pb.PhoneAuthRequest{
				Phone: "13800138000",
				Code:  "123456",
			},
			wantErr: false,
		},
		{
			name: "invalid code",
			req: &pb.PhoneAuthRequest{
				Phone: "13800138000",
				Code:  "000000",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "expired code",
			req: &pb.PhoneAuthRequest{
				Phone: "13900139000",
				Code:  "123456",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "empty phone",
			req: &pb.PhoneAuthRequest{
				Phone: "",
				Code:  "123456",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "empty code",
			req: &pb.PhoneAuthRequest{
				Phone: "13800138000",
				Code:  "",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.PhoneAuth(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok || st.Code() != tt.errCode {
						t.Errorf("expected error code %v, got %v", tt.errCode, st.Code())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resp == nil {
					t.Error("expected response")
				}
				if resp.TokenType != "Bearer" {
					t.Errorf("expected token type Bearer, got %s", resp.TokenType)
				}
				if resp.AccessToken == "" {
					t.Error("expected access token")
				}
			}
		})
	}
}

func TestAuthService_PhoneAuth_ExistingUser(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := jwt.NewJWTAuth(&jwt.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Create existing user
	authRepo.users[1] = &entity.User{
		Username: "existing_user",
		Phone:    "13800138001",
		Status:   entity.UserStatusActive,
	}

	// Pre-store SMS code
	tokenRepo.Set(context.Background(), "sms_code:13800138001", "654321", 5*time.Minute)

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, jwtAuth, mp).(*AuthService)

	resp, err := service.PhoneAuth(context.Background(), &pb.PhoneAuthRequest{
		Phone: "13800138001",
		Code:  "654321",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Error("expected response")
	}

	if resp.User == nil || resp.User.Username != "existing_user" {
		t.Error("expected existing user to be returned")
	}
}