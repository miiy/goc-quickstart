package router

import (
	"net/http"

	"github.com/miiy/goc-quickstart/nova-web/internal/config"
	authmid "github.com/miiy/goc-quickstart/nova-web/internal/middleware/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/service/about"
	"github.com/miiy/goc-quickstart/nova-web/internal/service/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/service/home"
	"github.com/miiy/goc-quickstart/nova-web/internal/service/page"
	"github.com/miiy/goc-quickstart/nova-web/internal/service/post"
	"github.com/miiy/goc-quickstart/nova-web/internal/service/user"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc-quickstart/nova-web/resources/static"
	resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/csrf"
	pkgTemplate "github.com/miiy/goc/gin/template"
	"github.com/miiy/goc/logger"
)

func Router(r *gin.Engine, sessionManager *websession.Manager, timezone string, authClient authmid.RefreshClient, log logger.Logger, storageRoot string) {
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
	// readable, just like /static. nova-file stores relative object keys; nova-web
	// expands those keys to /uploads/... URLs for browser rendering.
	if storageRoot != "" {
		r.Static("/uploads", storageRoot)
	}

	// Load the session before any auth-aware middleware or handlers run.
	r.Use(sessionManager.Middleware())

	r.Use(csrf.Middleware())

	// proactively refresh an expiring access token, then bridge the session-backed
	// identity into the gin/template + request contexts for downstream handlers.
	r.Use(authmid.SessionAuth(sessionManager, authClient, log))

	// template
	t := pkgTemplate.NewTemplate()
	t.AddFunc("config", func(key string) any {
		return config.GetConfig(key)
	})
	t.AddFunc("alertType", template.FlashLevelClass)
	t.AddFunc("formatTime", template.NewFormatTimeFunc(timezone))
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
	about.Router(r)
	home.Router(r)
	post.Router(r)
	auth.Router(r)
	user.Router(r)
}
