package health

type Module struct {
	healthAPI *HealthAPI
}

func NewModule() *Module {
	return &Module{healthAPI: NewHealthAPI()}
}
