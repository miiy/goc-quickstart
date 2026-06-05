package post

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(r gin.IRouter) {
	m.RegisterPublicRouter(r)
	m.RegisterProtectedRouter(r)
}

func (m *Module) RegisterPublicRouter(r gin.IRouter) {
	g := r.Group("/posts")
	g.GET("", m.list)
	g.GET("/:id", m.get)
}

func (m *Module) RegisterProtectedRouter(r gin.IRouter) {
	g := r.Group("/posts")
	g.POST("", m.create)
	g.PUT("/:id", m.update)
	g.DELETE("/:id", m.delete)
}
