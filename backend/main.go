package main

import (
	"faceid/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Роуты
	r.POST("/api/register", handlers.Register)
	r.POST("/api/verify", handlers.Verify)
	r.GET("/api/users", handlers.GetUsers)

	r.Run(":3000")
}
