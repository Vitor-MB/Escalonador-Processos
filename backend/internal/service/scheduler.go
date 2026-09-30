package service

import (
	"escalprocess/internal/dto"
	"escalprocess/internal/scheduler"
)

// Simulate converte a requisição da API para o formato interno e executa a simulação.
func Simulate(req dto.RequestSimulate) (dto.ResponseSimulate, error) {
	// Prepara a estrutura de entrada do núcleo do escalonador.
	s := scheduler.Simulate{
		Algorithm: req.Algorithm,
		Quantum:   req.Quantum,
		Aging:     req.Aging,
		Processes: make([]scheduler.Process, len(req.Processes)),
	}

	for i, p := range req.Processes {
		s.Processes[i] = scheduler.Process{
			Name:     p.Name,
			Arrival:  p.Arrival,
			Burst:    p.Burst,
			Priority: p.Priority,
		}
	}

	res, err := scheduler.Run(s)
	if err != nil {
		return dto.ResponseSimulate{}, err
	}

	out := dto.ResponseSimulate{
		Algorithm:       res.Algorithm,
		TotalTime:       res.TotalTime,
		ContextSwitches: res.ContextSwitches,
		Processes:       make([]dto.ResponseProcess, len(res.Processes)),
		Intervals:       make([]dto.ResponseIntervals, len(res.Intervals)),
		Timeline:        make([]dto.ResponseTimeline, len(res.Timeline)),
		Averages: dto.ResponseAverages{
			Turnaround: res.Averages.Turnaround,
			Waiting:    res.Averages.Waiting,
			Response:   res.Averages.Response,
		},
	}

	// Converte cada processo do resultado interno para o formato da resposta.
	for i, p := range res.Processes {
		out.Processes[i] = dto.ResponseProcess{
			Name:       p.Name,
			Arrival:    p.Arrival,
			Burst:      p.Burst,
			Priority:   p.Priority,
			Start:      p.Start,
			Finish:     p.Finish,
			Waiting:    p.Waiting,
			Turnaround: p.Turnaround,
			Response:   p.Response,
		}
	}

	// Converte os intervalos de execução em um formato serializável.
	for i, it := range res.Intervals {
		pname := it.ProcessName
		if pname == "" && it.Id > 0 && it.Id-1 < len(req.Processes) {
			pname = req.Processes[it.Id-1].Name
		}
		out.Intervals[i] = dto.ResponseIntervals{
			ProcessName: pname,
			Start:       it.Start,
			Finish:      it.Finish,
		}
	}

	// Converte a timeline interna.
	for i, t := range res.Timeline {
		out.Timeline[i] = dto.ResponseTimeline{
			From:   t.From,
			To:     t.To,
			States: t.States,
		}
	}

	return out, nil
}

// GetAlgorithms retorna todos os algoritmos disponíveis com suas características.
func GetAlgorithms() dto.ResponseAlgorithms {
	out := dto.ResponseAlgorithms{Algorithms: []dto.ResponseAlgorithm{}}
	algorithms := scheduler.ListAlgorithms()
	for _, m := range algorithms {
		out.Algorithms = append(out.Algorithms, dto.ResponseAlgorithm{
			Algorithm:    m.Algorithm,
			Preemptive:   m.Preemptive,
			UsesQuantum:  m.UsesQuantum,
			UsesAging:    m.UsesAging,
			UsesPriority: m.UsesPriority,
		})
	}
	return out
}
