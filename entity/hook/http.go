package hook

import (
	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils"
)

type Http struct {
	urlOnline string
	urlStatus string
	urlAlarm  string
	urlEvent  string
}

func NewHttp(o Option) *Http {
	return &Http{
		urlOnline: o.Address + "/" + o.Online,
		urlStatus: o.Address + "/" + o.Status,
		urlAlarm:  o.Address + "/" + o.Alarm,
		urlEvent:  o.Address + "/" + o.Event,
	}
}

func (h *Http) Online(deviceId uint, a *xproto.Access) {
	xutils.HttpPost(h.urlOnline, &Online{DeviceId: deviceId, Access: a}, nil)
}

func (h *Http) Status(deviceId uint, s *xproto.Status) {
	xutils.HttpPost(h.urlStatus, &Status{DeviceId: deviceId, Status: s}, nil)
}

func (h *Http) Alarm(deviceId uint, a *xproto.Alarm) {
	xutils.HttpPost(h.urlAlarm, &Alarm{DeviceId: deviceId, Alarm: a}, nil)
}

func (h *Http) Event(deviceId uint, e *xproto.Event) {
	xutils.HttpPost(h.urlEvent, &Event{DeviceId: deviceId, Event: e}, nil)
}

func (h *Http) Stop() {
}
