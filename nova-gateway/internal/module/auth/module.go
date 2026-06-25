package auth

import authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"

type Module struct {
	authClient authv1.AuthServiceClient
}

func NewModule(authClient authv1.AuthServiceClient) *Module {
	return &Module{authClient: authClient}
}
