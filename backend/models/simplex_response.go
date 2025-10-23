package models

// SimplexTableau representa la matriz de coeficientes del Simplex
type SimplexTableau [][]float64

// TableauStep combina la matriz numérica con sus encabezados de columna
type TableauStep struct {
	Headers []string       `json:"headers"`
	Matrix  SimplexTableau `json:"matrix"`
}

type SimplexResponse struct {
	Variables       map[string]float64 `json:"variables"`
	Optimal         float64            `json:"optimal"`
	Status          string             `json:"status"`
	TableauxHistory []TableauStep      `json:"tableaux_history,omitempty"`
}
