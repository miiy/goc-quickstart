package auth

import (
	"context"
	"encoding/gob"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	webclient "github.com/miiy/goc-quickstart/nova-web/client"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"github.com/miiy/goc/gin/middleware/sessionauth"
	gocsessions "github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/logger/zap"
)

// The cookie store used in these tests encodes session values with gob, which
// requires concrete interface{} types to be registered. The real app uses the
// Redis store with the JSON serializer, so this only affects tests.
func init() {
	gob.Register(map[string]any{})
}

// fakeRefreshClient is a TokenRefresher stub that counts calls and can simulate
// success, a terminal refresh failure, or a transient 5xx / network error.
type fakeRefreshClient struct {
	calls int32
	resp  *RefreshedTokens
	err   error
	delay time.Duration // make concurrent calls overlap to exercise singleflight
}

func (f *fakeRefreshClient) RefreshToken(ctx context.Context, refreshToken string) (*RefreshedTokens, error) {
	atomic.AddInt32(&f.calls, 1)
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
		}
	}
	return f.resp, f.err
}

type captured struct {
	hasUser bool
	token   string
}

type nopLogger struct{}

func (nopLogger) Debug(string, ...logger.Field) {}
func (nopLogger) Info(string, ...logger.Field)  {}
func (nopLogger) Warn(string, ...logger.Field)  {}
func (nopLogger) Error(string, ...logger.Field) {}
func (nopLogger) DPanic(string, ...logger.Field) {
}
func (nopLogger) Panic(string, ...logger.Field) {}
func (nopLogger) Fatal(string, ...logger.Field) {}
func (nopLogger) ZapLogger() *zap.Logger        { return zap.NewNop() }

// newEngine builds a router with the session + refresh + session user attachment
// middleware plus a priming /login route and a /check route that records what the
// middleware chain injected.
func newEngine(t *testing.T, client TokenRefresher) (*gin.Engine, *captured) {
	t.Helper()
	store := gocsessions.NewCookieStore("test-secret")
	mgr := websession.NewManager(store, "sess")
	cap := &captured{}

	r := gin.New()
	r.Use(mgr.Middleware())
	r.Use(RefreshSessionToken(mgr, client, nopLogger{}))
	r.Use(sessionauth.LoadSessionUser())

	// Prime the session. ?skew=<seconds> sets access-token expiry skew-seconds in
	// the future; ?refresh=<token> stores a refresh token (omit to store none).
	r.POST("/login", func(c *gin.Context) {
		skew, _ := time.ParseDuration(c.Query("skew"))
		if skew == 0 {
			skew = 30 * time.Second // expiring soon by default (within refreshSkew)
		}
		exp := time.Now().Add(skew).UTC().Format(time.RFC3339)
		refresh := c.Query("refresh") // empty => no refresh token
		if err := mgr.SaveLoginSession(c, map[string]any{"id": "7", "username": "alice"}, "access-old", exp, refresh); err != nil {
			t.Fatalf("SaveLoginSession: %v", err)
		}
		c.Status(http.StatusOK)
	})

	r.GET("/check", func(c *gin.Context) {
		_, cap.hasUser = authctx.CurrentUser(c)
		cap.token, _ = webclient.AccessTokenFromContext(c.Request.Context())
		c.Status(http.StatusOK)
	})

	return r, cap
}

// sessionCookieHeader extracts the session cookie from a response, keeping the
// last value for the session name. SaveLoginSession writes two Set-Cookie headers
// for the same name; a browser applies the last one, so the test must too.
func sessionCookieHeader(rec *httptest.ResponseRecorder, name string) string {
	var value string
	for _, c := range rec.Result().Cookies() {
		if c.Name == name {
			value = c.Value
		}
	}
	if value == "" {
		return ""
	}
	return name + "=" + value
}

// primeAndCheck performs the cookie round-trip: POST /login, then GET /check
// with that cookie so RefreshSessionToken runs on a real session.
func primeAndCheck(t *testing.T, r *gin.Engine, skew string, refresh string) {
	t.Helper()

	wLogin := httptest.NewRecorder()
	reqLogin := httptest.NewRequest(http.MethodPost, "/login?skew="+skew+"&refresh="+refresh, nil)
	r.ServeHTTP(wLogin, reqLogin)
	if wLogin.Code != http.StatusOK {
		t.Fatalf("prime /login: got %d", wLogin.Code)
	}

	wCheck := httptest.NewRecorder()
	reqCheck := httptest.NewRequest(http.MethodGet, "/check", nil)
	if h := sessionCookieHeader(wLogin, "sess"); h != "" {
		reqCheck.Header.Set("Cookie", h)
	}
	r.ServeHTTP(wCheck, reqCheck)
	if wCheck.Code != http.StatusOK {
		t.Fatalf("GET /check: got %d", wCheck.Code)
	}
}

func TestRefreshSessionToken_NoRefreshToken_NoRefresh(t *testing.T) {
	client := &fakeRefreshClient{resp: &RefreshedTokens{AccessToken: "access-new"}}
	r, cap := newEngine(t, client)

	primeAndCheck(t, r, "1h", "") // valid token, no refresh token stored

	if client.calls != 0 {
		t.Fatalf("expected no refresh call, got %d", client.calls)
	}
	if !cap.hasUser {
		t.Fatal("expected session user")
	}
	if cap.token != "access-old" {
		t.Fatalf("expected access-old, got %q", cap.token)
	}
}

