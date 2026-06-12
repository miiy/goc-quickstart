package router

import (
	"net/http"
	"strings"

	"github.com/miiy/goc-quickstart/nova-gateway/internal/app"
	authmiddleware "github.com/miiy/goc-quickstart/nova-gateway/internal/middleware/auth"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/service/auth"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/service/file"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/service/health"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/service/post"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/service/user"
	"github.com/miiy/goc/gin"
)

// Router returns a function that registers all routes and middleware onto the gin engine.
func Router(app *app.App) func(r *gin.Engine) {
	return func(r *gin.Engine) {
		if root := strings.TrimSpace(app.Config().Files.Root); root != "" {
			r.StaticFS("/uploads", http.Dir(root))
		}

		health.NewModule(r).RegisterRouter()

		authModule := auth.NewModule(app.AuthClient())
		authModule.RegisterPublicRouter(r)

		public := r.Group("/api/v1")
		postModule := post.NewModule(app.PostClient(), app.UserClient())
		postModule.RegisterPublicRouter(public)

		protected := r.Group("/api/v1")
		protected.Use(authmiddleware.JWTAuthenticationMiddleware(app.JWTAuth(), app.AuthClient()))

		authModule.RegisterProtectedRouter(protected)
		postModule.RegisterProtectedRouter(protected)
		file.NewModule(app.FileClient(), app.UserClient()).RegisterRouter(protected)
		user.NewModule(app.UserClient()).RegisterRouter(protected)
	}
}
