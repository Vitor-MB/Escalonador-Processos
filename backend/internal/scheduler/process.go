package scheduler

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
