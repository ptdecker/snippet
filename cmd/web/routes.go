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
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	// Initialize new server mux
	mux := pat.New()

	// Register page routes
	//TODO: The endpoints that should not be used by authenticated users (signup and login) should also be protected
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	// Create a file server to serve static content
	fileServer := http.FileServer(neuteredFileSystem{http.Dir(cfg.staticDir)})
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// log the request, add security headers, then handle the request
	// also provides panic recovery as first thing
	return standardMiddleware.Then(mux)
}
