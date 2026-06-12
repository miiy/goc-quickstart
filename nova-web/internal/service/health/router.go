package health

func (m *Module) RegisterRouter() {

	hG := m.router.Group("/health")
	{
		hG.GET("/liveness", livenessHandler)
		hG.GET("/readiness", readinessHandler)
	}
}
