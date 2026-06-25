//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"github.com/google/wire"
	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc-quickstart/nova-web/internal/app"
	"github.com/miiy/goc-quickstart/nova-web/internal/config"
	"github.com/miiy/goc-quickstart/nova-web/internal/module"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/post"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/user"
	"github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
)

func InitApp(conf string) (*app.App, func(), error) {
	panic(wire.Build(
		config.NewConfig,
		wire.NewSet(logger.NewLogger, provideLoggerOption),
		provideClients,
		providePostClient,
		provideAuthClient,
		provideUserClient,
		provideFileClient,
		provideSessionOptions,
		provideSessionStore,
		provideSessionManager,
		provideModules,
		app.NewApp,
	))
}

func provideLoggerOption() []logger.Option {
	return nil
}

func provideClients(config *config.Config) (*client.Clients, func(), error) {
	return client.NewClients(config.Gateway.Addr)
}

func providePostClient(clients *client.Clients) *client.PostClient {
	return clients.Post
}

func provideAuthClient(clients *client.Clients) *client.AuthClient {
	return clients.Auth
}

func provideUserClient(clients *client.Clients) *client.UserClient {
	return clients.User
}

func provideFileClient(clients *client.Clients) *client.FileClient {
	return clients.File
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

func provideSessionManager(store sessions.Store, config *config.Config) *session.Manager {
	return session.NewManager(store, config.Session.Name)
}

func provideModules(
	log logger.Logger,
	postClient *client.PostClient,
	authClient *client.AuthClient,
	userClient *client.UserClient,
	fileClient *client.FileClient,
	sessionManager *session.Manager,
) *module.Modules {
	return &module.Modules{
		Post: post.NewModule(log, postClient, fileClient, sessionManager),
		Auth: auth.NewModule(log, authClient, sessionManager),
		User: user.NewModule(log, authClient, userClient, fileClient, sessionManager),
	}
}
