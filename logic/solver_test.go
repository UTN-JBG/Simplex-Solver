package logic

import "testing"

// Test:problema caso basico
func TestSolveSimplex_CasoBasico(t *testing.T) {
	c := []float64{3, 5}
	A := [][]float64{
		{1, 0},
		{0, 2},
		{3, 2},
	}
	b := []float64{4, 12, 18}

	result := SolveSimplex(c, A, b)
	expected := "Solución óptima:\nx1 = 2.00\nx2 = 6.00\nValor óptimo = 36.00"

	if result != expected {
		t.Errorf("Resultado incorrecto, got: %v, want: %v", result, expected)
	}
}

// Test:problema con todos ceros
func TestSolveSimplex_TodosCeros(t *testing.T) {
	c := []float64{0, 0}
	A := [][]float64{
		{1, 1},
	}
	b := []float64{5}

	result := SolveSimplex(c, A, b)
	expected := "Solución óptima:\nx1 = 0.00\nx2 = 5.00\nValor óptimo = 0.00"

	if result != expected {
		t.Errorf("Resultado incorrecto, got: %v, want: %v", result, expected)
	}
}

// Test: problema sin solución factible
func TestSolveSimplex_ProblemaInviable(t *testing.T) {
	c := []float64{1, 1}
	A := [][]float64{
		{1, 1},
	}
	b := []float64{-1} // RHS negativa, problema inviables

	result := SolveSimplex(c, A, b)
	expected := "Problema sin solución factible"

	if result != expected {
		t.Errorf("Resultado incorrecto, got: %v, want: %v", result, expected)
	}
}

// Test: problema ilimitado
func TestSolveSimplex_ProblemaIlimitado(t *testing.T) {
	c := []float64{1, 1}
	A := [][]float64{
		{-1, 1},
	}
	b := []float64{1}

	result := SolveSimplex(c, A, b)
	expected := "Problema ilimitado"

	if result != expected {
		t.Errorf("Resultado incorrecto, got: %v, want: %v", result, expected)
	}
}
