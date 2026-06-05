package auth

import authv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/auth/v1"

type Module struct {
	client authv1.AuthServiceClient
}

func NewModule(client authv1.AuthServiceClient) *Module {
	return &Module{client: client}
}
