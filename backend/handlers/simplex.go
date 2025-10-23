package handlers

import (
	"net/http"

	"proyecto/simplex/logic"

	"proyecto/simplex/models"

	"github.com/gin-gonic/gin"
)

func SolveSimplexHandler(c *gin.Context) {
	var req models.SimplexRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar type
	if req.Type != "max" && req.Type != "min" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'type' debe ser 'max' o 'min'"})
		return
	}

	var result models.SimplexResponse

	if req.Type == "min" {
		result = logic.SolveSimplexMin(req.Objective, req.Constraints, req.RHS)
	} else {
		result = logic.SolveSimplexMax(req.Objective, req.Constraints, req.RHS)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
