// Idiomatic format for a middleware function:
//
// func newMiddlewareFunc(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		next.ServeHTTP(w, r)
// 	})
// }

package main

import "net/http"

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
