package user

import userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"

type Module struct {
	userClient userv1.UserServiceClient
}

func NewModule(userClient userv1.UserServiceClient) *Module {
	return &Module{userClient: userClient}
}
