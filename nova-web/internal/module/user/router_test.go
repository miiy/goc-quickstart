package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miiy/goc/gin"
)

func TestRouterAppliesAuthToUserPages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	protected := r.Group("")
	protected.Use(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
	})

	m := &Module{handler: &UserHandler{}}
	m.RegisterRouter(r.Group(""), protected)

	for _, path := range []string{
		"/profile",
		"/users/alice",
	} {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("GET %s status = %d, want %d", path, w.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestRouterDoesNotRegisterUserSectionPages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	protected := r.Group("")

	m := &Module{handler: &UserHandler{}}
	m.RegisterRouter(r.Group(""), protected)

	for _, path := range []string{
		"/users/alice/activity",
		"/users/alice/followers",
		"/users/alice/following",
	} {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))

			if w.Code != http.StatusNotFound {
				t.Fatalf("GET %s status = %d, want %d", path, w.Code, http.StatusNotFound)
			}
		})
	}
}
