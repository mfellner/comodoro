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

// HandleCORS handles a cross-origin preflight request.
func HandleCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			var origin string
			if origin = string(r.Header.Get("Origin")); origin == "" {
				origin = "null"
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "86400") // 86400 seconds = 1 day
			w.Header().Set("Content-Type", "text/plain")
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
