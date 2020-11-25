package server

import (
	"html/template"
	"path/filepath"
	"bytes"
	"io"
)

func (env *Env) GetTemplate(name string) (*template.Template, error) {
	pagePath := filepath.Join(env.TemplatePath, name)
	layoutsPath := filepath.Join(env.TemplatePath, "layouts", "*.html")
	layouts, err := filepath.Glob(layoutsPath)
	if err != nil {
		return nil, err
	}
	files := append(layouts, pagePath)
	t, err := template.New(name).Funcs(env.Helpers).ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (env *Env) RenderTemplate(wr io.Writer, t *template.Template, data interface{}) error {
	buf := new(bytes.Buffer)
	err := t.Funcs(env.Helpers).ExecuteTemplate(buf, "layout", data)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(wr)
	return err
}
