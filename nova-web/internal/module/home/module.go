package home

import (
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/logger"
)

type Module struct {
	router *gin.Engine
	logger logger.Logger
}

var module *Module

func NewModule(router *gin.Engine, logger logger.Logger) *Module {
	module = &Module{
		router: router,
		logger: logger,
	}
	return module
}
