package models

type SimplexResponse struct {
	Variables map[string]float64 `json:"variables"`
	Optimal   float64            `json:"optimal"`
	Status    string             `json:"status"`
}
