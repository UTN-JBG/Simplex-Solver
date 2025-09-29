package logic

import (
	"fmt"
	"strings"
)

// SimplexSolver resuelve problemas de maximización lineales
type SimplexSolver struct {
	tableau [][]float64
	basis   []int
	status  string // nuevo: guarda estado ("ok", "inviable", "ilimitado")
}

// NewSimplex crea el solver a partir de la función objetivo y restricciones
func NewSimplex(c []float64, A [][]float64, b []float64) *SimplexSolver {
	m := len(A) // cantidad de restricciones
	n := len(c) // cantidad de variables
	totalVars := n + m

	// Verificar inviabilidad inicial (b < 0)
	for _, bi := range b {
		if bi < 0 {
			return &SimplexSolver{status: "inviable"}
		}
	}

	// Crear tableau con holguras
	tableau := make([][]float64, m+1)
	for i := 0; i < m; i++ {
		tableau[i] = make([]float64, totalVars+1)
		copy(tableau[i], A[i])
		tableau[i][n+i] = 1 // variable de holgura
		tableau[i][totalVars] = b[i]
	}
	// Fila Z
	tableau[m] = make([]float64, totalVars+1)
	for j := 0; j < n; j++ {
		tableau[m][j] = -c[j]
	}

	basis := make([]int, m)
	for i := 0; i < m; i++ {
		basis[i] = n + i
	}

	return &SimplexSolver{tableau, basis, "ok"}
}

// Solve ejecuta el algoritmo del simplex
func (s *SimplexSolver) Solve() {
	if s.status != "ok" {
		return
	}

	m := len(s.tableau) - 1
	n := len(s.tableau[0]) - 1

	for {
		// Encontrar columna entrante (coeficiente más negativo en Z)
		pivotCol := -1
		minVal := 0.0
		for j := 0; j < n; j++ {
			if s.tableau[m][j] < minVal {
				minVal = s.tableau[m][j]
				pivotCol = j
			}
		}
		if pivotCol == -1 {
			break // óptimo alcanzado
		}

		// Encontrar fila saliente
		pivotRow := -1
		ratioMin := 1e18
		for i := 0; i < m; i++ {
			if s.tableau[i][pivotCol] > 0 {
				ratio := s.tableau[i][n] / s.tableau[i][pivotCol]
				if ratio < ratioMin {
					ratioMin = ratio
					pivotRow = i
				}
			}
		}

		if pivotRow == -1 {
			s.status = "ilimitado"
			return
		}

		// Pivotear
		s.pivot(pivotRow, pivotCol)
		s.basis[pivotRow] = pivotCol
	}
}

// pivot realiza una operación de pivoteo en el tableau
func (s *SimplexSolver) pivot(row, col int) {
	m := len(s.tableau)
	n := len(s.tableau[0])

	pivotVal := s.tableau[row][col]
	for j := 0; j < n; j++ {
		s.tableau[row][j] /= pivotVal
	}

	for i := 0; i < m; i++ {
		if i != row {
			factor := s.tableau[i][col]
			for j := 0; j < n; j++ {
				s.tableau[i][j] -= factor * s.tableau[row][j]
			}
		}
	}
}

func SolveSimplex(objective []float64, constraints [][]float64, rhs []float64) string {
	solver := NewSimplex(objective, constraints, rhs)
	if solver.status == "inviable" {
		return "Problema sin solución factible"
	}

	solver.Solve()
	if solver.status == "ilimitado" {
		return "Problema ilimitado"
	}

	// Armar resultado como string
	m := len(solver.tableau) - 1
	n := len(solver.tableau[0]) - 1
	sol := make([]float64, n)

	for i := 0; i < m; i++ {
		if solver.basis[i] < n {
			sol[solver.basis[i]] = solver.tableau[i][n]
		}
	}

	var sb strings.Builder
	sb.WriteString("Solución óptima:\n")
	for j := 0; j < len(objective); j++ {
		sb.WriteString(fmt.Sprintf("x%d = %.2f\n", j+1, sol[j]))
	}
	sb.WriteString(fmt.Sprintf("Valor óptimo = %.2f", solver.tableau[m][n]))

	return sb.String()
}
