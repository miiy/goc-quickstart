package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/miiy/goc/grpc/gateway"

	"github.com/miiy/goc-quickstart/api-gateway/internal/config"
	authv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/auth/v1"
	postv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/post/v1"
	userv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/user/v1"
)

func main() {
	conf := flag.String("c", "./configs/default.yaml", "config file")
	flag.Parse()

	config, err := config.Load(*conf)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Build services config
	services := []gateway.ServiceConfig{
		{Name: "auth", Endpoint: gateway.Endpoint{Addr: config.Services["auth"].Endpoint}, Register: authv1.RegisterAuthServiceHandler},
		{Name: "post", Endpoint: gateway.Endpoint{Addr: config.Services["post"].Endpoint}, Register: postv1.RegisterPostServiceHandler},
		{Name: "user", Endpoint: gateway.Endpoint{Addr: config.Services["user"].Endpoint}, Register: userv1.RegisterUserServiceHandler},
	}

	// Build options
	opts := gateway.Options{
		Addr:     config.Server.HTTP.Addr,
		Services: services,
	}

	// TLS config
	if config.TLS.Enabled {
		opts.TLSConfig = gateway.MTLSConfig(
			config.TLS.ServerName,
			config.TLS.CertFile,
			config.TLS.KeyFile,
			config.TLS.CaFile,
		)
	}

	// Context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := gateway.Run(ctx, opts); err != nil {
		log.Fatal(err)
	}
}
