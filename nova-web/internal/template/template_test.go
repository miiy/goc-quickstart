package template

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"github.com/miiy/goc/gin/sessions"
)

func TestNewViewDataIncludesSiteData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetDefaultSite(SiteData{
		Name:            "nova-web",
		URL:             "https://example.test",
		Locale:          "zh-CN",
		RegisterEnabled: true,
	})

	view := newTestViewData(t, nil)
	if view.Site.Name != "nova-web" {
		t.Fatalf("expected site name, got %q", view.Site.Name)
	}
	if view.Site.URL != "https://example.test" {
		t.Fatalf("expected site url, got %q", view.Site.URL)
	}
	if view.Site.Locale != "zh-CN" {
		t.Fatalf("expected site locale, got %q", view.Site.Locale)
	}
	if !view.Site.RegisterEnabled {
		t.Fatal("expected register enabled")
	}
	if view.Auth.IsLoggedIn {
		t.Fatal("expected logged out auth data")
	}
	if view.Auth.CurrentUser != nil {
		t.Fatalf("expected no current user, got %#v", view.Auth.CurrentUser)
	}
	if view.CSRFToken == "" {
		t.Fatal("expected csrf token")
	}
}

func TestNewViewDataIncludesAuthData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetDefaultSite(SiteData{})

	r := gin.New()
	r.Use(sessions.Middleware("sess", sessions.NewCookieStore("test-secret")))
	r.GET("/", func(c *gin.Context) {
		authctx.SetUser(c, &gocauth.AuthenticatedUser{ID: "7", Username: "alice"})

		view := NewViewData(c)
		if !view.Auth.IsLoggedIn {
			t.Fatal("expected logged in auth data")
		}
		if view.Auth.CurrentUser == nil {
			t.Fatal("expected current user")
		}
		if view.Auth.CurrentUser.ID != "7" || view.Auth.CurrentUser.Username != "alice" {
			t.Fatalf("unexpected current user: %#v", view.Auth.CurrentUser)
		}
		if view.CSRFToken == "" {
			t.Fatal("expected csrf token for logged in view")
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET / got %d", w.Code)
	}
}

func newTestViewData(t *testing.T, user *gocauth.AuthenticatedUser) ViewData {
	t.Helper()

	var view ViewData
	r := gin.New()
	r.Use(sessions.Middleware("sess", sessions.NewCookieStore("test-secret")))
	r.GET("/", func(c *gin.Context) {
		if user != nil {
			authctx.SetUser(c, user)
		}
		view = NewViewData(c)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET / got %d", w.Code)
	}
	return view
}

func TestNewFormatTimeFunc(t *testing.T) {
	formatTime := NewFormatTimeFunc("Asia/Shanghai")
	want := "2026-01-01 00:00"

	tests := []struct {
		name string
		in   any
	}{
		{
			name: "rfc3339 string",
			in:   "2025-12-31T16:00:00Z",
		},
		{
			name: "rfc3339 string with nanos",
			in:   "2025-12-31T16:00:00.000000000Z",
		},
		{
			name: "formatted string",
			in:   "2026-01-01 00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTime(tt.in); got != want {
				t.Fatalf("formatTime() = %q, want %q", got, want)
			}
		})
	}
}
