package health

type Module struct {
	handler *HealthHandler
}

func NewModule(gatewayAddr string) *Module {
	return &Module{
		handler: NewHealthHandler(gatewayAddr),
	}
}
