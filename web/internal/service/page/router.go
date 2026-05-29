package page

func Templates() map[string][]string {
	return map[string][]string{
		"pages/404": {"layout/layout.html", "layout/header.html", "layout/footer.html", "pages/404.html"},
		"pages/500": {"layout/layout.html", "layout/header.html", "layout/footer.html", "pages/500.html"},
	}
}
