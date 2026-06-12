package app

import (
	pb "github.com/miiy/goc-quickstart/upload-service/gen/go/blog/upload/v1"
	"github.com/miiy/goc-quickstart/upload-service/internal/config"
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/redis"
)

type App struct {
	config        *config.Config
	db            *db.DB
	redis         redis.UniversalClient
	logger        logger.Logger
	uploadService pb.UploadServiceServer
}

var app *App

func NewApp(c *config.Config, db *db.DB, rdb redis.UniversalClient, l logger.Logger, uploadService pb.UploadServiceServer) *App {
	app = &App{
		config:        c,
		db:            db,
		redis:         rdb,
		logger:        l,
		uploadService: uploadService,
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

func (a *App) UploadService() pb.UploadServiceServer {
	return a.uploadService
}
