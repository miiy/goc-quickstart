package main

import (
	"context"
	"flag"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	userv1 "github.com/miiy/goc-quickstart/user-service/gen/go/blog/user/v1"
	"github.com/miiy/goc-quickstart/user-service/internal/di"
	"github.com/miiy/goc/grpc/server"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var NoOpAuthFunc = func(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

var NoOpMatcher = func(ctx context.Context, callMeta interceptors.CallMeta) bool {
	return false
}

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

	serverOpts = append(serverOpts, server.DefaultInterceptor(
		logger,
		NoOpAuthFunc,
		selector.MatchFunc(NoOpMatcher),
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
