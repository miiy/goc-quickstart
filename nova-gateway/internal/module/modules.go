package module

import (
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/auth"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/file"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/post"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/module/user"
)

type Modules struct {
	Auth *auth.Module
	Post *post.Module
	File *file.Module
	User *user.Module
}
