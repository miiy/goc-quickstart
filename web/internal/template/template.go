package template

import (
	"net/http"

	"github.com/miiy/goc/gin"
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
}

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "pages/404", ViewData{IsLoggedIn: c.GetBool("isLoggedIn")})
}

func InternalError(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "pages/500", ViewData{IsLoggedIn: c.GetBool("isLoggedIn")})
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
