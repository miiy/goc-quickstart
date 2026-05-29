package app

import (
	pb "github.com/miiy/goc-quickstart/post-service/gen/go/blog/post/v1"
	"github.com/miiy/goc-quickstart/post-service/internal/config"
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/redis"
)

type App struct {
	config      *config.Config
	db          *db.DB
	redis       redis.UniversalClient
	logger      logger.Logger
	postService  pb.PostServiceServer
}

var app *App

func NewApp(c *config.Config, db *db.DB, rdb redis.UniversalClient, l logger.Logger, postService pb.PostServiceServer) *App {
	app = &App{
		config:     c,
		db:         db,
		redis:      rdb,
		logger:     l,
		postService: postService,
	}
	return app
}

func (a *App) Config() *config.Config {
	return a.config
}

func (a *App) DB() *db.DB {
	return a.db
}

func (a *App) Redis() redis.UniversalClient {
	return a.redis
}

func (a *App) Logger() logger.Logger {
	return a.logger
}

func (a *App) PostService() pb.PostServiceServer {
	return a.postService
}
