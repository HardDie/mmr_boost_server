package middleware

import (
	"net/http"

	"github.com/HardDie/mmr_boost_server/internal/logger"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug.Println("["+r.Method+"]", r.URL.String())
		next.ServeHTTP(w, r)
	})
}
