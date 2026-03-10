package main

import (
	"SysTrace_Agent/internal/agent"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	a := agent.Agent{}
	a.StartAgent()
}
