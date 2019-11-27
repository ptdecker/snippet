package main

import "net/http"

func (app *application) routes() http.Handler {
	// Initialize new server mux
	mux := http.NewServeMux()

	// Register home page route
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server to serve static content
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(cfg.staticDir)})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// log the request, add security headers, then handle the request
	// also provides panic recovery as first thing
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
