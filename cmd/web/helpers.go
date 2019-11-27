package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
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

	// Execute the template set passing in any dynamic data
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write out the page
	buf.WriteTo(w)
}
