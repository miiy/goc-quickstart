package module

import (
	"github.com/miiy/goc-quickstart/nova-web/internal/module/auth"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/post"
	"github.com/miiy/goc-quickstart/nova-web/internal/module/user"
)

type Modules struct {
	Post *post.PostModule
	Auth *auth.AuthModule
	User *user.Module
}
