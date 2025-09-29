package logic

import (
	"fmt"
	"strings"
	"testing"
)

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

	if !strings.Contains(result, "Valor óptimo = 0.00") {
		t.Errorf("Se esperaba valor óptimo 0.00, got: %v", result)
	}

	// Extraer las variables de la salida
	var x1, x2, valor float64
	fmt.Sscanf(result, "Solución óptima:\nx1 = %f\nx2 = %f\nValor óptimo = %f", &x1, &x2, &valor)

	// Verificar factibilidad
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
