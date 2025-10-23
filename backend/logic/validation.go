package logic

import (
	"errors"
	"math"
)

// ValidarEntrada verifica que los datos sean correctos para el método simplex
func ValidarEntrada(objective []float64, constraints [][]float64, rhs []float64) error {
	// Vacíos o nulos
	if len(objective) == 0 {
		return errors.New("vector objective no puede estar vacío")
	}
	if len(constraints) == 0 {
		return errors.New("matriz constraints no puede estar vacía")
	}
	if len(rhs) == 0 {
		return errors.New("vector rhs no puede estar vacío")
	}

	// Dimensiones inconsistentes
	rows := len(constraints)
	cols := len(constraints[0])
	for i := range constraints {
		if len(constraints[i]) != cols {
			return errors.New("todas las filas de constraints deben tener el mismo número de columnas")
		}
	}
	if len(rhs) != rows {
		return errors.New("la longitud de rhs no coincide con el número de filas de constraints")
	}
	if len(objective) != cols {
		return errors.New("la longitud de objective no coincide con el número de columnas de constraints")
	}

	// Valores no finitos
	for _, v := range objective {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return errors.New("objective contiene valores no finitos")
		}
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if math.IsNaN(constraints[i][j]) || math.IsInf(constraints[i][j], 0) {
				return errors.New("constraints contiene valores no finitos")
			}
		}
		if math.IsNaN(rhs[i]) || math.IsInf(rhs[i], 0) {
			return errors.New("rhs contiene valores no finitos")
		}
	}

	return nil
}
