package render

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

func writeDeckIndex(outDir string, files []string) error {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<meta charset=\"utf-8\">\n<title>Slides</title>\n<ul>\n")
	for _, file := range files {
		fmt.Fprintf(&b, "<li><a href=\"%s\">%s</a></li>\n", html.EscapeString(file), html.EscapeString(file))
	}
	b.WriteString("</ul>\n")
	return os.WriteFile(filepath.Join(outDir, "index.html"), []byte(b.String()), 0o644)
}

func writeIndex(outRoot string, decks []RenderedDeck) error {
	var b strings.Builder
	b.WriteString("<!doctype html>\n<meta charset=\"utf-8\">\n<title>Slides</title>\n<ul>\n")
	for _, deck := range decks {
		label := deck.Name
		if deck.Title != "" {
			label += ": " + deck.Title
		}
		fmt.Fprintf(&b, "<li><a href=\"%s/\">%s</a></li>\n", html.EscapeString(deck.Name), html.EscapeString(label))
	}
	b.WriteString("</ul>\n")
	return os.WriteFile(filepath.Join(outRoot, "index.html"), []byte(b.String()), 0o644)
}
