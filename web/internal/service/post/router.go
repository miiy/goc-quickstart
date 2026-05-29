package post

import (
	"github.com/miiy/goc-quickstart/web/internal/service/auth"
	"github.com/miiy/goc/gin"
)

func Router(r *gin.Engine) {
	r.GET("/posts", indexHandler)
	r.GET("/posts/pages/:page", pagesHandler)
	r.GET("/posts/:id", showHandler)

	// 需要登录的操作
	authRequired := auth.AuthRequired()
	r.GET("/posts/create", authRequired, createHandler)
	r.POST("/posts", authRequired, storeHandler)
	r.GET("/posts/:id/edit", authRequired, editHandler)
	r.POST("/posts/:id", authRequired, postHandler)
	r.PUT("/posts/:id", authRequired, updateHandler)
	r.DELETE("/posts/:id", authRequired, destroyHandler)
}

func Templates() map[string][]string {
	return map[string][]string{
		"post/detail": {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/detail.html"},
		"post/list":   {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/list.html"},
		"post/create": {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/create.html"},
		"post/edit":   {"layout/layout.html", "layout/header.html", "layout/footer.html", "post/edit.html"},
	}
}
