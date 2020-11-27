package server

import (
	"html/template"
	"path/filepath"
	"bytes"
	"io"
	"log"
)

// GetTemplate returns template by filename, fails if there is any error
// templates are located in env.TemplatePath and can have base layouts
func (env *Env) GetTemplate(name string, helpers template.FuncMap) *template.Template {
	pagePath := filepath.Join(env.TemplatePath, name)
	layoutsPath := filepath.Join(env.TemplatePath, "layouts", "*.html")
	layouts, err := filepath.Glob(layoutsPath)
	if err != nil {
		log.Fatal(err)
	}
	files := append(layouts, pagePath)
	t, err := template.New(name).Funcs(helpers).ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

// RenderTemplate renders precompiled template with data
// result is buffered, in the case of errors, nothing should be written to output
func (env *Env) RenderTemplate(wr io.Writer, t *template.Template, data interface{}) error {
	buf := new(bytes.Buffer)
	err := t.ExecuteTemplate(buf, "layout", data)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(wr)
	return err
}
