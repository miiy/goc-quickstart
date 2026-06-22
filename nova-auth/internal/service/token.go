package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/miiy/goc-quickstart/nova-auth/internal/entity"
	"github.com/miiy/goc/auth"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/logger/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	authTokenKey = "user_token:%s" // user_token:{sha256(access_token)}
	userCacheKey = "auth_user:%d"  // active user cache, keyed by user id
	userCacheTTL = 10 * time.Second

	refreshTokenKey      = "refresh_token:%s"  // refresh_token:{sha256(opaque)}
	refreshFamilyKey     = "refresh_family:%s" // refresh_family:{family}; expires after refresh TTL
	refreshStatusActive  = "active"
	refreshStatusRevoked = "revoked"

	refreshFamilyBytes = 16
	refreshTokenBytes  = 32
)

// tokenPair holds a freshly issued access + refresh token pair.
type tokenPair struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}

// refreshRecord is the value stored under refresh_token:{sha256(opaque)}.
type refreshRecord struct {
	UserID   int64     `json:"user_id"`
	Family   string    `json:"family"`
	Status   string    `json:"status"`
	IssuedAt time.Time `json:"issued_at"`
}

// cachedUser is the cached subset of entity.User for token verification.
// Only active users are cached, so a cache hit implies Status == active.
type cachedUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func tokenDigest(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func formatTokenKey(token string) string {
	return fmt.Sprintf(authTokenKey, tokenDigest(token))
}

func formatRefreshTokenKey(token string) string {
	return fmt.Sprintf(refreshTokenKey, tokenDigest(token))
}

func formatRefreshFamilyKey(family string) string {
	return fmt.Sprintf(refreshFamilyKey, family)
}

// newOpaqueToken returns n random bytes encoded as unpadded base64url
// (n=32 -> 43 chars).
func newOpaqueToken(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// createClaims builds access-token claims for the given user.
func (s *AuthService) createClaims(userID int64, username string) (*auth.UserClaims, error) {
	return s.jwtAuth.CreateClaims(userID, username), nil
}

// storeToken persists an access token in Redis (keyed by sha256(token)) with TTL = its
// remaining lifetime, enabling immediate revocation via Delete (Logout).
func (s *AuthService) storeToken(ctx context.Context, token string, expiresTime time.Time) error {
	ttl := time.Until(expiresTime)
	if ttl <= 0 {
		return errors.New("token already expired")
	}
	if err := s.tokenRepo.Set(ctx, formatTokenKey(token), token, ttl); err != nil {
		s.logger.Error("tokenRepo.Set", zap.Error(err))
		return err
	}
	return nil
}

func (s *AuthService) revokeRefreshFamily(ctx context.Context, family string) {
	ttl := s.refreshTTL
	if ttl < 0 {
		ttl = 0
	}
	if err := s.tokenRepo.Set(ctx, formatRefreshFamilyKey(family), refreshStatusRevoked, ttl); err != nil {
		s.logger.Error("tokenRepo.Set refresh family", zap.Error(err), zap.String("family", family))
	}
}

// issueTokenPair issues a new access token (JWT) and a new refresh token (opaque,
// new family) for the user. Used by Login/MpLogin/PhoneAuth. RefreshToken rotation
// reuses the existing family.
func (s *AuthService) issueTokenPair(ctx context.Context, userID int64, username string) (*tokenPair, error) {
	claims, err := s.createClaims(userID, username)
	if err != nil {
		return nil, err
	}
	access, err := s.jwtAuth.CreateTokenByClaims(claims)
	if err != nil {
		return nil, err
	}
	accessExp := claims.ExpiresAt.Time
	if err := s.storeToken(ctx, access, accessExp); err != nil {
		return nil, err
	}

	family, err := newOpaqueToken(refreshFamilyBytes)
	if err != nil {
		return nil, err
	}
	opaque, err := newOpaqueToken(refreshTokenBytes)
	if err != nil {
		return nil, err
	}
	refreshExp := time.Now().Add(s.refreshTTL)
	rec := refreshRecord{UserID: userID, Family: family, Status: refreshStatusActive, IssuedAt: time.Now()}
	data, err := json.Marshal(rec)
	if err != nil {
		return nil, err
	}
	if err := s.tokenRepo.Set(ctx, formatRefreshTokenKey(opaque), string(data), s.refreshTTL); err != nil {
		return nil, err
	}

	return &tokenPair{
		AccessToken:      access,
		AccessExpiresAt:  accessExp,
		RefreshToken:     opaque,
		RefreshExpiresAt: refreshExp,
	}, nil
}

func (s *AuthService) verifyAccessToken(ctx context.Context, accessToken string) (*entity.User, error) {
	token := strings.TrimSpace(accessToken)
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}

	claims, err := s.jwtAuth.ParseToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	}

	// fail-closed: a missing or unreachable revocation record rejects the token.
	if _, err := s.tokenRepo.Get(ctx, formatTokenKey(token)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "token revoked or expired")
	}

	user, err := s.loadActiveUser(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.Unauthenticated, ErrUserNotFound.Error())
		}
		s.logger.Error("loadActiveUser", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if user.Status != entity.UserStatusActive {
		return nil, status.Error(codes.Unauthenticated, ErrUserNotFound.Error())
	}

	return user, nil
}

