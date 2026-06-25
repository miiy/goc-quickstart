package client

import (
	"errors"
	"fmt"
	"strings"

	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	filev1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/file/v1"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/config"
	goccredentials "github.com/miiy/goc/grpc/credentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Auth authv1.AuthServiceClient
	Post postv1.PostServiceClient
	File filev1.FileServiceClient
	User userv1.UserServiceClient
}

func NewClients(cfg *config.Config) (_ *Clients, cleanup func(), err error) {
	creds, err := transportCredentials(cfg)
	if err != nil {
		return nil, nil, err
	}
	var conns []*grpc.ClientConn
	defer func() {
		if err != nil {
			_ = closeConns(conns)
		}
	}()

	authConn, err := dialService(cfg, "auth", creds)
	if err != nil {
		return nil, nil, err
	}
	conns = append(conns, authConn)

	postConn, err := dialService(cfg, "post", creds)
	if err != nil {
		return nil, nil, err
	}
	conns = append(conns, postConn)

	userConn, err := dialService(cfg, "user", creds)
	if err != nil {
		return nil, nil, err
	}
	conns = append(conns, userConn)

	fileConn, err := dialService(cfg, "file", creds)
	if err != nil {
		return nil, nil, err
	}
	conns = append(conns, fileConn)

	clients := &Clients{
		Auth: authv1.NewAuthServiceClient(authConn),
		Post: postv1.NewPostServiceClient(postConn),
		File: filev1.NewFileServiceClient(fileConn),
		User: userv1.NewUserServiceClient(userConn),
	}
	return clients, func() {
		_ = closeConns(conns)
	}, nil
}

func closeConns(conns []*grpc.ClientConn) error {
	var errs []error
	for _, conn := range conns {
		if err := conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func dialService(cfg *config.Config, name string, creds credentials.TransportCredentials) (*grpc.ClientConn, error) {
	svc, ok := cfg.Services[name]
	if !ok || strings.TrimSpace(svc.Endpoint) == "" {
		return nil, fmt.Errorf("missing %s service endpoint", name)
	}

	conn, err := grpc.NewClient(svc.Endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("dial %s service %s: %w", name, svc.Endpoint, err)
	}
	return conn, nil
}

func transportCredentials(cfg *config.Config) (credentials.TransportCredentials, error) {
	if !cfg.TLS.Enabled {
		return insecure.NewCredentials(), nil
	}

	return goccredentials.NewClientMTLS(
		cfg.TLS.ServerName,
		cfg.TLS.CertFile,
		cfg.TLS.KeyFile,
		cfg.TLS.CaFile,
	)
}
