package session

import (
	"net/http"

	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/sessionauth"
	"github.com/miiy/goc/gin/sessions"
)

// Token value keys. Unexported so callers use Manager's typed accessors instead
// of touching raw session keys.
const (
	keyAccessToken     = "access_token"
	keyAccessExpiresAt = "access_expires_at"
	keyRefreshToken    = "refresh_token"
)

type Manager struct {
	store sessions.Store
	name  string
}

func NewManager(store sessions.Store, name string) *Manager {
	return &Manager{
		store: store,
		name:  name,
	}
}

func (m *Manager) Middleware() gin.HandlerFunc {
	return sessions.Middleware(m.name, m.store)
}

// AccessToken returns the stored access token, or "" when none is set.
func (m *Manager) AccessToken(c *gin.Context) string {
	v, _ := sessions.Default(c).Get(keyAccessToken).(string)
	return v
}

// RefreshToken returns the stored refresh token, or "" when none is set.
func (m *Manager) RefreshToken(c *gin.Context) string {
	v, _ := sessions.Default(c).Get(keyRefreshToken).(string)
	return v
}

// AccessExpiresAt returns the stored access-token expiry (RFC3339), or "".
func (m *Manager) AccessExpiresAt(c *gin.Context) string {
	v, _ := sessions.Default(c).Get(keyAccessExpiresAt).(string)
	return v
}

// Tokens returns the access and refresh tokens. Kept for handlers that need both
// (e.g. Logout), which would otherwise call the two accessors separately.
func (m *Manager) Tokens(c *gin.Context) (accessToken, refreshToken string) {
	return m.AccessToken(c), m.RefreshToken(c)
}

// SaveRefreshedTokens stores tokens returned by a refresh response. Empty fields
// keep the existing session values.
func (m *Manager) SaveRefreshedTokens(c *gin.Context, accessToken, expiresAt, refreshToken string) {
	session := sessions.Default(c)
	session.Set(keyAccessToken, accessToken)
	if expiresAt != "" {
		session.Set(keyAccessExpiresAt, expiresAt)
	}
	if refreshToken != "" {
		session.Set(keyRefreshToken, refreshToken)
	}
	_ = session.Save()
}

// SaveLoginSession expires the old cookie and saves a fresh login session.
// The response first deletes the old cookie, then sends a new session cookie.
// Set-Cookie order:
//
//	Set-Cookie: session_name=...; Max-Age=0
//	Set-Cookie: session_name=new_value; Path=/; Max-Age=...
func (m *Manager) SaveLoginSession(c *gin.Context, user map[string]any, accessToken, expiresAt, refreshToken string) error {
	m.Clear(c)

	// Create a new session
	newSession, err := m.store.New(requestWithoutCookie(c.Request, m.name), m.name)
	if err != nil {
		return err
	}
	newSession.Values[sessionauth.SessionKeyUser] = user
	newSession.Values[keyAccessToken] = accessToken
	if expiresAt != "" {
		newSession.Values[keyAccessExpiresAt] = expiresAt
	}
	if refreshToken != "" {
		newSession.Values[keyRefreshToken] = refreshToken
	}
	return m.store.Save(c.Request, c.Writer, newSession)
}

// Clear wipes session values and expires the browser cookie.
func (m *Manager) Clear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	// MaxAge < 0 tells the browser to delete the cookie.
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	_ = session.Save()
}

func requestWithoutCookie(r *http.Request, name string) *http.Request {
	req := r.Clone(r.Context())
	req.Header.Del("Cookie")
	for _, cookie := range r.Cookies() {
		if cookie.Name != name {
			req.AddCookie(cookie)
		}
	}
	return req
}
