//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	filev1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/file/v1"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/app"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/client"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/config"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/auth"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/file"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/post"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/user"
	"github.com/miiy/goc/logger"
)

func InitApp(conf string) (*app.App, func(), error) {
	panic(wire.Build(
		config.NewConfig,
		wire.NewSet(provideLoggerOption, provideLogger),
		client.NewClients,
		provideAuthClient,
		providePostClient,
		provideFileClient,
		provideUserClient,
		provideModules,
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

func provideAuthClient(clients *client.Clients) authv1.AuthServiceClient {
	return clients.Auth
}

func providePostClient(clients *client.Clients) postv1.PostServiceClient {
	return clients.Post
}

func provideFileClient(clients *client.Clients) filev1.FileServiceClient {
	return clients.File
}

func provideUserClient(clients *client.Clients) userv1.UserServiceClient {
	return clients.User
}

func provideModules(
	authClient authv1.AuthServiceClient,
	postClient postv1.PostServiceClient,
	fileClient filev1.FileServiceClient,
	userClient userv1.UserServiceClient,
) *module.Modules {
	return &module.Modules{
		Auth: auth.NewModule(authClient),
		Post: post.NewModule(postClient, userClient),
		File: file.NewModule(fileClient, userClient),
		User: user.NewModule(userClient),
	}
}
