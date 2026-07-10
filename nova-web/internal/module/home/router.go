package home

import (
	resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"
	"github.com/miiy/goc/gin"
)

func RegisterRouter(r gin.IRouter) {
	handler := NewHomeHandler()

	r.GET("/", handler.index)
}

func Templates() map[string][]string {
	return map[string][]string{
		"home/index": resourceTemplate.Layout("home/index.html"),
	}
}
