package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// prints out requests as they come in - this function prints out something like:
/*
{
	"level":"info",
	"ts":1739070606.8582923,
	"caller":"middleware/middleware.go:14",
	"msg":"incoming request",
	"method":"GET",
	"path":"/api/v1/drivers"
	},"path":"/api/v1/trucks"}
*/
// this function helps debug output in real time
func Logging(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("incoming request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)
			next.ServeHTTP(w, r)
		})
	}
}
