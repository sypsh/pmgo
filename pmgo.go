/*
APM is a lightweight process manager written in Golang for Golang applications. It helps you keep all of your applications alive forever, if you want to. You can also reload, start, stop, delete and query status on the fly.

APM also provide a way to start a process by compiling a Golang project source code.

The main APM module is the Master module, it's the glue that keep everything running as it should be.

If you need to use the remote version of APM, take a look at RemoteMaster on Master package.

To use the remote version of APM, use:

- remoteServer := master.StartRemoteMasterServer(dsn, configFile)

It will start a remote master and return the instance.

To make remote requests, use the Remote Client by instantiating using:

- remoteClient, err := master.StartRemoteClient(dsn, timeout)

It will start the remote client and return the instance so you can use to initiate requests, such as:

- remoteClient.StartGoBin(sourcePath, name, keepAlive, args)
*/
package main

// import "github.com/kardianos/osext"
import "gopkg.in/alecthomas/kingpin.v2"
import "github.com/struCoder/pmgo/lib/cli"
import "github.com/struCoder/pmgo/lib/master"

import "github.com/sevlyar/go-daemon"

import "path"
import "path/filepath"
import "syscall"
import "os"
import "os/signal"

import "github.com/Sirupsen/logrus"

var (
	app     = kingpin.New("pmgo", "Aguia Process Manager.")
	dns     = app.Flag("dns", "TCP Dns host.").Default(":9876").String()
	timeout = app.Flag("timeout", "Timeout to connect to client").Default("30s").Duration()

	serveStop           = app.Command("kill", "Kill daemon pmgo.")
	serveStopConfigFile = serveStop.Flag("config-file", "Config file location").String()

	serve           = app.Command("serve", "Create pmgo daemon.")
	serveConfigFile = serve.Flag("config-file", "Config file location").String()

	resurrect     = app.Command("resurrect", "Resurrect all previously save processes.")

	start           = app.Command("start", "start and daemonize an app.")
	startSourcePath = start.Arg("start go file", "go file.").Required().String()
	startName       = start.Arg("name", "Process name.").Required().String()
	startKeepAlive  = start.Flag("keep-alive", "Keep process alive forever.").Required().Bool()
	startArgs       = start.Flag("args", "External args.").Strings()

	restart     = app.Command("restart", "Restart a process.")
	restartName = restart.Arg("name", "Process name.").Required().String()

	// start     = app.Command("start", "Start a process.")
	// startName = start.Arg("name", "Process name.").Required().String()

	stop     = app.Command("stop", "Stop a process.")
	stopName = stop.Arg("name", "Process name.").Required().String()

	delete     = app.Command("delete", "Delete a process.")
	deleteName = delete.Arg("name", "Process name.").Required().String()

	save = app.Command("save", "Save a list of processes onto a file.")

	status = app.Command("list", "Get pmgo list.")

	log = logrus.New()
)

func main() {
	log.Out = os.Stdout
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case serveStop.FullCommand():
		stopRemoteMasterServer()
	case serve.FullCommand():
		startRemoteMasterServer()
	case resurrect.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.Resurrect()
	case start.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.StartGoBin(*startSourcePath, *startName, *startKeepAlive, *startArgs)
	case restart.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.RestartProcess(*restartName)
	// case start.FullCommand():
	// 	cli := cli.InitCli(*dns, *timeout)
	// 	cli.StartProcess(*startName)
	case stop.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.StopProcess(*stopName)
	case delete.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.DeleteProcess(*deleteName)
	case save.FullCommand():
		cli := cli.InitCli(*dns, *timeout)
		cli.Save()
	case status.FullCommand():
		// checkRemoteMasterServer()
		cli := cli.InitCli(*dns, *timeout)
		cli.Status()
	}
}

func isDaemonRunning(ctx *daemon.Context) (bool, *os.Process, error) {
	d, err := ctx.Search()

	if err != nil {
		return false, d, err
	}

	if err := d.Signal(syscall.Signal(0)); err != nil {
		return false, d, err
	}

	return true, d, nil
}

func getCtx() *daemon.Context {
	if *serveConfigFile == "" {
		folderPath := os.Getenv("HOME")
		*serveConfigFile = folderPath + "/.pmgo/config.toml"
		os.MkdirAll(path.Dir(*serveConfigFile), 0755)
	}

	ctx := &daemon.Context{
		PidFileName: path.Join(filepath.Dir(*serveConfigFile), "main.pid"),
		PidFilePerm: 0644,
		LogFileName: path.Join(filepath.Dir(*serveConfigFile), "main.log"),
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}
	return ctx
}

// if RemoteMasterServer not running, just run
func checkRemoteMasterServer() {
	ctx := getCtx()
	if ok, _, _ := isDaemonRunning(ctx); !ok {
		startRemoteMasterServer()
	}	
}

func startRemoteMasterServer() {
	ctx := getCtx()
	if ok, _, _ := isDaemonRunning(ctx); ok {
		log.Info("pmgo daemon is already running.")
		return
	}

	log.Info("daemon started")
	d, err := ctx.Reborn()
	if err != nil {
		log.Fatalf("Failed to reborn daemon due to %+v.", err)
	}

	if d != nil {
		return
	}
	defer ctx.Release()
	
	log.Info("Starting remote master server...")
	remoteMaster := master.StartRemoteMasterServer(*dns, *serveConfigFile)

	sigsKill := make(chan os.Signal, 1)
	signal.Notify(sigsKill,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sigsKill
	log.Info("Received signal to stop...")
	err = remoteMaster.Stop()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func stopRemoteMasterServer() {
	log.Info("pmgo stopping...")
	ctx := getCtx()
	if ok, p, _ := isDaemonRunning(ctx); ok {
		if err := p.Signal(syscall.Signal(syscall.SIGQUIT)); err != nil {
			log.Fatalf("Failed to kill daemon %v", err)
		}
	} else {
		ctx.Release()
		log.Info("instance is not running.")
	}
	log.Info("pmgo daemon terminated")
}
