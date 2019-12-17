package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"ptodd.org/snippetbox/pkg/forms"
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

	// Render the template passing the snippet
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

	// Retrieve and validate relevant data fields
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "1", "7", "365")

	// Handle errors if any were encountered
	// If there are any errors, re-display the template passing to it the
	// validation errors and previously submitted form data
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	// Insert the record through our model and receive back the ID of the new record
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Add a flash confirmation to the user session
	app.session.Put(r, "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

// createSnippetForm handler
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// signupUserForm handler
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {

	// Render the page
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// signupUser handler
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Retrieve and validate relevant data fields
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	// Handle errors if any were encountered
	// If there are any errors, re-display the template passing to it the
	// validation errors and previously submitted form data
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	// Try to create the user record in the database
	// If an error occurs, handle the error
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) { // email exists already
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else { // all other errors
			app.serverError(w, err)
		}
		return
	}

	// Notify the user of a successful record creation
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	// Redirect back to the log-in page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// loginUserForm handler
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// loginUser handler
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic error // message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Either your email is incorrect, password is incorrect, or your account is inactive")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Add the ID of the current user to the session, so that they are now 'logged // in'.
	app.session.Put(r, "authenticatedUserID", id)

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

// logoutUser handler
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {

	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.session.Remove(r, "authenticatedUserID")

	// Add a flash message to the session to confirm to the user that they've been // logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ping handler
// TODO: Why isn't this a funciton of application
// TODO: Implement a route
// TODO: Expand to check for database connectivity
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
