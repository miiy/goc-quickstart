package auth

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/sync/singleflight"

	webclient "github.com/miiy/goc-quickstart/nova-web/client"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/sessionauth"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/logger/zap"
)

// refreshSkew is how long before access-token expiry a proactive refresh is triggered.
const refreshSkew = 5 * time.Minute

// RefreshClient is the subset of client.AuthClient used for proactive refresh.
type RefreshClient interface {
	Refresh(ctx context.Context, refreshToken string) (*webclient.LoginResponse, error)
}

// RefreshSessionToken proactively renews an expiring/expired access token via
// the refresh token. It keeps user attachment and route protection out of this
// middleware so those responsibilities stay in goc/sessionauth.
func RefreshSessionToken(m *websession.Manager, client RefreshClient, log logger.Logger) gin.HandlerFunc {
	var refreshGroup singleflight.Group
	return func(c *gin.Context) {
		maybeRefresh(c, m, client, log, &refreshGroup)
		injectAccessToken(c, m)
		c.Next()
	}
}

// maybeRefresh renews the access token when it is expiring soon or already expired.
//
//   - Transient failures (network, 5xx, ...) leave the session untouched; the
//     request proceeds with the stale token and the next request retries.
//   - Terminal failures (refresh token invalid/expired/revoked/reused -> HTTP 401)
//     clear the session so the user is logged out cleanly.
//
// singleflight collapses concurrent refreshes that share the same refresh token.
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
	m.SaveRefreshedTokens(c, resp.AccessToken, resp.ExpiresAt, resp.RefreshToken)
}

func injectAccessToken(c *gin.Context, m *websession.Manager) {
	if !hasSessionUser(c) {
		return
	}
	if token := m.AccessToken(c); token != "" {
		c.Request = c.Request.WithContext(webclient.WithAccessToken(c.Request.Context(), token))
	}
}

func hasSessionUser(c *gin.Context) bool {
	values, ok := sessions.Default(c).Get(sessionauth.SessionKeyUser).(map[string]any)
	if !ok {
		return false
	}
	id, _ := values["id"].(string)
	username, _ := values["username"].(string)
	return id != "" && username != ""
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
