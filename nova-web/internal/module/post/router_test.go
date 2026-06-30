package post

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miiy/goc/gin"
)

func TestRouterAppliesAuthToPrivateRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	public := r.Group("")
	protected := r.Group("")
	protected.Use(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
	})

	m := &Module{handler: &PostsHandler{}}
	m.RegisterRouter(public, protected)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/posts/create", nil))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
