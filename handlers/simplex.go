package handlers

import (
	"net/http"

	"proyecto/simplex/logic"

	"github.com/gin-gonic/gin"
)

// Estructura para recibir JSON
type SimplexRequest struct {
	Objective   []float64   `json:"objective"`   // coeficientes función objetivo
	Constraints [][]float64 `json:"constraints"` // matriz de restricciones
	RHS         []float64   `json:"rhs"`         // términos independientes
	Type        string      `json:"type"`        // "max" o "min"
}

func SolveSimplexHandler(c *gin.Context) {
	var req SimplexRequest
	var result []string

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar type
	if req.Type != "max" && req.Type != "min" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'type' debe ser 'max' o 'min'"})
		return
	}

	// Si es minimización, invertimos los coeficientes
	if req.Type == "min" {
		result = logic.SolveSimplexMin(req.Objective, req.Constraints, req.RHS)
	} else {
		result = logic.SolveSimplexMax(req.Objective, req.Constraints, req.RHS)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
