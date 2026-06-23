package render

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/present"
)

func slideTemplate(baseRoot string) (*template.Template, error) {
	tmpl := present.Template().Funcs(template.FuncMap{
		"playable": func(present.Code) bool { return false },
	})
	_, err := tmpl.ParseFS(os.DirFS(baseRoot), "templates/action.tmpl", "templates/slides.tmpl")
	return tmpl, err
}

func renderSlide(tmpl *template.Template, slide string) (*present.Doc, string, error) {
	f, err := os.Open(slide)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	ctx := present.Context{ReadFile: os.ReadFile}
	logOut := log.Writer()
	log.SetOutput(io.Discard)
	doc, err := ctx.Parse(f, slide, 0)
	log.SetOutput(logOut)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	if err := doc.Render(&buf, tmpl); err != nil {
		return nil, "", err
	}
	return doc, buf.String(), nil
}

func outputName(slides []string, slide string) string {
	if len(slides) == 1 {
		return "index.html"
	}
	return strings.TrimSuffix(filepath.Base(slide), filepath.Ext(slide)) + ".html"
}
