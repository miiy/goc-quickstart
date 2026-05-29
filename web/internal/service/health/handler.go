package health

import (
	"net/http"
	"time"

	"github.com/miiy/goc/gin"
)

// Liveness - 检查服务是否存活
func livenessHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}

// Readiness - 检查服务是否就绪
// 对于 web 服务，主要检查 gateway 连接是否正常
func readinessHandler(ctx *gin.Context) {
	gatewayAddr := module.gatewayAddr
	if gatewayAddr == "" {
		ctx.String(http.StatusOK, "Ok")
		return
	}

	// 检查 gateway 是否可用
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + gatewayAddr + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.String(http.StatusServiceUnavailable, "gateway unavailable")
		return
	}
	resp.Body.Close()

	ctx.String(http.StatusOK, "Ok")
}
