package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"ptodd.org/snippetbox/pkg/models"
)

// Home page handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Site-wide 404 handler
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create an instance of templateData to hold the retreived snippets
	data := &templateData{Snippets: s}

	// Parse templates
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write template as the respnose body
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// showSnippet handler
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// Extract expected 'id' parameter from query string
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the model's get method to receive a record based upon its ID then
	// return the record or a 404
	s, err := app.snippets.Get(id)
	if err != nil && errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create an instance of templateData struct to hold snippet data
	data := &templateData{Snippet: s}

	// Define templates
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parse the templates
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Execute the templates passing the snippet data
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

// createSnippet handler
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// Guard against non-POST calls to this endpoint
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly, slowly!\n\n- Kobayashi Issa"
	expires := "7"

	// Insert the record through our model and receive back the ID of the new record
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
