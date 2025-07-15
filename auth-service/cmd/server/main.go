package main

import (
	"context"
	"flag"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	authpb "github.com/miiy/goc-quickstart/auth-service/gen/go/shop/auth/v1"
	"github.com/miiy/goc-quickstart/auth-service/internal/auth"
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

	// set logger
	logger := app.Logger().ZapLogger()
	grpclog.SetLoggerV2(zapgrpc.NewLogger(logger))

	// grpc server options
	var serverOpts []grpc.ServerOption
	// mTLS
	serverOpts = append(serverOpts,
		server.WithMTLS(
			config.Server.Grpc.Tls.CertFile,
			config.Server.Grpc.Tls.KeyFile,
			config.Server.Grpc.Tls.CaFile,
		),
	)
	// interceptor
	serverOpts = append(serverOpts, server.DefaultInterceptor(
		logger,
		auth.AuthFunc(app.JWTAuth(), app.UserProvider()),
		selector.MatchFunc(auth.AuthMatchFunc),
	)...)

	// run server
	err = server.Run(ctx, server.Options{
		Network:      "tcp",
		Addr:         app.Config().Server.Grpc.Addr,
		ServerOption: serverOpts,
		RegisterService: func(s server.GRPCServer) {
			healthpb.RegisterHealthServer(s, health.NewServer())
			authpb.RegisterAuthServer(s, app.AuthServer())
			reflection.Register(s)
		},
	})
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
