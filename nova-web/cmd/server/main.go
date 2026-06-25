package main

import (
	"flag"
	"log"

	"github.com/miiy/goc-quickstart/nova-web/internal/di"
	"github.com/miiy/goc-quickstart/nova-web/internal/router"
	"github.com/miiy/goc/gin"
	ginzap "github.com/miiy/goc/gin/middleware/zap"
	httpserver "github.com/miiy/goc/http/server"
)

func main() {
	conf := flag.String("c", "./config.yaml", "config file")
	flag.Parse()

	app, cleanup, err := di.InitApp(*conf)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer cleanup()

	cfg := app.Config()
	clients := app.Clients()
	sessionManager := app.SessionManager()

	log.Printf("Starting server on %s", cfg.Server.HTTP.Addr)
	server := httpserver.New(
		httpserver.WithAddr(cfg.Server.HTTP.Addr),
		httpserver.WithLogger(app.Logger().ZapLogger()),
	)
	server.Use(ginzap.Ginzap(app.Logger().ZapLogger()), ginzap.RecoveryWithZap(app.Logger().ZapLogger(), true))
	server.RegisterRouter(func(r *gin.Engine) {
		router.Router(r, sessionManager, cfg.App, clients.Auth, app.Logger(), cfg.Storage.Root)
	})
	server.Run()
}