// loadActiveUser loads a user by id with a short-lived Redis cache. Only active
// users are cached; disabled or missing users always hit the DB so status changes
// take effect within userCacheTTL. Token revocation is checked separately by
// VerifyToken (tokenRepo.Get on the token key) and is never cached.
func (s *AuthService) loadActiveUser(ctx context.Context, userID int64) (*entity.User, error) {
	key := fmt.Sprintf(userCacheKey, userID)
	if cached, err := s.tokenRepo.Get(ctx, key); err == nil {
		var c cachedUser
		if json.Unmarshal([]byte(cached), &c) == nil {
			u := &entity.User{Username: c.Username, Status: entity.UserStatusActive}
			u.ID = c.ID
			return u, nil
		}
	}
	user, err := s.authRepo.First(ctx, uint64(userID), "id", "username", "status")
	if err != nil {
		return nil, err
	}
	if user.Status == entity.UserStatusActive {
		if data, err := json.Marshal(cachedUser{ID: user.ID, Username: user.Username}); err == nil {
			_ = s.tokenRepo.Set(ctx, key, string(data), userCacheTTL)
		}
	}
	return user, nil
}

// rotateRefreshToken exchanges a refresh token for a new access + refresh token
// pair. Reuse of an already-rotated refresh token revokes the entire family.
func (s *AuthService) rotateRefreshToken(ctx context.Context, refreshToken string) (*tokenPair, *entity.User, error) {
	opaque := strings.TrimSpace(refreshToken)
	if opaque == "" {
		return nil, nil, status.Error(codes.InvalidArgument, ErrInvalidArgument.Error())
	}
	key := formatRefreshTokenKey(opaque)

	raw, err := s.tokenRepo.Get(ctx, key)
	if err != nil {
		return nil, nil, status.Error(codes.Unauthenticated, "refresh token invalid or expired")
	}
	var rec refreshRecord
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		s.logger.Error("refresh record unmarshal", zap.Error(err))
		return nil, nil, status.Error(codes.Internal, "invalid refresh record")
	}

	// family-level revocation (set after a reuse was detected on any sibling token)
	if _, err := s.tokenRepo.Get(ctx, formatRefreshFamilyKey(rec.Family)); err == nil {
		return nil, nil, status.Error(codes.Unauthenticated, "refresh token family revoked")
	}

	// reuse detection: an already-rotated (revoked) token is presented again
	if rec.Status == refreshStatusRevoked {
		s.revokeRefreshFamily(ctx, rec.Family)
		s.logger.Warn("refresh token reuse detected, family revoked",
			zap.Int64("user_id", rec.UserID), zap.String("family", rec.Family))
		return nil, nil, status.Error(codes.Unauthenticated, "refresh token reuse detected")
	}

	user, err := s.loadActiveUser(ctx, rec.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, status.Error(codes.Unauthenticated, ErrUserNotFound.Error())
		}
		s.logger.Error("loadActiveUser", zap.Error(err))
		return nil, nil, status.Error(codes.Internal, err.Error())
	}
	if user.Status != entity.UserStatusActive {
		return nil, nil, status.Error(codes.Unauthenticated, ErrUserNotFound.Error())
	}

	// rotate: atomically mark current refresh revoked (keep key for reuse detection).
	// CAS guarantees only one concurrent rotation of this token succeeds; a loser
	// means another request already rotated it -> treat as reuse -> revoke family.
	rec.Status = refreshStatusRevoked
	revokedData, err := json.Marshal(rec)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}
	swapped, err := s.tokenRepo.CompareAndSet(ctx, key, raw, string(revokedData))
	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}
	if !swapped {
		s.revokeRefreshFamily(ctx, rec.Family)
		s.logger.Warn("refresh token reuse detected (concurrent rotation), family revoked",
			zap.Int64("user_id", rec.UserID), zap.String("family", rec.Family))
		return nil, nil, status.Error(codes.Unauthenticated, "refresh token reuse detected")
	}

	claims, err := s.createClaims(user.ID, user.Username)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}
	access, err := s.jwtAuth.CreateTokenByClaims(claims)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}
	accessExp := claims.ExpiresAt.Time
	if err := s.storeToken(ctx, access, accessExp); err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	newOpaque, err := newOpaqueToken(refreshTokenBytes)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}
	refreshExp := time.Now().Add(s.refreshTTL)
	newRec := refreshRecord{UserID: user.ID, Family: rec.Family, Status: refreshStatusActive, IssuedAt: time.Now()}
	newData, _ := json.Marshal(newRec)
	if err := s.tokenRepo.Set(ctx, formatRefreshTokenKey(newOpaque), string(newData), s.refreshTTL); err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	return &tokenPair{
		AccessToken:      access,
		AccessExpiresAt:  accessExp,
		RefreshToken:     newOpaque,
		RefreshExpiresAt: refreshExp,
	}, user, nil
}

// revokeTokens revokes the access token and the current refresh token
// best-effort. Refresh tokens are marked revoked, not deleted, so reuse
// detection still works for stolen copies.
func (s *AuthService) revokeTokens(ctx context.Context, accessToken, refreshToken string) {
	if accessToken != "" {
		if err := s.tokenRepo.Del(ctx, formatTokenKey(accessToken)); err != nil {
			s.logger.Error("tokenRepo.Del access", zap.Error(err))
		}
	}
	if refreshToken == "" {
		return
	}

	rkey := formatRefreshTokenKey(refreshToken)
	if raw, err := s.tokenRepo.Get(ctx, rkey); err == nil {
		var rec refreshRecord
		if json.Unmarshal([]byte(raw), &rec) == nil {
			rec.Status = refreshStatusRevoked
			if data, err := json.Marshal(rec); err == nil {
				_ = s.tokenRepo.SetKeepTTL(ctx, rkey, string(data))
			}
		}
	}
}
