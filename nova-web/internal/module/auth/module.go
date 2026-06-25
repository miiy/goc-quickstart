package auth

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/logger"
)

type AuthModule struct {
	logger         logger.Logger
	authClient     *client.AuthClient
	sessionManager *session.Manager
}

var authModule *AuthModule

func NewModule(logger logger.Logger, authClient *client.AuthClient, sessionManager *session.Manager) *AuthModule {
	authModule = &AuthModule{
		logger:         logger,
		authClient:     authClient,
		sessionManager: sessionManager,
	}
	return authModule
}
