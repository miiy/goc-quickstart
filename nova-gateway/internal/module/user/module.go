package user

import (
	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
)

type Module struct {
	usersAPI openapi.UsersAPI
}

func NewModule(userClient userv1.UserServiceClient) *Module {
	return &Module{usersAPI: NewUsersAPI(userClient)}
}
