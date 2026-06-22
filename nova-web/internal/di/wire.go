//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"github.com/google/wire"
	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc-quickstart/nova-web/internal/app"
	"github.com/miiy/goc-quickstart/nova-web/internal/config"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
)

func InitApp(conf string) (*app.App, func(), error) {
	panic(wire.Build(
		config.NewConfig,
		wire.NewSet(logger.NewLogger, provideLoggerOption),
		provideClients,
		provideSessionOptions,
		provideSessionStore,
		app.NewApp,
	))
}

func provideLoggerOption() []logger.Option {
	return nil
}

func provideClients(config *config.Config) (*client.Clients, func(), error) {
	return client.NewClients(config.Gateway.Addr)
}

func provideSessionOptions(config *config.Config) sessions.Options {
	return sessions.Options{
		Path:     "/",
		MaxAge:   config.Session.MaxAge,
		Secure:   config.Session.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func provideSessionStore(config *config.Config, options sessions.Options) (sessions.Store, error) {
	store, err := sessions.NewRedisStore(10, "tcp", config.Redis.Addr, config.Redis.Password, []byte(config.Session.Secret))
	if err != nil {
		return nil, err
	}
	if err := sessions.UseJSONSerializer(store); err != nil {
		return nil, err
	}
	store.Options(options)
	if config.Session.MaxAge > 0 {
		if err := sessions.SetMaxAge(store, config.Session.MaxAge); err != nil {
			return nil, err
		}
	}
	return store, nil
}
