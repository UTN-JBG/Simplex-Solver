package logic

import (
	"fmt"

	"proyecto/simplex/models"
)

// Constantes para el Simplex detallado
const (
	Z_ROW_INDEX = 0 // La primera fila es la función objetivo (Z)
)

// --- Funciones Auxiliares ---

// copyTableau es una función auxiliar para clonar la matriz
func copyTableau(tableau models.SimplexTableau) models.SimplexTableau {
	tableauCopy := make(models.SimplexTableau, len(tableau))
	for i, row := range tableau {
		tableauCopy[i] = make([]float64, len(row))
		copy(tableauCopy[i], row)
	}
	return tableauCopy
}

// generateColumnHeaders genera los nombres de las columnas: Z, x1...xn, s1...sn, RHS
func generateColumnHeaders(numVariables, numConstraints int) []string {
	headers := make([]string, 0, 2+numVariables+numConstraints)

	// Columna Z
	headers = append(headers, "Z")

	// Variables de decisión (x1, x2, ...)
	for i := 1; i <= numVariables; i++ {
		headers = append(headers, fmt.Sprintf("x%d", i))
	}

	// Variables de holgura (s1, s2, ...)
	for i := 1; i <= numConstraints; i++ {
		headers = append(headers, fmt.Sprintf("s%d", i))
	}

	// Columna RHS
	headers = append(headers, "RHS")

	return headers
}

// pivot realiza la operación de pivoteo
func pivot(tableau models.SimplexTableau, pivotRow, pivotCol int) models.SimplexTableau {
	numRows := len(tableau)
	numCols := len(tableau[0])
	newTableau := copyTableau(tableau)

	pivotVal := newTableau[pivotRow][pivotCol]

	// 1. Normalizar la fila pivote
	for j := 0; j < numCols; j++ {
		newTableau[pivotRow][j] /= pivotVal
	}

	// 2. Eliminar el resto de los elementos en la columna pivote
	for i := 0; i < numRows; i++ {
		if i != pivotRow {
			factor := newTableau[i][pivotCol]
			for j := 0; j < numCols; j++ {
				newTableau[i][j] -= factor * newTableau[pivotRow][j]
			}
		}
	}

	return newTableau
}
