package user

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(r gin.IRouter) {
	g := r.Group("/users")
	g.GET("", m.list)
	g.GET("/:id", m.get)
	g.PUT("/:id", m.update)
}
