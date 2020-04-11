package handlers

import (
	"errors"
	"html/template"
	"log"
	"path/filepath"
)

var (
	errNoTemplateSetup  = errors.New("No template setup")
	errTemplateNotFound = errors.New("Template not found")
)

type TemplateMap map[string]*template.Template

type TemplateConfig struct {
	Path  string
	cache TemplateMap
}

func (tc *TemplateConfig) GetTemplate(name string) (*template.Template, error) {
	if tc.cache == nil {
		return nil, errNoTemplateSetup
	}
	t, ok := tc.cache[name]
	if !ok {
		return nil, errTemplateNotFound
	}
	return t, nil
}

func (tc *TemplateConfig) Setup() {
	if tc.cache == nil {
		tc.cache = make(TemplateMap)
	}
	layoutsPath := filepath.Join(tc.Path, "layouts", "*.html")
	layouts, err := filepath.Glob(layoutsPath)
	if err != nil {
		log.Fatal(err)
	}

	pagesPath := filepath.Join(tc.Path, "*.html")
	pages, err := filepath.Glob(pagesPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, page := range pages {
		files := append(layouts, page)
		name := filepath.Base(page)
		t, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}
		tc.cache[name] = t
	}
}
