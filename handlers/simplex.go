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
}

func SolveSimplexHandler(c *gin.Context) {
	var req SimplexRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := logic.SolveSimplex(req.Objective, req.Constraints, req.RHS)

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
