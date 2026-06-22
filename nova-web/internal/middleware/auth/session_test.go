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
	gocsessions "github.com/miiy/goc/gin/sessions"
)

// The cookie store used in these tests encodes session values with gob, which
// requires concrete interface{} types to be registered. The real app uses the
// Redis store with the JSON serializer, so this only affects tests.
func init() {
	gob.Register(map[string]any{})
}

// fakeRefreshClient is a RefreshClient stub that counts calls and can simulate
// success, a terminal 401, or a transient 5xx / network error.
type fakeRefreshClient struct {
	calls int32
	resp  *webclient.LoginResponse
	err   error
	delay time.Duration // make concurrent calls overlap to exercise singleflight
}

func (f *fakeRefreshClient) Refresh(ctx context.Context, refreshToken string) (*webclient.LoginResponse, error) {
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
	loggedIn bool
	token    string
}

// newEngine builds a router with the session + SessionAuth middleware plus a
// priming /login route (query params control the access-token expiry skew and
// whether a refresh token is stored) and a /check route that records what the
// middleware injected.
func newEngine(t *testing.T, client RefreshClient) (*gin.Engine, *captured) {
	t.Helper()
	store := gocsessions.NewCookieStore("test-secret")
	mgr := websession.NewManager(store, "sess")
	cap := &captured{}

	r := gin.New()
	r.Use(mgr.Middleware())
	r.Use(SessionAuth(mgr, client, nil))

	// Prime the session. ?skew=<seconds> sets access-token expiry skew-seconds in
	// the future; ?refresh=<token> stores a refresh token (omit to store none).
	r.POST("/login", func(c *gin.Context) {
		skew, _ := time.ParseDuration(c.Query("skew"))
		if skew == 0 {
			skew = 30 * time.Second // expiring soon by default (within refreshSkew)
		}
		exp := time.Now().Add(skew).UTC().Format(time.RFC3339)
		refresh := c.Query("refresh") // empty => no refresh token
		if err := mgr.SaveLogin(c, map[string]any{"id": int64(7), "username": "alice"}, "access-old", exp, refresh); err != nil {
			t.Fatalf("SaveLogin: %v", err)
		}
		c.Status(http.StatusOK)
	})

	r.GET("/check", func(c *gin.Context) {
		cap.loggedIn = c.GetBool("isLoggedIn")
		cap.token, _ = webclient.AccessTokenFromContext(c.Request.Context())
		c.Status(http.StatusOK)
	})

	return r, cap
}

// sessionCookieHeader extracts the session cookie from a response, keeping the
// last value for the session name. SaveLogin writes two Set-Cookie headers for
// the same name (an expiry one from Clear, then the fresh session); a browser
// applies the last one, so the test must too.
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

// primeAndCheck performs the cookie round-trip: POST /login (storing the session
// cookie), then GET /check with that cookie so SessionAuth runs on a real session.
func primeAndCheck(t *testing.T, r *gin.Engine, skew string, refresh string) *captured {
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
	return nil
}

func TestSessionAuth_NoRefreshToken_NoRefresh(t *testing.T) {
	client := &fakeRefreshClient{resp: &webclient.LoginResponse{AccessToken: "access-new"}}
	r, cap := newEngine(t, client)

	primeAndCheck(t, r, "1h", "") // valid token, no refresh token stored

	if client.calls != 0 {
		t.Fatalf("expected no refresh call, got %d", client.calls)
	}
	if !cap.loggedIn {
		t.Fatal("expected isLoggedIn=true")
	}
	if cap.token != "access-old" {
		t.Fatalf("expected access-old, got %q", cap.token)
	}
}

func TestSessionAuth_TokenStillValid_NoRefresh(t *testing.T) {
	client := &fakeRefreshClient{resp: &webclient.LoginResponse{AccessToken: "access-new"}}
	r, cap := newEngine(t, client)

	primeAndCheck(t, r, "1h", "refresh-1") // expires in 1h, outside refreshSkew

	if client.calls != 0 {
		t.Fatalf("expected no refresh call, got %d", client.calls)
	}
	if cap.token != "access-old" {
		t.Fatalf("expected original access-old, got %q", cap.token)
	}
}

func TestSessionAuth_ExpiringSoon_RefreshSucceeds(t *testing.T) {
	client := &fakeRefreshClient{resp: &webclient.LoginResponse{
		AccessToken:  "access-new",
		ExpiresAt:    time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
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
	if !cap.loggedIn {
		t.Fatal("expected isLoggedIn=true after refresh")
	}
}

func TestSessionAuth_RefreshUnauthorized_ClearsSession(t *testing.T) {
	client := &fakeRefreshClient{err: &webclient.HTTPError{StatusCode: http.StatusUnauthorized, Message: "reuse detected"}}
	r, cap := newEngine(t, client)

	primeAndCheck(t, r, "30s", "refresh-1")

	if client.calls != 1 {
		t.Fatalf("expected 1 refresh call, got %d", client.calls)
	}
	if cap.loggedIn {
		t.Fatal("expected isLoggedIn=false after terminal refresh failure")
	}
	if cap.token != "" {
		t.Fatalf("expected no token injected after clear, got %q", cap.token)
	}
}

func TestSessionAuth_RefreshTransient_KeepsSession(t *testing.T) {
	client := &fakeRefreshClient{err: &webclient.HTTPError{StatusCode: http.StatusBadGateway, Message: "upstream"}}
	r, cap := newEngine(t, client)

	primeAndCheck(t, r, "30s", "refresh-1")

	if client.calls != 1 {
		t.Fatalf("expected 1 refresh call, got %d", client.calls)
	}
	// Stale access token kept; request proceeds with it.
	if cap.token != "access-old" {
		t.Fatalf("expected stale access-old kept, got %q", cap.token)
	}
	if !cap.loggedIn {
		t.Fatal("expected isLoggedIn=true kept on transient failure")
	}
}

func TestSessionAuth_ConcurrentRefresh_Singleflight(t *testing.T) {
	client := &fakeRefreshClient{
		resp: &webclient.LoginResponse{
			AccessToken:  "access-new",
			ExpiresAt:    time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			RefreshToken: "refresh-2",
		},
		delay: 50 * time.Millisecond, // force overlap so calls would race without dedup
	}
	r, _ := newEngine(t, client)

	// Prime once to obtain a session cookie holding "refresh-1".
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
