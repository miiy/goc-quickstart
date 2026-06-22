package main

import (
	"flag"
	"log"

	"github.com/miiy/goc-quickstart/nova-gateway/internal/di"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/router"
	ginzap "github.com/miiy/goc/gin/middleware/zap"
	httpserver "github.com/miiy/goc/http/server"
)

func main() {
	conf := flag.String("c", "./config.yaml", "config file")
	flag.Parse()

	app, cleanup, err := di.InitApp(*conf)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}
	defer cleanup()

	zapLogger := app.Logger().ZapLogger()
	server := httpserver.New(
		httpserver.WithAddr(app.Config().Server.HTTP.Addr),
		httpserver.WithLogger(zapLogger),
	)
	server.Use(ginzap.Ginzap(zapLogger), ginzap.RecoveryWithZap(zapLogger, true))
	server.RegisterRouter(router.Router(app))
	server.Run()
}
