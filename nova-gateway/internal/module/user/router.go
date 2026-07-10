package user

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(public, protected gin.IRouter) {
	api := m.usersAPI

	protectedGroup := protected.Group("/users")
	protectedGroup.GET("", api.ListUsers)
	protectedGroup.GET("/:username", api.GetUser)

	protected.GET("/profile", api.GetProfile)
	protected.PUT("/profile", api.UpdateProfile)
}
