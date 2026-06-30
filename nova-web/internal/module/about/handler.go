package about

import (
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc/gin"
)

type AboutHandler struct{}

func NewAboutHandler() *AboutHandler {
	return &AboutHandler{}
}

func (h *AboutHandler) index(c *gin.Context) {
	view := template.NewViewData(c)
	view.PageTitle = "About"
	c.HTML(200, "about/index", view)
}
