package main

import (
	"context"
	"flag"
	"log"

	postv1 "github.com/miiy/goc-quickstart/post-service/gen/go/shop/post/v1"
	"github.com/miiy/goc-quickstart/post-service/internal/di"
	postSrv "github.com/miiy/goc-quickstart/post-service/server/post"
	"github.com/miiy/goc/grpc/server"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
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
		nil,
		nil,
	)...)

	// run server
	err = server.Run(ctx, server.Options{
		Network:      "tcp",
		Addr:         app.Config().Server.Grpc.Addr,
		ServerOption: serverOpts,
		RegisterService: func(s server.GRPCServer) {
			postv1.RegisterPostServiceServer(s, postSrv.NewPostServiceServer())
			reflection.Register(s)
		},
	})
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
