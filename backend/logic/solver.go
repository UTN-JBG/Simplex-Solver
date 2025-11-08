package logic

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"proyecto/simplex/models"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
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
		if val < -1e-9 { // Si hay un RHS negativo, Primal no puede comenzar.
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
