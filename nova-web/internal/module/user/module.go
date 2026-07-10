package user

import (
	userv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/user/v1"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
)

type Module struct {
	handler *UserHandler
}

func NewModule(userClient userv1.UserServiceClient, sessionManager *websession.Manager) *Module {
	return &Module{
		handler: NewUserHandler(userClient, sessionManager),
	}
}
