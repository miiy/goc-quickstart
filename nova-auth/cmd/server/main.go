package main

import (
	"context"
	"flag"
	"log"

	authpb "github.com/miiy/goc-quickstart/nova-auth/gen/go/nova/auth/v1"
	"github.com/miiy/goc-quickstart/nova-auth/internal/di"
	grpcauth "github.com/miiy/goc/grpc/interceptor/auth"
	"github.com/miiy/goc/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf := flag.String("c", "./config.yaml", "config file")
	flag.Parse()

	ctx := context.Background()
	app, cleanup, err := di.InitApp(*conf)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	config := app.Config()

	logger := app.Logger().ZapLogger()

	var serverOpts []grpc.ServerOption

	// TLS
	if config.Server.Grpc.Tls.Enabled {
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

	serverOpts = append(serverOpts, server.DefaultInterceptor(
		logger,
		grpcauth.MetadataAuthFunc,
		grpcauth.MatchFullMethods(
			authpb.AuthService_RefreshToken_FullMethodName,
			authpb.AuthService_ChangePassword_FullMethodName,
			authpb.AuthService_Logout_FullMethodName,
		),
	)...)

	err = server.Run(ctx, server.Options{
		Network:      "tcp",
		Addr:         app.Config().Server.Grpc.Addr,
		Logger:       logger,
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
