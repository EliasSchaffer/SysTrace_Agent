package data

type WSEvent struct {
	Type   string `json:"type"`
	Device Device `json:"device"`
}

func (e WSEvent) Metricname() string {
	return "WSEvent"
}

func newWsEvent(typ string, device Device) *WSEvent {
	return &WSEvent{
		Type:   typ,
		Device: device,
	}
}

func (e WSEvent) GetType() string {
	return e.Type
}

func (e WSEvent) GetDevice() Device {
	return e.Device
}

func (e *WSEvent) SetType(t string) {
	e.Type = t
}

func (e *WSEvent) SetDevice(d Device) {
	e.Device = d
}
