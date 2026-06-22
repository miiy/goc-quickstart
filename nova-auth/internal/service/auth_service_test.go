package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	pb "github.com/miiy/goc-quickstart/nova-auth/gen/go/nova/auth/v1"
	"github.com/miiy/goc-quickstart/nova-auth/internal/entity"
	"github.com/miiy/goc/auth"
	"github.com/miiy/goc/contrib/wechat/miniprogram"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/logger/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	existing, ok := m.users[int64(id)]
	if !ok {
		return 0, nil
	}
	if len(columns) == 0 {
		m.users[int64(id)] = user
		return 1, nil
	}
	for _, column := range columns {
		switch column {
		case "password":
			existing.Password = user.Password
		case "nickname":
			existing.Nickname = user.Nickname
		case "avatar":
			existing.Avatar = user.Avatar
		case "email":
			existing.Email = user.Email
		case "phone":
			existing.Phone = user.Phone
		case "status":
			existing.Status = user.Status
		}
	}
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

func (m *MockTokenRepository) SetKeepTTL(ctx context.Context, key string, value interface{}) error {
	if m.err != nil {
		return m.err
	}
	if _, ok := m.tokens[key]; ok {
		m.tokens[key] = value.(string)
	}
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
	delete(m.ttl, key)
	return nil
}

// MockSMSCodeRepository implements repository.SMSCodeRepository for testing.
type MockSMSCodeRepository struct {
	codes map[string]string
	err   error
}

func NewMockSMSCodeRepository() *MockSMSCodeRepository {
	return &MockSMSCodeRepository{codes: make(map[string]string)}
}

func (m *MockSMSCodeRepository) Set(ctx context.Context, phone, code string, ttl time.Duration) error {
	if m.err != nil {
		return m.err
	}
	m.codes[phone] = code
	return nil
}

func (m *MockSMSCodeRepository) Get(ctx context.Context, phone string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	if v, ok := m.codes[phone]; ok {
		return v, nil
	}
	return "", errors.New("not found")
}

func (m *MockSMSCodeRepository) Del(ctx context.Context, phone string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.codes, phone)
	return nil
}

// CompareAndSet mirrors the Redis Lua semantics for tests (not truly concurrent,
// but exercises the compare-and-replace logic).
func (m *MockTokenRepository) CompareAndSet(ctx context.Context, key, oldVal, newVal string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	if cur, ok := m.tokens[key]; ok && cur == oldVal {
		m.tokens[key] = newVal
		return true, nil
	}
	return false, nil
}

func TestRegisterValidate(t *testing.T) {
	tests := []struct {
		name      string
		req       *pb.RegisterRequest
		wantErr   bool
		wantErrIs error
	}{
		{
			name: "valid request",
			req: &pb.RegisterRequest{
				Username:             "testuser",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr: false,
		},
		{
			name: "empty username",
			req: &pb.RegisterRequest{
				Username:             "",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr: true,
		},
		{
			name: "passwords differ",
			req: &pb.RegisterRequest{
				Username:             "testuser",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password456",
			},
			wantErr:   true,
			wantErrIs: ErrPasswordsDiffer,
		},
		{
			name: "password too short",
			req: &pb.RegisterRequest{
				Username:             "testuser",
				Email:                "test@example.com",
				Password:             "ab1",
				PasswordConfirmation: "ab1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registerValidate(tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("registerValidate() error = nil, want error")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Errorf("registerValidate() error = %v, wantErrIs %v", err, tt.wantErrIs)
				}
				return
			}
			if err != nil {
				t.Errorf("registerValidate() error = %v, want nil", err)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, mp, RefreshTTL(time.Hour)).(*AuthService)

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
				if resp.User.Id <= 0 {
					t.Error("expected user id")
				}
				if resp.User.Username != tt.req.Username {
					t.Errorf("expected username %s, got %s", tt.req.Username, resp.User.Username)
				}
				user := authRepo.users[resp.User.Id]
				if user == nil {
					t.Fatalf("expected stored user %d", resp.User.Id)
				}
				if user.Nickname != user.Username {
					t.Fatalf("expected nickname %q, got %q", user.Username, user.Nickname)
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Create a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	loginUser := &entity.User{
		Username: "testuser",
		Password: string(hashedPassword),
		Email:    "test@example.com",
		Status:   entity.UserStatusActive,
	}
	loginUser.ID = 1
	authRepo.users[1] = loginUser

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, mp, RefreshTTL(time.Hour)).(*AuthService)

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
				if resp.User == nil || resp.User.Id != 1 || resp.User.Username != "testuser" {
					t.Fatalf("unexpected login user: %+v", resp.User)
				}
				if ttl := tokenRepo.ttl[formatTokenKey(resp.AccessToken)]; ttl <= 0 {
					t.Fatalf("expected positive token ttl, got %v", ttl)
				}
			}
		})
	}
}

