package user

import (
	resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"
	"github.com/miiy/goc/gin"
)

func (m *Module) RegisterRouter(_ gin.IRouter, protected gin.IRouter) {
	handler := m.handler

	protectedUsers := protected.Group("/users")
	protectedUsers.GET("/:username", handler.Show)

	protectedProfile := protected.Group("/profile")
	protectedProfile.GET("", handler.Profile)
}

func Templates() map[string][]string {
	return map[string][]string{
		"user/profile": resourceTemplate.Layout("user/profile.html"),
		"user/show":    resourceTemplate.Layout("user/show.html"),
	}
}
