package auth

import (
	openapi "github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server/go"
	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
)

type Module struct {
	authAPI openapi.AuthAPI
}

func NewModule(authClient authv1.AuthServiceClient) *Module {
	return &Module{authAPI: NewAuthAPI(authClient)}
}
