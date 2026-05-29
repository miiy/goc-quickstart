package main

import (
	"context"
	"flag"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	authpb "github.com/miiy/goc-quickstart/auth-service/gen/go/blog/auth/v1"
	"github.com/miiy/goc-quickstart/auth-service/internal/middleware"
	"github.com/miiy/goc-quickstart/auth-service/internal/di"
	"github.com/miiy/goc/grpc/server"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf := flag.String("c", "./configs/default.yaml", "config file")
	flag.Parse()

	ctx := context.Background()
	app, cleanup, err := di.InitApp(*conf)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	config := app.Config()

	logger := app.Logger().ZapLogger()
	grpclog.SetLoggerV2(zapgrpc.NewLogger(logger))

	var serverOpts []grpc.ServerOption

	// TLS
	if config.Server.Grpc.Tls.CertFile != "" && config.Server.Grpc.Tls.KeyFile != "" {
		tlsOpt, err := server.WithMTLS(
			config.Server.Grpc.Tls.CertFile,
			config.Server.Grpc.Tls.KeyFile,
			config.Server.Grpc.Tls.CaFile,
		)
		if err != nil {
			log.Fatalf("failed to configure mTLS: %v", err)
		}
		serverOpts = append(serverOpts, tlsOpt)
	}

	// interceptor
	serverOpts = append(serverOpts, server.DefaultInterceptor(
		logger,
		middleware.AuthFunc(app.JWTAuth(), app.UserProvider()),
		selector.MatchFunc(middleware.AuthMatchFunc),
	)...)

	err = server.Run(ctx, server.Options{
		Network:      "tcp",
		Addr:         app.Config().Server.Grpc.Addr,
		ServerOption: serverOpts,
		RegisterService: func(s server.GRPCServer) {
			healthpb.RegisterHealthServer(s, health.NewServer())
			authpb.RegisterAuthServiceServer(s, app.AuthServiceServer())
			reflection.Register(s)
		},
	})
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
