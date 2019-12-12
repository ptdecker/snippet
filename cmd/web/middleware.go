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
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
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

// recoverPanic provides middleware that nicely recovers from if a panic is
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

		// If the user is not authenticated, redirect them to the login page and // return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or // other intermediary cache).
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
