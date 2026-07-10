package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/singleflight"

	authv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/auth/v1"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc-quickstart/nova-web/internal/transport"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/logger/zap"
)

// refreshSkew is how long before access-token expiry a proactive refresh is triggered.
const refreshSkew = 5 * time.Minute

// ErrInvalidRefreshToken marks a terminal refresh failure.
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

// RefreshedTokens is the normalized token payload stored in the session.
type RefreshedTokens struct {
	AccessToken  string
	ExpiresAt    time.Time
	RefreshToken string
}

// TokenRefresher refreshes session tokens without exposing a transport DTO.
type TokenRefresher interface {
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshedTokens, error)
}

// RefreshSessionToken proactively renews an expiring/expired access token via
// the refresh token. It keeps user attachment and route protection out of this
// middleware so those responsibilities stay in goc/sessionauth.
func RefreshSessionToken(m *websession.Manager, refresher TokenRefresher, log logger.Logger) gin.HandlerFunc {
	var refreshGroup singleflight.Group
	return func(c *gin.Context) {
		maybeRefresh(c, m, refresher, log, &refreshGroup)
		c.Next()
	}
}

// maybeRefresh renews the access token when it is expiring soon or already expired.
//
//   - Transient failures (network, 5xx, ...) leave the session untouched; the
//     request proceeds with the stale token and the next request retries.
//   - Terminal failures (refresh token invalid/expired/revoked/reused) clear the
//     session so the user is logged out cleanly.
//
// singleflight collapses concurrent refreshes that share the same refresh token.
func maybeRefresh(c *gin.Context, m *websession.Manager, refresher TokenRefresher, log logger.Logger, group *singleflight.Group) {
	if refresher == nil {
		return
	}
	refreshToken := m.RefreshToken(c)
	if refreshToken == "" || !shouldRefresh(m.AccessExpiresAt(c)) {
		return
	}

	v, err, _ := group.Do(refreshToken, func() (any, error) {
		return refresher.RefreshToken(c.Request.Context(), refreshToken)
	})
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			log.Warn("session: refresh token invalid or revoked; clearing session", zap.Error(err))
			m.Clear(c)
		} else {
			log.Debug("session: transient refresh failure; keeping session", zap.Error(err))
		}
		return
	}

	tokens, _ := v.(*RefreshedTokens)
	if tokens == nil {
		return
	}
	m.SaveRefreshedTokens(c, tokens.AccessToken, formatAPITime(tokens.ExpiresAt), tokens.RefreshToken)
}

func shouldRefresh(expiresAt string) bool {
	if expiresAt == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return false
	}
	return time.Now().Add(refreshSkew).After(t)
}

func formatAPITime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// NewAuthServiceRefresher adapts the generated auth RPC client to the
// middleware's transport-neutral TokenRefresher contract.
func NewAuthServiceRefresher(authClient authv1.AuthServiceClient) TokenRefresher {
	if authClient == nil {
		return nil
	}
	return &authClientTokenRefresher{authClient: authClient}
}

type authClientTokenRefresher struct {
	authClient authv1.AuthServiceClient
}

func (r *authClientTokenRefresher) RefreshToken(ctx context.Context, refreshToken string) (*RefreshedTokens, error) {
	resp, err := r.authClient.RefreshToken(ctx, &authv1.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		err = transport.FromGRPCError(err)
		if transport.IsStatus(err, http.StatusUnauthorized) {
			return nil, fmt.Errorf("%w: %v", ErrInvalidRefreshToken, err)
		}
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return &RefreshedTokens{
		AccessToken:  resp.GetAccessToken(),
		ExpiresAt:    resp.GetExpiresAt().AsTime(),
		RefreshToken: resp.GetRefreshToken(),
	}, nil
}
