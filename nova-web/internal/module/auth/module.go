package auth

import (
	authv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/auth/v1"
	"github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/logger"
)

type Module struct {
	handler *AuthHandler
}

func NewModule(log logger.Logger, authClient authv1.AuthServiceClient, sessionManager *session.Manager, registerEnabled bool) *Module {
	return &Module{
		handler: NewAuthHandler(log, authClient, sessionManager, registerEnabled),
	}
}
