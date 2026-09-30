package scheduler

// Simulate descreve a entrada da simulação de escalonamento.
type Simulate struct {
	Algorithm string
	Quantum   int
	Aging     int
	Processes []Process
}

// ProcessResult contém as métricas finais de um processo após a simulação.
type ProcessResult struct {
	Id         int
	Name       string
	Arrival    int
	Burst      int
	Priority   int
	Start      int
	Finish     int
	Turnaround int
	Waiting    int
	Response   int
}

// TimeLineRow representa um instante da linha do tempo da execução.
type TimeLineRow struct {
	From   int
	To     int
	States map[string]string
}

// Averages agrega as médias de tempo calculadas ao final da simulação.
type Averages struct {
	Turnaround float64
	Waiting    float64
	Response   float64
}

// Interval representa um intervalo contínuo de execução de um processo.
type Interval struct {
	Id          int
	ProcessName string
	Start       int
	Finish      int
}

// Result reúne todos os dados gerados pela simulação do escalonador.
type Result struct {
	Algorithm       string
	TotalTime       int
	ContextSwitches int
	Processes       []ProcessResult
	Intervals       []Interval
	Timeline        []TimeLineRow
	Averages        Averages
}
