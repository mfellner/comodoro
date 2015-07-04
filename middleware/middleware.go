package middleware

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// Chain is a chain of http.Handlers.
type Chain struct {
	handlers []func(http.Handler) http.Handler
}

// New creates a new chain of http.Handlers.
func New(handlers ...func(http.Handler) http.Handler) Chain {
	return Chain{handlers: handlers}
}

// Then connects the middlewares and returns the final http.Handler.
func (c Chain) Then(h http.Handler) http.Handler {
	for _, handler := range c.handlers {
		h = handler(h)
	}
	return h
}

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

// LogHTTP logs HTTP requests.
func LogHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"method": r.Method}).Debug("=>")
		next.ServeHTTP(w, r)
		log.WithFields(log.Fields{"content-type": w.Header().Get("content-type")}).Debug("<=")
	})
}
