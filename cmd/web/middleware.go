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

import "net/http"

import "fmt"

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
