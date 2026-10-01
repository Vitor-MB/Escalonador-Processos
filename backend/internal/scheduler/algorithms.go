package scheduler

// Estrutura que representa um algoritmo de escalonamento
type MetaData struct {
	Algorithm    string
	Preemptive   bool
	UsesQuantum  bool
	UsesAging    bool
	UsesPriority bool
}

// Algoritmos disponiveis
var Algorithms = map[string]MetaData{
	"fcfs":          {Algorithm: "fcfs", Preemptive: false, UsesQuantum: false, UsesAging: false, UsesPriority: false},
	"sjf":           {Algorithm: "sjf", Preemptive: false, UsesQuantum: false, UsesAging: false, UsesPriority: false},
	"srtf":          {Algorithm: "srtf", Preemptive: true, UsesQuantum: false, UsesAging: false, UsesPriority: false},
	"prio-np":       {Algorithm: "prio-np", Preemptive: false, UsesQuantum: false, UsesAging: false, UsesPriority: true},
	"prio-p":        {Algorithm: "prio-p", Preemptive: true, UsesQuantum: false, UsesAging: false, UsesPriority: true},
	"rr":            {Algorithm: "rr", Preemptive: true, UsesQuantum: true, UsesAging: false, UsesPriority: false},
	"rr-prio-aging": {Algorithm: "rr-prio-aging", Preemptive: true, UsesQuantum: true, UsesAging: true, UsesPriority: true},
}

// Ordem para saida padronizada em GET/algorithm
var AlgorithmOrder = []string{
	"fcfs", "sjf", "srtf", "prio-np", "prio-p", "rr", "rr-prio-aging",
}

func ListAlgorithms() []MetaData {
	list := make([]MetaData, 0, len(AlgorithmOrder))
	for _, name := range AlgorithmOrder {
		list = append(list, Algorithms[name])
	}
	return list
}

// Interface que define o comportamento de um algoritmo de escalonamento
type algorithm interface {
	// Escolhe o próximo processo a ser executado no tick atual
	pick(e *Engine) *Process

	// Chamado apos a execução de um tick, para atualizar o estado do algoritmo
	// running: processo que estava rodando no tick atual, ou nil se nenhum processo estava rodando
	// sliceEnded: indica se o slice de execução do processo acabou no tick atual
	after(e *Engine, Running *Process, sliceEnded bool)
}

func newAlgorithm(id string) algorithm {
	switch id {
	case "fcfs":
		return minorKey{key: func(p *Process) int { return p.Arrival }}
	case "sjf":
		return minorKey{key: func(p *Process) int { return p.Burst }}
	case "srtf":
		return minorKey{key: func(p *Process) int { return p.remaining }, preemptive: true}
	case "prio-np":
		return minorKey{key: func(p *Process) int { return p.base }}
	case "prio-p":
		return minorKey{key: func(p *Process) int { return p.base }, preemptive: true}
	case "rr":
		return roundRobin{}
	default: // "rr-prio-aging" (já validado em Run)
		return rrAging{}
	}
}

// minorKey implementa a interface algorithm para algoritmos que escolhem o
// próximo processo com base em uma chave menor (arrival, burst, remaining, priority).
// Se priority for escolhida pela maior, inverte o sinal da chave
type minorKey struct {
	key        func(p *Process) int
	preemptive bool
}

// pick escolhe o próximo processo a ser executado com base na chave menor.
func (m minorKey) pick(e *Engine) *Process {
	if !m.preemptive {
		if c := e.Current(); c != nil {
			return c
		}
	}

	return e.Best(m.key)
}

func (minorKey) after(*Engine, *Process, bool) {}

// roundRobin implementa a interface algorithm para o algoritmo Round Robin.
type roundRobin struct{}

func (roundRobin) pick(e *Engine) *Process {
	if c := e.Current(); c != nil && e.Slice > 0 {
		return c // Ainda no quantum
	}
	if len(e.Ready) == 0 {
		return nil // Execuçao encerrada
	}

	return e.Ready[0]
}

// after é chamado após a execução de um tick, para atualizar o estado do algoritmo Round Robin.
func (roundRobin) after(e *Engine, running *Process, sliceEnded bool) {
	if sliceEnded && running.remaining > 0 {
		e.removeFromReady(running)
		e.Ready = append(e.Ready, running)
	}
}

// rrAging implementa a interface algorithm para o algoritmo Round Robin com prioridade e Aging.
type rrAging struct{}

// pick escolhe o próximo processo a ser executado com base na prioridade e no Aging.
func (rrAging) pick(e *Engine) *Process {
	if c := e.Current(); c != nil && e.Slice > 0 {
		return c
	}

	return e.Best(func(p *Process) int { return p.key })
}

// after é chamado após a execução de um tick, para atualizar o estado do algoritmo Round Robin com prioridade e Aging.
func (rrAging) after(e *Engine, running *Process, sliceEnded bool) {
	if !sliceEnded {
		return
	}

	running.key = running.base
	for _, p := range e.Ready {
		// Só envelhece quem esperou durante o quantum (chegou antes do instante atual).
		if p != running && p.Arrival < e.T {
			// Sem piso: a prioridade efetiva pode passar do valor máximo estático,
			// garantindo que quem espera acabe vencendo (sem inanição).

			p.key -= e.Aging
		}
	}
}
