package post

import postv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/post/v1"

type Module struct {
	client postv1.PostServiceClient
}

func NewModule(client postv1.PostServiceClient) *Module {
	return &Module{client: client}
}
