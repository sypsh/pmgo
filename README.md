<div align="center">
     <a>
        <img src="http://7xjbiz.com1.z0.glb.clouddn.com/github/socJAdzByYtu5maI">
     </a>
     <br/>
     <b>PMGO</b>
     <br/><br/>
</div>


# PMGO 
PMGO is a lightweight process manager written in Golang for Golang applications. It helps you keep your applications alive forever, reload and start them from the source code.

[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) 
[![Join the chat at https://gitter.im/getpmgo/Lobby](https://badges.gitter.im/getpmgo/Lobby.svg)](https://gitter.im/getpmgo/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) 
[![Go Report Card](https://goreportcard.com/badge/github.com/struCoder/pmgo)](https://goreportcard.com/report/github.com/struCoder/pmgo) 
[![GoDoc](https://godoc.org/github.com/struCoder/pmgo?status.svg)](https://godoc.org/github.com/struCoder/pmgo)

Starting an application is easy:
```bash
$ pmgo start source app-name --keep-alive
```

This will basically compile your project source code and start it as a
daemon in the background.

You will probably be able to run anything in any directory, as long as
it is under `GOPATH`

## Install pmgo

```bash
$ go get github.com/struCoder/pmgo
```

## Start pmgo

In order to properly use APM, you always need to start a server. This will be changed in the next version, but in the meantime you need to run the command bellow to start using APM.
```bash
$ pmgo serve
```
If no config file is provided, it will default to a folder '~/.pmgo' where `pmgo` is first started.

## Stop pmgo

```bash
$ pmgo kill
```

## Starting a new application
If it's the first time you are starting a new golang application, you need to tell APM to first build its binary. Then you need to first run:
```bash
$ pmgo start source app-name --keep-alive
```

This will automatically compile, start and daemonize your application. If you need to later on, stop, restart or delete your app from PMGO, you can just run normal commands using the app-name you specified. Example:
```bash
$ pmgo stop app-name
$ pmgo restart app-name
$ pmgo delete app-name
```

## Main features

### Commands overview

```bash
$ pmgo serve
$ pmgo kill

$ pmgo start source app-name --keep-alive                    # Compile, start, daemonize and auto  restart application.
$ pmgo restart app-name                                      # Restart a previously saved process
$ pmgo stop app-name                                         # Stop application.
$ pmgo delete app-name                                       # Delete application forever.

$ pmgo save                                                  # Save current process list
$ pmgo resurrect                                             # Restore previously saved processes

$ pmgo status                                                # Display status for each app.
```

### Managing process via HTTP

You can also use all of the above commands via HTTP requests. Just set the flag ```--dns``` together with ```./pmgo serve``` and then you can use a remote client to start, stop, delete and query status for each app. 
