package app

import (
	"github.com/miiy/goc-quickstart/web/client"
	"github.com/miiy/goc-quickstart/web/internal/config"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
)

type App struct {
	config       *config.Config
	logger       logger.Logger
	clients      *client.Clients
	sessionStore sessions.Store
}

func NewApp(config *config.Config, logger logger.Logger, clients *client.Clients, sessionStore sessions.Store) *App {
	return &App{
		config:       config,
		logger:       logger,
		clients:      clients,
		sessionStore: sessionStore,
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
	return a.sessionStore
}
