package middleware

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

// NewCORSMiddleware tạo middleware CORS
func NewCORSMiddleware() iris.Handler {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return crs
}
