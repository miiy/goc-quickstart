//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"github.com/google/wire"
	authv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/auth/v1"
	postv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-web/internal/app"
	"github.com/miiy/goc-quickstart/nova-web/internal/client"
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
		providePostService,
		provideAuthService,
		provideUserService,
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
	return client.NewClients(config)
}

func providePostService(clients *client.Clients) postv1.PostServiceClient {
	return clients.Post
}

func provideAuthService(clients *client.Clients) authv1.AuthServiceClient {
	return clients.Auth
}

func provideUserService(clients *client.Clients) userv1.UserServiceClient {
	return clients.User
}

func provideSessionOptions(config *config.Config) sessions.Options {
	return sessions.Options{
		Path:     "/",
		Domain:   config.Session.Domain,
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
	return session.NewManager(store, config.Session.Name, provideSessionOptions(config))
}

func provideModules(
	cfg *config.Config,
	log logger.Logger,
	postClient postv1.PostServiceClient,
	authClient authv1.AuthServiceClient,
	userClient userv1.UserServiceClient,
	sessionManager *session.Manager,
) *module.Modules {
	return &module.Modules{
		Post: post.NewModule(postClient, userClient),
		Auth: auth.NewModule(log, authClient, sessionManager, cfg.App.RegisterEnabled),
		User: user.NewModule(userClient, sessionManager),
	}
}
