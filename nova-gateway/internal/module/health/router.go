package health

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(r gin.IRouter) {
	api := m.healthAPI

	r.GET("/healthz", api.Healthz)
}
