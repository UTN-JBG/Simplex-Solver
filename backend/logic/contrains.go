package logic

import (
	"errors"
)

// StandardizeConstraints convierte las restricciones a la forma canónica (solo <=)
// Esto multiplica por -1 las restricciones GE (>=) y devuelve error en EQ (=).
func StandardizeConstraints(constraints [][]float64, rhs []float64, types []string) ([][]float64, []float64, error) {
	if len(constraints) != len(types) || len(constraints) != len(rhs) {
		return nil, nil, errors.New("los tamaños de las restricciones, RHS y tipos no coinciden")
	}

	newConstraints := make([][]float64, len(constraints))
	newRHS := make([]float64, len(rhs))

	for i := 0; i < len(constraints); i++ {
		newConstraints[i] = make([]float64, len(constraints[i]))
		copy(newConstraints[i], constraints[i])
		newRHS[i] = rhs[i]

		constraintType := types[i]

		switch constraintType {
		case "le": // Menor o igual (<=): 	.
		case "ge": // Mayor o igual (>=): Multiplicar por -1 para convertir a <=.
			// Esto convierte A_i x >= b_i en -A_i x <= -b_i.
			for j := 0; j < len(newConstraints[i]); j++ {
				newConstraints[i][j] *= -1.0
			}
			newRHS[i] *= -1.0
		case "eq": // Igualdad (=): Requiere Gran M.
			return nil, nil, errors.New("solver simple no soporta restricciones de igualdad (=). Se necesita el método de Gran M")
		default:
			return nil, nil, errors.New("tipo de restricción no reconocido: " + constraintType)
		}
	}

	return newConstraints, newRHS, nil
}
