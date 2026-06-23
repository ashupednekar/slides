package render

import (
	"html"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	attrURLRE = regexp.MustCompile(`(?i)\b(?:src|href|poster)=["']([^"']+)["']`)
	srcsetRE  = regexp.MustCompile(`(?i)\bsrcset=["']([^"']+)["']`)
	cssURLRE  = regexp.MustCompile(`(?i)url\(([^)]+)\)`)
)

func patchStatic(outDir string) error {
	slidesJS := filepath.Join(outDir, "static", "slides.js")
	b, err := os.ReadFile(slidesJS)
	if err != nil {
		return err
	}
	s := strings.ReplaceAll(string(b), "var PERMANENT_URL_PREFIX = '/static/';", "var PERMANENT_URL_PREFIX = 'static/';")
	return os.WriteFile(slidesJS, []byte(s), 0o644)
}

func rewriteStaticHTML(s string) string {
	replacements := []struct {
		old string
		new string
	}{
		{`src='/static/`, `src='static/`},
		{`src="/static/`, `src="static/`},
		{`href='/static/`, `href='static/`},
		{`href="/static/`, `href="static/`},
		{`src='/play.js'`, `src='play.js'`},
		{`src="/play.js"`, `src="play.js"`},
		{`="../`, `="`},
		{`='../`, `='`},
		{`url("../`, `url("`},
		{`url('../`, `url('`},
		{`url(../`, `url(`},
	}
	for _, replacement := range replacements {
		s = strings.ReplaceAll(s, replacement.old, replacement.new)
	}
	return s
}

func copyReferencedAssets(repoRoot, deckDir, outDir, rendered string) error {
	for _, asset := range localAssets(rendered) {
		if strings.HasPrefix(asset, "static/") || asset == "play.js" {
			continue
		}
		if exists(filepath.Join(outDir, filepath.FromSlash(asset))) {
			continue
		}

		var src string
		if exists(filepath.Join(deckDir, filepath.FromSlash(asset))) {
			src = filepath.Join(deckDir, filepath.FromSlash(asset))
		} else if exists(filepath.Join(repoRoot, filepath.FromSlash(asset))) {
			src = filepath.Join(repoRoot, filepath.FromSlash(asset))
		} else {
			continue
		}

		info, err := os.Stat(src)
		if err != nil || info.IsDir() {
			continue
		}
		if err := copyFile(src, filepath.Join(outDir, filepath.FromSlash(asset)), info.Mode()); err != nil {
			return err
		}
	}
	return nil
}

func localAssets(rendered string) []string {
	seen := map[string]bool{}
	add := func(raw string) {
		raw = strings.TrimSpace(html.UnescapeString(raw))
		raw = strings.Trim(raw, `"'`)
		if raw == "" || isExternal(raw) {
			return
		}
		if idx := strings.IndexAny(raw, "?#"); idx >= 0 {
			raw = raw[:idx]
		}
		raw = pathCleanSlash(strings.TrimPrefix(raw, "/"))
		if raw == "." || strings.HasPrefix(raw, "../") {
			return
		}
		seen[raw] = true
	}

	for _, match := range attrURLRE.FindAllStringSubmatch(rendered, -1) {
		add(match[1])
	}
	for _, match := range srcsetRE.FindAllStringSubmatch(rendered, -1) {
		for _, candidate := range strings.Split(match[1], ",") {
			fields := strings.Fields(strings.TrimSpace(candidate))
			if len(fields) > 0 {
				add(fields[0])
			}
		}
	}
	for _, match := range cssURLRE.FindAllStringSubmatch(rendered, -1) {
		add(match[1])
	}

	assets := make([]string, 0, len(seen))
	for asset := range seen {
		assets = append(assets, asset)
	}
	sort.Strings(assets)
	return assets
}

func isExternal(raw string) bool {
	if strings.HasPrefix(raw, "#") ||
		strings.HasPrefix(raw, "mailto:") ||
		strings.HasPrefix(raw, "tel:") ||
		strings.HasPrefix(raw, "data:") ||
		strings.HasPrefix(raw, "//") {
		return true
	}
	u, err := url.Parse(raw)
	return err == nil && u.IsAbs()
}
