package auth

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
)

type AuthModule struct {
	logger       logger.Logger
	client       *client.AuthClient
	sessionStore sessions.Store
	sessionName  string
}

var authModule *AuthModule

func NewModule(logger logger.Logger, authClient *client.AuthClient, sessionStore sessions.Store, sessionName string) *AuthModule {
	authModule = &AuthModule{
		logger:       logger,
		client:       authClient,
		sessionStore: sessionStore,
		sessionName:  sessionName,
	}
	return authModule
}
