package router

import (
	"net/http"

	webclient "github.com/miiy/goc-quickstart/web/client"
	"github.com/miiy/goc-quickstart/web/internal/config"
	"github.com/miiy/goc-quickstart/web/internal/service/about"
	"github.com/miiy/goc-quickstart/web/internal/service/auth"
	"github.com/miiy/goc-quickstart/web/internal/service/home"
	"github.com/miiy/goc-quickstart/web/internal/service/page"
	"github.com/miiy/goc-quickstart/web/internal/service/post"
	"github.com/miiy/goc-quickstart/web/internal/template"
	"github.com/miiy/goc-quickstart/web/resources/static"
	resourceTemplate "github.com/miiy/goc-quickstart/web/resources/template"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	gocauthmid "github.com/miiy/goc/gin/middleware/auth"
	"github.com/miiy/goc/gin/middleware/csrf"
	"github.com/miiy/goc/gin/sessions"
	pkgTemplate "github.com/miiy/goc/gin/template"
)

func Router(r *gin.Engine, store sessions.Store, sessionName string) {
	// assets
	r.StaticFS("/static", http.FS(static.FS))

	// favicon
	faviconHandler := func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(static.FS))
	}
	r.HEAD("/favicon.ico", faviconHandler)
	r.GET("/favicon.ico", faviconHandler)

	r.Use(sessions.Middleware(sessionName, store))
	r.Use(csrf.Middleware())

	// inject login state into template context
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		user, ok := gocauthmid.SessionUser(session.Get(gocauthmid.SessionKeyAuthUser))
		c.Set("isLoggedIn", ok)
		if ok {
			c.Request = c.Request.WithContext(gocauth.InjectAuthenticatedUser(c.Request.Context(), user))
		}

		if token, ok := session.Get(auth.SessionKeyAccessToken).(string); ok {
			c.Request = c.Request.WithContext(webclient.WithAccessToken(c.Request.Context(), token))
		}

		c.Next()
	})

	// template
	t := pkgTemplate.NewTemplate()
	t.AddFunc("config", func(key string) any {
		return config.GetConfig(key)
	})
	t.AddFunc("alertType", template.FlashLevelClass)
	t.AddTemplate(resourceTemplate.FS, about.Templates())
	t.AddTemplate(resourceTemplate.FS, home.Templates())
	t.AddTemplate(resourceTemplate.FS, post.Templates())
	t.AddTemplate(resourceTemplate.FS, auth.Templates())
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
}
