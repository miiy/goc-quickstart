package auth

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/logger"
)

type AuthModule struct {
	logger         logger.Logger
	authClient     *client.AuthClient
	sessionManager *websession.Manager
}

var authModule *AuthModule

func NewModule(logger logger.Logger, authClient *client.AuthClient, sessionManager *websession.Manager) *AuthModule {
	authModule = &AuthModule{
		logger:         logger,
		authClient:     authClient,
		sessionManager: sessionManager,
	}
	return authModule
}
