package middleware

import (
	"net/http"

	"github.com/jwald3/waybill/internal/config"
)

// this is a SUPER basic API key-based authentication that uses a key stored on the server and requires that clients provide the key in their request headers. If the client key matches the server key, they can access the resources. Otherwise, they cannot.
// I'll expand on a more robust auth system later on.
func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we're going to access the API key as stored in config.go. Using the config means that we have the default value as a fallback.
		cfg := config.Load()
		validKey := cfg.App.APIKey
		if validKey == "" {
			// this should never happen outside of the default being an empty string or the user-defined value being empty
			http.Error(w, "server not configured with an API key", http.StatusInternalServerError)
			return
		}

		// access the client API key value from the header. If it doesn't match the server value, throw a 401. Otherwise, continue to the `next()` method (either the handler or additional middleware)
		clientKey := r.Header.Get("X-API-Key")
		if clientKey == "" || clientKey != validKey {
			http.Error(w, "invalid or missing API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
