package about

import (
	"github.com/miiy/goc-quickstart/web/internal/template"
	"github.com/miiy/goc/gin"
)

func indexHandler(c *gin.Context) {
	c.HTML(200, "about/index", template.ViewData{
		PageTitle:  "About",
		IsLoggedIn: c.GetBool("isLoggedIn"),
	})
}
