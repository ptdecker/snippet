package main

import (
	"path/filepath"
	"text/template"
	"time"

	"ptodd.org/snippetbox/pkg/forms"
	"ptodd.org/snippetbox/pkg/models"
)

// templateData acts as a holding structure for any dynamic data passed to
// HTML templates. 'CurrentYear' is an example of common dynamic data
type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	IsAuthenticated bool
}

// Initialize a tempate.FuncMap object for registering custom functions for
// use inside templates
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// humanDate returns a human-friendly formated string representation of a
// time.Time object. This function provides an example of how to use custom
// functions in a template
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// newTemplateCache creates a new template cache
func newTemplateCache(dir string) (map[string]*template.Template, error) {

	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Build a slice of all the 'page' templates
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages
	for _, page := range pages {

		// Extract the file name from the full path
		name := filepath.Base(page)

		// Parse the page tempate file
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add any layouts to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add any partial templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the teamplate set to the cache using the page name as the key
		cache[name] = ts
	}

	return cache, nil
}
