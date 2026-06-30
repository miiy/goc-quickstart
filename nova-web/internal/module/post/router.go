package post

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(public, protected gin.IRouter) {
	handler := m.handler

	publicGroup := public.Group("/posts")
	publicGroup.GET("", handler.index)
	publicGroup.GET("/pages/:page", handler.pages)
	publicGroup.GET("/:id", handler.show)

	protectedGroup := protected.Group("/posts")
	protectedGroup.GET("/create", handler.create)
	protectedGroup.POST("", handler.store)
	protectedGroup.GET("/:id/edit", handler.edit)
	protectedGroup.POST("/:id", handler.post)
	protectedGroup.PUT("/:id", handler.update)
	protectedGroup.DELETE("/:id", handler.destroy)
}

func Templates() map[string][]string {
	return map[string][]string{
		"post/detail": {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/detail.html"},
		"post/list":   {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/list.html"},
		"post/create": {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/create.html"},
		"post/edit":   {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/edit.html"},
	}
}