func TestRefreshSessionToken_NoSessionUser_NoAccessTokenInjection(t *testing.T) {
	store := gocsessions.NewCookieStore("test-secret")
	mgr := websession.NewManager(store, "sess")
	client := &fakeRefreshClient{resp: &RefreshedTokens{AccessToken: "access-new"}}

	r := gin.New()
	r.Use(mgr.Middleware())
	r.Use(RefreshSessionToken(mgr, client, nopLogger{}))
	r.Use(sessionauth.LoadSessionUser())

	r.POST("/prime", func(c *gin.Context) {
		mgr.SaveRefreshedTokens(c, "access-orphan", time.Now().Add(time.Hour).UTC().Format(time.RFC3339), "refresh-1")
		c.Status(http.StatusOK)
	})
	r.GET("/check", func(c *gin.Context) {
		if _, ok := webclient.AccessTokenFromContext(c.Request.Context()); ok {
			t.Fatal("expected no access token without session user")
		}
		if _, ok := authctx.CurrentUser(c); ok {
			t.Fatal("expected no session user")
		}
		c.Status(http.StatusOK)
	})

	wPrime := httptest.NewRecorder()
	r.ServeHTTP(wPrime, httptest.NewRequest(http.MethodPost, "/prime", nil))
	if wPrime.Code != http.StatusOK {
		t.Fatalf("prime: got %d", wPrime.Code)
	}

	wCheck := httptest.NewRecorder()
	reqCheck := httptest.NewRequest(http.MethodGet, "/check", nil)
	if h := sessionCookieHeader(wPrime, "sess"); h != "" {
		reqCheck.Header.Set("Cookie", h)
	}
	r.ServeHTTP(wCheck, reqCheck)
	if wCheck.Code != http.StatusOK {
		t.Fatalf("GET /check: got %d", wCheck.Code)
	}
}

func TestRefreshSessionToken_TokenStillValid_NoRefresh(t *testing.T) {
	client := &fakeRefreshClient{resp: &RefreshedTokens{AccessToken: "access-new"}}
	r, cap := newEngine(t, client)

	primeAndCheck(t, r, "1h", "refresh-1") // expires in 1h, outside refreshSkew

	if client.calls != 0 {
		t.Fatalf("expected no refresh call, got %d", client.calls)
	}
	if cap.token != "access-old" {
		t.Fatalf("expected original access-old, got %q", cap.token)
	}
}

func TestRefreshSessionToken_ExpiringSoon_RefreshSucceeds(t *testing.T) {
	client := &fakeRefreshClient{resp: &RefreshedTokens{
		AccessToken:  "access-new",
		ExpiresAt:    time.Now().Add(time.Hour).UTC(),
		RefreshToken: "refresh-2",
	}}
	r, cap := newEngine(t, client)

	primeAndCheck(t, r, "30s", "refresh-1") // expiring within refreshSkew

	if client.calls != 1 {
		t.Fatalf("expected 1 refresh call, got %d", client.calls)
	}
	if cap.token != "access-new" {
		t.Fatalf("expected refreshed access-new, got %q", cap.token)
	}
	if !cap.hasUser {
		t.Fatal("expected session user after refresh")
	}
}

func TestRefreshSessionToken_RefreshUnauthorized_ClearsSession(t *testing.T) {
	client := &fakeRefreshClient{err: ErrInvalidRefreshToken}
	r, cap := newEngine(t, client)

	primeAndCheck(t, r, "30s", "refresh-1")

	if client.calls != 1 {
		t.Fatalf("expected 1 refresh call, got %d", client.calls)
	}
	if cap.hasUser {
		t.Fatal("expected no session user after terminal refresh failure")
	}
	if cap.token != "" {
		t.Fatalf("expected no token injected after clear, got %q", cap.token)
	}
}

func TestRefreshSessionToken_RefreshTransient_KeepsSession(t *testing.T) {
	client := &fakeRefreshClient{err: webclient.NewHTTPError(http.StatusBadGateway, "upstream")}
	r, cap := newEngine(t, client)

	primeAndCheck(t, r, "30s", "refresh-1")

	if client.calls != 1 {
		t.Fatalf("expected 1 refresh call, got %d", client.calls)
	}
	if cap.token != "access-old" {
		t.Fatalf("expected stale access-old kept, got %q", cap.token)
	}
	if !cap.hasUser {
		t.Fatal("expected session user kept on transient failure")
	}
}

func TestRefreshSessionToken_ConcurrentRefresh_Singleflight(t *testing.T) {
	client := &fakeRefreshClient{
		resp: &RefreshedTokens{
			AccessToken:  "access-new",
			ExpiresAt:    time.Now().Add(time.Hour).UTC(),
			RefreshToken: "refresh-2",
		},
		delay: 50 * time.Millisecond, // force overlap so calls would race without dedup
	}
	r, _ := newEngine(t, client)

	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, httptest.NewRequest(http.MethodPost, "/login?skew=30s&refresh=refresh-1", nil))
	cookieHeader := sessionCookieHeader(wLogin, "sess")

	const n = 8
	done := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/check", nil)
			req.Header.Set("Cookie", cookieHeader)
			r.ServeHTTP(w, req)
		}()
	}
	for i := 0; i < n; i++ {
		<-done
	}

	if got := atomic.LoadInt32(&client.calls); got != 1 {
		t.Fatalf("expected singleflight to collapse %d concurrent refreshes to 1 call, got %d", n, got)
	}
}
