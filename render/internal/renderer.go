package render

import (
	"html/template"
	"os"
	"path/filepath"
)

func renderDeck(cfg Config, tmpl *template.Template, d deck) (RenderedDeck, error) {
	outDir := filepath.Join(cfg.OutRoot, filepath.FromSlash(d.name))
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return RenderedDeck{}, err
	}
	if err := copyDir(d.dir, outDir); err != nil {
		return RenderedDeck{}, err
	}
	if err := copyDir(filepath.Join(cfg.BaseRoot, "static"), filepath.Join(outDir, "static")); err != nil {
		return RenderedDeck{}, err
	}
	if err := patchStatic(outDir); err != nil {
		return RenderedDeck{}, err
	}

	var firstTitle string
	written := make([]string, 0, len(d.slides))
	for _, slide := range d.slides {
		doc, htmlBody, err := renderSlide(tmpl, slide)
		if err != nil {
			return RenderedDeck{}, err
		}
		if firstTitle == "" {
			firstTitle = doc.Title
		}

		htmlBody = rewriteStaticHTML(htmlBody)
		if err := copyReferencedAssets(cfg.RepoRoot, d.dir, outDir, htmlBody); err != nil {
			return RenderedDeck{}, err
		}

		outName := outputName(d.slides, slide)
		outPath := filepath.Join(outDir, outName)
		if err := os.WriteFile(outPath, []byte(htmlBody), 0o644); err != nil {
			return RenderedDeck{}, err
		}
		written = append(written, outName)
	}

	if len(written) > 1 {
		if err := writeDeckIndex(outDir, written); err != nil {
			return RenderedDeck{}, err
		}
	}
	if err := offloadLargeAssets(cfg, outDir); err != nil {
		return RenderedDeck{}, err
	}
	return RenderedDeck{Name: d.name, Title: firstTitle, OutDir: outDir}, nil
}
