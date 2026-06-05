package main

import (
	"flag"
	"log"

	"github.com/miiy/goc-quickstart/api-gateway/internal/app"
	"github.com/miiy/goc-quickstart/api-gateway/internal/config"
	"github.com/miiy/goc-quickstart/api-gateway/internal/router"
	ginzap "github.com/miiy/goc/gin/middleware/zap"
	httpserver "github.com/miiy/goc/http/server"
	"github.com/miiy/goc/logger"
)

func main() {
	conf := flag.String("c", "./config.yaml", "config file")
	flag.Parse()

	cfg, err := config.Load(*conf)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("failed to close grpc clients: %v", err)
		}
	}()

	l, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func() {
		if sync, ok := l.(interface{ Sync() error }); ok {
			_ = sync.Sync()
		}
	}()

	server := httpserver.New(
		httpserver.WithAddr(cfg.Server.HTTP.Addr),
		httpserver.WithLogger(l.ZapLogger()),
	)
	server.Use(ginzap.Ginzap(l.ZapLogger()), ginzap.RecoveryWithZap(l.ZapLogger(), true))
	server.RegisterRouter(router.Router(application))
	server.Run()
}
