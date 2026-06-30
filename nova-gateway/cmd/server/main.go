package main

import (
	"flag"
	"log"

	"github.com/miiy/goc-quickstart/nova-gateway/internal/di"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/router"
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

	cfg := app.Config()
	zapLogger := app.Logger().ZapLogger()
	opts := []httpserver.Option{
		httpserver.WithAddr(cfg.Server.HTTP.Addr),
		httpserver.WithLogger(zapLogger),
	}
	if cfg.App.Debug {
		opts = append(opts, httpserver.WithDebug())
	}

	server := httpserver.New(opts...)
	server.RegisterRouter(router.Router(app))
	server.Run()
}
