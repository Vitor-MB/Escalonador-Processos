package dto

// ResponseAlgorithm representa os metadados de um algoritmo de escalonamento disponível.
type ResponseAlgorithm struct {
	Algorithm    string `json:"algorithm"`
	Preemptive   bool   `json:"preemptive"`
	UsesQuantum  bool   `json:"uses_quantum"`
	UsesAging    bool   `json:"uses_aging"`
	UsesPriority bool   `json:"uses_priority"`
}

// ResponseAlgorithms é a resposta que lista todos os algoritmos suporteados pela API.
type ResponseAlgorithms struct {
	Algorithms []ResponseAlgorithm `json:"algorithms"`
}
