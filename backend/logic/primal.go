package logic

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"proyecto/simplex/models"
)

// --- Lógica del Método Simplex Primal (Usado para MAX) ---

// buildInitialTableau construye la tabla inicial del Simplex para un problema de MAXIMIZACIÓN.
func buildInitialTableau(objective []float64, constraints [][]float64, rhs []float64) models.SimplexTableau {
	numConstraints := len(constraints)
	numVariables := len(objective)
	numCols := numVariables + numConstraints + 2 // Z_col + Var_cols + Slack_cols + RHS_col

	tableau := make(models.SimplexTableau, numConstraints+1) // Fila Z + Filas de Restricción

	// 1. Fila Z (Índice 0)
	zRow := make([]float64, numCols)
	zRow[0] = 1.0 // Coeficiente de Z
	for j := 0; j < numVariables; j++ {
		zRow[j+1] = -objective[j] // Coeficientes de variables originales (negativos para MAX)
	}
	tableau[Z_ROW_INDEX] = zRow

	// 2. Filas de Restricción (Índice 1 en adelante)
	for i := 0; i < numConstraints; i++ {
		row := make([]float64, numCols)
		row[0] = 0.0 // Coeficiente de Z es 0

		// Coeficientes de variables originales
		for j := 0; j < numVariables; j++ {
			row[j+1] = constraints[i][j]
		}

		// Variables de holgura (identidad)
		slackColIndex := numVariables + 1 + i
		row[slackColIndex] = 1.0

		// RHS (columna final)
		row[numCols-1] = rhs[i]

		tableau[i+1] = row
	}

	return tableau
}

// findPivotColumn (Primal) encuentra la columna pivote (variable que entra a la base)
// En maximización, es el valor más negativo en la fila Z.
func findPivotColumn(tableau models.SimplexTableau) (int, error) {
	zRow := tableau[Z_ROW_INDEX]
	pivotCol := -1
	minVal := 0.0

	// Buscar el valor más negativo en la fila Z (índice 1 a len-2)
	for j := 1; j < len(zRow)-1; j++ {
		if zRow[j] < minVal {
			minVal = zRow[j]
			pivotCol = j
		}
	}

	if pivotCol == -1 {
		return -1, errors.New("solución óptima encontrada")
	}
	return pivotCol, nil
}

// findPivotRow (Primal) realiza la prueba del cociente mínimo para encontrar la fila pivote
func findPivotRow(tableau models.SimplexTableau, pivotCol int) (int, error) {
	numRows := len(tableau)
	rhsCol := len(tableau[0]) - 1

	minRatio := math.Inf(1)
	pivotRow := -1

	// Iterar sobre las filas de restricción (índice 1 en adelante)
	for i := 1; i < numRows; i++ {
		pivotElement := tableau[i][pivotCol]
		rhs := tableau[i][rhsCol]

		// Los divisores deben ser positivos
		if pivotElement > 1e-9 {
			ratio := rhs / pivotElement
			if ratio < minRatio {
				minRatio = ratio
				pivotRow = i
			}
		}
	}

	if pivotRow == -1 {
		return -1, errors.New("problema ilimitado (unbounded)")
	}

	return pivotRow, nil
}

// SolvePrimalSimplexDetailed implementa el algoritmo Simplex Primal (MAX).
func SolvePrimalSimplexDetailed(objective []float64, constraints [][]float64, rhs []float64) models.SimplexResponse {
	response := models.SimplexResponse{
		Variables: make(map[string]float64),
		Status:    "error: unknown",
	}

	// 1. Validar entrada
	if err := ValidarEntrada(objective, constraints, rhs); err != nil {
		response.Status = "error: " + err.Error()
		return response
	}
	// El Simplex Primal requiere que RHS >= 0 para comenzar.
	for _, val := range rhs {
		if val < -1e-9 {
			// Si hay un RHS negativo, el problema es inviable para el Primal simple.
			response.Status = "infeasible"
			return response
		}
	}

	numVariables := len(objective)
	numConstraints := len(constraints)

	headers := generateColumnHeaders(numVariables, numConstraints)

	// 2. Construir la tabla inicial
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

	// Si la función objetivo es constante (todos los coeficientes 0),
	// devolver la primera solución factible trivial.
	isZeroObjective := true
	for _, v := range objective {
		if math.Abs(v) > 1e-9 {
			isZeroObjective = false
			break
		}
	}
	if isZeroObjective {
		for j := 1; j <= numVariables; j++ {
			response.Variables[fmt.Sprintf("x%d", j)] = 0.0
		}
		response.Optimal = 0.0
		response.Status = "optimal (degenerate: multiple solutions)"
		return response
	}

	numCols := len(currentTableau[0])
	rhsCol := numCols - 1

	for iter := 0; iter < 100; iter++ { // Límite de iteraciones
		// 3. Encontrar columna pivote (Primal: más negativo en Z)
		pivotCol, err := findPivotColumn(currentTableau)
		if err != nil {
			if strings.Contains(err.Error(), "óptima") {
				response.Status = "optimal"
				break
			}
			response.Status = "error: " + err.Error()
			return response
		}

		// 4. Encontrar fila pivote (Primal: Cociente Mínimo)
		pivotRow, err := findPivotRow(currentTableau, pivotCol)
		if err != nil {
			response.Status = "unbounded"
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
	response.Optimal = currentTableau[Z_ROW_INDEX][rhsCol]
	response.Optimal = math.Trunc(response.Optimal*100) / 100 // Truncamiento

	// Extracción de valores de variables originales
	for j := 1; j <= numVariables; j++ {
		isBasic := false
		basicRow := -1

		// 1. Buscar el 1.0 en la columna
		for i := 0; i < len(currentTableau); i++ {
			if math.Abs(currentTableau[i][j]-1.0) < 1e-9 {
				basicRow = i
				break
			}
		}

		// 2. Si hay un 1.0, verificar que el resto de la columna sea 0.0 (columna canónica)
		if basicRow != -1 {
			isCanonical := true
			for i := 0; i < len(currentTableau); i++ {
				if i != basicRow && math.Abs(currentTableau[i][j]) > 1e-9 {
					isCanonical = false
					break
				}
			}
			if isCanonical && basicRow > 0 { // basicRow > 0 excluye la fila Z
				isBasic = true
			}
		}

		// 3. Asignar valor
		if isBasic {
			response.Variables[fmt.Sprintf("x%d", j)] = math.Trunc(currentTableau[basicRow][rhsCol]*100) / 100
		} else {
			response.Variables[fmt.Sprintf("x%d", j)] = 0.0
		}
	}

	return response
}