func TestAuthService_ChangePassword(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.ChangePasswordRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "successful change",
			req: &pb.ChangePasswordRequest{
				OldPassword:             "password123",
				NewPassword:             "newpass123",
				NewPasswordConfirmation: "newpass123",
			},
		},
		{
			name: "wrong old password",
			req: &pb.ChangePasswordRequest{
				OldPassword:             "wrong-password",
				NewPassword:             "newpass123",
				NewPasswordConfirmation: "newpass123",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "new password confirmation mismatch",
			req: &pb.ChangePasswordRequest{
				OldPassword:             "password123",
				NewPassword:             "newpass123",
				NewPasswordConfirmation: "otherpass123",
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zap.NewNop()
			authRepo := NewMockAuthRepository()
			tokenRepo := NewMockTokenRepository()
			jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})
			mp := &miniprogram.MiniProgram{}

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
			authRepo.users[1] = &entity.User{
				Username: "testuser",
				Password: string(hashedPassword),
				Status:   entity.UserStatusActive,
			}
			service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, mp, RefreshTTL(time.Hour)).(*AuthService)
			ctx := auth.InjectAuthenticatedUser(context.Background(), &auth.AuthenticatedUser{
				ID:       1,
				Username: "testuser",
			})

			_, err := service.ChangePassword(ctx, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if code := status.Code(err); code != tt.errCode {
					t.Fatalf("expected error code %v, got %v", tt.errCode, code)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if err := bcrypt.CompareHashAndPassword([]byte(authRepo.users[1].Password), []byte(tt.req.NewPassword)); err != nil {
				t.Fatal("expected stored password to match new password")
			}
		})
	}
}

func TestAuthService_UserExist(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Create a test user
	authRepo.users[1] = &entity.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "13800138000",
		Status:   entity.UserStatusActive,
	}

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, mp, RefreshTTL(time.Hour)).(*AuthService)

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
			var exist bool
			var err error

			switch tt.method {
			case "UsernameCheck":
				resp, checkErr := service.UsernameCheck(context.Background(), &pb.UsernameCheckRequest{Value: tt.value})
				err = checkErr
				exist = resp.GetExist()
			case "EmailCheck":
				resp, checkErr := service.EmailCheck(context.Background(), &pb.EmailCheckRequest{Value: tt.value})
				err = checkErr
				exist = resp.GetExist()
			case "PhoneCheck":
				resp, checkErr := service.PhoneCheck(context.Background(), &pb.PhoneCheckRequest{Value: tt.value})
				err = checkErr
				exist = resp.GetExist()
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if exist != tt.wantExist {
				t.Errorf("expected exist=%v, got %v", tt.wantExist, exist)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Store a test token
	testToken := "test-token-123"
	tokenRepo.Set(context.Background(), formatTokenKey(testToken), testToken, time.Hour)

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, mp, RefreshTTL(time.Hour)).(*AuthService)

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
			wantErr: true,
			errCode: codes.InvalidArgument,
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
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})

	activeUser := &entity.User{Username: "testuser", Status: entity.UserStatusActive}
	activeUser.ID = 1
	authRepo.users[1] = activeUser

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, &miniprogram.MiniProgram{}, RefreshTTL(time.Hour)).(*AuthService)

	// issue an initial token pair (as Login does)
	pair, err := service.issueTokenPair(context.Background(), 1, "testuser")
	if err != nil {
		t.Fatalf("issueTokenPair: %v", err)
	}
	oldRefreshKey := formatRefreshTokenKey(pair.RefreshToken)
	oldRefreshTTL := 30 * time.Minute
	tokenRepo.ttl[oldRefreshKey] = oldRefreshTTL

	// empty refresh token -> InvalidArgument
	if _, err := service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{}); err == nil {
		t.Fatal("expected error for empty refresh token")
	}
	// unknown refresh token -> Unauthenticated
	if _, err := service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: "unknown-opaque"}); err == nil {
		t.Fatal("expected error for unknown refresh token")
	}

	// successful refresh -> rotation
	resp, err := service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: pair.RefreshToken})
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatal("expected new access + refresh token")
	}
	if resp.User == nil || resp.User.Id != 1 || resp.User.Username != "testuser" {
		t.Fatalf("unexpected user: %+v", resp.User)
	}
	if resp.RefreshToken == pair.RefreshToken {
		t.Fatal("expected rotation to issue a new refresh token")
	}
	// old refresh now revoked (key retained for reuse detection)
	oldRec, _ := tokenRepo.Get(context.Background(), oldRefreshKey)
	if !strings.Contains(oldRec, refreshStatusRevoked) {
		t.Fatalf("expected old refresh revoked, got %s", oldRec)
	}
	if got := tokenRepo.ttl[oldRefreshKey]; got != oldRefreshTTL {
		t.Fatalf("expected old refresh ttl to be preserved, got %v want %v", got, oldRefreshTTL)
	}
	// new refresh active
	newRec, _ := tokenRepo.Get(context.Background(), formatRefreshTokenKey(resp.RefreshToken))
	if !strings.Contains(newRec, refreshStatusActive) {
		t.Fatalf("expected new refresh active, got %s", newRec)
	}

	// reuse the rotated refresh token -> reuse detection -> family revoked
	if _, err := service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: pair.RefreshToken}); err == nil {
		t.Fatal("expected reuse to fail")
	}
	// family now revoked: the rotated (active) token must also be rejected
	if _, err := service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: resp.RefreshToken}); err == nil {
		t.Fatal("expected family-revoked token to be rejected")
	}
}

