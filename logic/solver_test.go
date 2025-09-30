package logic

import (
	"fmt"
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
	expected := []string{
		"Solución óptima:",
		"x1 = 2.00",
		"x2 = 6.00",
		"Valor óptimo = 36.00",
	}
	if len(result) != len(expected) {
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

	valorOptimo := result[len(result)-1]
	if valorOptimo != "Valor óptimo = 0.00" {
		t.Errorf("Se esperaba valor óptimo 0.00, got: %v", valorOptimo)
	}

	// Extraer las variables de la salida
	var x1, x2 float64
	_, err1 := fmt.Sscanf(result[1], "x1 = %f", &x1)
	_, err2 := fmt.Sscanf(result[2], "x2 = %f", &x2)
	if err1 != nil || err2 != nil {
		t.Fatalf("Error al parsear variables de la salida")
	}

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
	expected := []string{"Problema sin solución factible"}

	if len(result) != len(expected) || result[0] != expected[0] {
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
	expected := []string{"Problema ilimitado"}

	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Resultado incorrecto, got: %v, want: %v", result, expected)
	}
}
