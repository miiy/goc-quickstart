package post

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/logger"
)

type PostModule struct {
	logger         logger.Logger
	postClient     *client.PostClient
	fileClient     *client.FileClient
	sessionManager *websession.Manager
}

var postModule *PostModule

func NewModule(logger logger.Logger, postClient *client.PostClient, fileClient *client.FileClient, sessionManager *websession.Manager) *PostModule {
	postModule = &PostModule{
		logger:         logger,
		postClient:     postClient,
		fileClient:     fileClient,
		sessionManager: sessionManager,
	}
	return postModule
}