func TestAuthService_RefreshToken_DisabledUser(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})

	user := &entity.User{Username: "u", Status: entity.UserStatusActive}
	user.ID = 1
	authRepo.users[1] = user

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, &miniprogram.MiniProgram{}, RefreshTTL(time.Hour)).(*AuthService)

	pair, err := service.issueTokenPair(context.Background(), 1, "u")
	if err != nil {
		t.Fatalf("issueTokenPair: %v", err)
	}
	// disable the user after a token pair was issued
	user.Status = entity.UserStatusDisabled

	_, err = service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: pair.RefreshToken})
	if err == nil {
		t.Fatal("expected disabled user to be rejected")
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got %v", err)
	}
}

func TestAuthService_Logout_RevokeAccessAndRefresh(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})

	user := &entity.User{Username: "u", Status: entity.UserStatusActive}
	user.ID = 1
	authRepo.users[1] = user

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, &miniprogram.MiniProgram{}, RefreshTTL(time.Hour)).(*AuthService)

	pair, err := service.issueTokenPair(context.Background(), 1, "u")
	if err != nil {
		t.Fatalf("issueTokenPair: %v", err)
	}
	if _, err := tokenRepo.Get(context.Background(), formatTokenKey(pair.AccessToken)); err != nil {
		t.Fatal("expected access token stored before logout")
	}
	if _, err := tokenRepo.Get(context.Background(), formatRefreshTokenKey(pair.RefreshToken)); err != nil {
		t.Fatal("expected refresh token stored before logout")
	}
	refreshKey := formatRefreshTokenKey(pair.RefreshToken)
	refreshTTL := 15 * time.Minute
	tokenRepo.ttl[refreshKey] = refreshTTL

	if _, err := service.Logout(context.Background(), &pb.LogoutRequest{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}); err != nil {
		t.Fatalf("logout: %v", err)
	}

	if _, err := tokenRepo.Get(context.Background(), formatTokenKey(pair.AccessToken)); err == nil {
		t.Fatal("expected access token revoked after logout")
	}
	rec, err := tokenRepo.Get(context.Background(), refreshKey)
	if err != nil {
		t.Fatal("expected refresh record retained after logout for reuse detection")
	}
	if !strings.Contains(rec, refreshStatusRevoked) {
		t.Fatalf("expected refresh revoked after logout, got %s", rec)
	}
	if got := tokenRepo.ttl[refreshKey]; got != refreshTTL {
		t.Fatalf("expected refresh ttl to be preserved after logout, got %v want %v", got, refreshTTL)
	}
}

