package agent

import (
	"fmt"
	"sync/atomic"

	"github.com/getlantern/systray"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var systemTrayOnClose func()
var settingsWindowOpen atomic.Bool

func InitSysTray(onClose func()) {
	systemTrayOnClose = onClose
	systray.Run(onReady, onExit)
}

func onReady() {
	//TODO: Icon setzen
	// systray.SetIcon(iconData)

	systray.SetTitle("SysTrace")
	systray.SetTooltip("SysTrace Settings")

	mSettings := systray.AddMenuItem("Settings", "Open settings")
	go func() {
		for {
			<-mSettings.ClickedCh
			openSettings()
		}
	}()

	mQuit := systray.AddMenuItem("Quit", "Exit app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		systemTrayOnClose()
	}()
}

func onExit() {
	fmt.Println("Application exiting...")
}

func openSettings() {
	if !settingsWindowOpen.CompareAndSwap(false, true) {
		fmt.Println("Settings-Fenster ist bereits offen")
		return
	}

	defer settingsWindowOpen.Store(false)

	var mw *walk.MainWindow
	window := MainWindow{
		AssignTo: &mw,
		Title:    "SysTrace Settings",
		Size:     Size{360, 180},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "Settings-Fenster aktiv."},
			Label{Text: "Hier koennen spaeter Konfigurationen bearbeitet werden."},
			PushButton{
				Text: "Schliessen",
				OnClicked: func() {
					if err := mw.Close(); err != nil {
						fmt.Println("Settings-Fenster konnte nicht geschlossen werden:", err)
					}
				},
			},
		},
	}

	if _, err := window.Run(); err != nil {
		fmt.Println("Settings-Fenster konnte nicht gestartet werden:", err)
	}
}
