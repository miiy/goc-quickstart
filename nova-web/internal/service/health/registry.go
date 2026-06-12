package health

import (
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/logger"
)

type Module struct {
	router      *gin.Engine
	logger      logger.Logger
	gatewayAddr string
}

var (
	module *Module
)

func NewModule(router *gin.Engine, logger logger.Logger, gatewayAddr string) *Module {
	module = &Module{
		router:      router,
		logger:      logger,
		gatewayAddr: gatewayAddr,
	}
	return module
}
