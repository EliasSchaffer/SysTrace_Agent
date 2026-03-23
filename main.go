package main

import (
	"SysTrace_Agent/internal/agent"
	"os"
	"runtime"
)

var appAgent = &agent.Agent{}

func main() {
	runtime.LockOSThread()
	go appAgent.StartAgent()
	agent.InitSysTray(Close)
}

func Close() {
	appAgent.StopAgent()
	os.Exit(0)
}
