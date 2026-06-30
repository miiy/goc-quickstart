package post

import (
	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
)

type Module struct {
	postsAPI openapi.PostsAPI
}

func NewModule(postClient postv1.PostServiceClient, userClient userv1.UserServiceClient) *Module {
	return &Module{
		postsAPI: NewPostsAPI(postClient, userClient),
	}
}
