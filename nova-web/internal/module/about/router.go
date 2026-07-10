package about

import (
	resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"
	"github.com/miiy/goc/gin"
)

func RegisterRouter(r gin.IRouter) {
	handler := NewAboutHandler()

	r.GET("/about", handler.index)
}

func Templates() map[string][]string {
	return map[string][]string{
		"about/index": resourceTemplate.Layout("about/index.html"),
	}
}
