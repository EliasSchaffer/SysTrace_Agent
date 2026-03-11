package handler

import (
	"fmt"
	"sync"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func handleNewMessage(message string) {
	showMessageBox(message)
}

func showMessageBox(message string) {

	var mw *walk.MainWindow
	var label *walk.Label
	var once sync.Once

	finish := func(execute bool) {
		once.Do(func() {
			if mw != nil {
				mw.Close()
			}
		})
	}

	window := MainWindow{
		AssignTo: &mw,
		Title:    "SysTrace Agent",
		Size:     Size{300, 150},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				AssignTo: &label,
				Text:     fmt.Sprintf("Server Message: '%s'", message),
			},
			PushButton{
				Text: "Ok",
				OnClicked: func() {
					finish(true)
				},
			},
		},
	}

	window.Run()
}
