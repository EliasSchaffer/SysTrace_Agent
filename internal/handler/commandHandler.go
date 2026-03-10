package handler

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func HandleCommand(command string) {
	switch command {
	case "shutdown":
		println("Shutting down the system...")
		showMessageBox("shutdown", shutdown)
	case "restart":
		println("Restarting the system...")
		restart()
	case "sleep":
		println("Putting the system to sleep...")
		sleep()
	default:
		println("Unknown command:", command)

	}
}

func shutdown() {
	cmd := exec.Command("shutdown", "/s", "/t", "0")
	cmd.Run()
}

func restart() {
	cmd := exec.Command("shutdown", "/r", "/t", "0")
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to restart system:", err)
	}
}

func sleep() {
	cmd := exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0")
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to suspend system:", err)
	}
}

func runHeadlessCountdown(message string, action func()) {
	fmt.Printf("UI unavailable, running headless countdown for command '%s'...\n", message)
	for i := 10; i >= 0; i-- {
		fmt.Printf("Action starting in %d seconds\n", i)
		time.Sleep(time.Second)
	}
	action()
}

func showMessageBox(message string, action func()) {
	var mw *walk.MainWindow
	var label *walk.Label
	var once sync.Once
	shouldExecute := false

	finish := func(execute bool) {
		once.Do(func() {
			shouldExecute = execute
			if mw != nil {
				if err := mw.Close(); err != nil {
					fmt.Println("Failed to close countdown window:", err)
				}
			}
		})
	}

	window := MainWindow{
		AssignTo: &mw,
		Title:    "SysTrace Agent",
		Size:     Size{300, 150},
		MinSize:  Size{250, 120},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				AssignTo: &label,
				Text:     fmt.Sprintf("Die Aktion '%s' wird in 10 Sekunden ausgeführt.", message),
			},
			PushButton{
				Text: "Jetzt ausführen",
				OnClicked: func() {
					finish(true)
				},
			},
			PushButton{
				Text: "Abbrechen",
				OnClicked: func() {
					finish(false)
				},
			},
		},
	}

	if err := window.Create(); err != nil {
		fmt.Println("Failed to create countdown window:", err)
		runHeadlessCountdown(message, action)
		return
	}

	go func() {
		for i := 10; i >= 0; i-- {
			remaining := i
			mw.Synchronize(func() {
				if err := label.SetText(fmt.Sprintf("Die Aktion '%s' wird in %d Sekunden ausgeführt.", message, remaining)); err != nil {
					fmt.Println("Failed to update countdown text:", err)
				}
			})

			if remaining == 0 {
				mw.Synchronize(func() {
					finish(true)
				})
				return
			}

			time.Sleep(time.Second)
		}
	}()

	mw.Run()

	if shouldExecute {
		action()
	}
}
