package cli

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/struCoder/pmgo/lib/master"
	"github.com/struCoder/pmgo/lib/utils"
)

// Cli is the command line client.
type Cli struct {
	remoteClient *master.RemoteClient
}

// InitCli initiates a remote client connecting to dsn.
// Returns a Cli instance.
func InitCli(dsn string, timeout time.Duration) *Cli {
	client, err := master.StartRemoteClient(dsn, timeout)
	if err != nil {
		log.Fatalf("Failed to start remote client due to: %+v\n", err)
	}
	return &Cli{
		remoteClient: client,
	}
}

// Save will save all previously saved processes onto a list.
// Display an error in case there's any.
func (cli *Cli) Save() {
	err := cli.remoteClient.Save()
	if err != nil {
		log.Fatalf("Failed to save list of processes due to: %+v\n", err)
	}
}

// Resurrect will restore all previously save processes.
// Display an error in case there's any.
func (cli *Cli) Resurrect() {
	err := cli.remoteClient.Resurrect()
	if err != nil {
		log.Fatalf("Failed to resurrect all previously save processes due to: %+v\n", err)
	}
}

// StartGoBin will try to start a go binary process.
// Returns a fatal error in case there's any.
func (cli *Cli) StartGoBin(sourcePath string, name string, keepAlive bool, args []string) {
	err := cli.remoteClient.StartGoBin(sourcePath, name, keepAlive, args)
	if err != nil {
		log.Fatalf("Failed to start go bin due to: %+v\n", err)
	}
}

// RestartProcess will try to restart a process with procName. Note that this process
// must have been already started through StartGoBin.
func (cli *Cli) RestartProcess(procName string) {
	err := cli.remoteClient.RestartProcess(procName)
	if err != nil {
		log.Fatalf("Failed to restart process due to: %+v\n", err)
	}
}

// StartProcess will try to start a process with procName. Note that this process
// must have been already started through StartGoBin.
func (cli *Cli) StartProcess(procName string) {
	err := cli.remoteClient.StartProcess(procName)
	if err != nil {
		log.Fatalf("Failed to start process due to: %+v\n", err)
	}
}

// StopProcess will try to stop a process named procName.
func (cli *Cli) StopProcess(procName string) {
	err := cli.remoteClient.StopProcess(procName)
	if err != nil {
		log.Fatalf("Failed to stop process due to: %+v\n", err)
	}
}

// DeleteProcess will stop and delete all dependencies from process procName forever.
func (cli *Cli) DeleteProcess(procName string) {
	err := cli.remoteClient.DeleteProcess(procName)
	if err != nil {
		log.Fatalf("Failed to delete process due to: %+v\n", err)
	}
}

// Status will display the status of all procs started through StartGoBin.
func (cli *Cli) Status() {
	procResponse, err := cli.remoteClient.MonitStatus()
	if err != nil {
		log.Fatalf("Failed to get status due to: %+v\n", err)
	}
	maxName := 0
	for id := range procResponse.Procs {
		proc := procResponse.Procs[id]
		maxName = int(math.Max(float64(maxName), float64(len(proc.Name))))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader([]string{"name", "pid", "status", "uptime", "restart", "CPU·%", "memory"})

	for id := range procResponse.Procs {
		proc := procResponse.Procs[id]
		table.Append([]string{
			proc.Name, fmt.Sprintf("%d", proc.Pid), proc.Status.Status, proc.Status.Uptime,
			strconv.Itoa(proc.Status.Restarts), strconv.Itoa(int(proc.Status.Sys.CPU)),
			utils.FormatMemory(int(proc.Status.Sys.Memory)),
		})
	}

	table.SetRowLine(true)
	table.Render()
}

// ProcInfo will display process information
func (cli *Cli) ProcInfo(procName string) {

}
