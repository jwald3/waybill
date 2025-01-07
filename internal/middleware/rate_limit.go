package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

// here's some very simple rate limiting to illustrate how an API can easily include the ability to restrict excessive API calls.
// this is just global limiting (all clients share the same pool of requests), but you'd realistically tie rate limits to individual consumers
func RateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// the limiter will determine if clients reached the max number of requests in the current window and either return an error or go to the next function/handler/middleware
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
