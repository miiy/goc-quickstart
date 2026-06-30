package router

import (
	htmltemplate "html/template"
	"net/http"

	"github.com/miiy/goc-quickstart/nova-web/internal/app"
	authmw "github.com/miiy/goc-quickstart/nova-web/internal/middleware/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/about"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/health"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/home"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/page"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/post"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/user"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc-quickstart/nova-web/resources/static"
	resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/csrf"
	"github.com/miiy/goc/gin/middleware/sessionauth"
	ginzap "github.com/miiy/goc/gin/middleware/zap"
	pkgTemplate "github.com/miiy/goc/gin/template"
)

func Router(app *app.App) func(r *gin.Engine) {
	return func(r *gin.Engine) {
		cfg := app.Config()
		sessionManager := app.SessionManager()
		logger := app.Logger()

		r.Use(ginzap.Ginzap(logger.ZapLogger()), ginzap.RecoveryWithZap(logger.ZapLogger(), true))

		health.NewModule(cfg.Gateway.Addr).RegisterRouter(r)

		// assets
		r.StaticFS("/static", http.FS(static.FS))

		// favicon
		faviconHandler := func(c *gin.Context) {
			c.FileFromFS("favicon.ico", http.FS(static.FS))
		}
		r.HEAD("/favicon.ico", faviconHandler)
		r.GET("/favicon.ico", faviconHandler)

		// uploaded media (avatars, post covers, ...) served straight from nova-file's
		// storage dir. Registered before session/csrf/auth middleware so it's publicly
		// readable, just like /static. nova-file stores relative object keys; handlers
		// expand those keys to /uploads/... URLs for browser rendering.
		if cfg.Storage.Root != "" {
			r.Static("/uploads", cfg.Storage.Root)
		}

		// Load the session before any auth-aware middleware or handlers run.
		r.Use(sessionManager.Middleware())

		r.Use(csrf.Middleware())

		// Proactively refresh an expiring access token, then bridge the session-backed
		// identity into the gin/template + request contexts for downstream handlers.
		r.Use(authmw.RefreshSessionToken(sessionManager, authmw.NewAuthClientTokenRefresher(app.Clients().Auth), logger))
		r.Use(sessionauth.LoadSessionUser())

		public := r.Group("")
		protected := r.Group("")
		protected.Use(sessionauth.Authenticate(sessionauth.WithRedirect("/login")))
		modules := app.Modules()

		// template
		template.SetDefaultSite(template.SiteData{
			Name:            cfg.App.Name,
			URL:             cfg.App.Url,
			Locale:          cfg.App.Locale,
			FooterCopyright: htmltemplate.HTML(cfg.App.FooterCopyright),
		})

		t := pkgTemplate.NewTemplate()
		t.AddFunc("alertType", template.FlashLevelClass)
		t.AddFunc("formatTime", template.NewFormatTimeFunc(cfg.App.Timezone))
		t.AddTemplate(resourceTemplate.FS, about.Templates())
		t.AddTemplate(resourceTemplate.FS, home.Templates())
		t.AddTemplate(resourceTemplate.FS, post.Templates())
		t.AddTemplate(resourceTemplate.FS, auth.Templates())
		t.AddTemplate(resourceTemplate.FS, user.Templates())
		t.AddTemplate(resourceTemplate.FS, page.Templates())

		r.HTMLRender = t.Render

		// fallback: render 404 page for unmatched routes and methods
		r.NoRoute(template.NotFound)
		r.NoMethod(template.NotFound)

		// modules router
		about.RegisterRouter(r)
		home.RegisterRouter(r)
		modules.Post.RegisterRouter(public, protected)
		modules.Auth.RegisterRouter(r)
		modules.User.RegisterRouter(protected)
	}
}
