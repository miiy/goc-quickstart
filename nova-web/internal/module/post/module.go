package post

import (
	"github.com/miiy/goc-quickstart/nova-web/client"
	websession "github.com/miiy/goc-quickstart/nova-web/internal/session"
	"github.com/miiy/goc/logger"
)

type Module struct {
	handler *PostsHandler
}

func NewModule(log logger.Logger, postClient *client.PostClient, fileClient *client.FileClient, sessionManager *websession.Manager) *Module {
	return &Module{
		handler: NewPostsHandler(log, postClient, fileClient, sessionManager),
	}
}
