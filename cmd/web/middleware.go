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
