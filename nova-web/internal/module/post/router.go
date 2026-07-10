package post

import (
	resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"
	"github.com/miiy/goc/gin"
)

func (m *Module) RegisterRouter(public, protected gin.IRouter) {
	handler := m.handler

	publicGroup := public.Group("/posts")
	publicGroup.GET("", handler.list)
	publicGroup.GET("/pages/:page", handler.list)
	publicGroup.GET("/:id", handler.show)

	protectedGroup := protected.Group("/posts")
	protectedGroup.GET("/create", handler.create)
	protectedGroup.GET("/:id/edit", handler.edit)
}

func Templates() map[string][]string {
	return map[string][]string{
		"post/detail": resourceTemplate.Layout("post/detail.html"),
		"post/list":   resourceTemplate.Layout("post/list.html"),
		"post/create": resourceTemplate.Layout("post/create.html"),
		"post/edit":   resourceTemplate.Layout("post/edit.html"),
	}
}
