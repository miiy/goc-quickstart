package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/miiy/goc-quickstart/web/client"
	"github.com/miiy/goc-quickstart/web/internal/config"
	"github.com/miiy/goc-quickstart/web/internal/router"
	"github.com/miiy/goc-quickstart/web/internal/service/auth"
	"github.com/miiy/goc-quickstart/web/internal/service/post"
	"github.com/miiy/goc/gin/sessions"
	"github.com/miiy/goc/logger"
)

func main() {
	conf := flag.String("c", "./configs/default.yaml", "config file")
	flag.Parse()

	cfg, err := config.NewConfig(*conf)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	l, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func() {
		if sync, ok := l.(interface{ Sync() error }); ok {
			sync.Sync()
		}
	}()

	// create HTTP client for gateway
	clients, cleanup, err := client.NewClients(cfg.Gateway.Addr)
	if err != nil {
		log.Fatalf("failed to create clients: %v", err)
	}
	defer cleanup()

	// init modules
	post.NewModule(l, clients.Post)
	auth.NewModule(l, clients.Auth)

	// init session store
	store, err := sessions.NewRedisStore(10, "tcp", cfg.Redis.Addr, cfg.Redis.Password, []byte(cfg.Session.Secret))
	if err != nil {
		log.Fatalf("failed to create session store: %v", err)
	}

	// create gin engine
	engine := gin.Default()
	router.Router(engine, store, cfg.Session.Name)

	log.Printf("Starting server on %s", cfg.Server.HTTP.Addr)
	if err := engine.Run(cfg.Server.HTTP.Addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
