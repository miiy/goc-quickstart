package router

import (
	"github.com/miiy/goc-quickstart/nova-gateway/internal/app"
	authmw "github.com/miiy/goc-quickstart/nova-gateway/internal/middleware/auth"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/health"
	"github.com/miiy/goc/gin"
)

// Router returns a function that registers all routes and middleware onto the gin engine.
func Router(app *app.App) func(r *gin.Engine) {
	return func(r *gin.Engine) {
		health.NewModule(r).RegisterRouter()

		modules := app.Modules()
		clients := app.Clients()
		modules.Auth.RegisterPublicRouter(r)

		public := r.Group("/api/v1")
		modules.Post.RegisterPublicRouter(public)

		protected := r.Group("/api/v1")
		protected.Use(authmw.AuthenticationMiddleware(clients.Auth))

		modules.Auth.RegisterProtectedRouter(protected)
		modules.Post.RegisterProtectedRouter(protected)
		modules.File.RegisterRouter(protected)
		modules.User.RegisterRouter(protected)
	}
}