func TestAuthService_RefreshToken_FamilyIsolatedAcrossLogins(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})

	user := &entity.User{Username: "u", Status: entity.UserStatusActive}
	user.ID = 1
	authRepo.users[1] = user

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, &miniprogram.MiniProgram{}, RefreshTTL(time.Hour)).(*AuthService)

	// two logins -> two independent families
	pair1, err := service.issueTokenPair(context.Background(), 1, "u")
	if err != nil {
		t.Fatalf("issueTokenPair 1: %v", err)
	}
	pair2, err := service.issueTokenPair(context.Background(), 1, "u")
	if err != nil {
		t.Fatalf("issueTokenPair 2: %v", err)
	}

	// rotate family A once, then reuse -> family A revoked
	if _, err := service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: pair1.RefreshToken}); err != nil {
		t.Fatalf("first rotate family A: %v", err)
	}
	if _, err := service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: pair1.RefreshToken}); err == nil {
		t.Fatal("expected reuse of family A to fail")
	}

	// family B must be unaffected: pair2 still rotates successfully
	resp, err := service.RefreshToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: pair2.RefreshToken})
	if err != nil {
		t.Fatalf("family B should be independent, got %v", err)
	}
	if resp.AccessToken == "" {
		t.Fatal("expected new access token for family B")
	}
}

func TestAuthService_VerifyToken(t *testing.T) {
	logger := zap.NewNop()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})

	validToken, err := jwtAuth.CreateTokenByClaims(jwtAuth.CreateClaims(1, "testuser"))
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	newActiveUser := func() *entity.User {
		u := &entity.User{Username: "testuser", Status: entity.UserStatusActive}
		u.ID = 1
		return u
	}

	tests := []struct {
		name     string
		token    string
		setup    func(*MockAuthRepository, *MockTokenRepository)
		wantErr  bool
		errCode  codes.Code
		wantUser *pb.User
	}{
		{
			name:    "empty token",
			token:   "",
			setup:   func(*MockAuthRepository, *MockTokenRepository) {},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name:    "invalid signature",
			token:   "invalid-token",
			setup:   func(*MockAuthRepository, *MockTokenRepository) {},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			name:  "valid active user",
			token: validToken,
			setup: func(ar *MockAuthRepository, tr *MockTokenRepository) {
				ar.users[1] = newActiveUser()
				_ = tr.Set(context.Background(), formatTokenKey(validToken), validToken, time.Hour)
			},
			wantUser: &pb.User{Id: 1, Username: "testuser"},
		},
		{
			name:  "revoked token not stored",
			token: validToken,
			setup: func(ar *MockAuthRepository, tr *MockTokenRepository) {
				ar.users[1] = newActiveUser()
				// token intentionally not stored -> simulates post-logout
			},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			name:  "disabled user",
			token: validToken,
			setup: func(ar *MockAuthRepository, tr *MockTokenRepository) {
				u := &entity.User{Username: "testuser", Status: entity.UserStatusDisabled}
				u.ID = 1
				ar.users[1] = u
				_ = tr.Set(context.Background(), formatTokenKey(validToken), validToken, time.Hour)
			},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
		{
			name:  "token repo error fails closed",
			token: validToken,
			setup: func(ar *MockAuthRepository, tr *MockTokenRepository) {
				ar.users[1] = newActiveUser()
				tr.err = errors.New("redis down")
			},
			wantErr: true,
			errCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authRepo := NewMockAuthRepository()
			tokenRepo := NewMockTokenRepository()
			tt.setup(authRepo, tokenRepo)

			service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, &miniprogram.MiniProgram{}, RefreshTTL(time.Hour)).(*AuthService)
			resp, err := service.VerifyToken(context.Background(), &pb.VerifyTokenRequest{AccessToken: tt.token})

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				st, ok := status.FromError(err)
				if !ok || st.Code() != tt.errCode {
					t.Errorf("expected error code %v, got %v (%v)", tt.errCode, st.Code(), err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil || resp.User == nil {
				t.Fatal("expected user in response")
			}
			if resp.User.Id != tt.wantUser.Id || resp.User.Username != tt.wantUser.Username {
				t.Errorf("expected user %+v, got %+v", tt.wantUser, resp.User)
			}
		})
	}
}

