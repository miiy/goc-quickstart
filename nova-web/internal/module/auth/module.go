package auth

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/logger"
)

type Module struct {
	handler *AuthHandler
}

func NewModule(log logger.Logger, authClient *client.AuthClient, sessionManager *session.Manager) *Module {
	return &Module{
		handler: NewAuthHandler(log, authClient, sessionManager),
	}
}
