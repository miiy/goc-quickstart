package post

import (
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
)

type Module struct {
	postClient postv1.PostServiceClient
	userClient authorUserClient
}

func NewModule(postClient postv1.PostServiceClient, userClient userv1.UserServiceClient) *Module {
	return &Module{
		postClient: postClient,
		userClient: userClient,
	}
}
