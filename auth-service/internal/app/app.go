package app

import (
	authpb "github.com/miiy/goc-quickstart/auth-service/gen/go/shop/auth/v1"
	"github.com/miiy/goc-quickstart/auth-service/internal/config"
	"github.com/miiy/goc/auth"
	"github.com/miiy/goc/auth/jwt"
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/redis"
)

type App struct {
	config       *config.Config
	db           *db.DB
	redis        redis.UniversalClient
	logger       logger.Logger
	jwtAuth      *jwt.JWTAuth
	authServer   authpb.AuthServer
	userProvider auth.UserProvider
}

var app *App

func NewApp(c *config.Config, db *db.DB, rdb redis.UniversalClient, l logger.Logger, j *jwt.JWTAuth,
	as authpb.AuthServer, up auth.UserProvider) *App {
	app = &App{
		config:       c,
		db:           db,
		redis:        rdb,
		logger:       l,
		jwtAuth:      j,
		authServer:   as,
		userProvider: up,
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

func (a *App) JWTAuth() *jwt.JWTAuth {
	return a.jwtAuth
}

func (a *App) AuthServer() authpb.AuthServer {
	return a.authServer
}

func (a *App) UserProvider() auth.UserProvider {
	return a.userProvider
}
