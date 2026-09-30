package scheduler

import "testing"

// Testes de regressão para validar os algoritmos de escalonamento e suas métricas.

func TestRun_FirstComeFirstServed(t *testing.T) {
	// FCFS executa os processos na ordem de chegada.
	sim := Simulate{
		Algorithm: "fcfs",
		Processes: []Process{
			{Name: "P1", Arrival: 0, Burst: 5, Priority: 3},
			{Name: "P2", Arrival: 1, Burst: 2, Priority: 1},
			{Name: "P3", Arrival: 5, Burst: 3, Priority: 2},
		},
	}

	res, err := Run(sim)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if got, want := res.Processes[0].Start, 0; got != want {
		t.Fatalf("P1 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[1].Start, 5; got != want {
		t.Fatalf("P2 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[2].Start, 7; got != want {
		t.Fatalf("P3 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[1].Response, 4; got != want {
		t.Fatalf("P2 response time = %d, want %d", got, want)
	}
}

func TestRun_ShortestJobFirst(t *testing.T) {
	// SJF escolhe o processo com menor burst disponível.
	sim := Simulate{
		Algorithm: "sjf",
		Processes: []Process{
			{Name: "P1", Arrival: 0, Burst: 5, Priority: 3},
			{Name: "P2", Arrival: 0, Burst: 2, Priority: 1},
			{Name: "P3", Arrival: 1, Burst: 1, Priority: 2},
		},
	}

	res, err := Run(sim)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if got, want := res.Processes[1].Start, 0; got != want {
		t.Fatalf("P2 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[2].Start, 2; got != want {
		t.Fatalf("P3 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[0].Finish, 8; got != want {
		t.Fatalf("P1 finish = %d, want %d", got, want)
	}
}

func TestRun_ShortestRemainingTimeFirst(t *testing.T) {
	// SRTF realiza preempção quando um processo com menor tempo restante chega.
	sim := Simulate{
		Algorithm: "srtf",
		Processes: []Process{
			{Name: "P1", Arrival: 0, Burst: 5, Priority: 3},
			{Name: "P2", Arrival: 1, Burst: 2, Priority: 1},
			{Name: "P3", Arrival: 2, Burst: 1, Priority: 2},
		},
	}

	res, err := Run(sim)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if got, want := res.Processes[1].Start, 1; got != want {
		t.Fatalf("P2 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[2].Start, 3; got != want {
		t.Fatalf("P3 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[0].Finish, 8; got != want {
		t.Fatalf("P1 finish = %d, want %d", got, want)
	}
}

func TestRun_PriorityNonPreemptive(t *testing.T) {
	// Prioridade não preemptiva executa o processo atual até terminar antes de escolher outro.
	sim := Simulate{
		Algorithm: "prio-np",
		Processes: []Process{
			{Name: "P1", Arrival: 0, Burst: 5, Priority: 3},
			{Name: "P2", Arrival: 1, Burst: 2, Priority: 1},
			{Name: "P3", Arrival: 2, Burst: 2, Priority: 2},
		},
	}

	res, err := Run(sim)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if got, want := res.Processes[0].Start, 0; got != want {
		t.Fatalf("P1 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[1].Start, 5; got != want {
		t.Fatalf("P2 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[2].Start, 7; got != want {
		t.Fatalf("P3 start = %d, want %d", got, want)
	}
}

func TestRun_PriorityPreemptive(t *testing.T) {
	// Prioridade preemptiva troca de processo quando uma prioridade maior chega.
	sim := Simulate{
		Algorithm: "prio-p",
		Processes: []Process{
			{Name: "P1", Arrival: 0, Burst: 5, Priority: 3},
			{Name: "P2", Arrival: 1, Burst: 2, Priority: 1},
			{Name: "P3", Arrival: 2, Burst: 4, Priority: 2},
		},
	}

	res, err := Run(sim)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if got, want := res.Processes[1].Start, 1; got != want {
		t.Fatalf("P2 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[2].Start, 3; got != want {
		t.Fatalf("P3 start = %d, want %d", got, want)
	}
	if got, want := res.Processes[0].Finish, 11; got != want {
		t.Fatalf("P1 finish = %d, want %d", got, want)
	}
}

func TestRun_RoundRobinQuantum(t *testing.T) {
	// RR divide o tempo do processador em fatias, respeitando o quantum configurado.
	sim := Simulate{
		Algorithm: "rr",
		Quantum:   2,
		Processes: []Process{
			{Name: "P1", Arrival: 0, Burst: 4, Priority: 1},
			{Name: "P2", Arrival: 0, Burst: 3, Priority: 2},
		},
	}

	res, err := Run(sim)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if got, want := res.Processes[0].Finish, 6; got != want {
		t.Fatalf("P1 finish = %d, want %d", got, want)
	}
	if got, want := res.Processes[1].Finish, 7; got != want {
		t.Fatalf("P2 finish = %d, want %d", got, want)
	}
}

func TestRun_RoundRobinPriorityAging(t *testing.T) {
	// A funcionalidade de aging ajusta a prioridade dos processos que esperam muito tempo.
	running := &Process{Name: "P1", Arrival: 0, Burst: 5, Priority: 1, base: 1, key: 1}
	waiting1 := &Process{Name: "P2", Arrival: 0, Burst: 3, Priority: 3, base: 3, key: 3}
	waiting2 := &Process{Name: "P3", Arrival: 0, Burst: 2, Priority: 5, base: 5, key: 5}
	e := &Engine{
		Ready: []*Process{waiting1, waiting2},
		T:     2,
		Aging: 1,
		Last:  running,
	}

	rrAging{}.after(e, running, true)

	if got, want := waiting1.key, 2; got != want {
		t.Fatalf("P2 key after aging = %d, want %d", got, want)
	}
	if got, want := waiting2.key, 4; got != want {
		t.Fatalf("P3 key after aging = %d, want %d", got, want)
	}
	if got, want := running.key, 1; got != want {
		t.Fatalf("P1 key after aging = %d, want %d", got, want)
	}
}
