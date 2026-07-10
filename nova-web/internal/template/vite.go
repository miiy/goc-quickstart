package template

import (
	"encoding/json"
	"fmt"
	"html"
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"strings"
)

type ViteAssetsConfig struct {
	Dev          bool
	DevServerURL string
	ManifestPath string
	StaticPrefix string
}

type ViteAssets struct {
	dev          bool
	devServerURL string
	manifest     map[string]viteManifestEntry
	staticPrefix string
}

type viteManifestEntry struct {
	File    string   `json:"file"`
	CSS     []string `json:"css"`
	Imports []string `json:"imports"`
}

func NewViteAssets(config ViteAssetsConfig) (*ViteAssets, error) {
	staticPrefix := config.StaticPrefix
	if staticPrefix == "" {
		staticPrefix = "/static/"
	}
	if !strings.HasSuffix(staticPrefix, "/") {
		staticPrefix += "/"
	}

	assets := &ViteAssets{
		dev:          config.Dev,
		devServerURL: defaultViteDevServerURL(config.DevServerURL),
		staticPrefix: staticPrefix,
	}
	if assets.dev {
		return assets, nil
	}

	manifestPath := config.ManifestPath
	if manifestPath == "" {
		manifestPath = filepath.Join("dist", ".vite", "manifest.json")
	}
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read vite manifest %s: %w", manifestPath, err)
	}
	if err := json.Unmarshal(data, &assets.manifest); err != nil {
		return nil, fmt.Errorf("parse vite manifest %s: %w", manifestPath, err)
	}

	return assets, nil
}

func (assets *ViteAssets) Client() htmltemplate.HTML {
	if !assets.dev {
		return ""
	}

	refreshURL := assets.devAssetURL("@react-refresh")
	clientURL := assets.devAssetURL("@vite/client")
	return htmltemplate.HTML(fmt.Sprintf(`<script type="module">
import RefreshRuntime from %s
RefreshRuntime.injectIntoGlobalHook(window)
window.$RefreshReg$ = () => {}
window.$RefreshSig$ = () => (type) => type
window.__vite_plugin_react_preamble_installed__ = true
</script>
<script type="module" src="%s"></script>`, jsString(refreshURL), attr(clientURL)))
}

func (assets *ViteAssets) Entry(entry string) (htmltemplate.HTML, error) {
	entry = strings.TrimPrefix(entry, "/")
	if assets.dev {
		return htmltemplate.HTML(fmt.Sprintf(`<script type="module" src="%s"></script>`, attr(assets.devAssetURL(entry)))), nil
	}

	manifestEntry, ok := assets.manifest[entry]
	if !ok {
		return "", fmt.Errorf("vite manifest entry %q not found", entry)
	}

	imports, err := assets.collectImports(manifestEntry, map[string]bool{})
	if err != nil {
		return "", err
	}

	var out strings.Builder
	seenFiles := map[string]bool{}
	for _, key := range imports {
		importEntry := assets.manifest[key]
		assets.writeModulePreload(&out, importEntry.File, seenFiles)
	}
	for _, key := range imports {
		importEntry := assets.manifest[key]
		assets.writeCSS(&out, importEntry.CSS, seenFiles)
	}
	assets.writeCSS(&out, manifestEntry.CSS, seenFiles)
	assets.writeScript(&out, manifestEntry.File)

	return htmltemplate.HTML(out.String()), nil
}

func (assets *ViteAssets) collectImports(entry viteManifestEntry, seen map[string]bool) ([]string, error) {
	var imports []string
	for _, key := range entry.Imports {
		if seen[key] {
			continue
		}
		seen[key] = true

		importEntry, ok := assets.manifest[key]
		if !ok {
			return nil, fmt.Errorf("vite manifest import %q not found", key)
		}

		nested, err := assets.collectImports(importEntry, seen)
		if err != nil {
			return nil, err
		}
		imports = append(imports, nested...)
		imports = append(imports, key)
	}
	return imports, nil
}

func (assets *ViteAssets) writeModulePreload(out *strings.Builder, file string, seen map[string]bool) {
	if file == "" || seen[file] {
		return
	}
	seen[file] = true
	fmt.Fprintf(out, `<link rel="modulepreload" href="%s">`, attr(assets.assetURL(file)))
}

func (assets *ViteAssets) writeCSS(out *strings.Builder, files []string, seen map[string]bool) {
	for _, file := range files {
		if file == "" || seen[file] {
			continue
		}
		seen[file] = true
		fmt.Fprintf(out, `<link rel="stylesheet" href="%s">`, attr(assets.assetURL(file)))
	}
}

func (assets *ViteAssets) writeScript(out *strings.Builder, file string) {
	if file == "" {
		return
	}
	fmt.Fprintf(out, `<script type="module" src="%s"></script>`, attr(assets.assetURL(file)))
}

func (assets *ViteAssets) assetURL(file string) string {
	return assets.staticPrefix + strings.TrimPrefix(file, "/")
}

func (assets *ViteAssets) devAssetURL(file string) string {
	return assets.devServerURL + "/" + strings.TrimPrefix(file, "/")
}

func defaultViteDevServerURL(configured string) string {
	if configured == "" {
		configured = os.Getenv("VITE_DEV_SERVER_URL")
	}
	if configured == "" {
		if port := os.Getenv("VITE_PORT"); port != "" {
			configured = "http://localhost:" + port
		}
	}
	if configured == "" {
		configured = "http://localhost:5173"
	}
	return strings.TrimRight(configured, "/")
}

func attr(value string) string {
	return html.EscapeString(value)
}

func jsString(value string) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return `""`
	}
	return string(encoded)
}
