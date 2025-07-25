// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/miiy/goc-quickstart/auth-service/internal/app"
	"github.com/miiy/goc-quickstart/auth-service/internal/config"
	"github.com/miiy/goc-quickstart/auth-service/internal/repository"
	"github.com/miiy/goc-quickstart/auth-service/server"
	"github.com/miiy/goc/auth/jwt"
	"github.com/miiy/goc/contrib/wechat/miniprogram"
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/redis"
	"go.uber.org/zap"
)

// Injectors from wire.go:

func InitApp(conf string) (*app.App, func(), error) {
	configConfig, err := config.NewConfig(conf)
	if err != nil {
		return nil, nil, err
	}
	dbConfig := provideDBConfig(configConfig)
	v := provideDBOption(configConfig)
	dbDB, err := db.NewDB(dbConfig, v...)
	if err != nil {
		return nil, nil, err
	}
	options := provideRedisOptions(configConfig)
	universalClient, err := redis.NewRedis(options)
	if err != nil {
		return nil, nil, err
	}
	v2 := provideLoggerOption()
	loggerLogger, err := logger.NewLogger(v2...)
	if err != nil {
		return nil, nil, err
	}
	jwtOptions := provideJwtAuthOptions(configConfig)
	jwtAuth := jwt.NewJWTAuth(jwtOptions)
	zapLogger := provideZap(loggerLogger)
	gormDB := provideGorm(dbDB)
	authRepository := repository.NewAuthRepository(gormDB)
	tokenRepository := repository.NewTokenRepository(universalClient)
	miniProgram, err := provideMiniProgram(configConfig)
	if err != nil {
		return nil, nil, err
	}
	authServer := server.NewAuthServiceServer(zapLogger, authRepository, tokenRepository, jwtAuth, miniProgram)
	appApp := app.NewApp(configConfig, dbDB, universalClient, loggerLogger, jwtAuth, authServer, authRepository)
	return appApp, func() {
	}, nil
}

// wire.go:

func provideLoggerOption() []logger.Option {
	return nil
}

func provideZap(logger2 logger.Logger) *zap.Logger {
	return logger2.ZapLogger()
}

func provideDBConfig(config2 *config.Config) db.Config {
	return db.Config{
		Driver:   config2.Database.Driver,
		Host:     config2.Database.Host,
		Port:     config2.Database.Port,
		Username: config2.Database.Username,
		Password: config2.Database.Password,
		Database: config2.Database.Database,
	}
}

func provideDBOption(config2 *config.Config) []db.Option {
	return nil
}

func provideRedisOptions(config2 *config.Config) *redis.Options {
	return &redis.Options{
		Addrs:    config2.Redis.Addrs,
		Password: config2.Redis.Password,
		DB:       config2.Redis.DB,
	}
}

func provideGorm(db2 *db.DB) *gorm.DB {
	return db2.Gorm()
}

func provideJwtAuthOptions(config2 *config.Config) *jwt.Options {
	return &jwt.Options{
		Secret:    config2.Jwt.Secret,
		Issuer:    config2.Jwt.Issuer,
		ExpiresIn: config2.Jwt.ExpiresIn,
	}
}

func provideMiniProgram(config2 *config.Config) (*miniprogram.MiniProgram, error) {
	return nil, nil
}