func TestVerifyTokenServesCachedUser(t *testing.T) {
	logger := zap.NewNop()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})

	token, err := jwtAuth.CreateTokenByClaims(jwtAuth.CreateClaims(1, "alice"))
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	tokenRepo := NewMockTokenRepository()
	// token is valid (not revoked)
	_ = tokenRepo.Set(context.Background(), formatTokenKey(token), token, time.Hour)
	// user cached as active; authRepo intentionally has NO user with id 1,
	// so success here proves VerifyToken served from cache without hitting the DB.
	_ = tokenRepo.Set(context.Background(), fmt.Sprintf(userCacheKey, 1), `{"id":1,"username":"alice"}`, time.Minute)

	service := NewAuthServiceServer(logger, NewMockAuthRepository(), tokenRepo, NewMockSMSCodeRepository(), jwtAuth, &miniprogram.MiniProgram{}, RefreshTTL(time.Hour)).(*AuthService)

	resp, err := service.VerifyToken(context.Background(), &pb.VerifyTokenRequest{AccessToken: token})
	if err != nil {
		t.Fatalf("expected cached success, got %v", err)
	}
	if resp.User == nil || resp.User.Id != 1 || resp.User.Username != "alice" {
		t.Fatalf("expected cached user {1 alice}, got %+v", resp.User)
	}
}

func TestAuthService_SendSmsCode(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, NewMockSMSCodeRepository(), jwtAuth, mp, RefreshTTL(time.Hour)).(*AuthService)

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
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Pre-store SMS code
	smsCodeRepo := NewMockSMSCodeRepository()
	smsCodeRepo.Set(context.Background(), "13800138000", "123456", 5*time.Minute)

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, smsCodeRepo, jwtAuth, mp, RefreshTTL(time.Hour)).(*AuthService)

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
				if resp.User == nil || resp.User.Id <= 0 || resp.User.Username == "" {
					t.Fatalf("expected user with id and username, got %+v", resp.User)
				}
				user := authRepo.users[resp.User.Id]
				if user == nil {
					t.Fatalf("expected stored user %d", resp.User.Id)
				}
				if user.Nickname != user.Username {
					t.Fatalf("expected nickname %q, got %q", user.Username, user.Nickname)
				}
			}
		})
	}
}

func TestAuthService_PhoneAuth_ExistingUser(t *testing.T) {
	logger := zap.NewNop()
	authRepo := NewMockAuthRepository()
	tokenRepo := NewMockTokenRepository()
	jwtAuth := auth.NewJWTAuth(&auth.Options{Secret: "test-secret", ExpiresIn: 3600})
	mp := &miniprogram.MiniProgram{}

	// Create existing user
	existingUser := &entity.User{
		Username: "existing_user",
		Phone:    "13800138001",
		Status:   entity.UserStatusActive,
	}
	existingUser.ID = 1
	authRepo.users[1] = existingUser

	// Pre-store SMS code
	smsCodeRepo := NewMockSMSCodeRepository()
	smsCodeRepo.Set(context.Background(), "13800138001", "654321", 5*time.Minute)

	service := NewAuthServiceServer(logger, authRepo, tokenRepo, smsCodeRepo, jwtAuth, mp, RefreshTTL(time.Hour)).(*AuthService)

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
	if resp.User.Id != 1 {
		t.Fatalf("expected existing user id 1, got %d", resp.User.Id)
	}
}

func TestAuthService_StoreTokenRejectsExpiredToken(t *testing.T) {
	service := &AuthService{
		logger:    zap.NewNop(),
		tokenRepo: NewMockTokenRepository(),
	}

	err := service.storeToken(context.Background(), "expired-token", time.Now().Add(-time.Second))
	if err == nil {
		t.Fatal("expected expired token error")
	}
}

func TestRandUserName(t *testing.T) {
	username, err := randUserName()
	if err != nil {
		t.Fatalf("randUserName() error = %v", err)
	}
	if !strings.HasPrefix(username, "user_") {
		t.Fatalf("expected user_ prefix, got %q", username)
	}
	if len(username) != len("user_")+16 {
		t.Fatalf("expected 16 hex chars suffix, got %q", username)
	}
	if _, err := hex.DecodeString(strings.TrimPrefix(username, "user_")); err != nil {
		t.Fatalf("expected hex suffix, got %q: %v", username, err)
	}
}
