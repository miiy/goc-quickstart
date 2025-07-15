//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/miiy/goc-quickstart/post-service/internal/app"
	"github.com/miiy/goc-quickstart/post-service/internal/config"
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/redis"
	"go.uber.org/zap"
)

func InitApp(conf string) (*app.App, func(), error) {
	panic(wire.Build(
		config.NewConfig,
		wire.NewSet(logger.NewLogger, provideLoggerOption, provideZap),
		wire.NewSet(db.NewDB, provideDBConfig, provideDBOption, provideGorm),
		wire.NewSet(redis.NewRedis, provideRedisOptions),
		app.NewApp,
	))
}

func provideLoggerOption() []logger.Option {
	return nil
}

func provideZap(logger logger.Logger) *zap.Logger {
	return logger.ZapLogger()
}

func provideDBConfig(config *config.Config) db.Config {
	return db.Config{
		Driver:   config.Database.Driver,
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Username: config.Database.Username,
		Password: config.Database.Password,
		Database: config.Database.Database,
	}
}

func provideDBOption(config *config.Config) []db.Option {
	return nil
}

func provideRedisOptions(config *config.Config) *redis.Options {
	return &redis.Options{
		Addrs:    config.Redis.Addrs,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	}
}

func provideGorm(db *db.DB) *gorm.DB {
	return db.Gorm()
}
