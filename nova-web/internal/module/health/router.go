package health

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(r gin.IRouter) {
	handler := m.handler

	g := r.Group("/health")
	g.GET("/liveness", handler.liveness)
	g.GET("/readiness", handler.readiness)
}
