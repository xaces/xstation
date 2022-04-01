package hook

import "github.com/wlgd/xproto"

type Option struct {
	Name    string
	Address string
	Online  string
	Status  string
	Alarm   string
	Event   string
}

type Online struct {
	DeviceId uint `json:"deviceId"`
	*xproto.Access
}

type Status struct {
	DeviceId uint `json:"deviceId"`
	*xproto.Status
}

type Alarm struct {
	DeviceId uint `json:"deviceId"`
	*xproto.Alarm
}

type Event struct {
	DeviceId uint `json:"deviceId"`
	*xproto.Event
}
