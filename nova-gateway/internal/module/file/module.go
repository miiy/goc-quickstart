package file

import (
	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	filev1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/file/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
)

type Module struct {
	filesAPI openapi.FilesAPI
}

func NewModule(fileClient filev1.FileServiceClient, userClient userv1.UserServiceClient) *Module {
	return &Module{filesAPI: NewFilesAPI(fileClient, userClient)}
}
