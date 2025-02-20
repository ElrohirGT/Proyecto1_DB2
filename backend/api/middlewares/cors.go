package middlewares

import (
	"net/http"

	"github.com/ElrohirGT/Proyecto1_DB2/config"
	"github.com/go-chi/cors"
)

// CreateCors creates and returns a CORS middleware with custom options.
func CreateCors(config config.CorsConfig) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   config.AllowedOrigins, // Specify allowed origins
		AllowedMethods:   config.AllowedMethods, // Specify allowed HTTP methods
		AllowedHeaders:   config.AllowedHeaders, // Specify allowed headers
		AllowCredentials: true,                  // Allow credentials like cookies, authorization headers, etc.
		MaxAge:           300,                   // Cache duration for preflight requests
	})
}
