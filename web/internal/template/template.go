package template

import (
	"net/http"

	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/csrf"
	gocTemplate "github.com/miiy/goc/gin/template"
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
	CSRFToken   string
}

func NewViewData(c *gin.Context) ViewData {
	view := ViewData{
		IsLoggedIn: c.GetBool("isLoggedIn"),
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

// FlashLevelClass maps flash levels to Bootstrap alert CSS classes.
func FlashLevelClass(level string) string {
	switch level {
	case "error":
		return "danger"
	default:
		return level
	}
}
