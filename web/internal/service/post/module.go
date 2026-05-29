package post

import (
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc-quickstart/web/client"
)

type PostModule struct {
	logger logger.Logger
	client *client.PostClient
}

var postModule *PostModule

func NewModule(logger logger.Logger, postClient *client.PostClient) *PostModule {
	postModule = &PostModule{
		logger: logger,
		client: postClient,
	}
	return postModule
}