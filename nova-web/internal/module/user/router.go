package user

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(protected gin.IRouter) {
	handler := m.handler

	protectedGroup := protected.Group("/user")
	protectedGroup.GET("/profile", handler.Profile)
	protectedGroup.POST("/profile", handler.UpdateProfile)
	protectedGroup.POST("/avatar", handler.UploadAvatar)
	protectedGroup.PUT("/password", handler.ChangePassword)
}

func Templates() map[string][]string {
	return map[string][]string{
		"user/profile": {"layout/layout.html", "layout/header.html", "layout/footer.html", "user/profile.html"},
	}
}
