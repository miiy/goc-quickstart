package health

import (
	"net/http"

	"github.com/miiy/goc/gin"
)

type HealthAPI struct{}

func NewHealthAPI() *HealthAPI {
	return &HealthAPI{}
}

func (api *HealthAPI) Healthz(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
