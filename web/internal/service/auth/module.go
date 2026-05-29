package auth

import (
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc-quickstart/web/client"
)

type AuthModule struct {
	logger logger.Logger
	client *client.AuthClient
}

var authModule *AuthModule

func NewModule(logger logger.Logger, authClient *client.AuthClient) *AuthModule {
	authModule = &AuthModule{
		logger: logger,
		client: authClient,
	}
	return authModule
}