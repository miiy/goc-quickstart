package user

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(protected gin.IRouter) {
	api := m.usersAPI

	protectedGroup := protected.Group("/users")
	protectedGroup.GET("", api.ListUsers)
	protectedGroup.GET("/batch", api.BatchGetUsers)
	protectedGroup.GET("/:id", api.GetUser)
	protectedGroup.PUT("/:id", api.UpdateUser)
}
