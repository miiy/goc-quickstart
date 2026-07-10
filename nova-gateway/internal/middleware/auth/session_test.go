package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/csrf"
	"github.com/miiy/goc/gin/sessions"
)

func newSessionProtectedEngine(t *testing.T, client verifyTokenClient) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(sessions.Middleware("sess", sessions.NewCookieStore("test-secret")))

	r.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(sessionAccessTokenKey, "session-token")
		session.Set(csrf.SessionKey, "csrf-token")
		if err := session.Save(); err != nil {
			t.Fatalf("save session: %v", err)
		}
		c.Status(http.StatusNoContent)
	})

	protected := r.Group("/protected")
	protected.Use(SessionTokenToBearer(), SessionCSRF(), AuthenticationMiddleware(client))
	protected.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	protected.POST("", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r
}

func TestSessionTokenToBearerAuthenticatesSessionCookie(t *testing.T) {
	client := &fakeVerifyTokenClient{user: &authv1.User{Id: 42, Username: "alice"}}
	engine := newSessionProtectedEngine(t, client)

	login := httptest.NewRecorder()
	engine.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/login", nil))
	if login.Code != http.StatusNoContent {
		t.Fatalf("login status = %d, want %d", login.Code, http.StatusNoContent)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	for _, cookie := range login.Result().Cookies() {
		req.AddCookie(cookie)
	}
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if client.gotToken != "session-token" {
		t.Fatalf("expected session token forwarded, got %q", client.gotToken)
	}
}

func TestSessionTokenToBearerKeepsExplicitAuthorization(t *testing.T) {
	client := &fakeVerifyTokenClient{user: &authv1.User{Id: 42, Username: "alice"}}
	engine := newSessionProtectedEngine(t, client)

	login := httptest.NewRecorder()
	engine.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/login", nil))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("Authorization", "Bearer header-token")
	for _, cookie := range login.Result().Cookies() {
		req.AddCookie(cookie)
	}
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if client.gotToken != "header-token" {
		t.Fatalf("expected explicit bearer token forwarded, got %q", client.gotToken)
	}
}

func TestSessionCSRFRejectsCookieAuthenticatedWriteWithoutToken(t *testing.T) {
	client := &fakeVerifyTokenClient{user: &authv1.User{Id: 42, Username: "alice"}}
	engine := newSessionProtectedEngine(t, client)

	login := httptest.NewRecorder()
	engine.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/login", nil))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	for _, cookie := range login.Result().Cookies() {
		req.AddCookie(cookie)
	}
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	if client.gotToken != "" {
		t.Fatalf("expected authentication not to run, got token %q", client.gotToken)
	}
}

func TestSessionCSRFRejectsCookieAuthenticatedWriteWithInvalidToken(t *testing.T) {
	client := &fakeVerifyTokenClient{user: &authv1.User{Id: 42, Username: "alice"}}
	engine := newSessionProtectedEngine(t, client)

	login := httptest.NewRecorder()
	engine.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/login", nil))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set(csrf.HeaderName, "wrong-token")
	for _, cookie := range login.Result().Cookies() {
		req.AddCookie(cookie)
	}
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
	if client.gotToken != "" {
		t.Fatalf("expected authentication not to run, got token %q", client.gotToken)
	}
}

func TestSessionCSRFAcceptsCookieAuthenticatedWriteWithToken(t *testing.T) {
	client := &fakeVerifyTokenClient{user: &authv1.User{Id: 42, Username: "alice"}}
	engine := newSessionProtectedEngine(t, client)

	login := httptest.NewRecorder()
	engine.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/login", nil))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set(csrf.HeaderName, "csrf-token")
	for _, cookie := range login.Result().Cookies() {
		req.AddCookie(cookie)
	}
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if client.gotToken != "session-token" {
		t.Fatalf("expected session token forwarded, got %q", client.gotToken)
	}
}

func TestSessionCSRFLeavesUnauthenticatedWriteForAuthentication(t *testing.T) {
	client := &fakeVerifyTokenClient{user: &authv1.User{Id: 42, Username: "alice"}}
	engine := newSessionProtectedEngine(t, client)

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/protected", nil))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
