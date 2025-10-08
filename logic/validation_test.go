package logic

import (
	"math"
	"testing"
)

func TestValidarEntrada(t *testing.T) {
	validObjective := []float64{1, 2}
	validConstraints := [][]float64{
		{1, 1},
		{2, 0},
	}
	validRHS := []float64{5, 4}

	// Caso válido
	if err := ValidarEntrada(validObjective, validConstraints, validRHS); err != nil {
		t.Errorf("Validación falló para datos válidos: %v", err)
	}

	// Dimensiones incorrectas: rhs
	badRHS := []float64{1}
	if err := ValidarEntrada(validObjective, validConstraints, badRHS); err == nil {
		t.Errorf("Esperaba error por dimensiones inconsistentes en rhs")
	}

	// Dimensiones incorrectas: objective
	badObjective := []float64{1, 2, 3}
	if err := ValidarEntrada(badObjective, validConstraints, validRHS); err == nil {
		t.Errorf("Esperaba error por dimensiones inconsistentes en objective")
	}

	// Valores no finitos: NaN
	nanObjective := []float64{1, math.NaN()}
	if err := ValidarEntrada(nanObjective, validConstraints, validRHS); err == nil {
		t.Errorf("Esperaba error por valor NaN en objective")
	}

	// Valores no finitos: Inf
	infConstraints := [][]float64{
		{1, math.Inf(1)},
		{2, 0},
	}
	if err := ValidarEntrada(validObjective, infConstraints, validRHS); err == nil {
		t.Errorf("Esperaba error por valor Inf en constraints")
	}
}
