package health

import (
	"net/http"
	"time"

	"github.com/miiy/goc/gin"
)

type HealthHandler struct {
	gatewayAddr string
}

func NewHealthHandler(gatewayAddr string) *HealthHandler {
	return &HealthHandler{gatewayAddr: gatewayAddr}
}

// Liveness - 检查服务是否存活
func (h *HealthHandler) liveness(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}

// Readiness - 检查服务是否就绪
// 对于 web 服务，主要检查 gateway 连接是否正常
func (h *HealthHandler) readiness(ctx *gin.Context) {
	gatewayAddr := h.gatewayAddr
	if gatewayAddr == "" {
		ctx.String(http.StatusOK, "Ok")
		return
	}

	// 检查 gateway 是否可用
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + gatewayAddr + "/healthz")
	if err != nil {
		ctx.String(http.StatusServiceUnavailable, "gateway unavailable")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ctx.String(http.StatusServiceUnavailable, "gateway unavailable")
		return
	}

	ctx.String(http.StatusOK, "Ok")
}
