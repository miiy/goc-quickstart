package main

import (
	"flag"
	"log"

	"github.com/miiy/goc-quickstart/nova-web/internal/di"
	"github.com/miiy/goc-quickstart/nova-web/internal/router"
	"github.com/miiy/goc-quickstart/nova-web/internal/service/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/service/post"
	"github.com/miiy/goc-quickstart/nova-web/internal/service/user"
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

	// init modules
	clients := app.Clients()
	post.NewModule(app.Logger(), clients.Post, clients.File)

	cfg := app.Config()
	auth.NewModule(app.Logger(), clients.Auth, app.SessionStore(), cfg.Session.Name)
	user.NewModule(app.Logger(), clients.Auth, clients.User, clients.File, app.SessionStore(), cfg.Session.Name)

	log.Printf("Starting server on %s", cfg.Server.HTTP.Addr)
	server := httpserver.New(
		httpserver.WithAddr(cfg.Server.HTTP.Addr),
		httpserver.WithLogger(app.Logger().ZapLogger()),
	)
	server.Use(ginzap.Ginzap(app.Logger().ZapLogger()), ginzap.RecoveryWithZap(app.Logger().ZapLogger(), true))
	server.RegisterRouter(func(r *gin.Engine) {
		router.Router(r, app.SessionStore(), cfg.Session.Name, cfg.App.Timezone)
	})
	server.Run()
}
