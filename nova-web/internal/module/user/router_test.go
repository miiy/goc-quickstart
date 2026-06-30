package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miiy/goc/gin"
)

func TestRouterAppliesAuthToUserRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	protected := r.Group("")
	protected.Use(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
	})

	m := &Module{handler: &UserHandler{}}
	m.RegisterRouter(protected)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/user/profile", nil))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
