package template

import "embed"

//go:embed *
var FS embed.FS

func Layout(files ...string) []string {
	layoutFiles := []string{
		"layout/layout.html",
		"layout/header.html",
		"layout/footer.html",
		"layout/footer_custom.html",
	}
	return append(layoutFiles, files...)
}
