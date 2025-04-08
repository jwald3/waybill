package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// CORS middleware function that returns a handler
func CORS() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check the origin of the request
			origin := r.Header.Get("Origin")

			allowedOrigins := []string{
				"https://getwaybill.com",
				"https://www.getwaybill.com",
				"http://localhost:5173",
			}

			// Set appropriate CORS headers based on origin
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			// Always set these headers for all responses
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
