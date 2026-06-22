package health

import "github.com/miiy/goc/gin"

type Module struct {
	router *gin.Engine
}

func NewModule(router *gin.Engine) *Module {
	return &Module{router: router}
}
