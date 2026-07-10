package router

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/miiy/goc-quickstart/nova-web/internal/app"
	"github.com/miiy/goc-quickstart/nova-web/internal/dev"
	authmw "github.com/miiy/goc-quickstart/nova-web/internal/middleware/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/about"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/health"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/home"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/page"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/post"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/user"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
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

		// live reload for template and static files in debug mode
		var templateFS fs.FS = resourceTemplate.FS
		if cfg.App.Debug {
			templateFS = os.DirFS("resources/template")
			dev.NewLiveReload("resources/template", "resources/static").RegisterRouter(r)
		}

		// assets
		staticRoot := cfg.Static.Root
		if cfg.App.Debug && staticRoot == "dist" {
			staticRoot = "resources/static"
		}
		if staticRoot == "" {
			staticRoot = "dist"
		}
		staticFS := http.Dir(staticRoot)
		r.StaticFS("/static", staticFS)

		// favicon
		faviconHandler := func(c *gin.Context) {
			c.FileFromFS("favicon.ico", staticFS)
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

		// Proactively refresh an expiring access token, then load the session-backed
		// identity for downstream handlers.
		r.Use(authmw.RefreshSessionToken(sessionManager, authmw.NewAuthServiceRefresher(app.Clients().Auth), logger))
		r.Use(sessionauth.LoadSessionUser())

		public := r.Group("")
		protected := r.Group("")
		protected.Use(sessionauth.Authenticate(sessionauth.WithRedirect("/login")))
		modules := app.Modules()

		// template
		template.SetDefaultSite(template.SiteData{
			Name:            cfg.App.Name,
			Description:     cfg.App.Description,
			URL:             cfg.App.Url,
			Locale:          cfg.App.Locale,
			RegisterEnabled: cfg.App.RegisterEnabled,
			LiveReload:      cfg.App.Debug,
		})

		t := pkgTemplate.NewTemplate()
		t.AddFunc("alertType", template.FlashLevelClass)
		t.AddFunc("formatTime", template.NewFormatTimeFunc(cfg.App.Timezone))
		viteAssets, err := template.NewViteAssets(template.ViteAssetsConfig{
			Dev:          cfg.App.Debug,
			DevServerURL: os.Getenv("VITE_DEV_SERVER_URL"),
			ManifestPath: filepath.Join(staticRoot, ".vite", "manifest.json"),
			StaticPrefix: "/static/",
		})
		if err != nil {
			panic(err)
		}
		t.AddFunc("viteClient", viteAssets.Client)
		t.AddFunc("viteEntry", viteAssets.Entry)
		t.AddTemplate(templateFS, about.Templates())
		t.AddTemplate(templateFS, home.Templates())
		t.AddTemplate(templateFS, post.Templates())
		t.AddTemplate(templateFS, auth.Templates())
		t.AddTemplate(templateFS, user.Templates())
		t.AddTemplate(templateFS, page.Templates())

		r.HTMLRender = t.Render

		// fallback: render 404 page for unmatched routes and methods
		r.NoRoute(template.NotFound)
		r.NoMethod(template.NotFound)

		// modules router
		about.RegisterRouter(r)
		home.RegisterRouter(r)
		modules.Post.RegisterRouter(public, protected)
		modules.Auth.RegisterRouter(r)
		modules.User.RegisterRouter(public, protected)
	}
}
