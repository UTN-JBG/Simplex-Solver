package test

import (
	"math"
	"proyecto/simplex/logic"
	"strings"
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
	if err := logic.ValidarEntrada(validObjective, validConstraints, validRHS); err != nil {
		t.Errorf("Validación falló para datos válidos: %v", err)
	}

	// Dimensiones incorrectas: rhs
	badRHS := []float64{1}
	if err := logic.ValidarEntrada(validObjective, validConstraints, badRHS); err == nil {
		t.Errorf("Esperaba error por dimensiones inconsistentes en rhs")
	}

	// Dimensiones incorrectas: objective
	badObjective := []float64{1, 2, 3}
	if err := logic.ValidarEntrada(badObjective, validConstraints, validRHS); err == nil {
		t.Errorf("Esperaba error por dimensiones inconsistentes en objective")
	}

	// Valores no finitos: NaN
	nanObjective := []float64{1, math.NaN()}
	if err := logic.ValidarEntrada(nanObjective, validConstraints, validRHS); err == nil {
		t.Errorf("Esperaba error por valor NaN en objective")
	}

	// Valores no finitos: Inf
	infConstraints := [][]float64{
		{1, math.Inf(1)},
		{2, 0},
	}
	if err := logic.ValidarEntrada(validObjective, infConstraints, validRHS); err == nil {
		t.Errorf("Esperaba error por valor Inf en constraints")
	}
}

// Test: Problema de Maximización (Requiere 2 iteraciones)
func TestSolveSimplexMax_Detallado_Pasos(t *testing.T) {
	// Problema de ejemplo (MAX): Z = 3x1 + 5x2
	// s.a. x1 <= 4, 2x2 <= 12, 3x1 + 2x2 <= 18
	c := []float64{3, 5}
	A := [][]float64{
		{1, 0},
		{0, 2},
		{3, 2},
	}
	b := []float64{4, 12, 18}
	types := []string{"le", "le", "le"}

	result := logic.SolveSimplexMaxWithTypes(c, A, b, types)

	expectedOptimal := 36.0
	// Se esperan 3 Tablas: Inicial (1) + Pivot 1 (1) + Pivot 2 (1) = 3
	expectedTableauxCount := 3

	// 1. Verificar Estado y Validación
	if !strings.Contains(result.Status, "optimal") && !strings.Contains(result.Status, "OK") {
		t.Errorf("La validación con gonum falló o el estado es incorrecto. Estado: %v", result.Status)
	}

	// 2. Verificar el Valor Óptimo
	if math.Abs(result.Optimal-expectedOptimal) > 1e-6 {
		t.Errorf("Valor óptimo incorrecto, got: %.4f, want: %.4f", result.Optimal, expectedOptimal)
	}

	// 3. Verificar la Cantidad de Tablas
	if len(result.TableauxHistory) != expectedTableauxCount {
		t.Errorf("Historial de tablas incorrecto. Se esperaban %d tablas (Inicial + 2 Pivotes), got: %d", expectedTableauxCount, len(result.TableauxHistory))
	}
}

// Test: Problema de Minimización (Resuelto por Dualidad MAX, requiere 2 iteraciones)
func TestSolveSimplexMin_Detallado_Pasos(t *testing.T) {
	// Problema de ejemplo (MIN):
	c := []float64{40, 50, 60}
	A := [][]float64{
		{4, 4, 5},
		{200, 300, 300},
		{600, 400, 500},
	}
	b := []float64{80, 6000, 5000}
	types := []string{"le", "le", "le"}

	result := logic.SolveSimplexMinWithTypes(c, A, b, types)

	// La solución MIN es -7.0 (ya que MAX Z' = 7.0)
	expectedOptimal := -7.0
	// CORREGIDO: Se esperan 3 Tablas: Inicial (1) + Pivot 1 (1) + Pivot 2 (1) = 3
	expectedTableauxCount := 3

	// 1. Verificar Estado y Validación
	if !strings.Contains(result.Status, "optimal") && !strings.Contains(result.Status, "OK") {
		t.Errorf("La validación con gonum falló o el estado es incorrecto. Estado: %v", result.Status)
	}

	// 2. Verificar el Valor Óptimo
	if math.Abs(result.Optimal-expectedOptimal) > 1e-6 {
		t.Errorf("Valor óptimo incorrecto, got: %.4f, want: %.4f", result.Optimal, expectedOptimal)
	}

	// 3. Verificar la Cantidad de Tablas
	if len(result.TableauxHistory) != expectedTableauxCount {
		t.Errorf("Historial de tablas incorrecto. Se esperaban %d tablas, got: %d", expectedTableauxCount, len(result.TableauxHistory))
	}
}
