package post

import (
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
)

type Module struct {
	client     postv1.PostServiceClient
	userClient authorUserClient
}

func NewModule(client postv1.PostServiceClient, userClient userv1.UserServiceClient) *Module {
	return &Module{
		client:     client,
		userClient: userClient,
	}
}
