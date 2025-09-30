package logic

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
)

// SolveSimplexMin resuelve minimización con holguras automáticas
func SolveSimplexMin(objective []float64, constraints [][]float64, rhs []float64) []string {
	// Paso 0: invertir el signo de restricciones y RHS
	flippedConstraints := make([][]float64, len(constraints))
	flippedRHS := make([]float64, len(rhs))
	for i := range constraints {
		flippedConstraints[i] = make([]float64, len(constraints[i]))
		for j := range constraints[i] {
			flippedConstraints[i][j] = -constraints[i][j]
		}
		flippedRHS[i] = -rhs[i]
	}

	A, obj := addSlackVariables(flippedConstraints, objective)

	// Resolver LP
	optVal, x, err := lp.Simplex(obj, A, flippedRHS, 0, nil)
	if err != nil {
		return []string{"Error: " + err.Error()}
	}

	// Preparar resultado
	result := []string{"Solución óptima (minimización):"}
	for i := 0; i < len(objective); i++ { // mostrar solo las variables originales
		result = append(result, fmt.Sprintf("x%d = %.4f", i+1, x[i]))
	}
	result = append(result, fmt.Sprintf("Valor óptimo = %.4f", optVal))
	return result
}

// SolveSimplexMax resuelve maximización con holguras automáticas
func SolveSimplexMax(objective []float64, constraints [][]float64, rhs []float64) []string {
	// Para maximización, invertimos el objetivo y luego revertimos el valor óptimo
	negObj := make([]float64, len(objective))
	copy(negObj, objective)
	for i := range negObj {
		negObj[i] = -negObj[i]
	}

	A, obj := addSlackVariables(constraints, negObj)

	optVal, x, err := lp.Simplex(obj, A, rhs, 0, nil)
	if err != nil {
		return []string{"Error: " + err.Error()}
	}

	optVal = -optVal // revertir signo

	// Preparar resultado
	result := []string{"Solución óptima:"}
	for i := 0; i < len(objective); i++ {
		result = append(result, fmt.Sprintf("x%d = %.4f", i+1, x[i]))
	}
	result = append(result, fmt.Sprintf("Valor óptimo = %.4f", optVal))
	return result
}

// addSlackVariables agrega variables de holgura para restricciones ≤
func addSlackVariables(constraints [][]float64, objective []float64) (*mat.Dense, []float64) {
	rows := len(constraints)
	cols := len(constraints[0])
	totalCols := cols + rows // agregamos una variable de holgura por fila

	data := make([]float64, 0, rows*totalCols)
	for i := 0; i < rows; i++ {
		// copiar coeficientes originales
		row := make([]float64, totalCols)
		copy(row, constraints[i])
		// agregar variable de holgura
		row[cols+i] = 1
		data = append(data, row...)
	}

	A := mat.NewDense(rows, totalCols, data)

	// ampliar objetivo con ceros para las variables de holgura
	obj := make([]float64, totalCols)
	copy(obj, objective)

	return A, obj
}
