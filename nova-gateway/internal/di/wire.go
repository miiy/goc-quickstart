//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/app"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/config"
	"github.com/miiy/goc/logger"
)

func InitApp(conf string) (*app.App, func(), error) {
	panic(wire.Build(
		config.NewConfig,
		wire.NewSet(provideLoggerOption, provideLogger),
		app.NewApp,
	))
}

func provideLoggerOption() []logger.Option {
	return nil
}

func provideLogger(options []logger.Option) (logger.Logger, func(), error) {
	l, err := logger.NewLogger(options...)
	if err != nil {
		return nil, nil, err
	}
	return l, func() {
		if sync, ok := l.(interface{ Sync() error }); ok {
			_ = sync.Sync()
		}
	}, nil
}
