package about

import "github.com/miiy/goc/gin"

func RegisterRouter(r gin.IRouter) {
	handler := NewAboutHandler()

	r.GET("/about", handler.index)
}

func Templates() map[string][]string {
	return map[string][]string{
		"about/index": {"layout/layout.html", "layout/header.html", "layout/footer.html", "about/index.html"},
	}
}
