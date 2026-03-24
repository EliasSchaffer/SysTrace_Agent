package agent

import (
	"SysTrace_Agent/internal/data/static"
	"SysTrace_Agent/internal/transport"
	"fmt"
	"strconv"
	"strings"
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

	env := transport.ENVLoader{}
	settings := env.GetSettings()

	defer settingsWindowOpen.Store(false)

	var mw *walk.MainWindow
	var masterServerEdit *walk.LineEdit
	var apiKeyEdit *walk.LineEdit
	var logPathEdit *walk.LineEdit
	var sendGPSCheck *walk.CheckBox
	var staticGPSCheck *walk.CheckBox
	var latitudeEdit *walk.LineEdit
	var longitudeEdit *walk.LineEdit
	var cityEdit *walk.LineEdit
	var regionEdit *walk.LineEdit
	var countryEdit *walk.LineEdit

	window := MainWindow{
		AssignTo: &mw,
		Title:    "SysTrace Settings",
		Size:     Size{560, 520},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "Configure connectivity, logging, and GPS behavior."},
			GroupBox{
				Title:  "General",
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "Master server URL"},
					LineEdit{AssignTo: &masterServerEdit, Text: settings.MASTER_SERVER_URL},
					Label{Text: "Geolocation API key"},
					LineEdit{AssignTo: &apiKeyEdit, Text: settings.GEOLOCATION_API_KEY},
					Label{Text: "Log file path"},
					LineEdit{AssignTo: &logPathEdit, Text: settings.LOGFILE_PATH},
				},
			},
			GroupBox{
				Title:  "GPS",
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "Send GPS"},
					CheckBox{AssignTo: &sendGPSCheck, Checked: settings.SENDGPS},
					Label{Text: "Use static GPS"},
					CheckBox{AssignTo: &staticGPSCheck, Checked: settings.STATICGPS},
					Label{Text: "Latitude"},
					LineEdit{AssignTo: &latitudeEdit, Text: formatFloatForInput(settings.GPS_LATITUDE)},
					Label{Text: "Longitude"},
					LineEdit{AssignTo: &longitudeEdit, Text: formatFloatForInput(settings.GPS_LONGITUDE)},
					Label{Text: "City"},
					LineEdit{AssignTo: &cityEdit, Text: settings.GPS_CITY},
					Label{Text: "Region"},
					LineEdit{AssignTo: &regionEdit, Text: settings.GPS_REGION},
					Label{Text: "Country"},
					LineEdit{AssignTo: &countryEdit, Text: settings.GPS_COUNTRY},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "Save",
						OnClicked: func() {
							newSettings := static.Settings{
								MASTER_SERVER_URL:   strings.TrimSpace(masterServerEdit.Text()),
								GEOLOCATION_API_KEY: strings.TrimSpace(apiKeyEdit.Text()),
								LOGFILE_PATH:        strings.TrimSpace(logPathEdit.Text()),
								SENDGPS:             sendGPSCheck.Checked(),
								STATICGPS:           staticGPSCheck.Checked(),
								GPS_CITY:            strings.TrimSpace(cityEdit.Text()),
								GPS_REGION:          strings.TrimSpace(regionEdit.Text()),
								GPS_COUNTRY:         strings.TrimSpace(countryEdit.Text()),
							}

							if newSettings.MASTER_SERVER_URL == "" {
								walk.MsgBox(mw, "Validation", "Master server URL is required.", walk.MsgBoxIconWarning)
								return
							}

							latitude, err := parseFloatField(latitudeEdit.Text(), "Latitude", -90, 90)
							if err != nil {
								walk.MsgBox(mw, "Validation", err.Error(), walk.MsgBoxIconWarning)
								return
							}

							longitude, err := parseFloatField(longitudeEdit.Text(), "Longitude", -180, 180)
							if err != nil {
								walk.MsgBox(mw, "Validation", err.Error(), walk.MsgBoxIconWarning)
								return
							}

							newSettings.GPS_LATITUDE = latitude
							newSettings.GPS_LONGITUDE = longitude

							if err := env.SetSettings(newSettings); err != nil {
								walk.MsgBox(mw, "Error", fmt.Sprintf("Could not save settings: %v", err), walk.MsgBoxIconError)
								return
							}

							walk.MsgBox(mw, "Saved", "Settings saved successfully.", walk.MsgBoxIconInformation)
						},
					},
					PushButton{
						Text: "Close",
						OnClicked: func() {
							if err := mw.Close(); err != nil {
								fmt.Println("Settings window could not be closed:", err)
							}
						},
					},
				},
			},
		},
	}

	if _, err := window.Run(); err != nil {
		fmt.Println("Settings window could not be started:", err)
	}
}

func formatFloatForInput(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func parseFloatField(rawValue, fieldName string, minValue, maxValue float64) (float64, error) {
	trimmed := strings.TrimSpace(rawValue)
	if trimmed == "" {
		return 0, nil
	}

	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid number", fieldName)
	}

	if parsed < minValue || parsed > maxValue {
		return 0, fmt.Errorf("%s must be between %.2f and %.2f", fieldName, minValue, maxValue)
	}

	return parsed, nil
}
