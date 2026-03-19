package agent

import (
	"fmt"

	"github.com/getlantern/systray"
)

var systemTrayOnClose func()

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
	fmt.Println("Settings clicked")
	//TODO: Settings-Fenster öffnen
}
