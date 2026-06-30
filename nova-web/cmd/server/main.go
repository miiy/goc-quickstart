package main

import (
	"flag"
	"log"

	"github.com/miiy/goc-quickstart/nova-web/internal/di"
	"github.com/miiy/goc-quickstart/nova-web/internal/router"
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

	log.Printf("Starting server on %s", cfg.Server.HTTP.Addr)
	opts := []httpserver.Option{
		httpserver.WithAddr(cfg.Server.HTTP.Addr),
		httpserver.WithLogger(app.Logger().ZapLogger()),
	}
	if cfg.App.Debug {
		opts = append(opts, httpserver.WithDebug())
	}

	server := httpserver.New(opts...)
	server.RegisterRouter(router.Router(app))
	server.Run()
}
