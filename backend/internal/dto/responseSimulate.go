package dto

// Estrutura para representar um processo na resposta de simulação
type ResponseProcess struct {
	Name       string `json:"name"`
	Arrival    int    `json:"arrival"`
	Burst      int    `json:"burst"`
	Priority   int    `json:"priority"`
	Start      int    `json:"start"`
	Finish     int    `json:"finish"`
	Waiting    int    `json:"waiting"`
	Turnaround int    `json:"turnaround"`
	Response   int    `json:"response"`
}

// Estrutura para representar os intervalos de execução de um processo na resposta de simulação
type ResponseIntervals struct {
	ProcessName string `json:"process_name"`
	Start       int    `json:"start"`
	Finish      int    `json:"finish"`
}

// Estrutura para representar a linha do tempo da simulação
type ResponseTimeline struct {
	From   int               `json:"from"`
	To     int               `json:"to"`
	States map[string]string `json:"states"`
}

// Estrutura para representar as médias de tempo na resposta de simulação
type ResponseAverages struct {
	Turnaround float64 `json:"turnaround"`
	Waiting    float64 `json:"waiting"`
	Response   float64 `json:"response"`
}

// Estrutura para representar a resposta de simulação
type ResponseSimulate struct {
	Algorithm       string              `json:"algorithm"`
	TotalTime       int                 `json:"total_time"`
	ContextSwitches int                 `json:"context_switches"`
	Processes       []ResponseProcess   `json:"processes"`
	Intervals       []ResponseIntervals `json:"intervals"`
	Timeline        []ResponseTimeline  `json:"timeline"`
	Averages        ResponseAverages    `json:"averages"`
}
