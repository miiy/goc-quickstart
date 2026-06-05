package health

import (
	"net/http"

	"github.com/miiy/goc/gin"
)

func (m *Module) RegisterRouter() {
	m.router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
}
