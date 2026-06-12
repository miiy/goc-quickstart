package user

import (
	"github.com/miiy/goc-quickstart/web/client"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
)

type Module struct {
	logger       logger.Logger
	authClient   *client.AuthClient
	userClient   *client.UserClient
	sessionStore sessions.Store
	sessionName  string
}

var userModule *Module

func NewModule(logger logger.Logger, authClient *client.AuthClient, userClient *client.UserClient, sessionStore sessions.Store, sessionName string) *Module {
	userModule = &Module{
		logger:       logger,
		authClient:   authClient,
		userClient:   userClient,
		sessionStore: sessionStore,
		sessionName:  sessionName,
	}
	return userModule
}
