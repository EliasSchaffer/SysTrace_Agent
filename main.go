package main

import (
	"SysTrace_Agent/internal/agent"
	"os"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	agent.InitSysTray(Close)
	a := agent.Agent{}
	a.StartAgent()
}

func Close() {
	a := agent.Agent{}
	a.StopAgent()
	os.Exit(0)
}
