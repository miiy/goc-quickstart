package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miiy/goc/gin"
)

func TestRegisterRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")

	NewModule(&fakeUsersUserClient{}).RegisterRouter(api, api)
}

func TestRegisterRouterProtectsUsernameRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	public := api.Group("")
	protected := api.Group("")
	protected.Use(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
	})

	NewModule(&fakeUsersUserClient{}).RegisterRouter(public, protected)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/users/alice", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
