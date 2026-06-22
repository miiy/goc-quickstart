package session

import (
	"net/http"

	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	gocauthmid "github.com/miiy/goc/gin/middleware/auth"
	"github.com/miiy/goc/gin/sessions"
)

// Session value keys. Unexported on purpose: Manager is the single owner of the
// session schema, so callers go through its typed methods instead of touching
// raw keys (which used to leak into the auth middleware).
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

// AuthUser returns the authenticated user stored in the session, if any.
func (m *Manager) AuthUser(c *gin.Context) (*gocauth.AuthenticatedUser, bool) {
	return gocauthmid.SessionUser(sessions.Default(c).Get(gocauthmid.SessionKeyAuthUser))
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

// SaveTokens writes a refreshed token pair back to the session. Empty fields are
// ignored so a response without a rotated refresh token keeps the existing one.
func (m *Manager) SaveTokens(c *gin.Context, accessToken, expiresAt, refreshToken string) {
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

func (m *Manager) SaveLogin(c *gin.Context, user map[string]any, accessToken, expiresAt, refreshToken string) error {
	m.Clear(c)

	session, err := m.store.New(requestWithoutCookie(c.Request, m.name), m.name)
	if err != nil {
		return err
	}
	session.Values[gocauthmid.SessionKeyAuthUser] = user
	session.Values[keyAccessToken] = accessToken
	if expiresAt != "" {
		session.Values[keyAccessExpiresAt] = expiresAt
	}
	if refreshToken != "" {
		session.Values[keyRefreshToken] = refreshToken
	}
	return m.store.Save(c.Request, c.Writer, session)
}

// Clear wipes the session values and expires its cookie. It operates on the
// request's context session (sessions.Default) so that reads later in the same
// request — e.g. the auth middleware injecting identity right after a failed
// refresh — observe the cleared state instead of the stale pre-clear snapshot.
func (m *Manager) Clear(c *gin.Context) {
	if m == nil {
		return
	}
	session := sessions.Default(c)
	session.Clear()
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
