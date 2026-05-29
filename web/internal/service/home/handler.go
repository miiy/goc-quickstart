package home

import (
	"github.com/miiy/goc-quickstart/web/internal/template"
	"github.com/miiy/goc/gin"
)

type HomeView struct {
	template.ViewData
	Header  string
	Content string
}

func indexHandler(c *gin.Context) {
	c.HTML(200, "home/index", HomeView{
		ViewData: template.ViewData{
			PageTitle:  "Home",
			IsLoggedIn: c.GetBool("isLoggedIn"),
		},
		Header:  "header.",
		Content: "Hello, world.",
	})
}
