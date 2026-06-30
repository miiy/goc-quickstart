package file

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(protected gin.IRouter) {
	api := m.filesAPI

	protectedGroup := protected.Group("/files")
	protectedGroup.POST("/upload", api.UploadFile)
	protectedGroup.POST("/upload/avatar", api.UploadAvatar)
}
