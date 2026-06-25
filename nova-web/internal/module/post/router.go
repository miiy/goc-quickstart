package post

import (
	"github.com/miiy/goc/gin"
	"github.com/miiy/goc/gin/middleware/sessionauth"
)

func Router(r *gin.Engine) {
	public := r.Group("")
	{
		public.GET("/posts", indexHandler)
		public.GET("/posts/pages/:page", pagesHandler)
		public.GET("/posts/:id", showHandler)
	}

	private := r.Group("")
	private.Use(sessionauth.Authenticate(sessionauth.WithRedirect("/login")))
	{
		private.GET("/posts/create", createHandler)
		private.POST("/posts", storeHandler)
		private.GET("/posts/:id/edit", editHandler)
		private.POST("/posts/:id", postHandler)
		private.PUT("/posts/:id", updateHandler)
		private.DELETE("/posts/:id", destroyHandler)
	}
}

func Templates() map[string][]string {
	return map[string][]string{
		"post/detail": {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/detail.html"},
		"post/list":   {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/list.html"},
		"post/create": {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/create.html"},
		"post/edit":   {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/edit.html"},
	}
}
