package home

import "github.com/miiy/goc/gin"

func RegisterRouter(r gin.IRouter) {
	handler := NewHomeHandler()

	r.GET("/", handler.index)
}

func Templates() map[string][]string {
	return map[string][]string{
		"home/index": {"layout/layout.html", "layout/header.html", "layout/footer.html", "home/index.html"},
	}
}
