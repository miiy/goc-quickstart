package template

import (
	"net/http"
	"time"

	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin"
	gocauthmid "github.com/miiy/goc/gin/middleware/auth"
	"github.com/miiy/goc/gin/middleware/csrf"
	gocTemplate "github.com/miiy/goc/gin/template"
	timeutil "github.com/miiy/goc/utils/time"
)

// Re-export goc template types for convenience
type Template = gocTemplate.Template

func NewTemplate() *Template {
	return gocTemplate.NewTemplate()
}

// ViewData is the common data passed to all templates.
type ViewData struct {
	PageTitle   string
	Keywords    string
	Description string
	IsLoggedIn  bool
	CurrentUser *gocauth.AuthenticatedUser
	CSRFToken   string
}

func NewViewData(c *gin.Context) ViewData {
	view := ViewData{
		IsLoggedIn: c.GetBool("isLoggedIn"),
	}
	if user, ok := c.Get("currentUser"); ok {
		if authUser, ok := user.(*gocauth.AuthenticatedUser); ok {
			view.CurrentUser = authUser
		} else if authUser, ok := gocauthmid.SessionUser(user); ok {
			view.CurrentUser = authUser
		}
	}
	if view.IsLoggedIn {
		view.CSRFToken = csrf.Token(c)
	}
	return view
}

func NewFormViewData(c *gin.Context) ViewData {
	view := NewViewData(c)
	view.CSRFToken = csrf.Token(c)
	return view
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

// FlashLevelClass maps flash levels to Bootstrap alert CSS classes.
func FlashLevelClass(level string) string {
	switch level {
	case "error":
		return "danger"
	default:
		return level
	}
}
