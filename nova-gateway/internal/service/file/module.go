package file

import (
	filev1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/file/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
)

type Module struct {
	fileClient filev1.FileServiceClient
	userClient userv1.UserServiceClient
}

func NewModule(fileClient filev1.FileServiceClient, userClient userv1.UserServiceClient) *Module {
	return &Module{
		fileClient: fileClient,
		userClient: userClient,
	}
}
