package main

import (
	"context"
	"flag"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	userv1 "github.com/miiy/goc-quickstart/user-service/gen/go/blog/user/v1"
	"github.com/miiy/goc-quickstart/user-service/internal/di"
	grpcauth "github.com/miiy/goc/grpc/interceptor/auth"
	"github.com/miiy/goc/grpc/server"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var protectedMethods = map[string]struct{}{
	userv1.UserService_GetUser_FullMethodName:    {},
	userv1.UserService_UpdateUser_FullMethodName: {},
	userv1.UserService_ListUsers_FullMethodName:  {},
}

func protectedMethodMatcher(ctx context.Context, callMeta interceptors.CallMeta) bool {
	_, ok := protectedMethods[callMeta.FullMethod()]
	return ok
}

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

	serverOpts = append(serverOpts, server.DefaultInterceptor(
		logger,
		grpcauth.MetadataAuthFunc,
		selector.MatchFunc(protectedMethodMatcher),
	)...)

	err = server.Run(ctx, server.Options{
		Network:      "tcp",
		Addr:         app.Config().Server.Grpc.Addr,
		ServerOption: serverOpts,
		RegisterService: func(s server.GRPCServer) {
			healthpb.RegisterHealthServer(s, health.NewServer())
			userv1.RegisterUserServiceServer(s, app.UserService())
			reflection.Register(s)
		},
	})
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
