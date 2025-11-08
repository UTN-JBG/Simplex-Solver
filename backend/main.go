package main

import (
	"os"
	"proyecto/simplex/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())
	// Ruta para probar que el server responde
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Endpoint del simplex
	r.POST("/api/simplex", handlers.SolveSimplexHandler)
	// Puerto dinámico para Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback local
	}
	r.Run(":" + port)

}
