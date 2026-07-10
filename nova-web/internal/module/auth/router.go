package auth

import (
	resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"
	"github.com/miiy/goc/gin"
)

func (m *Module) RegisterRouter(r gin.IRouter) {
	handler := m.handler

	r.GET("/register", handler.RegisterForm)
	r.POST("/register", handler.Register)
	r.GET("/login", handler.LoginForm)
	r.POST("/login", handler.Login)

	g := r.Group("/auth")
	{
		g.GET("/logout", handler.Logout)
		g.POST("/logout", handler.Logout)
	}
}

func Templates() map[string][]string {
	return map[string][]string{
		"auth/register": resourceTemplate.Layout("auth/register.html"),
		"auth/login":    resourceTemplate.Layout("auth/login.html"),
	}
}
