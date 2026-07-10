package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestViteAssetsDevOutput(t *testing.T) {
	assets, err := NewViteAssets(ViteAssetsConfig{
		Dev:          true,
		DevServerURL: "http://localhost:3000/",
	})
	if err != nil {
		t.Fatalf("NewViteAssets() error = %v", err)
	}

	client := string(assets.Client())
	if !strings.Contains(client, `http://localhost:3000/@react-refresh`) {
		t.Fatalf("expected react refresh url, got %s", client)
	}
	if !strings.Contains(client, `http://localhost:3000/@vite/client`) {
		t.Fatalf("expected vite client url, got %s", client)
	}

	entry, err := assets.Entry("frontend/src/features/post/entries/create.tsx")
	if err != nil {
		t.Fatalf("Entry() error = %v", err)
	}
	if got := string(entry); got != `<script type="module" src="http://localhost:3000/frontend/src/features/post/entries/create.tsx"></script>` {
		t.Fatalf("Entry() = %s", got)
	}
}

func TestViteAssetsManifestOutput(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "manifest.json")
	manifest := `{
		"frontend/src/app.ts": {
			"file": "assets/app-BbYx.js",
			"css": ["assets/app-AaXx.css"],
			"imports": ["_react.js"]
		},
		"_react.js": {
			"file": "assets/react-CcZz.js"
		}
	}`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	assets, err := NewViteAssets(ViteAssetsConfig{
		ManifestPath: manifestPath,
		StaticPrefix: "/static",
	})
	if err != nil {
		t.Fatalf("NewViteAssets() error = %v", err)
	}

	entry, err := assets.Entry("frontend/src/app.ts")
	if err != nil {
		t.Fatalf("Entry() error = %v", err)
	}
	got := string(entry)
	for _, want := range []string{
		`<link rel="modulepreload" href="/static/assets/react-CcZz.js">`,
		`<link rel="stylesheet" href="/static/assets/app-AaXx.css">`,
		`<script type="module" src="/static/assets/app-BbYx.js"></script>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("Entry() = %s, missing %s", got, want)
		}
	}
}
