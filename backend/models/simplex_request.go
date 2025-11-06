package models

// Estructura para recibir JSON
type SimplexRequest struct {
	Objective       []float64   `json:"objective"`   // coeficientes función objetivo
	Constraints     [][]float64 `json:"constraints"` // matriz de restricciones
	RHS             []float64   `json:"rhs"`         // términos independientes
	Type            string      `json:"type"`        // "max" o "min"
	ConstraintTypes []string    `json:"constraint_types"`
}
