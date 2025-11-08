package logic

import (
	"gonum.org/v1/gonum/mat"
)

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
