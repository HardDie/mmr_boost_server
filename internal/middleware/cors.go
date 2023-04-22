package middleware

import "net/http"

// CorsMiddleware CORS Headers middleware
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://editor.swagger.io")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,UPDATE,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
		if (*r).Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
