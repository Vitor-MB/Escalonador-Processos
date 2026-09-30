package scheduler

// Process representa um processo a ser escalonado, incluindo campos internos usados pela simulação.
type Process struct {
	Id        int
	Name      string
	Arrival   int
	Burst     int
	Priority  int
	remaining int
	start     int
	finish    int
	base      int
	key       int
}
