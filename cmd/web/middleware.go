// Idiomatic format for a middleware function:
//
// func newMiddlewareFunc(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		next.ServeHTTP(w, r)
// 	})
// }

//TODO: Consider leveraging https://github.com/justinas/alice to simplify the middleware
//      handler chain as suggested in "Let's Go" 6.5

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"ptodd.org/snippetbox/pkg/models"
)

// secureHeaders adds security improvment headers to help prevent XSS and
// clickjacking attacks
// c.f. https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
// c.f. https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode-block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

// logRequest provides a middleware component to log all API requests
// NOTE: Implementing this as a method under the application object grants
// the function access to application dependencies such as information logger
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// recoverPanic provides middleware that nicely recovers from a panic is
// encountered during runtime
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// requireAuthentication provides middleware that protects handlers that
// require an authorized user to use
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		// And call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// noSurf provides middleware that protects againt CSRF attacks when a user is using
// a browser that does not support SameSite cookie attributes
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true, Path: "/", Secure: true,
	})
	return csrfHandler
}

// authenticate provides middleware that gets the user's ID from the session, checks it
// against the database to see if it is valid, makes sure that the user is still marked
// as active, then updates the request context with this information for the use of other
// handlers
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check to see if the authenticatedUserID value exists in the sesison.  If it does
		// not exist, simply pass control to the next handler
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		// Retrieve the current user's information from the database.  If it either cannot
		// be found or has been deactivated, then remove the now invalid authenticatedUserID
		// from the session and pass control to the next handler.
		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active { // user does not exist or is inactive
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		}
		if err != nil { // handle all other errors
			app.serverError(w, err)
			return
		}

		// Having confirmed the request is from an active and authenticated user, create a new
		// request context that indicates so and call the next handler using this new context
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
