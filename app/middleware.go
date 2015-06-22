package app

import "net/http"

// AllowOrigin sets the Access-Control-Allow-Origin header for the current request.
func AllowOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var origin string
		if origin = string(r.Header.Get("Origin")); origin == "" {
			origin = "null"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		next.ServeHTTP(w, r)
	})
}
