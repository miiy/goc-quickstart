package post

import (
	postv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-web/gen/go/nova/user/v1"
)

type Module struct {
	handler *PostsHandler
}

func NewModule(postClient postv1.PostServiceClient, userClient userv1.UserServiceClient) *Module {
	return &Module{
		handler: NewPostsHandler(postClient, userClient),
	}
}
