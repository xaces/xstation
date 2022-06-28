package hook

import "github.com/xaces/xproto"

type Option struct {
	Address string
	Online  string
	Status  string
	Alarm   string
	Event   string
}

type online struct {
	MsgCode  int  `json:"msgCode"`
	DeviceId uint `json:"deviceId"`
	*xproto.Access
}

func newOnline(deviceId uint, a *xproto.Access) *online {
	return &online{MsgCode: 10000, DeviceId: deviceId, Access: a}
}

type status struct {
	MsgCode  int  `json:"msgCode"`
	DeviceId uint `json:"deviceId"`
	*xproto.Status
}

func newStatus(deviceId uint, s *xproto.Status) *status {
	return &status{MsgCode: 10001, DeviceId: deviceId, Status: s}
}

type alarm struct {
	MsgCode  int  `json:"msgCode"`
	DeviceId uint `json:"deviceId"`
	*xproto.Alarm
}

func newAlarm(deviceId uint, a *xproto.Alarm) *alarm {
	return &alarm{MsgCode: 10002, DeviceId: deviceId, Alarm: a}
}

type event struct {
	MsgCode  int  `json:"msgCode"`
	DeviceId uint `json:"deviceId"`
	*xproto.Event
}

func newEvent(deviceId uint, e *xproto.Event) *event {
	return &event{MsgCode: 10003, DeviceId: deviceId, Event: e}
}
