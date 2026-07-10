package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/sessions"
)

func TestRegisterFormRedirectsWhenRegisterDisabled(t *testing.T) {
	r := newDisabledRegisterRouter()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/register", nil))

	assertRedirectsToLogin(t, w)
}

func TestRegisterRedirectsWhenRegisterDisabled(t *testing.T) {
	r := newDisabledRegisterRouter()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/register", nil))

	assertRedirectsToLogin(t, w)
}

func TestRegisterDisabledFlashIsReadableOnLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAuthHandler(nil, nil, nil, false)
	r := gin.New()
	r.Use(sessions.Middleware("sess", sessions.NewCookieStore("test-secret")))
	r.GET("/register", h.RegisterForm)
	r.GET("/login", func(c *gin.Context) {
		flashes, err := sessions.Flashes(c)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if len(flashes) == 0 {
			c.String(http.StatusInternalServerError, "missing flash")
			return
		}
		c.String(http.StatusOK, flashes[0].Message)
	})

	registerResp := httptest.NewRecorder()
	r.ServeHTTP(registerResp, httptest.NewRequest(http.MethodGet, "/register", nil))
	assertRedirectsToLogin(t, registerResp)

	loginReq := httptest.NewRequest(http.MethodGet, "/login", nil)
	for _, cookie := range registerResp.Result().Cookies() {
		loginReq.AddCookie(cookie)
	}
	loginResp := httptest.NewRecorder()
	r.ServeHTTP(loginResp, loginReq)

	if loginResp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, loginResp.Code, loginResp.Body.String())
	}
	if got := loginResp.Body.String(); got != "注册暂未开放" {
		t.Fatalf("expected disabled register flash, got %q", got)
	}
}

func newDisabledRegisterRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	h := NewAuthHandler(nil, nil, nil, false)
	r := gin.New()
	r.Use(sessions.Middleware("sess", sessions.NewCookieStore("test-secret")))
	r.GET("/register", h.RegisterForm)
	r.POST("/register", h.Register)
	return r
}

func assertRedirectsToLogin(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	if w.Code != http.StatusFound {
		t.Fatalf("expected status %d, got %d", http.StatusFound, w.Code)
	}
	if location := w.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}
