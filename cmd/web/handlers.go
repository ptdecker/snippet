package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"ptodd.org/snippetbox/pkg/models"
)

// Home page handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// showSnippet handler
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	// Extract expected 'id' parameter from query string
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

// createSnippet handler
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// Add any data in POST (works for PUT and PATCH) bodies to the r.PostForm map
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Retrieve relevant data fields
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// Set up a map to hold validation errors
	errors := make(map[string]string)

	// *** Perform field validations ***

	// Title must not be blank and not more than 100 characters long
	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	// Content must not be blank
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	// Expires must not be blank and must have the value of 1, 7, or 365
	// TODO: Add support for leap years
	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	// *** End of field validations ***

	// Handle errors if any were encountered
	// If there are any errors, re-display the template passing to it the
	// validation errors and previously submitted form data
	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	// Insert the record through our model and receive back the ID of the new record
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

// createSnippetForm handler
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}
