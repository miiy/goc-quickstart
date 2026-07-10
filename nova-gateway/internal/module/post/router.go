package post

import "github.com/miiy/goc/gin"

func (m *Module) RegisterRouter(public, protected gin.IRouter) {
	api := m.postsAPI

	public.GET("/categories", api.ListCategories)

	publicGroup := public.Group("/posts")
	publicGroup.GET("", api.ListPosts)
	publicGroup.GET("/:id", api.GetPost)

	protectedGroup := protected.Group("/posts")
	protectedGroup.POST("", api.CreatePost)
	protectedGroup.PUT("/:id", api.UpdatePost)
	protectedGroup.DELETE("/:id", api.DeletePost)

	protected.GET("/users/:username/posts", api.ListUserPosts)
	protected.GET("/users/:username/posts/:id", api.GetUserPost)
}
