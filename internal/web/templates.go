package web

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"log"
	"path/filepath"
)

//go:embed templates/* templates/layouts/*
var content embed.FS

var templateCache = map[string]*template.Template{
	"git-diff.html":    requireTemplate("git-diff.html"),
	"git-tree.html":    requireTemplate("git-tree.html"),
	"git-summary.html": requireTemplate("git-summary.html"),
	"git-commits.html": requireTemplate("git-commits.html"),
}

func RenderTemplate(wr io.Writer, name string, data interface{}) error {
	t, ok := templateCache[name]
	if !ok {
		tt, err := getTemplate(name)
		if err != nil {
			return err
		}
		templateCache[name] = tt
		t = tt
	}
	return renderTemplate(wr, t, data)
}

// GetTemplate returns template by filename, fails if there is any error
// templates are located in env.TemplatePath and can have base layouts
func getTemplate(name string) (*template.Template, error) {
	layoutsPath := filepath.Join("templates", "layouts", "*.html")
	pagePath := filepath.Join("templates", name)
	helpers := TemplateHelpers()
	return template.New(name).Funcs(helpers).ParseFS(content, layoutsPath, pagePath)
}

func requireTemplate(name string) *template.Template {
	t, err := getTemplate(name)
	if err != nil {
		log.Fatal(err.Error())
	}
	return t
}

// RenderTemplate renders precompiled template with data
// result is buffered, in the case of errors, nothing should be written to output
func renderTemplate(wr io.Writer, t *template.Template, data interface{}) error {
	buf := new(bytes.Buffer)
	err := t.ExecuteTemplate(buf, "layout", data)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(wr)
	return err
}
