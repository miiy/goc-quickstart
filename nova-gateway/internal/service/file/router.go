package file

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(r gin.IRouter) {
	g := r.Group("/files")
	g.POST("/upload/avatar", m.avatar)
}
