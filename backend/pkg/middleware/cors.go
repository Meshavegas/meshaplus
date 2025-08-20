package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// CORS retourne un middleware CORS configuré
func CORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"https://localhost:3000",
			"https://localhost:8080",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Link",
			"X-Total-Count",
		},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	})
}
