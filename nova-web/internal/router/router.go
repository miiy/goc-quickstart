package router

import (
	htmltemplate "html/template"
	"net/http"

	"github.com/miiy/goc-quickstart/nova-web/internal/config"
	authmw "github.com/miiy/goc-quickstart/nova-web/internal/middleware/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/about"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/home"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/page"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/post"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/user"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc-quickstart/nova-web/resources/static"
	resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/csrf"
	"github.com/miiy/goc/gin/middleware/sessionauth"
	pkgTemplate "github.com/miiy/goc/gin/template"
	"github.com/miiy/goc/logger"
)

func Router(r *gin.Engine, sessionManager *websession.Manager, appConfig config.AppConfig, authClient authmw.RefreshClient, log logger.Logger, storageRoot string) {
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

	// Proactively refresh an expiring access token, then bridge the session-backed
	// identity into the gin/template + request contexts for downstream handlers.
	r.Use(authmw.RefreshSessionToken(sessionManager, authClient, log))
	r.Use(sessionauth.LoadSessionUser())

	// template
	template.SetDefaultSite(template.SiteData{
		Name:            appConfig.Name,
		URL:             appConfig.Url,
		Locale:          appConfig.Locale,
		FooterCopyright: htmltemplate.HTML(appConfig.FooterCopyright),
	})

	t := pkgTemplate.NewTemplate()
	t.AddFunc("alertType", template.FlashLevelClass)
	t.AddFunc("formatTime", template.NewFormatTimeFunc(appConfig.Timezone))
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
