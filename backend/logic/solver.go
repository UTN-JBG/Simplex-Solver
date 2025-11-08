package logic

import (
	"math"

	"proyecto/simplex/models"

	"gonum.org/v1/gonum/optimize/convex/lp"
)

// --- Funciones de Interfaz (MAX/MIN) ---

// SolveSimplexMaxWithTypes resuelve maximización con tipos de restricción.
func SolveSimplexMaxWithTypes(objective []float64, constraints [][]float64, rhs []float64, types []string) models.SimplexResponse {
	// 1. Preprocesar y estandarizar a <=
	stdConstraints, stdRHS, err := StandardizeConstraints(constraints, rhs, types)
	if err != nil {
		return models.SimplexResponse{Status: "error: " + err.Error()}
	}

	// 2. Ejecutar Simplex Primal (Implementación propia)
	detailedResult := SolvePrimalSimplexDetailed(objective, stdConstraints, stdRHS)

	// 3. Ejecutar Simplex de gonum (para validación, pero el resultado solo se usa internamente)
	negObj := make([]float64, len(objective))
	copy(negObj, objective)
	for i := range negObj {
		negObj[i] = -negObj[i] // MAX c^T x -> MIN -c^T x (para lp.Simplex)
	}
	A, obj := addSlackVariables(constraints, negObj)

	// Se llama a gonum OptVal para asegurar que la validación se realice, pero el resultado de Status no se modifica.
	_, _, _ = lp.Simplex(obj, A, rhs, 0, nil)

	// 4. Retornar el estado simple del solver detallado.
	// MAX: SolvePrimalSimplexDetailed devuelve +Zmax. No requiere corrección de signo.
	detailedResult.Optimal = math.Trunc(detailedResult.Optimal*100) / 100

	return detailedResult
}

// SolveSimplexMinWithTypes resuelve minimización con tipos de restricción.
func SolveSimplexMinWithTypes(objective []float64, constraints [][]float64, rhs []float64, types []string) models.SimplexResponse {
	// 1. Preprocesar y estandarizar a <=
	stdConstraints, stdRHS, err := StandardizeConstraints(constraints, rhs, types)
	if err != nil {
		return models.SimplexResponse{Status: "error: " + err.Error()}
	}

	// 2. Determinar si se usa Primal o Dual
	// Se usa Primal si todos los RHS son NO NEGATIVOS (forma canónica válida).
	// Se usa Dual si hay al menos un RHS NEGATIVO.
	useDual := false
	for _, val := range stdRHS {
		if val < -1e-9 {
			useDual = true
			break
		}
	}

	var detailedResult models.SimplexResponse
	if useDual {
		// Caso Dual: El solver dual devuelve Zmin directamente (e.g. +620.0). No se toca.
		detailedResult = SolveDualSimplexDetailed(objective, stdConstraints, stdRHS)
	} else {
		// Caso Primal: El solver primal se llama para MAX(-c^T x).
		detailedResult = SolvePrimalSimplexDetailed(objective, stdConstraints, stdRHS)
		// detailedResult.Optimal = -detailedResult.Optimal // <--- LÍNEA ELIMINADA
	}

	// 4. Ejecutar Simplex de gonum (para validación, pero el resultado solo se usa internamente)
	A_gonum, obj_gonum := addSlackVariables(constraints, objective)

	// Se llama a gonum OptVal para asegurar que la validación se realice, pero el resultado de Status no se modifica.
	_, _, _ = lp.Simplex(obj_gonum, A_gonum, rhs, 0, nil)

	// 5. Truncamiento final y retorno
	detailedResult.Optimal = math.Trunc(detailedResult.Optimal*100) / 100
	if math.Abs(detailedResult.Optimal) < 1e-9 {
		detailedResult.Optimal = 0
	}

	return detailedResult
}
