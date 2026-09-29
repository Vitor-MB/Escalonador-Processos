package dto

// Estrutura para representar um processo na requisição de simulação
type RequestProcess struct {
	Name     string `json:"name"`
	Arrival  int    `json:"arrival"`
	Burst    int    `json:"burst"`
	Priority int    `json:"priority"`
}

// Estrutura para representar a requisição de simulação
type RequestSimulate struct {
	Algorithm string           `json:"algorithm"`
	Quantum   int              `json:"quantum"`
	Aging     int              `json:"aging"`
	Processes []RequestProcess `json:"processes"`
}
