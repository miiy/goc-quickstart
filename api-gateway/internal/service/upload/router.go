package upload

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(r gin.IRouter) {
	g := r.Group("/uploads")
	g.POST("/avatar", m.avatar)
}
