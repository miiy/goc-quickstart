package router

import (
	"github.com/miiy/goc-quickstart/nova-gateway/internal/app"
	authmw "github.com/miiy/goc-quickstart/nova-gateway/internal/middleware/auth"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/health"
	"github.com/miiy/goc/gin"
	corsmw "github.com/miiy/goc/gin/middleware/cors"
	ginzap "github.com/miiy/goc/gin/middleware/zap"
	"github.com/miiy/goc/gin/sessions"
)

// Router returns a function that registers all routes and middleware onto the gin engine.
func Router(app *app.App) func(r *gin.Engine) {
	return func(r *gin.Engine) {
		cfg := app.Config()
		logger := app.Logger().ZapLogger()
		r.Use(ginzap.Ginzap(logger), ginzap.RecoveryWithZap(logger, true))
		if len(cfg.CORS.AllowOrigins) > 0 {
			corsConfig := corsmw.DefaultConfig()
			corsConfig.AllowOrigins = cfg.CORS.AllowOrigins
			corsConfig.AllowCredentials = cfg.CORS.AllowCredentials
			corsConfig.AddAllowHeaders("Authorization", "X-CSRF-Token")
			r.Use(corsmw.NewWithConfig(corsConfig))
		}

		modules := app.Modules()
		clients := app.Clients()

		health.NewModule().RegisterRouter(r)

		api := r.Group("/api/v1")

		public := api.Group("")
		protected := api.Group("")
		if store := app.SessionStore(); store != nil && cfg.Session.Name != "" {
			protected.Use(
				sessions.Middleware(cfg.Session.Name, store),
				authmw.SessionTokenToBearer(),
				authmw.SessionCSRF(),
			)
		}
		protected.Use(authmw.AuthenticationMiddleware(clients.Auth))

		modules.Auth.RegisterRouter(public, protected)
		modules.Post.RegisterRouter(public, protected)
		modules.File.RegisterRouter(protected)
		modules.User.RegisterRouter(public, protected)
	}
}
