package main

import (
	"proyecto/simplex/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"


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

	r.Run(":8080") // localhost:8080
}
