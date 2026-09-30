package dto

// RequestProcess representa um processo enviado na requisição de simulação.
type RequestProcess struct {
	Name     string `json:"name"`
	Arrival  int    `json:"arrival"`
	Burst    int    `json:"burst"`
	Priority int    `json:"priority"`
}

// RequestSimulate é o corpo da requisição para iniciar uma simulação de escalonamento.
type RequestSimulate struct {
	Algorithm string           `json:"algorithm"`
	Quantum   int              `json:"quantum"`
	Aging     int              `json:"aging"`
	Processes []RequestProcess `json:"processes"`
}
