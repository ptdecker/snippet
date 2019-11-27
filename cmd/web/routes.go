// c.f. https://vluxe.io/golang-router.html for a reference implementation of a basic custom router

package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	// Initialize new server mux

	mux := pat.New()

	// Register home page route
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	// Create a file server to serve static content
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(cfg.staticDir)})
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// log the request, add security headers, then handle the request
	// also provides panic recovery as first thing
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
