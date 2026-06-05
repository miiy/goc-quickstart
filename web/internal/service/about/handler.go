package about

import (
	"github.com/miiy/goc-quickstart/web/internal/template"
	"github.com/miiy/goc/gin"
)

func indexHandler(c *gin.Context) {
	view := template.NewViewData(c)
	view.PageTitle = "About"
	c.HTML(200, "about/index", view)
}
