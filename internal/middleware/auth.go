package middleware

import (
	"net/http"

	"github.com/jwald3/go_rest_template/internal/config"
)

func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Load()

		validKey := cfg.App.APIKey
		if validKey == "" {
			http.Error(w, "server not configured with an API key", http.StatusInternalServerError)
			return
		}

		clientKey := r.Header.Get("X-API-Key")
		if clientKey == "" || clientKey != validKey {
			http.Error(w, "invalid or missing API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
