package main

import (
	"path/filepath"
	"text/template"

	"ptodd.org/snippetbox/pkg/models"
)

// templateData acts as a holding structure for any dynamic data passed to
// HTML templates.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
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
		ts, err := template.ParseFiles(page)
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
