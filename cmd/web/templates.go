package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.alberttseng.net/internal/models"
	"snippetbox.alberttseng.net/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	//return t.Format("2006-01-02 at 15:04:05")
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02 at 15:04:05")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		// if err != nil {
		// 	return nil, err
		// }

		// ts, err = ts.ParseFiles(page)
		// if err != nil {
		// 	return nil, err
		// }

		cache[name] = ts
	}
	return cache, nil
}
