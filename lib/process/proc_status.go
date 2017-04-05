package process

import (
	"time"
)
// ProcStatus is a wrapper with the process current status.
type ProcStatus struct {
	Status		string
	Restarts	int
	Time			time.Time
	Uptime		string
}

// SetStatus will set the process string status.
func (proc_status *ProcStatus) SetStatus(status string) {
	proc_status.Status = status
}

// AddRestart will add one restart to the process status.
func (proc_status *ProcStatus) AddRestart() {
	proc_status.Restarts++
}

func (proc_status *ProcStatus) InitUptime() {
	proc_status.Time = time.Now()
}

// process uptime 
func (proc_status *ProcStatus) SetUptime() {
  proc_status.Uptime = time.Since(proc_status.Time).String()
}