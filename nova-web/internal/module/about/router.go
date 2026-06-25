package about

import "github.com/miiy/goc/gin"

func Router(r *gin.Engine) {
	r.GET("/about", indexHandler)
}

func Templates() map[string][]string {
	return map[string][]string{
		"about/index": {"layout/layout.html", "layout/header.html", "layout/footer.html", "about/index.html"},
	}
}
