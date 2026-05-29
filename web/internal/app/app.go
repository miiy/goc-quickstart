package app

import (
	"github.com/miiy/goc-quickstart/web/internal/config"
	"github.com/miiy/goc/logger"
)

type App struct {
	config *config.Config
	logger logger.Logger
}

func NewApp(config *config.Config, logger logger.Logger) *App {
	return &App{
		config: config,
		logger: logger,
	}
}

func (a *App) Config() *config.Config {
	return a.config
}

func (a *App) Logger() logger.Logger {
	return a.logger
}
