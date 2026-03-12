package ws

import "SysTrace_Agent/internal/data/static"

type WSEvent struct {
	Type   string        `json:"type"`
	Device static.Device `json:"device"`
}

func (e WSEvent) Metricname() string {
	return "WSEvent"
}

func newWsEvent(typ string, device static.Device) *WSEvent {
	return &WSEvent{
		Type:   typ,
		Device: device,
	}
}

func (e WSEvent) GetType() string {
	return e.Type
}

func (e WSEvent) GetDevice() static.Device {
	return e.Device
}

func (e *WSEvent) SetType(t string) {
	e.Type = t
}

func (e *WSEvent) SetDevice(d static.Device) {
	e.Device = d
}
