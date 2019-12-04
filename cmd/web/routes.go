// c.f. https://vluxe.io/golang-router.html for a reference implementation of a basic custom router

package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// Set-up middleware chain for all requests
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Set-up dynamic middleware chain for specific requests that need
	// session state
	dynamicMiddleware := alice.New(app.session.Enable)

	// Initialize new server mux
	mux := pat.New()

	// Register home page route
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// Create a file server to serve static content
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(cfg.staticDir)})
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// log the request, add security headers, then handle the request
	// also provides panic recovery as first thing
	return standardMiddleware.Then(mux)
}
