package upload

import (
	uploadv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/upload/v1"
	userv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/user/v1"
)

type Module struct {
	client     uploadv1.UploadServiceClient
	userClient userv1.UserServiceClient
}

func NewModule(client uploadv1.UploadServiceClient, userClient userv1.UserServiceClient) *Module {
	return &Module{
		client:     client,
		userClient: userClient,
	}
}
