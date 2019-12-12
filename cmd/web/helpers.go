package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

// serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError helper sends a specific status code and corresponding description
// to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound helper is a convenience wrapper around clientError to send a 404
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// addDefaultData adds default common data to the passed template structure
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	// Populate common data
	td.CurrentYear = time.Now().Year()

	// Retreive flash message from user session (if one)
	td.Flash = app.session.PopString(r, "flash")

	// Determine authentication status
	td.IsAuthenticated = app.isAuthenticated(r)

	// Add the CSRF protection token
	td.CSRFToken = nosurf.Token(r)

	return td
}

// render helper renders a page based upon templates from our cache.  It first
// renders the page to a buffer to trap render errors.  If succesful displays
// the page; otherwise, gracefully fails
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	// Retrieve template set from cache based upon teh page name.  If no entry
	// exists, error out gracefully
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// Initialize a buffer to hold the trial rendered page
	buf := new(bytes.Buffer)

	// Execute the template set passing in any dynamic data.  Inject into
	// the dynamic data any common data
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write out the page
	buf.WriteTo(w)
}

// isAuthenticated determines if the current request is from an authenticated user
// by checking the request context
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool) // cast the {interface} into the expected type
	return ok && isAuthenticated
}
