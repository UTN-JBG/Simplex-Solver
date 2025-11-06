package logic

import (
	"math"
	"testing"
)

// Test: problema caso basico
func TestSolveSimplex_CasoBasico(t *testing.T) {
	c := []float64{3, 5}
	A := [][]float64{
		{1, 0},
		{0, 2},
		{3, 2},
	}
	b := []float64{4, 12, 18}
	// Creamos los tipos de restricción, asumiendo que todas son <= ("le")
	types := []string{"le", "le", "le"}

	result := SolveSimplexMaxWithTypes(c, A, b, types)

	if result.Status != "optimal" {
		t.Errorf("Se esperaba estado 'optimal', got: %v", result.Status)
	}

	if math.Abs(result.Optimal-36.0) > 1e-6 {
		t.Errorf("Valor óptimo incorrecto, got: %v, want: 36.0", result.Optimal)
	}

	if math.Abs(result.Variables["x1"]-2.0) > 1e-6 || math.Abs(result.Variables["x2"]-6.0) > 1e-6 {
		t.Errorf("Variables incorrectas, got: %v", result.Variables)
	}
}

// Test: problema con todos ceros
func TestSolveSimplex_TodosCeros(t *testing.T) {
	c := []float64{0, 0}
	A := [][]float64{
		{1, 1},
	}
	b := []float64{5}
	types := []string{"le"}

	result := SolveSimplexMaxWithTypes(c, A, b, types)

	if result.Status != "optimal" {
		t.Errorf("Se esperaba estado 'optimal', got: %v", result.Status)
	}

	if math.Abs(result.Optimal-0.0) > 1e-6 {
		t.Errorf("Valor óptimo incorrecto, got: %v, want: 0.0", result.Optimal)
	}

	x1 := result.Variables["x1"]
	x2 := result.Variables["x2"]

	if x1 < -1e-6 || x2 < -1e-6 || x1+x2 > 5+1e-6 {
		t.Errorf("Solución no factible: x1=%v, x2=%v", x1, x2)
	}
}

// Test: problema sin solución factible
func TestSolveSimplex_ProblemaInviable(t *testing.T) {
	c := []float64{1, 1}
	A := [][]float64{
		{1, 1},
	}
	b := []float64{-1} // RHS negativa, problema inviables
	types := []string{"le"}

	result := SolveSimplexMaxWithTypes(c, A, b, types)

	if result.Status != "infeasible" {
		t.Errorf("Resultado incorrecto, se esperaba infeasible, got: %v", result.Status)
	}
}

// Test: problema ilimitado
func TestSolveSimplex_ProblemaIlimitado(t *testing.T) {
	c := []float64{1, 1}
	A := [][]float64{
		{-1, 1},
	}
	b := []float64{1}
	types := []string{"le"}

	result := SolveSimplexMaxWithTypes(c, A, b, types)

	if result.Status != "unbounded" {
		t.Errorf("Resultado incorrecto, se esperaba unbounded, got: %v", result.Status)
	}
}
