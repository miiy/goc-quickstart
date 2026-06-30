package user

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/logger"
)

type Module struct {
	handler *UserHandler
}

func NewModule(log logger.Logger, authClient *client.AuthClient, userClient *client.UserClient, fileClient *client.FileClient, sessionManager *websession.Manager) *Module {
	return &Module{
		handler: NewUserHandler(log, authClient, userClient, fileClient, sessionManager),
	}
}
