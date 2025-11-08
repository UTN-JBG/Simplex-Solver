package logic

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"proyecto/simplex/models"
)

// --- Lógica del Método Simplex Dual (Usado para MIN con RHS negativo) ---

// findDualPivotRow encuentra la fila pivote (variable que sale)
// Simplex Dual: Fila con el valor RHS más NEGATIVO.
func findDualPivotRow(tableau models.SimplexTableau) (int, error) {
	numRows := len(tableau)
	rhsCol := len(tableau[0]) - 1

	minRHS := 0.0
	pivotRow := -1

	// Iterar sobre las filas de restricción (índice 1 en adelante)
	for i := 1; i < numRows; i++ {
		rhs := tableau[i][rhsCol]
		if rhs < minRHS {
			minRHS = rhs
			pivotRow = i
		}
	}

	if pivotRow == -1 {
		return -1, errors.New("solución factible encontrada")
	}
	return pivotRow, nil
}

// findDualPivotColumn encuentra la columna pivote (variable que entra)
// Simplex Dual: Mínimo cociente (Fila Z / Elemento Pivote), solo para elementos negativos en la Fila Pivote.
func findDualPivotColumn(tableau models.SimplexTableau, pivotRow int) (int, error) {
	numCols := len(tableau[0])
	zRow := tableau[Z_ROW_INDEX]
	pivotElements := tableau[pivotRow]

	minRatio := math.Inf(1)
	pivotCol := -1

	// Iterar sobre las columnas de variables (índice 1 hasta len-2)
	for j := 1; j < numCols-1; j++ {
		pivotElement := pivotElements[j]
		if pivotElement < -1e-9 { // Solo considerar coeficientes NEGATIVOS en la fila pivote
			// Cociente Z[j] / |Pivot[j]|
			ratio := math.Abs(zRow[j] / pivotElement)

			if ratio < minRatio {
				minRatio = ratio
				pivotCol = j
			}
		}
	}

	if pivotCol == -1 {
		return -1, errors.New("problema infactible (infeasible)")
	}
	return pivotCol, nil
}

// SolveDualSimplexDetailed implementa el algoritmo Simplex Dual (MIN).
func SolveDualSimplexDetailed(objective []float64, constraints [][]float64, rhs []float64) models.SimplexResponse {
	response := models.SimplexResponse{
		Variables: make(map[string]float64),
		Status:    "error: unknown",
	}

	// 1. Validar entrada
	if err := ValidarEntrada(objective, constraints, rhs); err != nil {
		response.Status = "error: " + err.Error()
		return response
	}

	numVariables := len(objective)
	numConstraints := len(constraints)

	headers := generateColumnHeaders(numVariables, numConstraints)

	// 2. Construir la tabla inicial (es la misma lógica que Primal)
	currentTableau := buildInitialTableau(objective, constraints, rhs)

	// --- Lógica de truncamiento en línea para la Tabla Inicial ---
	// Guardar una copia truncada para la visualización (Tabla 0)
	truncatedInitial := copyTableau(currentTableau)
	for i, row := range truncatedInitial {
		for j, val := range row {
			truncatedInitial[i][j] = math.Trunc(val*100) / 100
		}
	}
	// Guardar la tabla inicial ENCAPSULADA
	response.TableauxHistory = append(response.TableauxHistory, models.TableauStep{
		Headers: headers,
		Matrix:  truncatedInitial,
	})

	numCols := len(currentTableau[0])
	rhsCol := numCols - 1

	for iter := 0; iter < 100; iter++ { // Límite de iteraciones
		// 3. Encontrar fila pivote (Dual: RHS más negativo)
		pivotRow, err := findDualPivotRow(currentTableau)
		if err != nil {
			if strings.Contains(err.Error(), "factible") {
				response.Status = "optimal" // Óptimo y Factible (solución encontrada)
				break
			}
			response.Status = "error: " + err.Error()
			return response
		}

		// 4. Encontrar columna pivote (Dual: Cociente Mínimo Z/|Pivot|)
		pivotCol, err := findDualPivotColumn(currentTableau, pivotRow)
		if err != nil {
			response.Status = "infeasible" // Infactibilidad detectada
			return response
		}

		// 5. Pivoteo
		currentTableau = pivot(currentTableau, pivotRow, pivotCol)
		// --- Lógica de truncamiento en línea para las Tablas Intermedias ---
		// Guardar una copia truncada para la visualización
		truncatedStep := copyTableau(currentTableau)
		for i, row := range truncatedStep {
			for j, val := range row {
				truncatedStep[i][j] = math.Trunc(val*100) / 100
			}
		}
		// Guardar el resultado del pivoteo ENCAPSULADO
		response.TableauxHistory = append(response.TableauxHistory, models.TableauStep{
			Headers: headers,
			Matrix:  truncatedStep,
		})
	}

	if response.Status != "optimal" {
		return response
	}

	// 6. Extracción de resultados
	// En Dual, la tabla final ya está en estado óptimo/factible, y el valor
	// RHS de la Fila Z es directamente el valor de Z_min.
	response.Optimal = currentTableau[Z_ROW_INDEX][rhsCol]

	for j := 1; j <= numVariables; j++ {
		isBasic := false
		basicRow := -1

		for i := 0; i < len(currentTableau); i++ {
			if math.Abs(currentTableau[i][j]-1.0) < 1e-9 {
				basicRow = i
				break
			}
		}

		if basicRow != -1 {
			isCanonical := true
			for i := 0; i < len(currentTableau); i++ {
				if i != basicRow && math.Abs(currentTableau[i][j]) > 1e-9 {
					isCanonical = false
					break
				}
			}
			if isCanonical && basicRow > 0 {
				isBasic = true
			}
		}

		if isBasic {
			response.Variables[fmt.Sprintf("x%d", j)] = math.Trunc(currentTableau[basicRow][rhsCol]*100) / 100
		} else {
			response.Variables[fmt.Sprintf("x%d", j)] = 0.0
		}
	}
	// Truncación final
	response.Optimal = math.Trunc(response.Optimal*100) / 100
	if math.Abs(response.Optimal) < 1e-9 {
		response.Optimal = 0
	}

	return response
}
