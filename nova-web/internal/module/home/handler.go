package home

import (
	"github.com/miiy/goc-quickstart/nova-web/internal/template"
	"github.com/miiy/goc/gin"
)

type HomeView struct {
	template.ViewData
	Header  string
	Content string
}

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) index(c *gin.Context) {
	view := template.NewViewData(c)
	view.PageTitle = "Home"
	c.HTML(200, "home/index", HomeView{
		ViewData: view,
		Header:   "header.",
		Content:  "Hello, world.",
	})
}
