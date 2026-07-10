package app

import (
	"github.com/miiy/goc-quickstart/nova-gateway/internal/client"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/config"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
)

type App struct {
	config *config.Config
	logger logger.Logger

	clients *client.Clients
	session sessions.Store
	modules *module.Modules
}

func NewApp(cfg *config.Config, logger logger.Logger, clients *client.Clients, session sessions.Store, modules *module.Modules) *App {
	return &App{
		config:  cfg,
		logger:  logger,
		clients: clients,
		session: session,
		modules: modules,
	}
}

func (a *App) Config() *config.Config {
	return a.config
}

func (a *App) Logger() logger.Logger {
	return a.logger
}

func (a *App) Clients() *client.Clients {
	return a.clients
}

func (a *App) SessionStore() sessions.Store {
	return a.session
}

func (a *App) Modules() *module.Modules {
	return a.modules
}
