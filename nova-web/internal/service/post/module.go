package post

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	"github.com/miiy/goc/logger"
)

type PostModule struct {
	logger     logger.Logger
	postClient *client.PostClient
	fileClient *client.FileClient
}

var postModule *PostModule

func NewModule(logger logger.Logger, postClient *client.PostClient, fileClient *client.FileClient) *PostModule {
	postModule = &PostModule{
		logger:     logger,
		postClient: postClient,
		fileClient: fileClient,
	}
	return postModule
}
