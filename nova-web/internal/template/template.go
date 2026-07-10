package template

import (
	"net/http"
	"time"

	"github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/authctx"
	"github.com/miiy/goc/gin/middleware/csrf"
	timeutil "github.com/miiy/goc/utils/time"
)

// SiteData is the public, display-safe site metadata available to templates.
type SiteData struct {
	Name            string
	Description     string
	URL             string
	Locale          string
	RegisterEnabled bool
	LiveReload      bool
}

// AuthData is the authenticated request identity available to templates.
type AuthData struct {
	IsLoggedIn  bool
	CurrentUser *auth.AuthenticatedUser
}

// ViewData is the common data passed to all templates.
type ViewData struct {
	PageTitle   string
	Keywords    string
	Description string
	Site        SiteData
	Auth        AuthData
	CSRFToken   string
}

var defaultSite SiteData

func SetDefaultSite(site SiteData) {
	defaultSite = site
}

func NewViewData(c *gin.Context) ViewData {
	view := ViewData{
		Site: defaultSite,
	}
	if user, ok := authctx.CurrentUser(c); ok {
		view.Auth.IsLoggedIn = true
		view.Auth.CurrentUser = user
	}
	view.CSRFToken = csrf.Token(c)
	return view
}

func NewFormViewData(c *gin.Context) ViewData {
	return NewViewData(c)
}

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "pages/404", NewViewData(c))
}

func InternalError(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "pages/500", NewViewData(c))
}

// NewFormatTimeFunc binds template time formatting to the configured display timezone.
func NewFormatTimeFunc(timezone string) func(any) string {
	loc, err := timeutil.LoadLocation(timezone)
	if err != nil {
		loc = time.Local
	}
	return func(v any) string {
		return timeutil.FormatTime(v, loc, timeutil.DateMinuteLayout)
	}
}

// FlashLevelClass maps flash levels to alert CSS class suffixes.
func FlashLevelClass(level string) string {
	switch level {
	case "error":
		return "error"
	default:
		return level
	}
}
