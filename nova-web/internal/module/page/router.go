package page

import resourceTemplate "github.com/miiy/goc-quickstart/nova-web/resources/template"

func Templates() map[string][]string {
	return map[string][]string{
		"pages/404": resourceTemplate.Layout("pages/404.html"),
		"pages/500": resourceTemplate.Layout("pages/500.html"),
	}
}
