package post

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miiy/goc/gin"
)

func TestRouterRegistersPublicReadAndProtectedEditorRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	public := r.Group("")
	protected := r.Group("")
	protected.Use(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
	})

	m := &Module{handler: &PostsHandler{}}
	m.RegisterRouter(public, protected)

	routes := registeredRoutes(r)
	for _, route := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/posts"},
		{http.MethodGet, "/posts/pages/2"},
		{http.MethodGet, "/posts/example-id"},
	} {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			if !routes[route.method+" "+route.path] {
				t.Fatalf("route %s %s was not registered", route.method, route.path)
			}
		})
	}

	for _, path := range []string{"/posts/create", "/posts/example-id/edit"} {
		t.Run("protected "+path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("GET %s status = %d, want %d", path, w.Code, http.StatusUnauthorized)
			}
		})
	}

	for _, route := range []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/posts"},
		{http.MethodPost, "/posts/example-id"},
		{http.MethodPut, "/posts/example-id"},
		{http.MethodDelete, "/posts/example-id"},
	} {
		t.Run("unregistered "+route.method+" "+route.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(route.method, route.path, nil))

			if w.Code == http.StatusUnauthorized {
				t.Fatalf("route %s %s should not be registered", route.method, route.path)
			}
		})
	}
}

func registeredRoutes(r *gin.Engine) map[string]bool {
	result := make(map[string]bool)
	for _, route := range r.Routes() {
		for _, path := range expandedRoutePaths(route.Path) {
			result[route.Method+" "+path] = true
		}
	}
	return result
}

func expandedRoutePaths(path string) []string {
	switch path {
	case "/posts/pages/:page":
		return []string{path, "/posts/pages/2"}
	case "/posts/:id":
		return []string{path, "/posts/example-id"}
	case "/posts/:id/edit":
		return []string{path, "/posts/example-id/edit"}
	default:
		return []string{path}
	}
}
