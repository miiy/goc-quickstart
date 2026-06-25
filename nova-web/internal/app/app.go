package app

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc-quickstart/nova-web/internal/config"
	"github.com/miiy/goc-quickstart/nova-web/internal/module"
	"github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/logger"
)

type App struct {
	config         *config.Config
	logger         logger.Logger
	clients        *client.Clients
	sessionManager *session.Manager
	modules        *module.Modules
}

func NewApp(config *config.Config, logger logger.Logger, clients *client.Clients, sessionManager *session.Manager, modules *module.Modules) *App {
	return &App{
		config:         config,
		logger:         logger,
		clients:        clients,
		sessionManager: sessionManager,
		modules:        modules,
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

func (a *App) SessionManager() *session.Manager {
	return a.sessionManager
}

func (a *App) Modules() *module.Modules {
	return a.modules
}
