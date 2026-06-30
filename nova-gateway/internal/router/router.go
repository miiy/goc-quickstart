package router

import (
	"github.com/miiy/goc-quickstart/nova-gateway/internal/app"
	authmw "github.com/miiy/goc-quickstart/nova-gateway/internal/middleware/auth"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/health"
	"github.com/miiy/goc/gin"
	ginzap "github.com/miiy/goc/gin/middleware/zap"
)

// Router returns a function that registers all routes and middleware onto the gin engine.
func Router(app *app.App) func(r *gin.Engine) {
	return func(r *gin.Engine) {
		logger := app.Logger().ZapLogger()
		r.Use(ginzap.Ginzap(logger), ginzap.RecoveryWithZap(logger, true))

		modules := app.Modules()
		clients := app.Clients()

		health.NewModule().RegisterRouter(r)

		api := r.Group("/api/v1")

		public := api.Group("")
		protected := api.Group("")
		protected.Use(authmw.AuthenticationMiddleware(clients.Auth))

		modules.Auth.RegisterRouter(public, protected)
		modules.Post.RegisterRouter(public, protected)
		modules.File.RegisterRouter(protected)
		modules.User.RegisterRouter(protected)
	}
}
