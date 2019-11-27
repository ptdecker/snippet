package main

import "ptodd.org/snippetbox/pkg/models"

// templateData acts as a holding structure for any dynamic data passed to
// HTML templates.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
