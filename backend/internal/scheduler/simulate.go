package scheduler

type Simulate struct {
	Algorithm string
	Quantum   int
	Aging     int
	Processes []Process
}

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

type TimeLineRow struct {
	From   int
	To     int
	States map[string]string
}

type Averages struct {
	Turnaround float64
	Waiting    float64
	Response   float64
}

type Interval struct {
	Id          int
	ProcessName string
	Start       int
	Finish      int
}

type Result struct {
	Algorithm       string
	TotalTime       int
	ContextSwitches int
	Processes       []ProcessResult
	Intervals       []Interval
	Timeline        []TimeLineRow
	Averages        Averages
}
