package user

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/logger"
)

type Module struct {
	logger         logger.Logger
	authClient     *client.AuthClient
	userClient     *client.UserClient
	fileClient     *client.FileClient
	sessionManager *websession.Manager
}

var userModule *Module

func NewModule(logger logger.Logger, authClient *client.AuthClient, userClient *client.UserClient, fileClient *client.FileClient, sessionManager *websession.Manager) *Module {
	userModule = &Module{
		logger:         logger,
		authClient:     authClient,
		userClient:     userClient,
		fileClient:     fileClient,
		sessionManager: sessionManager,
	}
	return userModule
}
