package user

import userv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/user/v1"

type Module struct {
	client userv1.UserServiceClient
}

func NewModule(client userv1.UserServiceClient) *Module {
	return &Module{client: client}
}
