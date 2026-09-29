package dto

type ResponseAlgorithm struct {
	Algorithm    string `json:"algorithm"`
	Preemptive   bool   `json:"preemptive"`
	UsesQuantum  bool   `json:"uses_quantum"`
	UsesAging    bool   `json:"uses_aging"`
	UsesPriority bool   `json:"uses_priority"`
}
type ResponseAlgorithms struct {
	Algorithms []ResponseAlgorithm `json:"algorithms"`
}
