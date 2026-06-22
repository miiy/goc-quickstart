package auth

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/sync/singleflight"

	webclient "github.com/miiy/goc-quickstart/nova-web/client"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/logger/zap"
)

// refreshSkew is how long before access-token expiry a proactive refresh is triggered.
const refreshSkew = 5 * time.Minute

// RefreshClient is the subset of client.AuthClient used for proactive refresh.
type RefreshClient interface {
	Refresh(ctx context.Context, refreshToken string) (*webclient.LoginResponse, error)
}

// SessionAuth establishes a fresh, valid login identity for each request. It first
// proactively renews an expiring/expired access token via the refresh token, then
// bridges the session-backed identity into the gin/template and request contexts.
//
// It must be registered after the generic session middleware (Manager.Middleware),
// which is what makes sessions.Default(c) usable. Renewal happens before injection
// so downstream handlers always see the freshly rotated access token.
func SessionAuth(m *websession.Manager, client RefreshClient, log logger.Logger) gin.HandlerFunc {
	// One group per middleware instance collapses concurrent refreshes.
	var refreshGroup singleflight.Group
	return func(c *gin.Context) {
		maybeRefresh(c, m, client, log, &refreshGroup)
		injectIdentity(c, m)
		c.Next()
	}
}

// maybeRefresh renews the access token when it is expiring soon or already expired.
//
//   - Transient failures (network, 5xx, …) leave the session untouched; the request
//     proceeds with the stale token and the next request retries.
//   - Terminal failures (refresh token invalid/expired/revoked/reused → HTTP 401)
//     clear the session so the user is logged out cleanly, instead of looping on a
//     dead refresh token and repeatedly tripping the auth service's reuse detection.
//
// singleflight collapses concurrent refreshes that share the same refresh token.
// With the auth service's rotation + reuse detection, two simultaneous Refresh RPCs
// race the server-side rotation CAS and the loser revokes the whole token family;
// collapsing them to a single RPC avoids that for in-process concurrency.
// (Cross-process concurrency still needs a distributed lock; out of scope here.)
func maybeRefresh(c *gin.Context, m *websession.Manager, client RefreshClient, log logger.Logger, group *singleflight.Group) {
	if client == nil {
		return
	}
	refreshToken := m.RefreshToken(c)
	if refreshToken == "" || !shouldRefresh(m.AccessExpiresAt(c)) {
		return
	}

	v, err, _ := group.Do(refreshToken, func() (any, error) {
		return client.Refresh(c.Request.Context(), refreshToken)
	})
	if err != nil {
		if webclient.IsStatus(err, http.StatusUnauthorized) {
			log.Warn("session: refresh token invalid or revoked; clearing session", zap.Error(err))
			m.Clear(c)
		} else {
			log.Debug("session: transient refresh failure; keeping session", zap.Error(err))
		}
		return
	}

	resp, _ := v.(*webclient.LoginResponse)
	if resp == nil {
		return
	}
	m.SaveTokens(c, resp.AccessToken, resp.ExpiresAt, resp.RefreshToken)
}

// injectIdentity bridges the session into the gin context (for templates) and the
// request context (for handlers and downstream gRPC clients). The access token is
// only injected for an authenticated session.
func injectIdentity(c *gin.Context, m *websession.Manager) {
	user, ok := m.AuthUser(c)
	c.Set("isLoggedIn", ok)
	if !ok {
		return
	}
	c.Set("currentUser", user)
	ctx := gocauth.InjectAuthenticatedUser(c.Request.Context(), user)
	if token := m.AccessToken(c); token != "" {
		ctx = webclient.WithAccessToken(ctx, token)
	}
	c.Request = c.Request.WithContext(ctx)
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
