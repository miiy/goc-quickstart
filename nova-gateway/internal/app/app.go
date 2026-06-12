package app

import (
	"errors"
	"fmt"
	"strings"

	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	filev1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/file/v1"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/config"
	"github.com/miiy/goc/auth"
	goccredentials "github.com/miiy/goc/grpc/credentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	config  *config.Config
	jwtAuth *auth.JWTAuth

	authClient authv1.AuthServiceClient
	postClient postv1.PostServiceClient
	fileClient filev1.FileServiceClient
	userClient userv1.UserServiceClient

	conns []*grpc.ClientConn
}

func NewApp(cfg *config.Config) (_ *App, err error) {
	creds, err := transportCredentials(cfg)
	if err != nil {
		return nil, err
	}
	var conns []*grpc.ClientConn
	defer func() {
		if err != nil {
			_ = closeConns(conns)
		}
	}()

	authConn, err := dialService(cfg, "auth", creds)
	if err != nil {
		return nil, err
	}
	conns = append(conns, authConn)

	postConn, err := dialService(cfg, "post", creds)
	if err != nil {
		return nil, err
	}
	conns = append(conns, postConn)

	userConn, err := dialService(cfg, "user", creds)
	if err != nil {
		return nil, err
	}
	conns = append(conns, userConn)

	fileConn, err := dialService(cfg, "file", creds)
	if err != nil {
		return nil, err
	}
	conns = append(conns, fileConn)

	return &App{
		config: cfg,
		jwtAuth: auth.NewJWTAuth(&auth.Options{
			Secret: cfg.JWT.Secret,
			Issuer: cfg.JWT.Issuer,
		}),
		authClient: authv1.NewAuthServiceClient(authConn),
		postClient: postv1.NewPostServiceClient(postConn),
		fileClient: filev1.NewFileServiceClient(fileConn),
		userClient: userv1.NewUserServiceClient(userConn),
		conns:      conns,
	}, nil
}

func (a *App) Config() *config.Config {
	return a.config
}

func (a *App) JWTAuth() *auth.JWTAuth {
	return a.jwtAuth
}

func (a *App) AuthClient() authv1.AuthServiceClient {
	return a.authClient
}

func (a *App) PostClient() postv1.PostServiceClient {
	return a.postClient
}

func (a *App) UserClient() userv1.UserServiceClient {
	return a.userClient
}

func (a *App) FileClient() filev1.FileServiceClient {
	return a.fileClient
}

func (a *App) Close() error {
	if a == nil {
		return nil
	}
	return closeConns(a.conns)
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
